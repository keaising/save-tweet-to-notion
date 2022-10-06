package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dstotijn/go-notion"
)

type Config struct {
	Secret     string `toml:"secret"`
	RootPageID string `toml:"root_page_id"`
}

//go:embed config.toml
var configStr string

var (
	config   Config
	clt      *notion.Client
	pageTree *TreeNode
)

func init() {
	_, err := toml.Decode(configStr, &config)
	if err != nil {
		log.Panicln("decode config.toml", err)
	}
	log.Println("init notion start")
	clt = notion.NewClient(config.Secret)
	pageTree, err = getRootTree()
	if err != nil {
		log.Panicln("get root tree failed", err)
	}
	log.Println("init notion finished")
}

func main() {
	log.Println("all tweets", len(tweets))
	for _, tweet := range tweets {
		ti, err := tweet.GetCreatedAt()
		if err != nil {
			break
		}
		err = createPageOnDate(ti)
		if err != nil {
			break
		}
	}
	for _, tweet := range tweets {
		ti, err := tweet.GetCreatedAt()
		if err != nil {
			break
		}
		if err != nil {
			break
		}
		log.Println(tweet.Tweet.ID, ti.Format(time.RFC3339))
	}
}
