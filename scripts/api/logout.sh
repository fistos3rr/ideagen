#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "$SCRIPT_DIR/../common.sh"

CONFIG_FILE="${CONFIG_FILE:-$SCRIPT_DIR/../../user.env}"
load_env "$CONFIG_FILE" || exit 1

REQUIRED_VARS=("API_URL")
check_required_vars "${REQUIRED_VARS[@]}" || exit 1

TOKEN_FILE="$SCRIPT_DIR/../../.jwt_token"
COOKIE_FILE="$SCRIPT_DIR/../../.cookies.txt"

rm $COOKIE_FILE
rm $TOKEN_FILE
