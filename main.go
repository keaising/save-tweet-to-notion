package main

import (
	_ "embed"
	"log"

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
	config Config
	clt    *notion.Client
	tree   *TreeNode
)

func init() {
	_, err := toml.Decode(configStr, &config)
	if err != nil {
		log.Panicln("decode config.toml", err)
	}
	clt = notion.NewClient(config.Secret)
	tree, err = getRootTree()
	if err != nil {
		log.Panicln("get root tree failed", err)
	}
}

func main() {
	demo()
	// _ = createPageOnDate(time.Now().Unix())
	// t, err := getRootTree()
	// if err != nil {
	//     log.Println("iter tree failed", err)
	//     return
	// }
	// t.Print(0)
	_ = addTweetToCallout("3f605dc646364bbf894671f4830ad7e1", nil)
}
