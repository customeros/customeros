# Test: Enrich Organization by domain
GET {{cos_url}}/enrich/v1/organization
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
domain: customeros.ai

HTTP *
[Asserts]
status == 200
header "Content-Type" == "application/json"
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.data" exists
jsonpath "$.data.name" ==  "CustomerOS"
jsonpath "$.data.domain" == "customeros.ai"
jsonpath "$.data.website" == "customeros.ai"
jsonpath "$.data.description" exists
jsonpath "$.data.longDescription" exists
jsonpath "$.data.employees" exists
jsonpath "$.data.foundedYear"  == 2022
jsonpath "$.data.logos" isCollection
jsonpath "$.data.icons" isCollection
jsonpath "$.data.industry.industry" == "Computers Electronics and Technology"
jsonpath "$.data.socials" isCollection
jsonpath "$.data.location.isHeadquarter" == false
jsonpath "$.data.location.country" == "United Kingdom"
jsonpath "$.data.location.countryCodeA2" == "GB"
jsonpath "$.data.location.city" == "London"
jsonpath "$.data.location.region" == "Europe"
