#!/bin/sh
rm -rf tmp-sources
mkdir tmp-sources
mkdir tmp-sources/customer-os-common-module
mkdir tmp-sources/file-storage-api

cp -rv `ls -A | grep -vE "tmp-sources|tmp-sources/file-storage-api"` tmp-sources/file-storage-api
cp -r ../customer-os-common-module/* tmp-sources/customer-os-common-module

cp Dockerfile tmp-sources/Dockerfile

cd tmp-sources
docker build -t aa .
