package desc

// Constraint on version requirement
type Constraint int

// Constraints on package deps
// Least to most constraining
const (
	None Constraint = iota
	GTE
	GT
	Equals
	LTE
	LT
)

// Version represents a package version
// want to support package version found on CRAN
// pkgstring <- as.character(package_version(
// unique(as.data.frame(available.packages())$Version)))
// max(sapply(stringr::str_split(pkgstring, "\\."), length))
type Version struct {
	Major int
	Minor int
	Patch int
	Dev   int
	// max amount detected on CRAN was 5
	Other int
}

// Dep represents a dependency
type Dep struct {
	Name       string
	Version    Version
	Constraint Constraint
}

// Desc represents a package description
type Desc struct {
	Package     string
	Source      string
	Version     string
	Maintainer  string
	Description string
	MD5sum      string
	Remotes     []string
	Repository  string
	Imports     map[string]Dep
	Suggests    map[string]Dep
	Depends     map[string]Dep
	LinkingTo   []string
}

// TODO figure out unmarshalling pattern so can
// implement that on Desc so don't need intermediate
// desc struct
type desc struct {
	Package     string
	Source      string
	Version     string
	Maintainer  string
	Description string
	MD5sum      string
	Remotes     []string `delim:"," strip:"\n\r\t "`
	Repository  string
	Imports     []string `delim:"," strip:"\n\r\t "`
	Suggests    []string `delim:"," strip:"\n\r\t "`
	Depends     []string `delim:"," strip:"\n\r\t "`
	LinkingTo   []string `delim:"," strip:"\n\r\t "`
}
