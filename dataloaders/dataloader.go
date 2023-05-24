package dataloader

import (
	"10/graph/model"
	"context"
	database "10/db"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type ctxKeyType struct{ name string }

var ctxKey = ctxKeyType{"userCtx"}

type loaders struct {
	Users *UserLoader
}

func LoadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ldrs := loaders{}
		ctx := r.Context()
		// crConn := ctx.Value("crConn").(*sql.DB)
		crConn := database.Connect()
		wait := 350 * time.Millisecond

		// Todo loader
		ldrs.Users = &UserLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(ids []string) ([]*model.User, []error) {
				var sqlQuery string
				if len(ids) == 1 {
					sqlQuery = "SELECT id, name, phone FROM users WHERE id = ?"
				} else {
					sqlQuery = "SELECT id, name, phone FROM users WHERE id IN (?)"
				}
				sqlQuery, arguments, err := sqlx.In(sqlQuery, ids)
				if err != nil {
					log.Println("Error to fetch user query")
				}
				sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
				rows, err := crConn.Query(sqlQuery, arguments...)
				if err != nil {
					log.Println("Error to read user data")
				}
				userById := map[int]*model.User{}
				defer func() {
					if err := rows.Close(); err != nil {
						log.Println("Error at defer function of user")
					}
				}()

				for rows.Next() {
					var user model.User
					err := rows.Scan(&user.ID, &user.Name, &user.Phone)
					if err != nil {
						log.Println("Error to scan todo data")
					}
					// Store user.ID in key & model of user in value
					userById[user.ID] = &user
				}

				users := make([]*model.User, len(ids))
				for i, id := range ids {
					x, _ := strconv.Atoi(id)
					users[i] = userById[x]
					i++
				}
				return users, nil
			},
		}

		ctx = context.WithValue(ctx, ctxKey, ldrs)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func CtxLoaders(ctx context.Context) loaders {
	return ctx.Value(ctxKey).(loaders)
}
