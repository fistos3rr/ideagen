#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "$SCRIPT_DIR/../common.sh"

CONFIG_FILE="${CONFIG_FILE:-$SCRIPT_DIR/../../user.env}"
load_env "$CONFIG_FILE" || exit 1

REQUIRED_VARS=("API_URL" "LOGIN_ENDPOINT" "USER_EMAIL" "USER_PASSWORD")
check_required_vars "${REQUIRED_VARS[@]}" || exit 1

TOKEN_FILE="$SCRIPT_DIR/../../.jwt_token"
COOKIE_FILE="$SCRIPT_DIR/../../.cookies.txt"

EMAIL_ESC=$(escape_json "$USER_EMAIL")
PASSWORD_ESC=$(escape_json "$USER_PASSWORD")

JSON_BODY="{"
JSON_BODY+="\"email\":\"$EMAIL_ESC\","
JSON_BODY+="\"password\":\"$PASSWORD_ESC\""
JSON_BODY+="}"

echo "JSON: $JSON_BODY"

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_URL}${LOGIN_ENDPOINT}" \
    -H "Content-Type: application/json" \
    -d "$JSON_BODY" \
    -c "$COOKIE_FILE" \
    -b "$COOKIE_FILE")

HTTP_BODY=$(echo "$RESPONSE" | head -n -1)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

echo "HTTP status: $HTTP_CODE"

if [ "$HTTP_CODE" -eq 200 ]; then
    TOKEN=$(extract_json_field "$HTTP_BODY" "access_token")
    [ -z "$TOKEN" ] && TOKEN=$(extract_json_field "$HTTP_BODY" "token")
    if [ -n "$TOKEN" ]; then
        echo "$TOKEN" > "$TOKEN_FILE"
    else
        echo "Access-token not found"
    fi
    pretty_json "$HTTP_BODY"
else
    echo "Login error with code $HTTP_CODE"
    pretty_json "$HTTP_BODY"
fi
