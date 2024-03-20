package httpserver

import (
	"fmt"
	"net/http"

	"chat-api/pkg/mysqlrepo"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func StartHTTPServer() {
	// initialise redis
	// redisClient := mysqlrepo.InitialiseRedis()
	// defer redisClient.Close()

	// create indexes
	client := mysqlrepo.ConnectDB()
	defer client.Close()
	mysqlrepo.CreateFetchChatBetweenIndex()

	r := mux.NewRouter()
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	}).Methods(http.MethodGet)

	r.HandleFunc("/register", registerHandler).Methods(http.MethodPost)
	r.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	r.HandleFunc("/verify-contact", verifyContactHandler).Methods(http.MethodPost)
	r.HandleFunc("/chat-history", chatHistoryHandler).Methods(http.MethodGet)
	r.HandleFunc("/contact-list", contactListHandler).Methods(http.MethodGet)

	// Use default options
	// handler := cors.Default().Handler(r)
	// http.ListenAndServe(":8080", handler)
}
