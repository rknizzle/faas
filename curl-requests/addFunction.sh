url="localhost:${PORT:-8080}/functions"

curl -s -X POST \
  -H "content-type: application/json" \
  -w " %{http_code}" \
  "$url" \
  -d @./api-examples/function.json
