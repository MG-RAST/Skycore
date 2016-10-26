#!/bin/sh

set -x
set -e

echo "####### go fmt Skycore     #######"
go fmt github.com/MG-RAST/Skycore/skycore

echo "####### go fix Skycore     #######"
go fix github.com/MG-RAST/Skycore/skycore

echo "####### go install Skycore #######"
go install -v github.com/MG-RAST/Skycore/skycore
