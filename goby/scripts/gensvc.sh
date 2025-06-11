#!/bin/bash

svcName=${1}


cd app/rpc && kitex -module github.com/bitdance-panic/gobuy/app/rpc -service ${svcName} idl/${svcName}.thrift
