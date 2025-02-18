#!/bin/bash

# Check for required environment variable
if [ -z "${HURL_TENANT_API_KEY}" ]; then
    echo "❌ Error: HURL_TENANT_API_KEY environment variable is required"
    exit 1
fi

COS_URL="https://api.customeros.ai"

# Initialize arrays for final summary
declare -a ALL_TEST_NAMES
declare -a ALL_TEST_STATUSES
declare -a ORGANIZATION_IDS
declare -a CONTACT_IDS

# Initialize overall test status
TESTS_FAILED=0

# Create/clear the test-output.txt file
: > test-output.txt

# Function to extract organization ID from test output
extract_org_id() {
    local output_file=$1
    local org_id=""

    # Look for organization ID in the JSON response
    if grep -q "\"organization\":{\"id\":\"[^\"]*\"" "$output_file"; then
        org_id=$(grep -o "\"organization\":{\"id\":\"[^\"]*\"" "$output_file" | grep -o "\"id\":\"[^\"]*\"" | cut -d'"' -f4)
        if [ ! -z "$org_id" ]; then
            ORGANIZATION_IDS+=("$org_id")
            echo "Captured organization ID: $org_id" >> test-output.txt
        fi
    fi
}

# Function to fetch and hide contacts
fetch_and_hide_contacts() {
    echo "Fetching contacts to hide..." | tee -a test-output.txt

    # First query - search by email
    cat > contact_search_email.hurl << EOF
POST ${COS_URL}/query
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: ${HURL_TENANT_API_KEY}
{
    "query": "query searchContacts(\$limit: Int, \$where: Filter, \$sort: SortBy) { ui_contacts_search(limit: \$limit, where: \$where, sort: \$sort) { ids totalElements totalAvailable } }",
    "variables": {
        "where": {
            "AND": [{
                "filter": {
                    "property": "CONTACTS_PRIMARY_EMAIL",
                    "operation": "CONTAINS",
                    "value": "hurl-"
                }
            }]
        },
        "sort": {
            "by": "CONTACTS_UPDATED_AT",
            "direction": "DESC"
        }
    },
    "operationName": "searchContacts"
}

HTTP 200
[Captures]
contact_ids: jsonpath "$.data.ui_contacts_search.ids"
EOF

    # Second query - search by LinkedIn URL
    cat > contact_search_linkedin.hurl << EOF
POST ${COS_URL}/query
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: ${HURL_TENANT_API_KEY}
{
    "query": "query searchContacts(\$limit: Int, \$where: Filter, \$sort: SortBy) { ui_contacts_search(limit: \$limit, where: \$where, sort: \$sort) { ids totalElements totalAvailable } }",
    "variables": {
        "where": {
            "AND": [{
                "filter": {
                    "property": "CONTACTS_LINKEDIN",
                    "operation": "CONTAINS",
                    "includeEmpty": false,
                    "value": "https://linkedin.com/in/hurl"
                }
            }]
        },
        "sort": {
            "by": "CONTACTS_UPDATED_AT",
            "direction": "DESC"
        }
    },
    "operationName": "searchContacts"
}

HTTP 200
[Captures]
linkedin_contact_ids: jsonpath "$.data.ui_contacts_search.ids"
EOF

    # Run both contact searches
    echo "Searching for contacts with email containing 'hurl-'..." | tee -a test-output.txt
    contact_search_output=$(hurl --very-verbose contact_search_email.hurl)

    echo "Searching for contacts with LinkedIn URLs containing 'hurl'..." | tee -a test-output.txt
    linkedin_search_output=$(hurl --very-verbose contact_search_linkedin.hurl)

    # Clean up temporary files
    rm -f contact_search_email.hurl contact_search_linkedin.hurl

    # Process both sets of contact IDs
    declare -a all_contact_ids=()

    # Process email-based contacts
    if [[ $contact_search_output =~ \"ui_contacts_search\":[[:space:]]*{[[:space:]]*\"ids\":[[:space:]]*\[([^\]]*)\] ]]; then
        IFS=',' read -ra email_contacts <<< "${BASH_REMATCH[1]//\"/}"
        for id in "${email_contacts[@]}"; do
            id=${id// /}  # Remove any whitespace
            if [ ! -z "$id" ]; then
                all_contact_ids+=("$id")
            fi
        done
    fi

    # Process LinkedIn-based contacts
    if [[ $linkedin_search_output =~ \"ui_contacts_search\":[[:space:]]*{[[:space:]]*\"ids\":[[:space:]]*\[([^\]]*)\] ]]; then
        IFS=',' read -ra linkedin_contacts <<< "${BASH_REMATCH[1]//\"/}"
        for id in "${linkedin_contacts[@]}"; do
            id=${id// /}  # Remove any whitespace
            if [ ! -z "$id" ]; then
                all_contact_ids+=("$id")
            fi
        done
    fi

    # Remove duplicates from all_contact_ids
    all_contact_ids=($(echo "${all_contact_ids[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' '))

    echo "Found ${#all_contact_ids[@]} total unique contacts to hide" | tee -a test-output.txt

    # Hide each contact
    for contact_id in "${all_contact_ids[@]}"; do
        if [ ! -z "$contact_id" ]; then
            echo "Hiding contact: $contact_id" | tee -a test-output.txt

            # Create temporary hide contact file
            cat > hide_contact.hurl << EOF
POST ${COS_URL}/query
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: ${HURL_TENANT_API_KEY}
{
    "query": "mutation HideContact(\$contactId: ID!) { contact_Hide(contactId: \$contactId) { accepted } }",
    "variables": {
        "contactId": "${contact_id}"
    }
}

HTTP 200
[Asserts]
jsonpath "$.data.contact_Hide.accepted" == true
EOF

            # Run hide contact operation
            hurl --very-verbose hide_contact.hurl >> test-output.txt 2>&1
            rm -f hide_contact.hurl
        fi
    done
}

# Process each test file
while IFS= read -r test_file
do
    echo "Running tests from file: $test_file" | tee -a test-output.txt

    # Generate timestamp for this specific test
    TIMESTAMP=$(date +%s)
    capitalized_custom_id="hurl-${TIMESTAMP}"
    RANDOM_STRING=$(openssl rand -base64 12 | tr -dc 'a-z' | fold -w 10 | head -n 1)

    # Get the test name from the file
    test_name=""
    while IFS= read -r line
    do
        if [[ $line =~ ^"# Test"[[:space:]]*":"[[:space:]]*(.+) ]]; then
            test_name="${BASH_REMATCH[1]}"
            break
        fi
    done < "$test_file"

    # Run hurl command with verbose output to capture all details
    hurl --very-verbose --test --continue-on-error \
        --variable "custom_id=$capitalized_custom_id" \
        --variable "random_str=$RANDOM_STRING" \
        --variable "cos_url=$COS_URL" \
        --variable "api_key=$HURL_TENANT_API_KEY" \
        "$test_file" > temp_output.txt 2>&1

    TEST_EXIT_CODE=$?

    # Extract organization ID if present
    extract_org_id "temp_output.txt"

    # Determine test status
    if [ $TEST_EXIT_CODE -eq 0 ]; then
        test_status="✅ Passed"
    else
        test_status="❌ Failed"
        TESTS_FAILED=1

        # Process the output file line by line for error details
        in_error_block=0
        error_location=""
        actual_value=""
        expected_value=""

        while IFS= read -r line
        do
            # Capture error location
            if [[ $line == *"--> "* && $line == *".hurl:"* ]]; then
                # If we were processing a previous error, print it
                if [ -n "$error_location" ]; then
                    {
                        echo "🔍 Error location: $error_location"
                        [ -n "$actual_value" ] && echo "   Actual: $actual_value"
                        [ -n "$expected_value" ] && echo "   Expected: $expected_value"
                        echo ""
                    } | tee -a test-output.txt
                fi

                error_location="$line"
                actual_value=""
                expected_value=""
                in_error_block=1
                continue
            fi

            # Capture actual and expected values when in an error block
            if [ $in_error_block -eq 1 ]; then
                # Check for status code error with carets (^^^)
                if [[ $line =~ "^^^ actual value is" ]]; then
                    actual_value="$(echo "$line" | sed -n 's/.*actual value is <\([^>]*\)>.*/\1/p')"
                elif [[ $line =~ "HTTP "* ]]; then
                    expected_value="$(echo "$line" | grep -o '[0-9]\{3\}')"
                # Check for regular value comparisons
                elif [[ $line == *"actual: "* ]]; then
                    actual_value="${line#*actual: }"
                elif [[ $line == *"expected: "* ]]; then
                    expected_value="${line#*expected: }"
                fi
            fi

            # Reset error block flag when we hit a blank line
            if [[ -z "$line" && $in_error_block -eq 1 ]]; then
                in_error_block=0
            fi
        done < temp_output.txt

        # Print the last error if there is one
        if [ -n "$error_location" ]; then
            {
                echo "🔍 Error location: $error_location"
                [ -n "$actual_value" ] && echo "   Actual: $actual_value"
                [ -n "$expected_value" ] && echo "   Expected: $expected_value"
                echo ""
            } | tee -a test-output.txt
        fi
    fi

    # Add results to global arrays
    if [ -n "$test_name" ]; then
        ALL_TEST_NAMES+=("$test_name")
        ALL_TEST_STATUSES+=("$test_status")
    fi

    # Cleanup temporary file
    rm -f temp_output.txt
done < <(find . -name "*.hurl")

# Run teardown for organization and contact cleanup
echo "Running teardown operations..." | tee -a test-output.txt

# First, handle contacts
fetch_and_hide_contacts

# Then handle organizations
if [ ${#ORGANIZATION_IDS[@]} -gt 0 ]; then
    # Convert organization IDs array to JSON array format
    printf -v org_ids_json '"%s",' "${ORGANIZATION_IDS[@]}"
    org_ids_json="[${org_ids_json%,}]"

    echo "Cleaning up organizations: $org_ids_json" | tee -a test-output.txt

    # Create temporary teardown.hurl file
    cat > teardown.hurl << EOF
# Test: Teardown - Hide test organizations
POST ${COS_URL}/query
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: ${HURL_TENANT_API_KEY}
{
    "query": "mutation OrganizationHideAll(\$ids: [ID!]!) { organization_HideAll(ids: \$ids) { result } }",
    "variables": {
        "ids": ${org_ids_json}
    }
}

HTTP 200
[Asserts]
jsonpath "$.data.organization_HideAll.result" == true
EOF

    # Run teardown
    hurl --very-verbose --test --continue-on-error \
        --variable "cos_url=$COS_URL" \
        --variable "api_key=$HURL_TENANT_API_KEY" \
        teardown.hurl >> test-output.txt 2>&1

    TEARDOWN_EXIT_CODE=$?
    if [ $TEARDOWN_EXIT_CODE -ne 0 ]; then
        echo "⚠️ Warning: Teardown operations completed with some errors" | tee -a test-output.txt
    fi

    # Clean up temporary teardown file
    rm -f teardown.hurl
else
    echo "No organization IDs collected, skipping organization cleanup" | tee -a test-output.txt
fi

# Print final aggregated summary
{
    echo "Final Test Summary"
    echo "=================="
    echo "Total Tests: ${#ALL_TEST_NAMES[@]}"
} | tee -a test-output.txt

# Count passed and failed tests
PASSED_TESTS=0
FAILED_TESTS=0
i=0
while [ $i -lt "${#ALL_TEST_STATUSES[@]}" ]
do
    if [[ ${ALL_TEST_STATUSES[$i]} == "✅ Passed" ]]; then
        ((PASSED_TESTS++))
    else
        ((FAILED_TESTS++))
    fi
    ((i++))
done

{
    echo "Passed: $PASSED_TESTS"
    echo "Failed: $FAILED_TESTS"
    echo ""
    echo "Detailed Results:"
    echo "----------------"
} | tee -a test-output.txt

i=0
while [ $i -lt "${#ALL_TEST_NAMES[@]}" ]
do
    echo "${ALL_TEST_STATUSES[$i]}: ${ALL_TEST_NAMES[$i]}" | tee -a test-output.txt
    ((i++))
done

{
    echo "===================="
    echo "Test execution completed"
} | tee -a test-output.txt

# Exit with failure if any tests failed
exit $TESTS_FAILED