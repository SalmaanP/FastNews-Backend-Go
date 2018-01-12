package main

import "time"

type Article struct {
	//Text  string `bson:"article_text"`
	Title string `bson:"title"`
	Summary string `bson:"summary"`
	Domain string `bson:"domain"`
	//Url   string `bson:"url"`
	Score int `bson:"score"`
	Comments int `bson:"num_comments"`
	Date time.Time `bson:"date_added"`
	Id string `bson:"reddit_id"`
	KeyPoints []string `bson:"keypoints"`
	Permalink string `bson:"permalink"`
	Url string `bson:"url"`
}

type ErrorMap struct {
	InvalidCategory string
	InvalidPageNo   string
	ServerError     string
	InvalidArticle  string
}

type DbMap struct {
	url      string
	authDB   string
	username string
	password string
	dbName   string
}

type RedisMap struct {
	url string
	password string
	db int
}

type Articles []Article