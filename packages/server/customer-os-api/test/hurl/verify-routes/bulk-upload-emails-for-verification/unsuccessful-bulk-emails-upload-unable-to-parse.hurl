# Test: Unsuccessful bulk emails upload missing emailColumn
POST {{cos_url}}/verify/v1/email/bulk
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,unsuccessful-bulk-emails-upload-unable-to-parse.csv;

HTTP 400
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.message" == "Unable to parse csv file"
