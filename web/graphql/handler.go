package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/labstack/echo/v4"
)

func GraphQLHandler(resolver *Resolver) echo.HandlerFunc {
	graphQLHandler := handler.NewDefaultServer(
		NewExecutableSchema(Config{Resolvers: resolver}),
	)

	return func(c echo.Context) error {
		// Обработка GraphQL-запроса.
		graphQLHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
