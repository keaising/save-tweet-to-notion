# save-tweet-to-notion

Save tweet archive into your notion pages

## 1. Results

### 1.1 Page Tree

Page tree will be created according to the date of all your tweet archive

![page-tree](/fixture/image/page-tree.png)

As you can see, first level is year, and month follows.

### 1.2 Tweets

Every tweet will be put into a notion `Callout`, just like:

![tweets-every-day](/fixture/image/tweets-every-day.png)

That's it.

## 2. Usage

### 2.1 Config

Change `config.toml.example` to `config.toml`

`secret`: Notion integration's secret
`root_page_id`: The root page which you want to write all tweets to

You can find these config in [Notion API Doc](https://developers.notion.com/docs), go through guide and things will be clear.

### 2.2 Twitter Archive

Download your twitter archive and put it into `fixture` directory like

![twitter-archive-directory](/fixture/image/twitter-archive-directory.png)

`tweet.js` and `profile.js` is needed.

### 2.3 Run

Install `go` in `go.dev`

Run

```
go run .
```

### 2.4 Time cost

For my over 30k tweets from 2012 to 2022, creating page tree will cost 30 to 40 minutes, writing all tweets into proper pages will cost over 10 hours.

Notions API rate limit is 3 requests per second, but their respond is very slow, please be patient.
