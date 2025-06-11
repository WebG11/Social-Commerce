#!/bin/bash

. scripts/list_mod.sh
get_mod_list

# 遍历每个模块并运行 go test
for module in $MODULES; do
    echo "Running tests for module: $module"
    (cd "$module" && go test -v ./...)
done