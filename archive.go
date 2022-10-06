package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"sort"
)

var (
	//go:embed fixture/twitter/data/profile.js
	profileStr string
	//go:embed fixture/twitter/data/tweet.js
	tweetStr string

	profile Profile
	tweets  []Tweet
)

func init() {
	profileStr = profileStr[28:]
	profileStr = profileStr[:len(profileStr)-1]

	tweetStr = tweetStr[25:]

	log.Println("init archive")
	err := json.Unmarshal([]byte(profileStr), &profile)
	if err != nil {
		log.Println("unmarshal profile", err)
	}

	err = json.Unmarshal([]byte(tweetStr), &tweets)
	if err != nil {
		log.Println("unmarshal tweets", err)
	}

	log.Println("sort tweets")
	sort.Slice(tweets, func(i, j int) bool {
		ti, _ := tweets[i].GetCreatedAt()
		tj, _ := tweets[j].GetCreatedAt()
		return ti.Unix() < tj.Unix()
	})

	log.Println(tweets[0].GetCreatedAtMonth())
	log.Println(tweets[len(tweets)-1].GetCreatedAtMonth())
}
