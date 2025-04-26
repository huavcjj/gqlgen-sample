package main

import (
	"context"
	_ "embed"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5"
	"gqlgen-todos/internal/db"

	"github.com/vektah/gqlparser/v2/ast"
	"gqlgen-todos/graph"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8081"

var (
	//go:embed sql/schema/init.sql
	ddl     string
	queries *db.Queries
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://sampleuser:samplepass@localhost:5432/sampledb?sslmode=disable")
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err = conn.Close(ctx)
		if err != nil {

		}
	}(conn, ctx)

	if _, err = conn.Exec(ctx, ddl); err != nil {
		log.Fatalf("failed to initialize schema: %v", err)
	}

	queries = db.New(conn)

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
	http.Handle("/query", graph.Middleware(queries, srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
