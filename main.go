package main

import (
	_ "embed"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dstotijn/go-notion"
	"github.com/samber/lo"
	"golang.org/x/time/rate"
)

type Config struct {
	Secret     string `toml:"secret"`
	RootPageID string `toml:"root_page_id"`
}

//go:embed config.toml
var configStr string

var (
	config   Config
	pageTree *TreeNode
)

type notionClient struct {
	RateLimiter *rate.Limiter
	*notion.Client
}

var clt *notionClient

var timeZone = "Asia/Shanghai"

var loc *time.Location

func init() {
	_, err := toml.Decode(configStr, &config)
	if err != nil {
		log.Panicln("decode config.toml", err)
	}

	loc, err = time.LoadLocation(timeZone)
	if err != nil {
		log.Panicln("get timezone failed", err)
	}

	log.Println("init notion start")
	nClt := notion.NewClient(config.Secret)
	clt = &notionClient{
		rate.NewLimiter(rate.Every(time.Second), 1),
		nClt,
	}
	pageTree, err = getRootTree()
	if err != nil {
		log.Panicln("get root tree failed", err)
	}
	log.Println("init notion finished")
}

func main() {
	createdPages()

	// _ = addTweetToCallout(tweets[0], true)
	addTweets()
}

func createdPages() {
	// get all month
	months := make(map[int64]string)
	for _, tweet := range tweets {
		month := time.Date(tweet.CreatedAt.Year(), tweet.CreatedAt.Month(), 1, 0, 0, 0, 0, loc).Unix()
		months[month] = ""
	}
	firstDays := lo.MapToSlice(months, func(k int64, v string) time.Time {
		return time.Unix(k, 0)
	})
	sort.Slice(firstDays, func(i, j int) bool {
		return firstDays[i].Unix() < firstDays[j].Unix()
	})

	// create all pages
	for _, month := range firstDays {
		err := createPageOnDate(month)
		if err != nil {
			break
		}
	}
}

type tweetGroup struct {
	month  int64
	tweets []*Tweet
}

func addTweets() {
	// group tweets by month
	var tweetGroups []tweetGroup

	for k, v := range lo.GroupBy(tweets, func(tweet *Tweet) int64 {
		return time.Date(tweet.CreatedAt.Year(), tweet.CreatedAt.Month(), 1, 0, 0, 0, 0, loc).Unix()
	}) {
		tweetGroups = append(tweetGroups, tweetGroup{k, v})
	}
	sort.Slice(tweetGroups, func(i, j int) bool {
		return tweetGroups[i].month < tweetGroups[j].month
	})

	numJobs := 20
	jobs := make(chan tweetGroup, numJobs)
	var wg sync.WaitGroup

	for w := 0; w < 10; w++ {
		go worker(w, jobs, &wg)
	}
	for i := range tweetGroups {
		jobs <- tweetGroups[i]
		wg.Add(1)
	}
	wg.Wait()

	close(jobs)
}

func worker(id int, jobs <-chan tweetGroup, wg *sync.WaitGroup) {
	// append tweet to correct month
	for g := range jobs {
		today := g.tweets[0].CreatedAt.YearDay() - 1
		var isFirst bool
		log.Printf("%d begin  month %d.%d", id, time.Unix(g.month, 0).Year(), time.Unix(g.month, 0).Month())
		for _, tweet := range g.tweets {
			if tweet.CreatedAt.YearDay() != today {
				isFirst = true
				today = tweet.CreatedAt.YearDay()
			} else {
				isFirst = false
			}
			err := addTweetToCallout(tweet, isFirst)
			if err != nil {
				break
			}
		}
		log.Printf("%d finish month %d.%d", id, time.Unix(g.month, 0).Year(), time.Unix(g.month, 0).Month())
		wg.Done()
	}
}
