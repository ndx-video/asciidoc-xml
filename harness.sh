#!/bin/bash

# AsciiDoc XML Converter - Development Harness Manager
# Manages the web server for testing XML generation and XSLT development

DEFAULT_PORT=8005
PORT=${PORT:-$DEFAULT_PORT}
PIDFILE="/tmp/asciidoc-xml-harness-${PORT}.pid"
LOGFILE="/tmp/asciidoc-xml-harness-${PORT}.log"
SKIP_TESTS=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Find process using a port
find_port_process() {
    local port=$1
    if command -v lsof >/dev/null 2>&1; then
        lsof -ti:$port 2>/dev/null
    elif command -v netstat >/dev/null 2>&1; then
        netstat -tlnp 2>/dev/null | grep ":$port " | awk '{print $7}' | cut -d'/' -f1 | head -1
    elif command -v ss >/dev/null 2>&1; then
        ss -tlnp 2>/dev/null | grep ":$port " | awk '{print $6}' | cut -d',' -f2 | cut -d'(' -f2 | cut -d')' -f1 | head -1
    else
        echo ""
    fi
}

# Get process info
get_process_info() {
    local pid=$1
    if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
        ps -p "$pid" -o pid,cmd --no-headers 2>/dev/null
    else
        echo ""
    fi
}

# Check if server is running
is_running() {
    if [ -f "$PIDFILE" ]; then
        local pid=$(cat "$PIDFILE" 2>/dev/null)
        if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
            return 0
        else
            rm -f "$PIDFILE"
            return 1
        fi
    fi
    return 1
}

