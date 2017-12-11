package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
	"github.com/salmaanp/common"
	"gopkg.in/mgo.v2/bson"
)

type dbMap struct {
	url string
	authDB string
	username string
	password string
	dbName string
}

type article struct {
	//Text  string `bson:"article_text"`
	//Title string `bson:"title"`
	//Url   string `bson:"url"`
	Score int `bson:"score"`
}

type errorMap struct {
	InvalidCategory string
	InvalidPageNo string
	ServerError string
	InvalidArticle string
}

var dbConfig dbMap = getDbDetails()

var errors = errorMap{"Invalid Category", "Invalid Page Number", "Server Error", "Invalid Article Id"}

var categories = map[string]bool{
	"india":         true,
	"worldnews":     true,
	"science":       true,
	"technology":    true,
	"news":          true,
	"canada":        true,
	"unitedkingdom": true,
	"europe":        true,
	"china":         true,
	"upliftingnews": true}

func getData(w http.ResponseWriter, req *http.Request) {

	var result []article

	vars := mux.Vars(req)

	_, isPresent := categories[vars["category"]]
	pageNo, pageErr := strconv.Atoi(vars["pageNo"])

	if !isPresent {
		w.WriteHeader(400)
		w.Write([]byte(errors.InvalidCategory))
		return
	}

	if pageNo < 0 || pageErr != nil {
		w.WriteHeader(400)
		w.Write([]byte(errors.InvalidPageNo))
		return
	}


	session := common.GetMongoSession(
		dbConfig.url,
		dbConfig.authDB,
		dbConfig.username,
		dbConfig.password)

	collection := session.DB(dbConfig.dbName).C(vars["category"])

	err := collection.
		Find(nil).
		Sort("-date_added", "-score").
		Skip(pageNo * 18).
		Limit(18).
		All(&result)

	jData, err := json.Marshal(result)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		w.Write([]byte(errors.ServerError))
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
		return
	}

}

func getArticle(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	result := article{}


	_, isPresent := categories[vars["category"]]

	if !isPresent {
		w.WriteHeader(400)
		w.Write([]byte(errors.InvalidCategory))
	}

	session := common.GetMongoSession(dbConfig.url, dbConfig.authDB, dbConfig.username, dbConfig.password)

	collection := session.DB(dbConfig.dbName).C(vars["category"])

	err := collection.Find(bson.M{"reddit_id": vars["articleId"]}).One(&result)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(errors.InvalidArticle))
		return
	}

	jData, err := json.Marshal(result)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(errors.ServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	return

}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/getData/{category}/{pageNo}", getData).Methods("GET")
	router.HandleFunc("/getArticle/{category}/{articleId}", getArticle).Methods("GET")

	log.Fatal(http.ListenAndServe(":12345", router))
	fmt.Print("Hi")
}
