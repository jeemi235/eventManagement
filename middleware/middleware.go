// This is the context for database
package middlewares

import (
	database "10/db"
	"context"
	"log"
	"net/http"
)

func DbContext(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//This will call function from file db and connect to the database
		db := database.Connect()
		ctx := context.WithValue(r.Context(), "database", db)
		defer db.Close()
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

var err error
var ctx context.Context

func DbMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		crConn := database.Connect()
		if err != nil {
			log.Fatal(err)
		}
		ctx = context.WithValue(request.Context(), "crConn", crConn)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
