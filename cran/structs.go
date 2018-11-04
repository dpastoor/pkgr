package cran

import (
	"time"

	"github.com/dpastoor/rpackagemanager/desc"
)

// RepoURL represents the URL and name for a repo
// to match the R convention of specifying a repository name
// CRAN = https://cran.rstudio.com would be
// RepoUrl{URL: "https://cran.rstudio.com", Name: "CRAN"}
type RepoURL struct {
	URL  string
	Name string
}

// RepoDb represents a Db
type RepoDb struct {
	Db   map[string]desc.Desc
	Time time.Time
	Repo RepoURL
}

// PkgDb represents a package database
type PkgDb struct {
	Db []*RepoDb
}