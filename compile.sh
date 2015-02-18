#!/bin/sh

set -x
set -e

echo "####### go fmt Skycore     #######"
go fmt github.com/wgerlach/Skycore/skycore

echo "####### go fix Skycore     #######"
go fix github.com/wgerlach/Skycore/skycore

echo "####### go install Skycore #######"
go install -v github.com/wgerlach/Skycore/skycore
