# Test: Successful existing risky email verification
GET {{cos_url}}/verify/v1/email
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
address: admin@customeros.ai

HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.email.emailAddress" == "admin@customeros.ai"
jsonpath "$.email.deliverable" == "false"
jsonpath "$.email.provider" == "google workspace"
jsonpath "$.email.secureGatewayProvider" == ""
jsonpath "$.email.isRisky" == true
jsonpath "$.email.isCatchAll" == false
jsonpath "$.email.risk" exists
jsonpath "$.email.risk.isFirewalled" == false
jsonpath "$.email.risk.isRoleMailbox" == true
jsonpath "$.email.risk.isSystemGenerated" == false
jsonpath "$.email.risk.isFreeProvider" == false
jsonpath "$.email.risk.isMailboxFull" == false
jsonpath "$.email.risk.isPrimaryDomain" == true
jsonpath "$.email.syntax.isValid" == true
jsonpath "$.email.syntax.domain" == "customeros.ai"
jsonpath "$.email.syntax.user" == "admin"
