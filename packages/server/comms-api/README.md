


To build the project run

```
make
```

## Development

This service uses the environment variables described below. The env files have a default value if not provided ( check .env file )

| Variable                        | Meaning                                                                                |
|---------------------------------|----------------------------------------------------------------------------------------|
| COMMS_API_SERVER_ADDRESS        | ip:port to bind for the rest api, normally should be ":8013"                           |
| COMMS_API_CORS_URL              | url of the frontend, needed to allow cros-site scripting                               |
| FILE_STORE_API                  | url of the file store api                                                              |
| FILE_STORE_API_KEY              | api key used to validate requests to the file store api                                |
| COMMS_API_MAIL_API_KEY          | The api key used to validated received emails, must mach what is set in the AWS Lambda |
| COMMS_API_VCON_API_KEY          | The api key used to validated received by the vcon endpoint                            |
| WEBCHAT_API_KEY                 | The api key used to validated received messages and login requests                     |
| WEBSOCKET_PING_INTERVAL         | Ping interval in seconds to monitor websocket connections                              |
| WEBRTC_AUTH_SECRET              | Secret used to sign the auth tokens                                                    |
| WEBRTC_AUTH_TTL                 | Validity time of Ephemeral auth tokens                                                 |
| REDIS_HOST                      | Redis host                                                                             |
| POSTGRES_HOST                   | Postgres host                                                                          |
| POSTGRES_PORT                   | Postgres port                                                                          |
| POSTGRES_USER                   | Postgres user                                                                          |
| POSTGRES_PASSWORD               | Postgres password                                                                      |
| POSTGRES_DB                     | Postgres database                                                                      |
| POSTGRES_DB_MAX_CONN            | Postgres max connections                                                               |
| POSTGRES_DB_MAX_IDLE_CONN       | Postgres max idle connections                                                          |
| POSTGRES_DB_CONN_MAX_LIFETIME   | Postgres max connection lifetime                                                       |
| CALCOM_SECRET                   | Secret used to validate cal.com webhooks                                               |





## Setting up google email forwarding in dev environment
1. Go to your gmail account settings. Click on the "Forwarding and POP/IMAP" tab.
2. Click on the "Add a forwarding address" button.
3. Enter dev@getopenline.com and click "Next".
4. Login into oasis and look for an email from forwarding-noreply@google.com
5. Get the verification code, go back to gmail settings, input the code, click "Verify" and then "Proceed".
6. Click "Save Changes" and you're done.

## Setting up google email forwarding in prod environment
1. Go to your gmail account settings. Click on the "Forwarding and POP/IMAP" tab.
2. Click on the "Add a forwarding address" button.
3. Enter openline@getopenline.com and click "Next".
4. Login into oasis and look for an email from forwarding-noreply@google.com.
5. Get the verification code, go back to gmail settings, input the code, click "Verify" and then "Proceed".
6. Click "Save Changes" and you're done.

## Setting up google email forwarding in local (ninja) environment
Ngrok and aws lambda are needed for this to work:
1. start ngrok to tunnel to channel-api: `ngrok http 8013`
2. copy the ngrok url and set it to the environment variable `OL_MAIL_CB_HOST` of the lambda function "openline-local-sender" by doing the following:
   1. go to the AWS console and go to the lambda service
   2. select the function "openline-local-sender"
   3. click on the "Configuration" tab
   4. click on the "Environment variables" section
   5. For the variable with the name "OL_MAIL_CB_HOST" and the value of your ngrok url

On gmail:
1. Go to your gmail account settings. Click on the "Forwarding and POP/IMAP" tab.
2. Click on the "Add a forwarding address" button.
3. Enter local@getopenline.com and click "Next".
4. Login into oasis and look for an email from forwarding-noreply@google.com
5. Get the verification code, go back to gmail settings, input the code, click "Verify" and then "Proceed".
6. Click "Save Changes" and you're done.

## SES - LAMBDA - S3
Naming convention:
1. LAMBDA names: $tenant-$domain-sender: 
   * openline-ai-sender,
   * openline-dev-sender, 

2. S3 bucket names: ses-$emailaddress: 
   * ses-dev-getopenline-com, 
   * ses-openline-getopenline-com
3. SES rules: $emailaddress:
    * dev-getopenline-com
    * openline-getopenline-com

Region: All of our forwarding infrastructure is in ireland eu-west-2

## CAL.COM integration
1. go to https://app.cal.com/settings/developer/api-keys and create an API KEY
2. Create a webhook
```
curl \
-X POST \
-H "Content-Type: application/json" \
-d '{"eventTriggers": ["BOOKING_CANCELLED", "BOOKING_CREATED", "BOOKING_RESCHEDULED"], "active": true, "subscriberUrl": "https://api.customeros.ai/calcom"}' \
'https://api.cal.com/v1/webhooks?apiKey=cal_live_a9e7877b0dc073de54d6c94fc86a4732'
```
3. Set the secret for the newly created webhook. 
- Atm this operation is done via the cal.com ui. unfortunately cal.com api does not support this operation yet.
- https://app.cal.com/settings/developer/webhooks/{web-hook-id}
- Find the secret for prod here: https://start.1password.com/open/i?a=247WXGWKQJDK7FK5GLXPC6EMZM&v=qs5pqywxpd3yhuypwk24kj5k4e&i=kvqs44yv26x37vihioj2bv3pxq&h=openline.1password.com


## call_progress / recording integration

if running out of k8s, need to set up the port forward on the postgres and the redis

kubectl port-forward --namespace openline svc/customer-db-redis-master 6379:6379 &
kubectl port-forward --namespace openline svc/customer-db-postgresql 5432:5432 &

### sending a call_progress webhook
curl -X 'POST' \
'http://127.0.0.1:8013/call_progress' \
-H 'accept: application/json' \
-H 'Content-Type: application/json' \
-H 'X-API-KEY: 44b4086d-a5d6-4954-b62d-bf2c78e6bb36' \
--data '{
    "version": "1.0",
    "correlation_id": "my_awesome_call",
    "event": "CALL_START",
    "from": {"tel": "+32485000000", "type": "pstn"},
    "to": {"mailto": "AgentSmith@openline.ai", "type": "webrtc"},
    "start_time": "2023-08-01T12:34:56Z"
}'


### sending a recording
curl -X 'POST' \
'http://127.0.0.1:8013/recording' \
-H 'accept: application/json' \
-H 'Content-Type: multipart/form-data' \
-H 'X-API-KEY: 44b4086d-a5d6-4954-b62d-bf2c78e6bb36' \
-F 'correlationId=my_awesome_call' \
-F 'audio=@test.mp3'
