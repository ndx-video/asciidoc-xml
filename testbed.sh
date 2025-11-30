#!/bin/bash

# Testbed test runner script
# Usage: ./testbed.sh [-all] [-log]
#   -all: Run tests on all subdirectories (uses converter_testbed_all_test.go)
#   -log: Use main logging system instead of just writing to results.log

set -e

# Parse arguments
USE_ALL=false
USE_LOG=false
RESULTS_LOG="./testbed/results.log"
FAILS_LOG="./testbed/fails.log"

while [[ $# -gt 0 ]]; do
    case $1 in
        -all|--all)
            USE_ALL=true
            shift
            ;;
        -log|--log)
            USE_LOG=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [-all] [-log]"
            echo "  -all: Test all subdirectories"
            echo "  -log: Use main logging system"
            exit 1
            ;;
    esac
done

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Ensure testbed directory exists
if [ ! -d "./testbed" ]; then
    echo "Error: ./testbed directory not found"
    exit 1
fi

# Create results.log and fails.log files (truncate if exists)
> "$RESULTS_LOG"
> "$FAILS_LOG"

# Determine which test to run
if [ "$USE_ALL" = true ]; then
    TEST_NAME="TestConvertMarkdownToAsciiDoc_TestbedAllFiles"
    echo "Running testbed tests on ALL files (including subdirectories)..."
    echo "Test: $TEST_NAME" | tee -a "$RESULTS_LOG"
else
    TEST_NAME="TestConvertMarkdownToAsciiDoc_TestbedFiles"
    echo "Running testbed tests on root directory files only..."
    echo "Test: $TEST_NAME" | tee -a "$RESULTS_LOG"
fi

echo "Results will be written to: $RESULTS_LOG"
echo "Failure details will be written to: $FAILS_LOG"
echo ""

