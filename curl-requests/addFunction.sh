url="localhost:5555/functions"

curl -s -X POST \
  -H "content-type: application/json" \
  -w " %{http_code}" \
  "$url" \
  -d @./api-examples/function.json
