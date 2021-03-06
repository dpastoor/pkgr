package pacman

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/dpastoor/goutils"
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/desc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// GetPriorInstalledPackages ...
func GetPriorInstalledPackages(fileSystem afero.Fs, libraryPath string) map[string]desc.Desc {
	installed := make(map[string]desc.Desc)
	installedLibrary, err := fileSystem.Open(libraryPath)

	if err != nil {
		log.WithField("libraryPath", libraryPath).Warn("package library not found at given library path")
		return installed
	}
	defer installedLibrary.Close()

	fileInfo, _ := installedLibrary.Readdir(0)
	installedPackageFolders := goutils.ListDirNames(fileInfo)

	for _, pkgFolder := range installedPackageFolders {
		descriptionFilePath := filepath.Join(libraryPath, pkgFolder, "DESCRIPTION")
		installedPackage, err := scanInstalledPackage(descriptionFilePath, fileSystem)

		if err != nil {
			log.Error(err)
			err = nil
		} else {
			log.WithFields(log.Fields{
				"package":     installedPackage.Package,
				"version":     installedPackage.Version,
				"source repo": installedPackage.Repository,
			}).Debug("found installed package")

			installed[installedPackage.Package] = installedPackage
		}
	}

	return installed
}

// GetInstallers returns the installers for the installed packages
func GetInstallers(ip map[string]desc.Desc) InstalledFromPkgs {
	var pkgr, packrat, unknown []string
	for k, v := range ip {
		if v.PkgrVersion == "" {
			packrat = append(packrat, k)
		} else {
			pkgr = append(pkgr, k)
		}
	}
	return InstalledFromPkgs{
		Pkgr:    pkgr,
		Packrat: packrat,
		Unknown: unknown,
	}

}

// GetPackagesByInstalledFrom returns InstalledFromPkgs structure
// single location where business rule of "not pkgr" is applied
func GetPackagesByInstalledFrom(fileSystem afero.Fs, libraryPath string) InstalledFromPkgs {
	var pkgr, packrat, unknown []string
	ip := GetPriorInstalledPackages(fileSystem, libraryPath)
	for k, v := range ip {
		if v.PkgrVersion == "" {
			packrat = append(packrat, k)
		} else {
			pkgr = append(pkgr, k)
		}
	}
	return InstalledFromPkgs{
		Pkgr:    pkgr,
		Packrat: packrat,
		Unknown: unknown,
	}
}

func scanInstalledPackage(
	descriptionFilePath string, fileSystem afero.Fs) (desc.Desc, error) {

	descriptionFile, err := fileSystem.Open(descriptionFilePath)

	if err != nil {
		log.WithField("file", descriptionFilePath).Warn("DESCRIPTION missing from installed package.")
		return desc.Desc{}, err
	}
	defer descriptionFile.Close()

	log.WithField("description_file", descriptionFilePath).Trace("scanning DESCRIPTION file")

	installedPackage, err := desc.ParseDesc(descriptionFile)

	if installedPackage.Package == "" {
		err = errors.New(fmt.Sprintf("encountered a description file without package name: %s", descriptionFile))
		log.WithField("description file", descriptionFilePath).Error(err)
		return desc.Desc{}, err
	}

	return installedPackage, nil
}

// GetOutdatedPackages returns a list of outdated packages
func GetOutdatedPackages(installed map[string]desc.Desc, availablePackages []cran.PkgDl) []cran.OutdatedPackage {
	var outdatedPackages []cran.OutdatedPackage

	for _, pkgDl := range availablePackages {

		pkgName := pkgDl.Package.Package
		availableVersion := pkgDl.Package.Version

		if installedPkg, found := installed[pkgName]; found {

			// If available version is later than currently installed version
			if desc.CompareVersionStrings(availableVersion, installedPkg.Version) > 0 {
				outdatedPackages = append(outdatedPackages, cran.OutdatedPackage{
					Package:    pkgName,
					OldVersion: installed[pkgName].Version,
					NewVersion: pkgDl.Package.Version,
				})
			}
		}
	}
	return outdatedPackages
}

// InstalledFromPkgs ...
type InstalledFromPkgs struct {
	Pkgr    []string `json:"pkgr"`
	Packrat []string `json:"packrat"`
	Unknown []string `json:"unknown"`
}

// NotPkgr returns a list of packages not installed by Pkgr
func (ip *InstalledFromPkgs) NotFromPkgr() []string {
	var list []string
	for _, p := range ip.Packrat {
		list = append(list, p)
	}
	for _, p := range ip.Unknown {
		list = append(list, p)
	}
	return list
}

// IsPkgr returns a list of packages installed by Pkgr
func (ip *InstalledFromPkgs) FromPkgr() []string {
	return ip.Pkgr
}
