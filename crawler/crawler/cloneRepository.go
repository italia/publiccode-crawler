package crawler

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/italia/developers-italia-backend/crawler/metrics"
	"github.com/spf13/viper"
)

// CloneRepository clone the repository into DATADIR/repos/<hostname>/<vendor>/<repo>/gitClone
func CloneRepository(domain Domain, hostname, name, gitURL, gitBranch, index string) error {
	if domain.Host == "" {
		return errors.New("cannot save a file without domain host")
	}
	if name == "" {
		return errors.New("cannot save a file without name")
	}
	if gitURL == "" {
		return errors.New("cannot clone a repository without git URL")
	}

	vendor, repo := splitFullName(name)
	path := filepath.Join(viper.GetString("CRAWLER_DATADIR"), "repos", hostname, vendor, repo, "gitClone")

	// If folder already exists it will do a fetch instead of a clone.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		//	Command is: git fetch --all
		cmd := exec.Command("git", "-C", path, "fetch", "--all") // nolint: gas
		err := cmd.Run()
		if err != nil {
			return errors.New("cannot git pull the repository: " + err.Error())
		}
		// Command is: git reset --hard origin/<branch_name>
		cmd = exec.Command("git", "-C", path, "reset", "--hard", "origin/"+gitBranch) // nolint: gas
		err = cmd.Run()
		if err != nil {
			return errors.New("cannot git pull the repository: " + err.Error())
		}
		return err
	}

	// Clone the repository using the external command "git".
	// Command is: git clone -b <branch> <remote_repo>
	cmd := exec.Command("git", "clone", "-b", gitBranch, gitURL, path) // nolint: gas
	err := cmd.Run()
	if err != nil {
		return errors.New("cannot git clone the repository: " + err.Error())
	}

	metrics.GetCounter("repository_cloned", index).Inc()
	return err
}
