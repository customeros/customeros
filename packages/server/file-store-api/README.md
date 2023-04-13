# File-store-api


## api
An example http request to send to upload a file to the server:

```
curl -X 'POST' \
'http://127.0.0.1:10001/file' \
-H 'accept: application/json' \
-H 'Content-Type: multipart/form-data' \
-H 'X-Openline-API-KEY: 9eb87aa2-75e7-45b2-a1e6-53ed297d0ba8' \
-H 'X-Openline-USERNAME: alex@openline.ai' \
-F 'file=@test1.pdf'
```

```
curl -X 'GET' \
'http://127.0.0.1:10001/file/598abeb1-2979-4a0a-b2fd-f7441c2c4366' \
-H 'accept: application/json' \
-H 'Content-Type: multipart/form-data' \
-H 'X-Openline-API-KEY: 9eb87aa2-75e7-45b2-a1e6-53ed297d0ba8' \
-H 'X-Openline-USERNAME: alex@openline.ai' 
```