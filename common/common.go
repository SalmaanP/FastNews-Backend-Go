package common

import (
	"gopkg.in/mgo.v2"
	"log"
	"crypto/tls"
	"net"
)

var mgoSession *mgo.Session

// Creates a new session if mgoSession is nil i.e there is no active mongo session.
//If there is an active mongo session it will return a Clone
func GetMongoSession(url string, authDB string, username string, password string) *mgo.Session {
	if mgoSession == nil {
		var err error

		tlsConfig := &tls.Config{}

		dialInfo := &mgo.DialInfo{

			Addrs:    []string{url},
			Database: authDB,
			Username: username,
			Password: password,
		}

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		mgoSession, err = mgo.DialWithInfo(dialInfo)

		if err != nil {
			log.Fatal("Failed to start the Mongo session")
		}
	}
	return mgoSession.Clone()
}
