package main

import (
	"net/http"
	"strconv"
	"encoding/json"
	"log"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"github.com/go-redis/redis"
	"time"
)

var dbConfig = getDbDetails()
var redisConfig = getRedisDetails()

var client = redis.NewClient(&redis.Options{
	Addr:     redisConfig.url,
	Password: redisConfig.password, // no password set
	DB:       redisConfig.db,  // use default DB
})

var errors = ErrorMap{"Invalid Category", "Invalid Page Number", "Server Error", "Invalid Article Id"}

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

	redisValue, redisError := client.Get(vars["category"] + vars["pageNo"]).Result()

	if redisError == nil {
		//From Cache

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(redisValue))
		return

	} else {

		//From DB

		var result Articles

		session := GetMongoSession(
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

			redisTest := vars["category"] + vars["pageNo"]

			client.Set(redisTest, jData, time.Hour)

			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}
	}

}

func getArticle(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	result := Article{}

	_, isPresent := categories[vars["category"]]

	if !isPresent {
		w.WriteHeader(400)
		w.Write([]byte(errors.InvalidCategory))
		return
	}

	redisValue, redisError := client.Get(vars["articleId"]).Result()

	if redisError == nil {

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(redisValue))
		return

	} else {
		session := GetMongoSession(dbConfig.url, dbConfig.authDB, dbConfig.username, dbConfig.password)

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

		client.Set(vars["articleId"], jData, time.Hour)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
		return
	}

}

func searchArticles(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	var result []Article

	_, isPresent := categories[vars["category"]]

	if !isPresent {
		w.WriteHeader(400)
		w.Write([]byte(errors.InvalidCategory))
		return
	}

	pageNo, pageErr := strconv.Atoi(vars["pageNo"])

	if pageNo < 0 || pageErr != nil {
		w.WriteHeader(400)
		w.Write([]byte(errors.InvalidPageNo))
		return
	}

	redisSearchString := "search" + vars["category"] + vars["searchString"] + vars["pageNo"]

	redisValue, redisError := client.Get(redisSearchString).Result()

	if redisError == nil {

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(redisValue))
		return

	} else {

		session := GetMongoSession(dbConfig.url, dbConfig.authDB, dbConfig.username, dbConfig.password)
		collection := session.DB(dbConfig.dbName).C(vars["category"])

		err := collection.
			Find(bson.M{"article_text": bson.RegEx{Pattern: vars["searchString"], Options: "i"}}).
			Sort("-date_added", "-score").
			Skip(pageNo * 18).
			Limit(18).
			All(&result)

		if err != nil {
			fmt.Print(err)
			w.WriteHeader(500)
			w.Write([]byte(errors.ServerError))
			return
		}

		jData, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(errors.ServerError))
			return
		}

		client.Set(redisSearchString, jData, time.Hour)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
		return

	}



}

func getRoot(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Fastnews by Salmaan Pehlari"))
	w.Write([]byte("\n\n"))
	w.Write([]byte("Supported endpoints: /getData/<category>/<pageNo>, /getArticle/<category>/<articleId>, /search/<searchString>/<category>/<pageNo>"))
}
