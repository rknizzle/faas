url="localhost:${PORT:-8080}/functions/rkneills/jsexample"

curl -s -X POST \
  -H "content-type: application/json" \
  -w " %{http_code}" \
  "$url"
