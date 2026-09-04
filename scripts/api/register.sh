#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../common.sh"

CONFIG_FILE="${CONFIG_FILE:-$SCRIPT_DIR/../../user.env}"
load_env "$CONFIG_FILE" || exit 1

if [ $# -ge 1 ]; then
    EMAIL="$1"
fi
if [ $# -ge 2 ]; then
    PASSWORD="$2"
fi

REQUIRED_VARS=("API_URL" "REGISTER_ENDPOINT" "EMAIL" "PASSWORD")
check_required_vars "${REQUIRED_VARS[@]}" || exit 1

EMAIL_ESC=$(escape_json "$EMAIL")
PASSWORD_ESC=$(escape_json "$PASSWORD")

JSON_BODY="{"
JSON_BODY+="\"email\":\"$EMAIL_ESC\","
JSON_BODY+="\"password\":\"$PASSWORD_ESC\""
JSON_BODY+="}"

echo "Registering: $EMAIL"
echo "JSON: $JSON_BODY"

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_URL}${REGISTER_ENDPOINT}" \
    -H "Content-Type: application/json" \
    -d "$JSON_BODY")

HTTP_BODY=$(echo "$RESPONSE" | head -n -1)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

echo "HTTP status: $HTTP_CODE"
echo "HTTP_BODY: "
pretty_json "$HTTP_BODY"
