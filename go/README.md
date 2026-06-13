Tip: Any time you modify those // @Summary comments in the future, you will just need to run swag init -g cmd/server/main.go from your go directory to update the docs!

swag init -g cmd/server/main.go --parseDependency --parseInternal

echo 'JWT_SECRET='$(openssl rand -hex 32) >> /home/steve/dev/me/neneloga/go/.env && echo "Added JWT_SECRET"


TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"me@stevenene.top","password":"'"$(grep ADMIN_PASSWORD /home/steve/dev/me/neneloga/go/.env | cut -d= -f2)"'"}' | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4) && \
echo "Token: ${TOKEN:0:40}..." && \
curl -s -o /dev/null -w "GET /servers with token: %{http_code}\n" \
  -H "Authorization: Bearer $TOKEN" http://localhost:8080/servers


go run cmd/seed/main.go --users=2 --servers=3 --logs=5
