#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "$SCRIPT_DIR/../common.sh"

CONFIG_FILE="${CONFIG_FILE:-$SCRIPT_DIR/../../user.env}"
load_env "$CONFIG_FILE" || exit 1

auth

SERVICE_ENDPOINT="/v1/types"
if [ $# -ge 1 ]; then
    NAME="$1"
else
    echo "no NAME found"
    exit 1
fi
NAME_ESC=$(escape_json "$NAME")
JSON_BODY="{"
JSON_BODY+="\"name\":\"$NAME_ESC\""
JSON_BODY+="}"


RESULT_URL="${API_URL}${SERVICE_ENDPOINT}"
echo "Requesting $RESULT_URL"
curl_args=(-s -w "\n%{http_code}" -X POST "$RESULT_URL")
curl_args+=(-H "Content-Type: application/json")
if [ -n "$TOKEN" ]; then
    curl_args+=(-H "Authorization: Bearer $TOKEN")
fi
curl_args+=(-d "$JSON_BODY")

RESPONSE=$(curl "${curl_args[@]}")

HTTP_BODY=$(echo "$RESPONSE" | head -n -1)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

echo "HTTP status: $HTTP_CODE"
pretty_json "$HTTP_BODY"
