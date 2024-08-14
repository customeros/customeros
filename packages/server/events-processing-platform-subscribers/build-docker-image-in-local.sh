#!/bin/sh
rm -rf tmp-sources
mkdir tmp-sources
mkdir tmp-sources/customer-os-common-module
mkdir tmp-sources/customer-os-common-ai
mkdir tmp-sources/customer-os-neo4j-repository
mkdir tmp-sources/events-processing-proto
mkdir tmp-sources/events-processing-platform
mkdir tmp-sources/events-processing-platform-subscribers
mkdir tmp-sources/enrichment-api

rsync -av --progress --exclude="tmp-sources" * tmp-sources/events-processing-platform-subscribers
cp .env tmp-sources/events-processing-platform-subscribers/.env

cp -r ../customer-os-common-ai/* tmp-sources/customer-os-common-ai
cp -r ../customer-os-common-module/* tmp-sources/customer-os-common-module
cp -r ../customer-os-neo4j-repository/* tmp-sources/customer-os-neo4j-repository
cp -r ../events-processing-proto/* tmp-sources/events-processing-proto
cp -r ../events-processing-platform/* tmp-sources/events-processing-platform
cp -r ../enrichment-api/* tmp-sources/enrichment-api

cp Dockerfile tmp-sources/Dockerfile

docker build -t events-processing-platform-subscribers-local tmp-sources/.

rm -rf tmp-sources