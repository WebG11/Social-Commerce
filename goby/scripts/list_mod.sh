#!/bin/bash

# service_list=()

# ls $PWD/app/services
# get_app_list(){
#     local idx=0
#     for d in app/services/*; do
#         if [ -d "$d" ] ; then
#             service_list[idx]=$d
#             idx+=1
#         fi
#     done
# }

# echo $service_list

MODULES=();

get_mod_list(){
    MODULES=$(find . -name "go.mod" -exec dirname {} \;)
}
