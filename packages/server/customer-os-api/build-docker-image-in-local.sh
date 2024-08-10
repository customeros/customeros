#!/bin/sh
rm -rf tmp-sources
mkdir tmp-sources
mkdir tmp-sources/customer-os-common-module
mkdir tmp-sources/customer-os-neo4j-repository
mkdir tmp-sources/validation-api
mkdir tmp-sources/enrichment-api
mkdir tmp-sources/customer-os-api

rsync -av --progress --exclude="tmp-sources" * tmp-sources/customer-os-api
cp .env tmp-sources/customer-os-api/.env
cp -r ../customer-os-common-module/* tmp-sources/customer-os-common-module
cp -r ../customer-os-neo4j-repository/* tmp-sources/customer-os-neo4j-repository
cp -r ../validation-api/* tmp-sources/validation-api
cp -r ../enrichment-api/* tmp-sources/enrichment-api

cp Dockerfile tmp-sources/Dockerfile

docker build -t aa tmp-sources/.

rm -rf tmp-sources