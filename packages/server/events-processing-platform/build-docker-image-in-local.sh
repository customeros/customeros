#!/bin/sh
rm -rf tmp-sources
mkdir tmp-sources
mkdir tmp-sources/customer-os-common-module
mkdir tmp-sources/customer-os-common-ai
mkdir tmp-sources/customer-os-neo4j-repository
mkdir tmp-sources/events-processing-proto
mkdir tmp-sources/events-processing-platform

rsync -av --progress --exclude="tmp-sources" * tmp-sources/events-processing-platform
cp .env tmp-sources/events-processing-platform/.env
cp -r ../customer-os-common-ai/* tmp-sources/customer-os-common-ai
cp -r ../customer-os-common-module/* tmp-sources/customer-os-common-module
cp -r ../customer-os-neo4j-repository/* tmp-sources/customer-os-neo4j-repository
cp -r ../events-processing-proto/* tmp-sources/events-processing-proto

cp Dockerfile tmp-sources/Dockerfile

docker build -t events-processing-platform-local tmp-sources/.

rm -rf tmp-sources