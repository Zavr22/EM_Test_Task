package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
)

func NewGraphqlHandler(graphqlH *handler.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		graphqlH.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func PlaygroundHandler() echo.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
