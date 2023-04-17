#!/bin/bash

# This script is used to build the environment for the transcription-api
# service. It is run by the ami packer script.

echo "CUSTOMER_OS_API_KEY=$CUSTOMER_OS_API_KEY" > /etc/transcription/environment
echo "CUSTOMER_OS_API_URL=$CUSTOMER_OS_API_URL" >> /etc/transcription/environment
echo "TRANSCRIPTION_KEY=$TRANSCRIPTION_KEY" >> /etc/transcription/environment
echo "VCON_API_KEY=$VCON_API_KEY" >> /etc/transcription/environment
echo "VCON_API_URL=$VCON_API_URL" >> /etc/transcription/environment
echo "FILE_STORE_API_KEY=$FILE_STORE_API_KEY" >> /etc/transcription/environment
echo "FILE_STORE_API_URL=$FILE_STORE_API_URL" >> /etc/transcription/environment
echo "REPLICATE_API_TOKEN=$REPLICATE_API_TOKEN" >> /etc/transcription/environment
