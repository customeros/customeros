# Test: Enrich Person by linkedinUrl
GET {{cos_url}}/enrich/v1/person
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
linkedinUrl: https://www.linkedin.com/in/mateocafe/
includeMobileNumber: true

HTTP 200
[Asserts]
header "Content-Type" contains "application/json; charset=utf-8"
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.isComplete" == true
jsonpath "$.data" exists
jsonpath "$.data.emails" isCollection
jsonpath "$.data.jobs" isCollection
jsonpath "$.data.phoneNumbers" isCollection
jsonpath "$.data.emails[0].address" == "matt@customeros.ai"
jsonpath "$.data.emails[0].deliverable" exists
jsonpath "$.data.emails[0].isRisky" == false
jsonpath "$.data.emails[0].type" == "work"
jsonpath "$.data.jobs[?(@.company == 'CustomerOS (YC S22)')].company" exists
jsonpath "$.data.location.region" == "London, England, United Kingdom"
jsonpath "$.data.name.firstName" == "Matt"
jsonpath "$.data.name.lastName" == "Brown"
jsonpath "$.data.name.fullName" == "Matt Brown"
jsonpath "$.data.phoneNumbers[0].number" == "+44 21 2614 6927"
jsonpath "$.data.phoneNumbers[0].type" == "mobile"
jsonpath "$.data.profilePic" exists
jsonpath "$.data.social.linkedin.publicId" == "mateocafe"
jsonpath "$.data.social.linkedin.url" == "https://www.linkedin.com/in/mateocafe"
jsonpath "$.data.social.x" exists
jsonpath "$.data.social.github" exists
jsonpath "$.data.social.discord" exists

