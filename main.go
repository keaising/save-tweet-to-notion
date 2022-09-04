package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/dstotijn/go-notion"
)

type Config struct {
	Secret string `toml:"secret"`
}

//go:embed config.toml
var configStr string

func main() {
	var config Config
	_, err := toml.Decode(configStr, &config)
	if err != nil {
		log.Panicln("decode config.toml", err)
	}
	fmt.Println("secret:", config.Secret)

	ctx := context.TODO()
	clt := notion.NewClient(config.Secret)
	result, err := clt.Search(ctx, &notion.SearchOpts{
		Query:    "2022.09",
		PageSize: 100,
	})
	// page, err := clt.CreatePage(ctx, notion.CreatePageParams{
	//     ParentType: notion.ParentTypePage,
	//     ParentID:   "cbfb36d7b05b4f74aca29b89b7c838be",
	//     Title: []notion.RichText{
	//         {
	//             Text: &notion.Text{
	//                 Content: "Create Page Example",
	//             },
	//         },
	//     },
	// })
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(len(result.Results))
}
