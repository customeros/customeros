An example http request to send to the server:

```
curl -X 'POST' \
'http://127.0.0.1:8014/transcribe' \
-H 'accept: application/json' \
-H 'Content-Type: multipart/form-data' \
-H 'X-Openline-API-KEY: b1ced267-43b9-4be1-a5ef-8d054e6f84c1' \
-H 'X-Openline-USERNAME: torrey@openline.ai' \
-F 'users=["d4acc4f7-be03-453d-8444-2ce842e721e0", "5ca3e332-8246-4a4a-99bb-27a7eec6d412", "2620f73c-ce68-4c28-ad06-26d78848d031", "255225d8-4b0b-44aa-b4bf-87e5eb4309ac"]' \
-F 'contacts=["echotest"]' \
-F 'file=@vuy-wxso-sik (2022-06-27 02_12 GMT-7).mp3'
```