package main

import (
	"github.com/eecscord/workspace/src/go/auth"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {
	//create http mux
	mux := http.NewServeMux()

	//register routes
	registerAllRoutes(mux)

	http.ListenAndServe(":8080", mux)

}

func registerAllRoutes(mux *http.ServeMux) {
	err := auth.RegisterRoutes(mux)
	logError(err)
}

func logError(err error) {
	if err != nil {
		log.Fatal("Error registering auth routes")
	}
}
