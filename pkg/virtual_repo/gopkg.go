package virtual_repo

import (
	"regexp"
	"sort"

	"github.com/Masterminds/semver/v3"
)

var packageRefRegex = regexp.MustCompile(`^refs/tags/(.*)/v(.*)$`)

func (v *VirtualRepo) ListPackages() (map[string]*GoPackage, error) {
	refs, err := v.FilterRefs("refs/tags/")
	if err != nil {
		return nil, err
	}

	packages := map[string]*GoPackage{}
	for _, ref := range refs {
		tag := ref.Name().String()
		matches := packageRefRegex.FindStringSubmatch(tag)

		if len(matches) != 3 {
			continue
		}

		packageName := matches[1]
		versionText := matches[2]

		if _, ok := packages[packageName]; !ok {
			packages[packageName] = &GoPackage{
				Name:     packageName,
				Ref:      tag,
				Versions: []*semver.Version{},
			}
		}

		if version, err := semver.NewVersion(versionText); err != nil {
			continue
		} else {
			packages[packageName].Versions = append(packages[packageName].Versions, version)
		}
	}

	for _, pkg := range packages {
		sort.Sort(semver.Collection(pkg.Versions))
	}

	return packages, nil
}

type GoPackage struct {
	Name     string
	Ref      string
	Versions []*semver.Version
}

func (p *GoPackage) GetLatestVersion() *semver.Version {
	return p.Versions[len(p.Versions)-1]
}

func (p *GoPackage) NextVersion(upgradeType ...string) *semver.Version {
	if len(upgradeType) > 1 {
		panic("only one upgrade type is allowed")
	}

	latestVersion := p.GetLatestVersion()

	nextVersion := latestVersion.IncPatch()

	switch upgradeType[0] {
	case "patch":
		nextVersion = latestVersion.IncPatch()
	case "minor":
		nextVersion = latestVersion.IncMinor()
	case "major":
		nextVersion = latestVersion.IncMajor()
	}

	return &nextVersion
}
