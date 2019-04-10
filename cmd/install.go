// Copyright © 2018 Devin Pastoor <devin.pastoor@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/spf13/afero"
	"path/filepath"
	"time"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/logger"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)


var updateArgument bool

// installCmd represents the R CMD install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install a package",
	Long: `
	install a package
 `,
	RunE: rInstall,
}

func init() {
	installCmd.Flags().BoolVar(&updateArgument, "update", false, "Update outdated packages during installation.")
	RootCmd.AddCommand(installCmd)
}

func rInstall(cmd *cobra.Command, args []string) error {

	//Init install-specific log, if one has been set. This overwrites the default log.
	if cfg.Logging.Install != "" {
		logger.AddLogFile(cfg.Logging.Install, cfg.Logging.Overwrite)
	} else {
		logger.AddLogFile(cfg.Logging.All, cfg.Logging.Overwrite)
	}

	startTime := time.Now()
	rs := rcmd.NewRSettings()
	rVersion := rcmd.GetRVersion(&rs)
	log.Infoln("R Version " + rVersion.ToFullString())
	cdb, installPlan := planInstall(rVersion)

	var packageUpdateInfo []UpdateAttempt
	if updateArgument {
		log.Info("update argument passed. staging packages for update...")
		packageUpdateInfo = tagOldInstallations(fs, cfg.Library, installPlan.OutdatedPackages)
	}

	var toDl []cran.PkgDl
	// starting packages
	for _, p := range installPlan.StartingPackages {
		pkg, cfg, _ := cdb.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Config: cfg})
	}
	// all other packages
	for p := range installPlan.DepDb {
		pkg, cfg, _ := cdb.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Config: cfg})
	}
	// // want to download the packages and return the full path of any downloaded package
	pc := rcmd.NewPackageCache(userCache(cfg.Cache), false)
	dl, err := cran.DownloadPackages(fs, toDl, pc.BaseDir, rVersion)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	ia := rcmd.NewDefaultInstallArgs()
	ia.Library, _ = filepath.Abs(cfg.Library)
	nworkers := getWorkerCount()

	// leave at least 1 thread open for coordination, given more than 2 threads available.
	// if only 2 available, will let the OS hypervisor coordinate some else would drop the
	// install time too much for the little bit of additional coordination going on.
	pkgCustomizations := cfg.Customizations.Packages
	for n, v := range pkgCustomizations {
		if v.Env != nil {
			rs.PkgEnvVars[n] = v.Env
		}
	}
	err = rcmd.InstallPackagePlan(fs, installPlan, dl, pc, ia, rs, rcmd.ExecSettings{}, nworkers)
	if err != nil {
		fmt.Println("failed package install")
		fmt.Println(err)
	}

	restoreUnupdatedPackages(fs, packageUpdateInfo)

	fmt.Println("duration:", time.Since(startTime))
	return nil
}

func tagOldInstallations(fileSystem afero.Fs, libraryPath string, outdatedPackages []gpsr.OutdatedPackage) []UpdateAttempt {
	var updateAttempts []UpdateAttempt

	//Tag each outdated pkg directory in library with "__OLD__"
	for _, pkg := range outdatedPackages {
		updateAttempts = append(updateAttempts, tagOldInstallation(fileSystem, libraryPath, pkg))
	}

	return updateAttempts
}

func tagOldInstallation(fileSystem afero.Fs, libraryPath string, outdatedPackage gpsr.OutdatedPackage) UpdateAttempt {
	packageDir := filepath.Join(libraryPath, outdatedPackage.Package)
	taggedPackageDir := filepath.Join(libraryPath, "__OLD__" + outdatedPackage.Package)
	err := fileSystem.Rename(packageDir, taggedPackageDir)

	if err != nil {
		log.WithField("package dir", packageDir).Warn("error when backing up outdated package")
		log.Error(err)
	}

	return UpdateAttempt{
		Package: outdatedPackage.Package,
		ActivePackageDirectory: packageDir,
		BackupPackageDirectory: taggedPackageDir,
		OldVersion: outdatedPackage.OldVersion,
		NewVersion: outdatedPackage.NewVersion,
	}
}

func restoreUnupdatedPackages(fileSystem afero.Fs, packageBackupInfo []UpdateAttempt) {

	if len(packageBackupInfo) == 0 {
		return
	}

	//libraryDirectoryFsObject, _ := fs.Open(libraryPath)
	//packageFolderObjects, _ := libraryDirectoryFsObject.Readdir(0)

	for _, info := range packageBackupInfo {
		_, err := fileSystem.Stat(info.ActivePackageDirectory) //Checking existence
		if err == nil {

			fileSystem.RemoveAll(info.BackupPackageDirectory)

			log.WithFields(log.Fields{
				"pkg": info.Package,
				"old_version": info.OldVersion,
				"new_version": info.NewVersion,
			}).Info("successfully updated package")

		} else {
			log.WithFields(log.Fields{
				"pkg": info.Package,
				"old_version": info.OldVersion,
				"new_version": info.NewVersion,
			}).Warn("could not update package, restoring last-installed version")
			fileSystem.Rename(info.BackupPackageDirectory, info.ActivePackageDirectory)
		}
	}
}

type UpdateAttempt struct {
	Package string
	ActivePackageDirectory string
	BackupPackageDirectory string
	OldVersion string
	NewVersion string
}


