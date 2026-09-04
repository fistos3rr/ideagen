#!/bin/bash
# common.sh - common utility functions

load_env() {
    local config_file="$1"
    if [ -f "$config_file" ]; then
        source "$config_file"
    else
        echo "File $config_file not found!"
        return 1
    fi
}

check_required_vars() {
    local missing=()
    for var in "$@"; do
        if [ -z "${!var:-}" ]; then
            missing+=("$var")
        fi
    done
    if [ ${#missing[@]} -ne 0 ]; then
        echo "No vars were found: ${missing[*]}"
        return 1
    fi
    return 0
}

extract_json_field() {
    local json="$1"
    local field="$2" 
    echo "$json" | grep -o "\"${field}\":\"[^\"]*\"" | head -1 | sed 's/.*:"//;s/"//'
}

pretty_json() {
    local json="$1"
    echo "$json"
}

escape_json() {
    printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}

auth() {
    local token_file="${AUTH_TOKEN_FILE:-$SCRIPT_DIR/../../.jwt_token}"
    local cookie_file="${AUTH_COOKIE_FILE:-$SCRIPT_DIR/../../.cookies.txt}"

    if [ ! -f "$token_file" ]; then
        echo "No jwt token file"
        export TOKEN=
    else
        export TOKEN=$(cat "$token_file")
    fi

    if [ ! -f "$cookie_file" ]; then
        echo "No cookies file"
        export COOKIE_FILE=
    else
        export COOKIE_FILE="$cookie_file"
    fi

    return 0
}
