package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	//"github.com/go-redis/redis"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/getData/{category}/{pageNo}", getData).Methods("GET")
	router.HandleFunc("/api/getArticle/{category}/{articleId}", getArticle).Methods("GET")
	router.HandleFunc("/api/search/{searchString}/{category}/{pageNo}", searchArticles).Methods("GET")
	router.HandleFunc("/api", getRoot).Methods("GET")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":12345", handler))
}
