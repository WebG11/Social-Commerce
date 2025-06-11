#!/bin/bash

. scripts/list_mod.sh

get_mod_list

readonly root_path=${PWD}

for module in $MODULES; do
    cd "${root_path}/${module}" && go mod tidy
done