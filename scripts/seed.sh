#!/usr/bin/env bash
# scripts/seed.sh
# ──────────────────────────────────────────────────────────────────────────────
# CLI helper that uses mongoimport to load seed_data.json into MongoDB.
#
# Usage:
#   bash scripts/seed.sh [--uri <mongo-uri>] [--db <database>] [--drop]
#
# Defaults:
#   --uri  mongodb://localhost:27017
#   --db   restaurant
#   --drop (flag) if passed, drops each collection before importing
#
# Examples:
#   bash scripts/seed.sh
#   bash scripts/seed.sh --uri "mongodb://user:pass@host:27017" --db mydb --drop
# ──────────────────────────────────────────────────────────────────────────────

set -euo pipefail

# ── Default values ────────────────────────────────────────────────────────────
MONGO_URI="mongodb://localhost:27017"
DB_NAME="restaurant"
DROP_FLAG=""

# ── Argument parsing ──────────────────────────────────────────────────────────
while [[ $# -gt 0 ]]; do
  case "$1" in
    --uri)  MONGO_URI="$2"; shift 2 ;;
    --db)   DB_NAME="$2";   shift 2 ;;
    --drop) DROP_FLAG="--drop"; shift ;;
    *) echo "Unknown argument: $1"; exit 1 ;;
  esac
done

SEED_FILE="$(dirname "$0")/seed_data.json"

if [[ ! -f "$SEED_FILE" ]]; then
  echo "❌  Seed file not found: $SEED_FILE"
  exit 1
fi

if ! command -v mongoimport &>/dev/null; then
  echo "❌  mongoimport not found in PATH."
  echo "    Install it via: https://www.mongodb.com/docs/database-tools/mongoimport/"
  exit 1
fi

if ! command -v jq &>/dev/null; then
  echo "❌  jq not found in PATH."
  echo "    Install it via: https://stedolan.github.io/jq/download/"
  exit 1
fi

echo "🌱  Seeding database '${DB_NAME}' at ${MONGO_URI}"
echo ""

# ── Helper: import one collection ─────────────────────────────────────────────
import_collection() {
  local collection="$1"
  local tmp_file
  tmp_file="$(mktemp /tmp/seed_${collection}_XXXXXX.json)"

  # Extract the array for this collection and write one document per line (NDJSON)
  jq -c ".${collection}[]" "$SEED_FILE" > "$tmp_file"

  local count
  count=$(wc -l < "$tmp_file" | tr -d ' ')

  echo "📦  Importing ${count} document(s) → collection '${collection}'"

  mongoimport \
    --uri="${MONGO_URI}" \
    --db="${DB_NAME}" \
    --collection="${collection}" \
    --file="$tmp_file" \
    --type=json \
    ${DROP_FLAG}

  rm -f "$tmp_file"
  echo "✅  '${collection}' imported successfully."
  echo ""
}

# ── Import each collection ────────────────────────────────────────────────────
import_collection "menus"
import_collection "foods"
import_collection "tables"

echo "🎉  Seed complete! Database '${DB_NAME}' is ready."
