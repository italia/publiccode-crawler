package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/italia/publiccode-crawler/v3/common"
	"github.com/italia/publiccode-crawler/v3/metrics"
	"github.com/spf13/viper"
)

// CloneRepository clone the repository into DATADIR/repos/<hostname>/<vendor>/<repo>/gitClone.
func CloneRepository(hostname, name, gitURL, index string) error {
	if name == "" {
		return errors.New("cannot save a file without name")
	}
	if gitURL == "" {
		return errors.New("cannot clone a repository without git URL")
	}

	vendor, repo := common.SplitFullName(name)
	path := filepath.Join(viper.GetString("CRAWLER_DATADIR"), "repos", hostname, vendor, repo, "gitClone")

	// If folder already exists it will do a fetch instead of a clone.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		out, err := exec.Command("git", "-C", path, "fetch", "--all").CombinedOutput()
		if err != nil {
			return fmt.Errorf("cannot git pull the repository: %s: %s", err.Error(), out)
		}

		return nil
	}

	out, err := exec.Command("git", "clone", "--filter=blob:none", "--mirror", gitURL, path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot git clone the repository: %s: %s", err.Error(), out)
	}

	metrics.GetCounter("repository_cloned", index).Inc()

	return err
}