# Run tests
run_tests() {
    if [ "$SKIP_TESTS" = true ]; then
        echo -e "${YELLOW}Skipping tests (--no-test flag provided)${NC}"
        return 0
    fi

    # Get the directory where this script is located
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    cd "$SCRIPT_DIR" || exit 1

    echo -e "${GREEN}Running all tests in project...${NC}"
    
    # Find all packages with test files
    echo -e "${YELLOW}Test packages found:${NC}"
    local test_packages
    test_packages=$(go list ./... | while read -r pkg; do
        if go list -f '{{len .TestGoFiles}}{{len .XTestGoFiles}}' "$pkg" 2>/dev/null | grep -q '[1-9]'; then
            echo "$pkg"
        fi
    done)
    
    if [ -z "$test_packages" ]; then
        echo -e "${YELLOW}  No test packages found${NC}"
    else
        echo "$test_packages" | while read -r pkg; do
            echo -e "  ${GREEN}✓${NC} $pkg"
        done
    fi
    echo ""
    
    # Run all tests from project root (includes web, cli, lib, etc.)
    # Capture output for processing and potential logging
    local test_exit_code
    local temp_output
    temp_output=$(mktemp)
    
    # Run tests and capture output
    echo -e "${GREEN}Executing: go test -v ./...${NC}"
    echo ""
    go test -v ./... > "$temp_output" 2>&1
    test_exit_code=$?
    
    # Count tests and packages from temp file before processing
    local test_count
    test_count=$(grep -E "^--- PASS:|^    --- PASS:" "$temp_output" 2>/dev/null | wc -l | tr -d ' ')
    local package_count
    package_count=$(grep -E "^ok[[:space:]]+[^[:space:]]" "$temp_output" 2>/dev/null | wc -l | tr -d ' ')
    if [ -z "$test_count" ] || [ "$test_count" = "0" ]; then
        test_count="0"
    fi
    if [ -z "$package_count" ] || [ "$package_count" = "0" ]; then
        package_count="0"
    fi
    
    # Process output and highlight FAIL lines in red, collect failure lines
    local failure_lines=()
    local in_failure=false
    while IFS= read -r line; do
        if [[ "$line" =~ ^---[[:space:]]FAIL: ]]; then
            echo -e "${RED}$line${NC}"
            failure_lines+=("$line")
            in_failure=true
        elif [[ "$line" =~ ^FAIL$ ]]; then
            echo -e "${RED}$line${NC}"
            failure_lines+=("$line")
            in_failure=false
        elif [ "$in_failure" = true ]; then
            # Capture error details that follow a FAIL line
            if [[ "$line" =~ ^[[:space:]]+(Error|FAIL|error|assertion|expected|got|want) ]] || \
               [[ "$line" =~ ^[[:space:]]+.*\.go:[0-9]+: ]] || \
               [[ "$line" =~ ^[[:space:]]+[0-9]+ ]]; then
                echo "$line"
                failure_lines+=("$line")
            else
                echo "$line"
                # Stop collecting if we hit a blank line or non-error line after several error lines
                if [[ -z "$line" ]]; then
                    in_failure=false
                fi
            fi
        else
            echo "$line"
        fi
    done < "$temp_output"
    
    # Clean up temp file
    rm -f "$temp_output"
    
    if [ $test_exit_code -eq 0 ]; then
        echo ""
        echo -e "${GREEN}✓ All tests passed${NC}"
        echo -e "${GREEN}  Packages tested: $package_count${NC}"
        echo -e "${GREEN}  Tests passed: $test_count${NC}"
        return 0
    else
        echo ""
        echo -e "${RED}✗ Tests failed${NC}"
        
        # Ask user if they want to save failure lines to log file
        if [ ${#failure_lines[@]} -gt 0 ]; then
            echo ""
            read -p "Save failure lines to failed-tests.log? (y/n): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                # Write failure lines to log file (overwrite existing)
                printf '%s\n' "${failure_lines[@]}" > failed-tests.log
                echo -e "${GREEN}Failure lines written to failed-tests.log${NC}"
            fi
        fi
        
        return 1
    fi
}

# Start the server
start_server() {
    local port=$1
    
    # Run tests first unless --no-test is provided
    if ! run_tests; then
        echo -e "${RED}Tests failed. Use --no-test to skip tests and start server anyway.${NC}"
        return 1
    fi
    
    # Check if port is already in use
    local existing_pid=$(find_port_process "$port")
    
    if [ -n "$existing_pid" ]; then
        local process_info=$(get_process_info "$existing_pid")
        echo -e "${YELLOW}Port $port is already in use by process:${NC}"
        echo "  PID: $existing_pid"
        echo "  Info: $process_info"
        echo ""
        read -p "Do you want to kill this process? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            kill "$existing_pid" 2>/dev/null
            sleep 1
            if kill -0 "$existing_pid" 2>/dev/null; then
                kill -9 "$existing_pid" 2>/dev/null
            fi
            echo -e "${GREEN}Process killed.${NC}"
        else
            # Try next port
            local next_port=$((port + 1))
            echo -e "${YELLOW}Trying port $next_port instead...${NC}"
            start_server "$next_port"
            return
        fi
    fi

    if is_running; then
        local pid=$(cat "$PIDFILE")
        echo -e "${YELLOW}Server is already running (PID: $pid)${NC}"
        return 1
    fi

    # Get the directory where this script is located
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    cd "$SCRIPT_DIR" || exit 1

    # Check if web/main.go exists
    if [ ! -f "web/main.go" ]; then
        echo -e "${RED}Error: web/main.go not found${NC}"
        exit 1
    fi

    echo -e "${GREEN}Starting server on port $port...${NC}"
    
    # Build and run the server
    cd web || exit 1
    PORT=$port go run main.go > "$LOGFILE" 2>&1 &
    local server_pid=$!
    
    # Wait a moment to see if it starts successfully
    sleep 1
    
    if kill -0 "$server_pid" 2>/dev/null; then
        echo "$server_pid" > "$PIDFILE"
        echo -e "${GREEN}Server started successfully (PID: $server_pid)${NC}"
        echo -e "${GREEN}Server URL: http://localhost:$port${NC}"
        echo -e "${GREEN}Log file: $LOGFILE${NC}"
        return 0
    else
        echo -e "${RED}Failed to start server${NC}"
        echo "Check log file: $LOGFILE"
        cat "$LOGFILE" 2>/dev/null
        return 1
    fi
}

# Stop the server
stop_server() {
    if ! is_running; then
        echo -e "${YELLOW}Server is not running${NC}"
        return 1
    fi

    local pid=$(cat "$PIDFILE")
    local port=$PORT
    echo -e "${GREEN}Stopping server (PID: $pid)...${NC}"
    
    kill "$pid" 2>/dev/null
    sleep 1
    
    if kill -0 "$pid" 2>/dev/null; then
        echo -e "${YELLOW}Server did not stop gracefully, forcing...${NC}"
        kill -9 "$pid" 2>/dev/null
        sleep 1
    fi
    
    # Wait for the process to fully terminate
    local waited=0
    while kill -0 "$pid" 2>/dev/null && [ $waited -lt 5 ]; do
        sleep 0.5
        waited=$((waited + 1))
    done
    
    rm -f "$PIDFILE"
    echo -e "${GREEN}Server stopped${NC}"
    return 0
}

# Restart the server (without port check)
restart_server() {
    echo -e "${YELLOW}Restarting server...${NC}"
    stop_server
    
    # Wait for port to be released (TCP TIME_WAIT can take a moment)
    local port=$PORT
    local max_wait=5
    local waited=0
    
    echo -e "${YELLOW}Waiting for port $port to be released...${NC}"
    while [ $waited -lt $max_wait ]; do
        local existing_pid=$(find_port_process "$port")
        if [ -z "$existing_pid" ]; then
            break
        fi
        sleep 1
        waited=$((waited + 1))
        echo -e "${YELLOW}  Waiting... (${waited}/${max_wait})${NC}"
    done
    
    # Final check - if port is still in use, try to kill it
    local existing_pid=$(find_port_process "$port")
    if [ -n "$existing_pid" ]; then
        echo -e "${YELLOW}Port still in use, attempting to free it...${NC}"
        kill "$existing_pid" 2>/dev/null
        sleep 1
        if kill -0 "$existing_pid" 2>/dev/null; then
            kill -9 "$existing_pid" 2>/dev/null
            sleep 1
        fi
    fi
    
    # Verify port is free before starting
    existing_pid=$(find_port_process "$port")
    if [ -n "$existing_pid" ]; then
        echo -e "${RED}Port $port is still in use by PID $existing_pid. Cannot restart.${NC}"
        return 1
    fi
    
    if is_running; then
        local pid=$(cat "$PIDFILE")
        echo -e "${YELLOW}Server is already running (PID: $pid)${NC}"
        return 1
    fi

    # Get the directory where this script is located
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    cd "$SCRIPT_DIR" || exit 1

    # Check if web/main.go exists
    if [ ! -f "web/main.go" ]; then
        echo -e "${RED}Error: web/main.go not found${NC}"
        exit 1
    fi

    echo -e "${GREEN}Starting server on port $port...${NC}"
    
    # Build and run the server
    cd web || exit 1
    PORT=$port go run main.go > "$LOGFILE" 2>&1 &
    local server_pid=$!
    
    # Wait a moment to see if it starts successfully
    sleep 1
    
    if kill -0 "$server_pid" 2>/dev/null; then
        echo "$server_pid" > "$PIDFILE"
        echo -e "${GREEN}Server started successfully (PID: $server_pid)${NC}"
        echo -e "${GREEN}Server URL: http://localhost:$port${NC}"
        echo -e "${GREEN}Log file: $LOGFILE${NC}"
        return 0
    else
        echo -e "${RED}Failed to start server${NC}"
        echo "Check log file: $LOGFILE"
        cat "$LOGFILE" 2>/dev/null
        return 1
    fi
}

# Reload the server (alias for restart)
reload_server() {
    restart_server
}

# Show server status
show_status() {
    if is_running; then
        local pid=$(cat "$PIDFILE")
        local process_info=$(get_process_info "$pid")
        echo -e "${GREEN}Server is running${NC}"
        echo "  PID: $pid"
        echo "  Port: $PORT"
        echo "  Process: $process_info"
        echo "  URL: http://localhost:$PORT"
        echo "  Log: $LOGFILE"
        
        # Check if port is actually listening
        local port_pid=$(find_port_process "$PORT")
        if [ -n "$port_pid" ] && [ "$port_pid" = "$pid" ]; then
            echo -e "  Status: ${GREEN}Port $PORT is listening${NC}"
        else
            echo -e "  Status: ${YELLOW}Port $PORT may not be listening${NC}"
        fi
    else
        echo -e "${RED}Server is not running${NC}"
        
        # Check if port is in use by another process
        local port_pid=$(find_port_process "$PORT")
        if [ -n "$port_pid" ]; then
            local process_info=$(get_process_info "$port_pid")
            echo -e "${YELLOW}Port $PORT is in use by another process:${NC}"
            echo "  PID: $port_pid"
            echo "  Info: $process_info"
        fi
    fi
}

# Show usage
show_usage() {
    echo "Usage: $0 {start|stop|restart|reload|status|test-all} [--no-test]"
    echo ""
    echo "Commands:"
    echo "  start   - Start the development server (default port: $DEFAULT_PORT)"
    echo "  stop    - Stop the development server"
    echo "  restart - Restart the development server (no port check)"
    echo "  reload  - Reload the development server (alias for restart)"
    echo "  status  - Show server status"
    echo "  test-all - Run all tests without starting the server"
    echo ""
    echo "Options:"
    echo "  --no-test  - Skip running tests before starting server"
    echo ""
    echo "Environment variables:"
    echo "  PORT    - Port number to use (default: $DEFAULT_PORT)"
    echo ""
    echo "Examples:"
    echo "  $0 start"
    echo "  $0 start --no-test"
    echo "  PORT=8080 $0 start"
    echo "  $0 status"
    echo "  $0 test-all"
}

# Parse arguments
ARGS=()
for arg in "$@"; do
    case $arg in
        --no-test)
            SKIP_TESTS=true
            ;;
        *)
            ARGS+=("$arg")
            ;;
    esac
done

# Main
case "${ARGS[0]:-}" in
    start)
        start_server "$PORT"
        ;;
    stop)
        stop_server
        ;;
    restart)
        restart_server
        ;;
    reload)
        reload_server
        ;;
    status)
        show_status
        ;;
    test-all)
        run_tests
        exit $?
        ;;
    *)
        show_usage
        exit 1
        ;;
esac

exit $?

