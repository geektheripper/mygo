package virtual_repo

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
)

type VirtualRepo struct {
	remoteURL string
	branch    string

	auth transport.AuthMethod

	repo   *git.Repository
	remote *git.Remote
	wt     *git.Worktree

	preCloseFuncs []func()
}

func NewVirtualRepo(remoteURL, branchName string) (*VirtualRepo, error) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	v := &VirtualRepo{
		remoteURL:     remoteURL,
		branch:        branchName,
		preCloseFuncs: []func(){},
	}

	mfs := memfs.New()

	repo, err := git.InitWithOptions(memory.NewStorage(), mfs, git.InitOptions{
		DefaultBranch: plumbing.NewBranchReferenceName(v.branch),
	})

	if err != nil {
		return nil, err
	}

	wt, err := repo.Worktree()

	if err != nil {
		return nil, err
	}

	remote, err := repo.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{v.remoteURL}})

	if err != nil {
		return nil, err
	}

	v.repo = repo
	v.remote = remote
	v.wt = wt

	err = v.EnsureAuth()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (v *VirtualRepo) Close() {
	for _, f := range v.preCloseFuncs {
		f()
	}
}

func (v *VirtualRepo) SetAuth(auth transport.AuthMethod) *VirtualRepo {
	v.auth = auth
	return v
}

func (v *VirtualRepo) FilterRefs(prefixes ...string) ([]*plumbing.Reference, error) {
	refs, err := v.remote.List(&git.ListOptions{Auth: v.auth})
	if err != nil {
		return nil, err
	}

	if len(prefixes) == 0 {
		return refs, nil
	}

	versions := make([]*plumbing.Reference, 0)
	for _, ref := range refs {
		rfsName := ref.Name().String()
		for _, prefix := range prefixes {
			if strings.HasPrefix(rfsName, prefix) {
				versions = append(versions, ref)
			}
		}
	}

	return versions, nil
}

func (v *VirtualRepo) CreateFile(path string, reader io.Reader, perm os.FileMode) error {
	wt, err := v.repo.Worktree()

	if err != nil {
		return err
	}

	file, err := wt.Filesystem.OpenFile(path, os.O_CREATE|os.O_WRONLY, perm)

	if err != nil {
		return err
	}

	_, err = io.Copy(file, reader)

	if err != nil {
		return err
	}

	err = file.Close()

	if err != nil {
		return err
	}

	return nil
}

func (v *VirtualRepo) CopyFile(from, to string, perm os.FileMode) error {
	sourceFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	return v.CreateFile(to, sourceFile, perm)
}

type ImportResult struct {
	Count int
	Size  uint64
}

func (r *ImportResult) Add(file os.FileInfo) {
	r.Count++
	r.Size += uint64(file.Size())
}

func (v *VirtualRepo) Import(from, to string) (*ImportResult, error) {
	result := &ImportResult{Count: 0, Size: 0}

	info, err := os.Stat(from)
	if err != nil {
		return result, err
	}

	if !info.IsDir() {
		err := v.CopyFile(from, to, info.Mode())
		if err != nil {
			return result, err
		}

		result.Add(info)
		return result, nil
	}

	err = filepath.Walk(from, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(from, path)
		if err != nil {
			return err
		}

		err = v.CopyFile(path, filepath.Join(to, rel), info.Mode())
		if err != nil {
			return err
		}

		result.Add(info)
		return nil
	})

	return result, err
}

func (v *VirtualRepo) Publish(tagName string, message string) error {
	err := v.wt.AddGlob(".")

	if err != nil {
		return err
	}

	commit, err := v.wt.Commit(message, &git.CommitOptions{All: true})

	if err != nil {
		return err
	}

	// tag
	tag, err := v.repo.CreateTag(tagName, commit, nil)

	if err != nil {
		return err
	}

	ref := tag.Name().String()
	err = v.repo.Push(&git.PushOptions{
		RemoteURL: v.remoteURL,
		Progress:  os.Stdout,
		RefSpecs:  []config.RefSpec{config.RefSpec(ref + ":" + ref)},
		Auth:      v.auth,
	})

	if err != nil {
		return err
	}

	return nil
}
