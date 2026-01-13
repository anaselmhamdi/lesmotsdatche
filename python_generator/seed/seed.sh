#!/bin/bash
# Seed the API with a sample puzzle

API_URL="${API_URL:-http://localhost:8080}"

echo "Seeding puzzle to $API_URL..."

curl -X POST "$API_URL/admin/v1/puzzles" \
  -H "Content-Type: application/json" \
  -d @"$(dirname "$0")/puzzle.json"

echo ""
echo "Done!"
