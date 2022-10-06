package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dstotijn/go-notion"
)

type Profile struct {
	Profile struct {
		AvatarMediaUrl string `json:"avatarMediaUrl"`
	}
}

type Tweet struct {
	Tweet struct {
		ID            string `json:"id"`
		CreatedAt     string `json:"created_at"`
		FullText      string `json:"full_text"`
		FavoriteCount string `json:"favorite_count"`
		RetweetCount  string `json:"retweet_count"`
		Entities      struct {
			Media []struct {
				ExpandedURL string `json:"expanded_url"`
				// if this struct exists, meaning this tweet contains image
				MediaUrlHttps string `json:"media_url_https"`
			} `json:"media"`
			UserMentions []struct {
				Name       string `json:"name"`
				ScreenName string `json:"screen_name"`
				Id         string `json:"id"`
			} `json:"user_mentions"`
			Urls []struct {
				Url         string `json:"url"`
				ExpandedUrl string `json:"expanded_url"`
				DisplayUrl  string `json:"display_url"`
			} `json:"urls"`
			Hashtags []struct {
				Text string `json:"text"`
			} `json:"hashtags"`
		} `json:"entities"`
	} `json:"tweet"`
}

func (t *Tweet) GetCreatedAt() (time.Time, error) {
	ti, err := time.Parse(time.RubyDate, t.Tweet.CreatedAt)
	if err != nil {
		log.Println("parse time fail", t.Tweet.ID, err)
		return time.Now(), err
	}
	return ti, nil
}

func (t *Tweet) GetCreatedAtMonth() (string, error) {
	ti, err := t.GetCreatedAt()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d", ti.Year(), ti.Month()), nil
}

func addTweetToCallout(tweet *Tweet) error {
	ctx := context.TODO()

	month, err := tweet.GetCreatedAtMonth()
	if err != nil {
		return err
	}
	page := pageTree.GetByTitle(month)
	if page == nil {
		log.Println("page not found", month)
		return nil
	}

	_, err = clt.AppendBlockChildren(ctx, page.ID, convertTweetToBlock(tweet))
	if err != nil {
		log.Println("append callout to page failed", err)
	}
	return err
}

func convertTweetToBlock(tweet *Tweet) []notion.Block {
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
				Content: "@" + tweet.Tweet.ID,
				Link: &notion.Link{
					URL: fmt.Sprintf("https://twitter.com/anyone/status/%s", tweet.Tweet.ID),
				},
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
				Content: tweet.Tweet.FullText,
			},
		},
	}
	var childrenBlocks []notion.Block
	// if content contains image
	// FIXME:
	// image's urls in twitter archive are broken. this should be replace with other OSS URL
	// original images can be found in `twitter_media` directory.
	// {
	//     Object: "block",
	//     Type:   notion.BlockTypeImage,
	//     Image: &notion.FileBlock{
	//         Type:     notion.FileTypeExternal,
	//         External: &notion.FileExternal{URL: "https://pbs.twimg.com/media/CafaY3CUsAASD9O.jpg"},
	//     },
	// },

	// TODO:
	// detect quote tweet and pull content from twitter.com
	// this work nedd to import twitter bot and sdk
	// if this is a quote tweet
	// {
	//     Object: "block",
	//     Type:   notion.BlockTypeCallout,
	//     Callout: &notion.Callout{
	//         Icon: &notion.Icon{
	//             Type: notion.IconTypeExternal,
	//             External: &notion.FileExternal{
	//                 URL: "https://shuxiao.wang/favicon.ico",
	//             },
	//         },
	//         RichTextBlock: notion.RichTextBlock{
	//             Text: []notion.RichText{
	//                 {
	//                     // name
	//                     Type: notion.RichTextTypeText,
	//                     Text: &notion.Text{
	//                         Content: "keaising",
	//                     },
	//                     Annotations: &notion.Annotations{Bold: true},
	//                 },
	//                 {
	//                     // id link
	//                     Type: notion.RichTextTypeText,
	//                     Text: &notion.Text{
	//                         Content: "@this_is_a_tweet_id",
	//                         Link:    &notion.Link{URL: "https://twitter.com/a/xxxxx"},
	//                     },
	//                 },
	//             },
	//         },
	//     },
	// },

	blocks := []notion.Block{
		{
			Object: "block",
			Type:   notion.BlockTypeCallout,
			Callout: &notion.Callout{
				Icon: &notion.Icon{
					Type: notion.IconTypeExternal,
					External: &notion.FileExternal{
						URL: profile.Profile.AvatarMediaUrl,
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
