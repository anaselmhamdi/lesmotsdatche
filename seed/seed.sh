#!/bin/sh
# Seed script to load test puzzle into the API

API_URL="${API_URL:-http://api:8080}"
TODAY=$(date +%Y-%m-%d)

echo "Waiting for API to be ready..."
until curl -sf "$API_URL/health" > /dev/null 2>&1; do
  sleep 1
done

echo "API is ready. Seeding puzzle for date: $TODAY"

# Update the puzzle JSON with today's date and post it
sed "s/2026-01-13/$TODAY/g" /seed/puzzle.json | \
  curl -sf -X POST "$API_URL/admin/v1/puzzles" \
    -H "Content-Type: application/json" \
    -d @-

if [ $? -eq 0 ]; then
  echo ""
  echo "Puzzle seeded successfully!"
else
  echo "Failed to seed puzzle"
  exit 1
fi
