Tip: Any time you modify those // @Summary comments in the future, you will just need to run swag init -g cmd/server/main.go from your go directory to update the docs!

swag init -g cmd/server/main.go --parseDependency --parseInternal
