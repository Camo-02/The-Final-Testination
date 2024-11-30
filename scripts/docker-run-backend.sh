#! /usr/bin/env sh

set -xe
#TODO: actually wait for the db to be ready, as this may still fail sometimes.
echo "Waiting for the database to be ready..."
sleep 10

/bin/backend
