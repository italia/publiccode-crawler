#!/bin/bash

echo "Building plugins... "

PLUGINS_FOLDER=plugins

#go build -buildmode=plugin -o plugins/out/github.so plugins/github/plugin.go
#go build -buildmode=plugin -o plugins/out/gitlab.so plugins/gitlab/plugin.go
#go build -buildmode=plugin -o plugins/out/bitbucket.so plugins/bitbucket/plugin.go

for i in $(find plugins/ -maxdepth 1 -type d); do
    if [[ $i != "plugins/out" ]] && [[ $i != "plugins/" ]]; then
        go build -buildmode=plugin -o plugins/out/${i//plugins\//}.so $i/plugin.go
    fi
done

echo "Plugin build: end."
