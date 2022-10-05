package main

import (
	"context"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dstotijn/go-notion"
)

type Tweet struct {
	ID      string
	Time    time.Time
	Content string
}

func addTweetToCallout(pageID string, tweet *twitter.Tweet) error {
	ctx := context.TODO()
	_, err := clt.AppendBlockChildren(ctx, pageID, convertTweetToBlock(tweet))
	if err != nil {
		log.Println("append callout to page failed", err)
	}
	return err
}

func convertTweetToBlock(tweet *twitter.Tweet) []notion.Block {
	texts := []notion.RichText{
		{
			// name
			Type: notion.RichTextTypeText,
			Text: &notion.Text{
				Content: "keaising",
			},
			Annotations: &notion.Annotations{Bold: true},
		},
		{
			// id link
			Type: notion.RichTextTypeText,
			Text: &notion.Text{
				Content: "@this_is_a_tweet_id",
				Link:    &notion.Link{URL: "https://twitter.com/a/xxxxx"},
			},
		},
		{
			// \r\n
			Type: notion.RichTextTypeText,
			Text: &notion.Text{
				Content: "\r\n",
			},
		},
		{
			// content
			Type: notion.RichTextTypeText,
			Text: &notion.Text{
				Content: "This is a content paragraph",
			},
		},
	}
	childrenBlocks := []notion.Block{
		{
			// if content contains image
			Object: "block",
			Type:   notion.BlockTypeImage,
			Image: &notion.FileBlock{
				Type:     notion.FileTypeExternal,
				External: &notion.FileExternal{URL: "https://shuxiao.wang/images/king-in-the-dark/pre3.png"},
			},
		},
		{
			// if this is a quote tweet
			Object: "block",
			Type:   notion.BlockTypeCallout,
			Callout: &notion.Callout{
				Icon: &notion.Icon{
					Type: notion.IconTypeExternal,
					External: &notion.FileExternal{
						URL: "https://shuxiao.wang/favicon.ico",
					},
				},
				RichTextBlock: notion.RichTextBlock{
					Text: []notion.RichText{
						{
							// name
							Type: notion.RichTextTypeText,
							Text: &notion.Text{
								Content: "keaising",
							},
							Annotations: &notion.Annotations{Bold: true},
						},
						{
							// id link
							Type: notion.RichTextTypeText,
							Text: &notion.Text{
								Content: "@this_is_a_tweet_id",
								Link:    &notion.Link{URL: "https://twitter.com/a/xxxxx"},
							},
						},
					},
				},
			},
		},
	}

	blocks := []notion.Block{
		{
			Object: "block",
			Type:   notion.BlockTypeCallout,
			Callout: &notion.Callout{
				Icon: &notion.Icon{
					Type: notion.IconTypeExternal,
					External: &notion.FileExternal{
						URL: "https://shuxiao.wang/favicon.ico",
					},
				},
				RichTextBlock: notion.RichTextBlock{
					Text:     texts,
					Children: childrenBlocks,
				},
			},
		},
	}
	return blocks
}
