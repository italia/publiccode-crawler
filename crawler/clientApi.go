package crawler

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"plugin"

	log "github.com/sirupsen/logrus"
)

// Crawler is the interface for crawler plugins.
type Crawler interface {
	Register() Handler
	GetId() string
}

var (
	clientApis map[string]Handler
)

const (
	pluginsDir = "plugins/out"
)

// RegisterCrawlers registers all founded crawler plugins.
func RegisterCrawlers() {
	clientApis = make(map[string]Handler)

	files, err := getPluginFiles()
	if err != nil {
		log.Error(err)
		panic(err)
	}

	for _, file := range files {
		err := registerPlugin(file)
		if err != nil {
			log.Error(err)
		}
	}
}

// GetClientApiCrawler returns the handler func to process domain.
func GetClientApiCrawler(clientApi string) (Handler, error) {
	if crawler, ok := clientApis[clientApi]; ok {
		return crawler, nil
	} else {
		return nil, errors.New(fmt.Sprintf("no client found for %s", clientApi))
	}
}

// GetPlugins returns a list of all registered plugins.
func GetPlugins() map[string]Handler {
	return clientApis
}

func getPluginFiles() ([]string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var files []string
	absolutePluginsDir := path.Join(workDir, pluginsDir)

	filepath.Walk(absolutePluginsDir, func(p string, f os.FileInfo, _ error) error {
		if filepath.Ext(p) == ".so" {
			files = append(files, path.Join(absolutePluginsDir, f.Name()))
		}
		return nil
	})

	return files, nil
}

func registerPlugin(file string) error {
	plug, err := plugin.Open(file)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to open file %s: %v", file, err))
	}

	symbol, err := plug.Lookup("Plugin")
	if err != nil {
		return errors.New(fmt.Sprintf("unable to lookup Plugin symbol: %v", err))
	}

	var crawler Crawler
	crawler, ok := symbol.(Crawler)
	if !ok {
		return errors.New("unexpected type from module symbol")
	}

	clientApis[crawler.GetId()] = crawler.Register()
	return nil
}
