package rcmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/stretchr/testify/assert"
)

func TestInstallArgs(t *testing.T) {
	assert := assert.New(t)

	var installArgsTests = []struct {
		in       InstallArgs
		expected []string
	}{
		{
			InstallArgs{},
			[]string{},
		},
		{
			InstallArgs{Clean: true},
			[]string{"--clean"},
		},
		{
			InstallArgs{Library: "path/to/lib"},
			[]string{"--library=path/to/lib"},
		},
		{
			NewDefaultInstallArgs(),
			[]string{"--build", "--install-tests", "--no-multiarch", "--with-keep.source"},
		},
	}
	for i, tt := range installArgsTests {
		actual := tt.in.CliArgs()
		assert.Equal(tt.expected, actual, fmt.Sprintf("test num: %v", i+1))
	}
}

// ATTENTION:
// This test is misconfigured, it shouldn't be pointing to the integration test folders as those folders are only valid
// after make test-install has been run from the integration_tests folder.
// This test might fail falsey depending on the current state of pkgr/integration_tests/simple/test-library
func TestUpdateDcfFile(t *testing.T) {
	var tests = []struct {
		filename     string
		version      string
		installType  string
		repoURL      string
		repo         string
		expectedRepo string
		message      string
	}{
		{
			filename:     "../integration_tests/simple/test-library/R6/Description",
			version:      "version",
			installType:  "binary",
			repoURL:      "myURL",
			repo:         "CRAN",
			expectedRepo: "CRAN",
			message:      "R6 test",
		},
		{
			filename:     "../integration_tests/simple/test-library/pillar/Description",
			version:      "1.2.3",
			installType:  "binary",
			repoURL:      "www.myURL.com",
			repo:         "GitHub",
			expectedRepo: "CRAN GitHub",
			message:      "pillar test",
		},
	}

	for _, tt := range tests {

		dcf, err := updateDescriptionInfo(tt.filename, tt.version, tt.installType, tt.repoURL, tt.repo, false)
		installedPackage, _ := desc.ParseDesc(bytes.NewReader(dcf))

		assert.Equal(t, nil, err, fmt.Sprintf("Error: %s", err))
		assert.Equal(t, tt.expectedRepo, installedPackage.Repository, fmt.Sprintf("Failed: %s", tt.message))
		assert.Equal(t, tt.version, installedPackage.PkgrVersion, fmt.Sprintf("Failed: %s", tt.message))
		assert.Equal(t, tt.repoURL, installedPackage.PkgrRepositoryURL, fmt.Sprintf("Failed: %s", tt.message))
		assert.Equal(t, tt.installType, installedPackage.PkgrInstallType, fmt.Sprintf("Failed: %s", tt.message))
	}
}
