##  validateEmail Endpoint

This endpoint is responsible for validating an email address using the POST method. It checks if the email is valid and returns a response in JSON format.

### Request

Method: POST
URL: /validateEmail

Headers
````
X-Openline-TENANT: The tenant value for authentication
X-Openline-API-KEY: The API key value for authentication
````
Body
The request body should be a JSON object with the following structure:
````
{
    "email": "sample@example.com" 
}
````
### Example

````
curl -X \
POST -H "Content-Type: application/json" \
-H "X-Openline-TENANT: <YOUR_TENANT_UUID>" \
-H "X-Openline-API-KEY: <YOUR_API_KEY_UUID>" \
-d '{"email": "example@example.com"}' \
https://validation.openline.ai/validateEmail
````