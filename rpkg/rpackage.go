package rpkg

import (
	"crypto/md5"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"io"
)

// Hash a tarball
func Hash(fs afero.Fs, tbp string) (string, error) {
	f, err := fs.Open(tbp)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), err
}
