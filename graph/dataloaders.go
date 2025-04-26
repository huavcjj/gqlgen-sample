package graph

import (
	"context"
	"gqlgen-todos/graph/model"
	"time"

	"gqlgen-todos/internal/db"

	"net/http"
	"strconv"
	"strings"
)

const loadersKey = "dataLoaders"

type Loaders struct {
	UsersByIDs     *UserLoader
	TodosByUserIDs *TodoLoader
}

func Middleware(queries *db.Queries, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		loaders := &Loaders{
			UsersByIDs: NewUserLoader(UserLoaderConfig{
				Wait:     1 * time.Millisecond,
				MaxBatch: 100,
				Fetch: func(keys []string) ([]*model.User, []error) {
					users, err := queries.ListUsersByIDs(ctx, keys)
					if err != nil {
						errs := make([]error, len(keys))
						for i := range errs {
							errs[i] = err
						}
						return nil, errs
					}

					userByID := make(map[string]*model.User)
					for _, u := range users {
						userByID[u.ID] = &model.User{
							ID:   u.ID,
							Name: u.Name,
						}
					}

					result := make([]*model.User, len(keys))
					for i, id := range keys {
						result[i] = userByID[id]
					}

					return result, nil
				},
			}),
		}

		ctx = context.WithValue(ctx, loadersKey, loaders)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

func toPKs(ids []int64) string {

	var pks []string
	for _, id := range ids {
		pks = append(pks, strconv.FormatInt(id, 10))
	}
	return strings.Join(pks, ",")
}
