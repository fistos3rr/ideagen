#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "$SCRIPT_DIR/../common.sh"

CONFIG_FILE="${CONFIG_FILE:-$SCRIPT_DIR/../../user.env}"
load_env "$CONFIG_FILE" || exit 1

auth

SERVICE_ENDPOINT="/v1/service/me"

curl_args=(-s -w "\n%{http_code}" -X GET "${API_URL}${SERVICE_ENDPOINT}")

if [ -n "$TOKEN" ]; then
    curl_args+=(-H "Authorization: Bearer $TOKEN")
fi

RESPONSE=$(curl "${curl_args[@]}")

HTTP_BODY=$(echo "$RESPONSE" | head -n -1)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

echo "HTTP status: $HTTP_CODE"
pretty_json "$HTTP_BODY"
