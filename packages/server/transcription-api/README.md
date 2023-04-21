# Transcription-api


## api
An example http request to send to the server:

### for transcription and summarisation
```
curl -X 'POST' \
'http://127.0.0.1:8014/transcribe' \
-H 'accept: application/json' \
-H 'Content-Type: multipart/form-data' \
-H 'X-Openline-API-KEY: b1ced267-43b9-4be1-a5ef-8d054e6f84c1' \
-H 'X-Openline-USERNAME: torrey@openline.ai' \
-F 'users=["d4acc4f7-be03-453d-8444-2ce842e721e0", "5ca3e332-8246-4a4a-99bb-27a7eec6d412", "2620f73c-ce68-4c28-ad06-26d78848d031", "255225d8-4b0b-44aa-b4bf-87e5eb4309ac"]' \
-F 'contacts=["echotest"]' \
-F 'topic=Discussion about a new call routing platform using Jambonz and Node-RED' \
-F 'start=2022-06-27T02:12-07:00' \
-F 'type=meeting' \
-F 'group_id=1234567890' \
-F 'file_id=f8d623b4-5cf0-417e-88bc-e2c2547eb1da'
```

### for summary only
```
curl -X 'POST' \
'http://127.0.0.1:8014/summary' \
-H 'accept: application/json' \
-H 'Content-Type: multipart/form-data' \
-H 'X-Openline-API-KEY: b1ced267-43b9-4be1-a5ef-8d054e6f84c1' \
-H 'X-Openline-USERNAME: torrey@openline.ai' \
-F "transcript=$(cat tests/data/transcription.json)"
```

### for action-points only
```
curl -X 'POST' \
'http://127.0.0.1:8014/action-items' \
-H 'accept: application/json' \
-H 'Content-Type: multipart/form-data' \
-H 'X-Openline-API-KEY: b1ced267-43b9-4be1-a5ef-8d054e6f84c1' \
-H 'X-Openline-USERNAME: torrey@openline.ai' \
-F "transcript=$(cat tests/data/transcription.json)"
```
## packer ami
to an aws ami image

to build the api image you must have AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_REGION properly set with your AWS credentials

then you can run the following commands to build the image

```
packer init aws-ubuntu.pkr.hcl
packer validate aws-ubuntu.pkr.hcl
packer build -var 'environment=openline-dev' aws-ubuntu.pkr.hcl
```

for production builds you also need to specify the region

```
export AWS_REGION=eu-west-1
packer init aws-ubuntu.pkr.hcl
packer validate aws-ubuntu.pkr.hcl
packer build -var 'region=eu-west-1' -var 'environment=openline-production' aws-ubuntu.pkr.hcl
