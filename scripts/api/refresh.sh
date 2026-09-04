#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "$SCRIPT_DIR/../common.sh"

CONFIG_FILE="${CONFIG_FILE:-$SCRIPT_DIR/../../user.env}"
load_env "$CONFIG_FILE" || exit 1

REQUIRED_VARS=("API_URL" "REFRESH_ENDPOINT")
check_required_vars "${REQUIRED_VARS[@]}" || exit 1

TOKEN_FILE="$SCRIPT_DIR/../../.jwt_token"
COOKIE_FILE="$SCRIPT_DIR/../../.cookies.txt"

if [ ! -f "$COOKIE_FILE" ]; then
    echo "No cookie-file $COOKIE_FILE provided, login pls" 
    exit 1
fi

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_URL}${REFRESH_ENDPOINT}" \
    -H "Content-Type: application/json" \
    -b "$COOKIE_FILE" \
    -c "$COOKIE_FILE")

HTTP_BODY=$(echo "$RESPONSE" | head -n -1)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

echo "HTTP status: $HTTP_CODE"

if [ "$HTTP_CODE" -eq 200 ]; then
    echo "Token refreshed successfully"
    NEW_TOKEN=$(extract_json_field "$HTTP_BODY" "access_token")
    [ -z "$NEW_TOKEN" ] && NEW_TOKEN=$(extract_json_field "$HTTP_BODY" "token")

    if [ -n "$NEW_TOKEN" ]; then
        echo "$NEW_TOKEN" > "$TOKEN_FILE"
    else
        echo "Access-token not found"
        echo "Server response:"
        pretty_json "$HTTP_BODY"
        exit 1
    fi

    if [ ! -s "$COOKIE_FILE" ]; then
        echo "Cookie-file empty, error while refreshing"
    fi

    pretty_json "$HTTP_BODY"
fi
