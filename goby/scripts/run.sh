#!/bin/bash

svcName=${1}


if [ -d "app/services/${svcName}" ];then
    cd app/services/${svcName} && go run .
fi