#!/bin/bash

echo "Plugin build: start..."

# Build plugins for the current architecture.
for i in $(find plugins/ -maxdepth 1 -type d); do
    if [[ $i != "plugins/out" ]] && [[ $i != "plugins/" ]]; then
        go build -buildmode=plugin -o plugins/out/${i//plugins\//}.so $i/plugin.go
    fi
done

echo "Plugin build: end."