# Function to show progress bar
show_progress() {
    local current=$1
    local total=$2
    local file=$3
    local status=$4
    
    # Calculate percentage
    local percent=0
    if [ $total -gt 0 ]; then
        percent=$((current * 100 / total))
    fi
    
    # Create progress bar (20 chars wide)
    local bar_width=20
    local filled=$((percent * bar_width / 100))
    local bar=""
    for ((i=0; i<filled; i++)); do
        bar+="="
    done
    for ((i=filled; i<bar_width; i++)); do
        bar+=" "
    done
    
    # Truncate filename if too long
    local short_file=$(basename "$file")
    if [ ${#short_file} -gt 40 ]; then
        short_file="${short_file:0:37}..."
    fi
    
    # Print progress line (overwrite previous with \r)
    printf "\r[%s] %3d%% (%d/%d) %s [%s]" "$bar" "$percent" "$current" "$total" "$short_file" "$status"
}

# Run the test and capture output
# We'll parse the output to show progress
TEMP_OUTPUT=$(mktemp)
trap "rm -f $TEMP_OUTPUT" EXIT

# Count total files first (for progress bar)
TOTAL_FILES=0
if [ "$USE_ALL" = true ]; then
    TOTAL_FILES=$(find ./testbed -name "*.md" -not -name "README.md" | wc -l)
else
    TOTAL_FILES=$(find ./testbed -maxdepth 1 -name "*.md" -not -name "README.md" | wc -l)
fi

echo "Found $TOTAL_FILES test files"
echo ""

# Run go test and capture output
# Use a temporary file to capture all output, then process it
TEMP_TEST_OUTPUT=$(mktemp)
trap "rm -f $TEMP_TEST_OUTPUT" EXIT

# Track failure state for detailed logging
IN_FAILURE=false
FAILURE_CONTEXT=""
FAILURE_LINES=()

# Run test and capture all output
go test ./lib -run "^${TEST_NAME}$" -v 2>&1 | tee "$TEMP_TEST_OUTPUT" | while IFS= read -r line; do
    # Write all output to results log
    echo "$line" >> "$RESULTS_LOG"
    
    # Parse test output for progress and failure tracking
    if [[ $line =~ ^===.*RUN.*TestConvertMarkdownToAsciiDoc.*/(.+)$ ]]; then
        # Test started - extract filename
        FILE_NAME="${BASH_REMATCH[1]}"
        # Count how many tests have run so far
        CURRENT_FILE=$(grep -c "^===.*RUN.*TestConvertMarkdownToAsciiDoc.*/" "$RESULTS_LOG" 2>/dev/null || echo "0")
        show_progress "$CURRENT_FILE" "$TOTAL_FILES" "$FILE_NAME" "RUN"
        
        # Reset failure tracking for new test
        IN_FAILURE=false
        FAILURE_CONTEXT=""
        FAILURE_LINES=()
        
    elif [[ $line =~ ^---.*PASS.*TestConvertMarkdownToAsciiDoc.*/(.+)$ ]]; then
        # Test passed
        FILE_NAME="${BASH_REMATCH[1]}"
        CURRENT_FILE=$(grep -c "^---.*(PASS|FAIL).*TestConvertMarkdownToAsciiDoc.*/" "$RESULTS_LOG" 2>/dev/null || echo "0")
        show_progress "$CURRENT_FILE" "$TOTAL_FILES" "$FILE_NAME" "PASS"
        echo ""  # New line after progress update
        
        # Reset failure tracking
        IN_FAILURE=false
        FAILURE_CONTEXT=""
        FAILURE_LINES=()
        
    elif [[ $line =~ ^---.*FAIL.*TestConvertMarkdownToAsciiDoc.*/(.+)$ ]]; then
        # Test failed - start capturing failure details
        FILE_NAME="${BASH_REMATCH[1]}"
        CURRENT_FILE=$(grep -c "^---.*(PASS|FAIL).*TestConvertMarkdownToAsciiDoc.*/" "$RESULTS_LOG" 2>/dev/null || echo "0")
        show_progress "$CURRENT_FILE" "$TOTAL_FILES" "$FILE_NAME" "FAIL"
        echo ""  # New line after progress update
        
        # Start failure capture
        IN_FAILURE=true
        FAILURE_CONTEXT="$FILE_NAME"
        FAILURE_LINES=("$line")
        
        # Write failure header to fails.log
        echo "========================================" >> "$FAILS_LOG"
        echo "FAILURE: $FILE_NAME" >> "$FAILS_LOG"
        echo "Test: ${TEST_NAME}" >> "$FAILS_LOG"
        echo "Time: $(date)" >> "$FAILS_LOG"
        echo "========================================" >> "$FAILS_LOG"
        echo "$line" >> "$FAILS_LOG"
        
    elif [[ $line =~ ^---.*(PASS|FAIL).*TestConvertMarkdownToAsciiDoc[^/]*$ ]]; then
        # Overall test result (no subtest)
        echo ""
        echo "$line"
        
        # If this is a FAIL for the overall test, capture it
        if [[ $line =~ ^---.*FAIL ]]; then
            echo "========================================" >> "$FAILS_LOG"
            echo "OVERALL TEST FAILURE: ${TEST_NAME}" >> "$FAILS_LOG"
            echo "Time: $(date)" >> "$FAILS_LOG"
            echo "========================================" >> "$FAILS_LOG"
            echo "$line" >> "$FAILS_LOG"
        fi
        
        # End any ongoing failure capture
        if [ "$IN_FAILURE" = true ] && [ ${#FAILURE_LINES[@]} -gt 0 ]; then
            echo "" >> "$FAILS_LOG"
            IN_FAILURE=false
        fi
        
    elif [ "$IN_FAILURE" = true ]; then
        # Capture failure details - include error messages, stack traces, etc.
        FAILURE_LINES+=("$line")
        echo "$line" >> "$FAILS_LOG"
        
        # Stop capturing if we hit a blank line followed by non-error content
        # But continue if it looks like error content (indented, contains error keywords, file paths, etc.)
        if [[ -z "$line" ]]; then
            # Blank line - might be end of error, but check next line
            continue
        elif [[ $line =~ ^(FAIL|ok|===|---) ]]; then
            # Hit next test or summary - stop capturing
            IN_FAILURE=false
            echo "" >> "$FAILS_LOG"
        elif [[ $line =~ ^[[:space:]]+(Error|error|FAIL|assertion|expected|got|want|panic|runtime|goroutine|\.go:[0-9]+) ]] || \
             [[ $line =~ ^[[:space:]]+[0-9]+ ]] || \
             [[ $line =~ \.go:[0-9]+ ]] || \
             [[ $line =~ ^[[:space:]]+at ]] || \
             [[ $line =~ ^[[:space:]]+.*\(.*\) ]]; then
            # Looks like error content - continue capturing
            continue
        fi
    fi
done

# Also write the full output to results log
cat "$TEMP_TEST_OUTPUT" >> "$RESULTS_LOG"

# Extract any additional failure information from the full output
# This catches any failures we might have missed in the streaming parse
if grep -q "^--- FAIL" "$TEMP_TEST_OUTPUT"; then
    # Extract failure sections with context - look for patterns that indicate detailed errors
    FAILURE_COUNT=$(grep -c "^--- FAIL" "$TEMP_TEST_OUTPUT" 2>/dev/null || echo "0")
    if [ "$FAILURE_COUNT" -gt 0 ]; then
        echo "" >> "$FAILS_LOG"
        echo "========================================" >> "$FAILS_LOG"
        echo "COMPLETE FAILURE EXTRACTION FROM FULL OUTPUT" >> "$FAILS_LOG"
        echo "Found $FAILURE_COUNT failure(s)" >> "$FAILS_LOG"
        echo "========================================" >> "$FAILS_LOG"
        
        # Extract failure blocks with surrounding context using awk
        awk '
            BEGIN {
                in_failure=0
                context_lines=0
                max_context=100
            }
            /^--- FAIL/ { 
                if (in_failure) {
                    # End previous failure block
                    print ""
                    print "---"
                    print ""
                }
                in_failure=1
                context_lines=0
                print ""
                print "========================================"
                print $0
                next
            }
            in_failure {
                print $0
                context_lines++
                if (/^=== RUN/ || /^--- PASS/ || /^ok[[:space:]]+[^[:space:]]/ || /^FAIL[[:space:]]*$/) {
                    in_failure=0
                    print ""
                    print "---"
                    print ""
                } else if (context_lines > max_context && /^[[:space:]]*$/) {
                    in_failure=0
                    print ""
                    print "---"
                    print ""
                }
            }
            END {
                if (in_failure) {
                    print ""
                    print "---"
                    print ""
                }
            }
        ' "$TEMP_TEST_OUTPUT" >> "$FAILS_LOG"
        
        echo "" >> "$FAILS_LOG"
        echo "========================================" >> "$FAILS_LOG"
        echo "END OF FAILURE EXTRACTION" >> "$FAILS_LOG"
        echo "========================================" >> "$FAILS_LOG"
    fi
fi

# Final newline
echo ""

# Check if we should use main logging system
if [ "$USE_LOG" = true ]; then
    echo "Logging to main logging system..."
    
    # Try to load logging config from cli/adc.json
    LOG_CONFIG=""
    if [ -f "./cli/adc.json" ]; then
        # Extract logging config from adc.json
        LOG_CONFIG=$(cat ./cli/adc.json | grep -A 20 '"logging"' || echo "")
    fi
    
    # If no config found, use defaults
    if [ -z "$LOG_CONFIG" ]; then
        echo "No logging config found in cli/adc.json, using defaults"
        # Could initialize logger here, but for now just note it
    fi
    
    # For now, we'll just note that logging was requested
    # Full integration would require Go code to initialize the logger
    echo "Note: Full logging system integration requires Go code initialization"
    echo "Results are available in: $RESULTS_LOG"
fi

# Show summary
echo ""
echo "========================================="
echo "Test Summary"
echo "========================================="
PASS_COUNT=$(grep -c "^--- PASS" "$RESULTS_LOG" 2>/dev/null | head -1 || echo "0")
FAIL_COUNT=$(grep -c "^--- FAIL" "$RESULTS_LOG" 2>/dev/null | head -1 || echo "0")
# Clean up any whitespace/newlines
PASS_COUNT=$(echo "$PASS_COUNT" | tr -d '\n\r ' | head -c 10)
FAIL_COUNT=$(echo "$FAIL_COUNT" | tr -d '\n\r ' | head -c 10)
# Default to 0 if empty
if [ -z "$PASS_COUNT" ]; then
    PASS_COUNT=0
fi
if [ -z "$FAIL_COUNT" ]; then
    FAIL_COUNT=0
fi
echo "Passed: $PASS_COUNT"
echo "Failed: $FAIL_COUNT"
echo ""
echo "Full results: $RESULTS_LOG"
if [ "$FAIL_COUNT" != "0" ] && [ "$FAIL_COUNT" -gt 0 ] 2>/dev/null; then
    echo "Failure details: $FAILS_LOG"
fi

# Check if we have failures (use numeric comparison)
if [ "$FAIL_COUNT" != "0" ] && [ "$FAIL_COUNT" -gt 0 ] 2>/dev/null; then
    echo ""
    echo "Failed tests:"
    grep "^--- FAIL" "$RESULTS_LOG" 2>/dev/null | sed 's/^--- FAIL: /  - /' || true
    echo ""
    echo "Detailed failure information has been written to: $FAILS_LOG"
    echo "Use this file for debugging failed tests."
    exit 1
fi

exit 0

