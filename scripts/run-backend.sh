#! /usr/bin/env bash

set -xe

GIT_ROOT=$(git rev-parse --show-toplevel)

cd $GIT_ROOT/backend
env $(xargs < $GIT_ROOT/.env) go run main.go

