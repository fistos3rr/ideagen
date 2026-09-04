#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "$SCRIPT_DIR/../common.sh"

CONFIG_FILE="${CONFIG_FILE:-$SCRIPT_DIR/../../user.env}"
load_env "$CONFIG_FILE" || exit 1

auth

SERVICE_ENDPOINT="/v1/ideas"
if [ $# -ge 1 ]; then
    TYPE_ID="$1"
else
    echo "no type_id found"
    exit 1
fi
TYPE_ID_ESC=$(escape_json "$TYPE_ID")
if [ $# -ge 2 ]; then
    TEXT="$2"
else
    echo "no text found"
    exit 1
fi
TEXT_ESC=$(escape_json "$TEXT")
JSON_BODY="{"
JSON_BODY+="\"type_id\":$TYPE_ID_ESC,"
JSON_BODY+="\"text\":\"$TEXT_ESC\""
JSON_BODY+="}"


RESULT_URL="${API_URL}${SERVICE_ENDPOINT}"
echo "Requesting $RESULT_URL"
echo "JSON request: $JSON_BODY"
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
