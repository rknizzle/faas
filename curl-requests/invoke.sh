url="localhost:${PORT:-8080}/functions/fn"

curl -s -X POST \
  -H "content-type: application/json" \
  -w " %{http_code}" \
  "$url"
