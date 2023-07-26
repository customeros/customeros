

### generate
https://grpc.io/docs/protoc-installation/
https://grpc.io/docs/languages/go/quickstart/

To build the project run

```
make
```

## Development

This service uses the environment variables described below. The env files have a default value if not provided ( check .env file )

| Variable                      | Meaning                                                                                |
|-------------------------------|----------------------------------------------------------------------------------------|
| MESSAGE_STORE_URL             | url of the GRPC interface of the message store                                         |
| MESSAGE_STORE_API_KEY         | message store API key                                                                  |
| CHANNELS_API_SERVER_ADDRESS   | ip:port to bind for the rest api, normally should be ":8013"                           |
| CHANNELS_GRPC_PORT            | port used for the channel-api grpc interface, should be 9013                           |
| MAIL_API_KEY                  | The api key used to validated received emails, must mach what is set in the AWS Lambda |
| OASIS_API_URL                 | IP & port of the GRPC interface of oasis api                                           |
| CHANNELS_API_CORS_URL         | url of the frontend, needed to allow cros-site scripting                               |
| WEBCHAT_API_KEY               | The api key used to validated received messages and login requests                     |
| WEBSOCKET_PING_INTERVAL       | Ping interval in seconds to monitor websocket connections                              |


## Setting up gmail in local environment

follow the procedure in https://developers.google.com/gmail/api/quickstart/go
start ngrok to tunnel to channel-api
```
ngrok http 8013
```

* create a credential of type oauth client id
* select web application as application type
* add http://localhost:3006 as authorized javascript origin
* add https://(your ngrok url)/auth as authorized redirect uri
* set GMAIL_CLIENT_ID to the client id
* set GMAIL_CLIENT_SECRET to the client secret
* set GMAIL_REDIRECT_URIS to the redirect url specified by ory

Additionally, you need to set up ory
* set up the ory tunnel as described in the oasis-frontent README
* create an API key in the ory admin console
* set ORY_API_KEY to the api key
* set ORY_SERVER_URL to http:://localhost:4000


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


