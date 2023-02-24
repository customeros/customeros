#!/bin/sh
rm -rf tmp-sources
mkdir tmp-sources
mkdir tmp-sources/customer-os-common-module
mkdir tmp-sources/validation-api

rsync -av --progress --exclude="tmp-sources" * tmp-sources/validation-api
cp .env tmp-sources/validation-api/.env
cp -r ../customer-os-common-module/* tmp-sources/customer-os-common-module

cp Dockerfile tmp-sources/Dockerfile

docker build -t aa tmp-sources/.

rm -rf tmp-sources