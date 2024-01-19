#!/bin/sh
rm -rf tmp-sources
mkdir tmp-sources
mkdir tmp-sources/customer-os-common-module
mkdir tmp-sources/customer-os-neo4j-repository
mkdir tmp-sources/events-processing-proto
mkdir tmp-sources/customer-os-data-upkeeper

rsync -av --progress --exclude="tmp-sources" * tmp-sources/customer-os-data-upkeeper
cp .env tmp-sources/customer-os-data-upkeeper/.env
cp -r ../../server/customer-os-common-module/* tmp-sources/customer-os-common-module
cp -r ../../server/customer-os-neo4j-repository/* tmp-sources/customer-os-neo4j-repository
cp -r ../../server/events-processing-proto/* tmp-sources/events-processing-proto

cp Dockerfile tmp-sources/Dockerfile

docker build -t aa tmp-sources/.

rm -rf tmp-sources