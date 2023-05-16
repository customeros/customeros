#!/bin/sh
rm -rf tmp-sources
mkdir tmp-sources
mkdir tmp-sources/customer-os-common-module
mkdir tmp-sources/events-processing-common
mkdir tmp-sources/events-processing-platform

rsync -av --progress --exclude="tmp-sources" * tmp-sources/events-processing-platform
cp .env tmp-sources/events-processing-platform/.env
cp -r ../customer-os-common-module/* tmp-sources/customer-os-common-module
cp -r ../events-processing-common/* tmp-sources/events-processing-common

cp Dockerfile tmp-sources/Dockerfile

docker build -t events-processing-platform-local tmp-sources/.

rm -rf tmp-sources