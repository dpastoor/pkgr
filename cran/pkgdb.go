package cran

import (
	"errors"
	"fmt"

	"github.com/dpastoor/rpackagemanager/desc"
)

// NewPkgDb returns a new package database
func NewPkgDb(urls []RepoURL, dst SourceType, cfgdb map[string]PkgConfig) (*PkgDb, error) {
	db := PkgDb{
		Config:            cfgdb,
		DefaultSourceType: dst,
	}
	if len(urls) == 0 {
		return &db, errors.New("Package database must contain at least one RepoUrl")
	}
	for _, url := range urls {
		rdb, err := NewRepoDb(url, dst)
		if err != nil {
			return &db, err
		}
		db.Db = append(db.Db, rdb)
	}
	return &db, nil
}

// SetPackageRepo sets a package repository so querying the package will
// pull from that repo
func (p *PkgDb) SetPackageRepo(pkg string, repo string) error {
	for _, r := range p.Db {
		if r.Repo.Name == repo {
			cfg := p.Config[pkg]
			cfg.Repo = r.Repo
			p.Config[pkg] = cfg
			return nil
		}
	}
	return fmt.Errorf("no repo: %s, detected containing package: %s", repo, pkg)
}

func pkgExists(pkg string, db map[string]desc.Desc) bool {
	_, exists := db[pkg]
	return exists
}
func pkgExistsInRepo(pkg string, dbs map[SourceType]map[string]desc.Desc) bool {
	exists := false
	for _, db := range dbs {
		_, exists = db[pkg]
		if exists {
			return exists
		}
	}
	return exists
}

func isCorrectRepo(pkg string, r RepoURL, cfg map[string]PkgConfig) bool {
	pkgcfg, exists := cfg[pkg]
	if exists {
		if pkgcfg.Repo.Name == r.Name {
			return true
		} else {
			return false
		}
	}
	return true
}

// GetPackage gets a package from the package database, returning the first match
func (p *PkgDb) GetPackage(pkg string) (desc.Desc, PkgConfig, bool) {
	st := p.Config[pkg].Type
	for _, db := range p.Db {
		// For now package existence is checked exactly as the package is specified
		// in the config. Eg, if specifies binary, will only check binary version
		// the checking if also exists as source or otherwise should occur upstream
		// then be set as part of the explicit configuration.
		if pkgExists(pkg, db.Dbs[st]) && isCorrectRepo(pkg, db.Repo, p.Config) {
			return db.Dbs[st][pkg], PkgConfig{Repo: db.Repo, Type: st}, true
		}
	}
	return desc.Desc{}, PkgConfig{}, false
}

// GetPackageFromRepo gets a package from a repo in the package database
func (p *PkgDb) GetPackageFromRepo(pkg string, repo string) (desc.Desc, PkgConfig, bool) {
	st := p.Config[pkg].Type
	for _, db := range p.Db {
		if repo != "" && db.Repo.Name != repo {
			continue
		}
		if pkgExists(pkg, db.Dbs[st]) {
			return db.Dbs[st][pkg], PkgConfig{Repo: db.Repo, Type: st}, true
		}
	}
	return desc.Desc{}, PkgConfig{}, false
}

// GetPackages returns all packages and the repo that they
// will be acquired from, as well as any missing packages
func (p *PkgDb) GetPackages(pkgs []string) AvailablePkgs {
	ap := AvailablePkgs{}
	for _, pkg := range pkgs {
		pd, cfg, found := p.GetPackage(pkg)
		ap.Packages = append(ap.Packages, PkgDl{
			Package: pd,
			Config:  cfg,
		})
		if !found {
			ap.Missing = append(ap.Missing, pkg)
		}
	}
	return ap
}

// CheckAllAvailable returns whether all requested packages
// are available in the package database. It is a simple wrapper
// around GetPackages
func (p *PkgDb) CheckAllAvailable(pkgs []string) bool {
	ap := p.GetPackages(pkgs)
	if len(ap.Missing) > 0 {
		return false
	}
	return true
}

// GetAllPkgsByName returns all packages in the database
func (p *PkgDb) GetAllPkgsByName() []string {
	// use map so will remove duplicate packages
	pkgMap := make(map[string]bool)
	for _, db := range p.Db {
		for st := range db.Dbs {
			for pkg := range db.Dbs[st] {
				pkgMap[pkg] = true
			}
		}
	}
	pkgs := []string{}
	for pkg := range pkgMap {
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}
