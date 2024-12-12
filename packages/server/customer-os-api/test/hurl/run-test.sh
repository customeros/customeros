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

# Initialize overall test status
TESTS_FAILED=0

# Create/clear the test-output.txt file
: > test-output.txt

# Process each test file
while IFS= read -r test_file
do
    echo "Running tests from file: $test_file" | tee -a test-output.txt

    # Generate a new UUID for this specific test
    NEW_UUID=$(uuidgen | tr '[:upper:]' '[:lower:]')  # Generate lowercase UUID
    capitalized_custom_id=$(echo "$NEW_UUID" | awk 'BEGIN{FS=OFS="-"} {for (i=1; i<=NF; i++) { sub(/[a-z]/, toupper(substr($i, match($i, /[a-z]/), 1)), $i) }} 1')
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
    hurl --very-verbose --test --continue-on-error --variable "custom_id=$capitalized_custom_id" --variable "random_str=$RANDOM_STRING" --variable "cos_url=$COS_URL" --variable "api_key=$HURL_TENANT_API_KEY" "$test_file" > temp_output.txt 2>&1

    TEST_EXIT_CODE=$?

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