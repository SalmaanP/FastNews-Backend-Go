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

type EchoRequest struct {
	Version string      `json:"version"`
	Session EchoSession `json:"session"`
	Request EchoReqBody `json:"request"`
	Context EchoContext `json:"context"`
}

type EchoSession struct {
	New         bool   `json:"new"`
	SessionID   string `json:"sessionId"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	Attributes map[string]interface{} `json:"attributes"`
	User       struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"user"`
}

type EchoContext struct {
	System struct {
		Device struct {
			DeviceId string `json:"deviceId,omitempty"`
		} `json:"device,omitempty"`
		Application struct {
			ApplicationID string `json:"applicationId,omitempty"`
		} `json:"application,omitempty"`
	} `json:"System,omitempty"`
}

type EchoReqBody struct {
	Type      string     `json:"type"`
	RequestID string     `json:"requestId"`
	Timestamp string     `json:"timestamp"`
	Intent    EchoIntent `json:"intent,omitempty"`
	Reason    string     `json:"reason,omitempty"`
	Locale    string     `json:"locale,omitempty"`
}

type EchoIntent struct {
	Name  string              `json:"name"`
	Slots map[string]EchoSlot `json:"slots"`
}

type EchoSlot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Response Types

type EchoResponse struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          EchoRespBody           `json:"response"`
}

type EchoRespBody struct {
	OutputSpeech     *EchoRespPayload `json:"outputSpeech,omitempty"`
	Card             *EchoRespPayload `json:"card,omitempty"`
	Reprompt         *EchoReprompt    `json:"reprompt,omitempty"` // Pointer so it's dropped if empty in JSON response.
	ShouldEndSession bool             `json:"shouldEndSession"`
}

type EchoReprompt struct {
	OutputSpeech EchoRespPayload `json:"outputSpeech,omitempty"`
}

type EchoRespImage struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

type EchoRespPayload struct {
	Type    string        `json:"type,omitempty"`
	Title   string        `json:"title,omitempty"`
	Text    string        `json:"text,omitempty"`
	SSML    string        `json:"ssml,omitempty"`
	Content string        `json:"content,omitempty"`
	Image   EchoRespImage `json:"image,omitempty"`
}