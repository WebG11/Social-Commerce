#!/bin/bash

. scripts/list_mod.sh

get_mod_list

readonly root_path=${PWD}

for module in $MODULES; do
    echo "lint ${root_path}/${module}"
    cd "${root_path}/${module}" && golangci-lint run -E gofumpt --path-prefix=. --fix --timeout=5m
done