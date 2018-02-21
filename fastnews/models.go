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

type Alexa struct {
	Version string `json:"version"`
	Response AlexaResponse `json:"response"`
}

type AlexaResponse struct {
	OutputSpeech AlexaOutputSpeech `json:"outputSpeech"`
	Card AlexaCard `json:"card"`
	ShouldEndSession string `json:"shouldEndSession"`
}

type AlexaOutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type AlexaCard struct {
	Type string `json:"type"`
	Title string `json:"title"`
	Content string `json:"content"`
}
