# Test: Unsuccessful email verification - bad address format
GET {{cos_url}}/verify/v1/email
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
address: silviucustomeros.ai

HTTP 200
[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.email.emailAddress" == "silviucustomeros.ai"
jsonpath "$.email.deliverable" == ""
jsonpath "$.email.provider" == ""
jsonpath "$.email.secureGatewayProvider" == ""
jsonpath "$.email.isRisky" == false
jsonpath "$.email.isCatchAll" == false
jsonpath "$.email.risk" exists
jsonpath "$.email.risk.isFirewalled" == false
jsonpath "$.email.risk.isRoleMailbox" == false
jsonpath "$.email.risk.isSystemGenerated" == false
jsonpath "$.email.risk.isFreeProvider" == false
jsonpath "$.email.risk.isMailboxFull" == false
jsonpath "$.email.risk.isPrimaryDomain" == false
jsonpath "$.email.syntax.isValid" == false
jsonpath "$.email.syntax.domain" == ""
jsonpath "$.email.syntax.user" == ""
