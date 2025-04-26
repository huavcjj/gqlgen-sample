package main

import (
	"context"
	"database/sql"
	_ "embed"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"gqlgen-todos/graph"
	"gqlgen-todos/internal/db"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
)

const defaultPort = "8081"

//go:embed sql/schema/todos.sql
var ddl string
var queries *db.Queries

func init() {
	ctx := context.Background()

	sqlDB, err := sql.Open("sqlite", "./sqlite.db")
	if err != nil {
		log.Fatalf("failed to open in-memory database: %v", err)
	}

	if _, err := sqlDB.ExecContext(ctx, ddl); err != nil {
		log.Fatalf("failed to initialize schema: %v", err)
	}

	queries = db.New(sqlDB)
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		DB: queries,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
