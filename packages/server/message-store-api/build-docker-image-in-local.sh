#!/bin/sh
rm -rf tmp-sources
mkdir tmp-sources
mkdir tmp-sources/customer-os-common-module
mkdir tmp-sources/events-processing-common
mkdir tmp-sources/customer-os-api
mkdir tmp-sources/message-store-api

rsync -av --progress --exclude="tmp-sources" * tmp-sources/message-store-api
cp .env tmp-sources/message-store-api/.env
cp -r ../customer-os-common-module/* tmp-sources/customer-os-common-module
cp -r ../events-processing-common/* tmp-sources/events-processing-common
cp -r ../customer-os-api/* tmp-sources/customer-os-api

cp Dockerfile tmp-sources/Dockerfile

docker build -t aa tmp-sources/.

rm -rf tmp-sources