package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"sort"
	"time"
)

var (
	//go:embed fixture/twitter/data/profile.js
	profileStr string
	//go:embed fixture/twitter/data/tweet.js
	tweetStr string

	profile Profile
	tweets  []*Tweet
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

	for i := range tweets {
		ti, err := time.Parse(time.RubyDate, tweets[i].Tweet.CreatedAt)
		if err != nil {
			log.Println("parse time fail", tweets[i].Tweet.ID, err)
		}
		tweets[i].CreatedAt = ti
	}

	log.Println("sort tweets")
	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.Unix() < tweets[j].CreatedAt.Unix()
	})

	log.Println(tweets[0].GetCreatedAtMonth())
	log.Println(tweets[len(tweets)-1].GetCreatedAtMonth())
}
