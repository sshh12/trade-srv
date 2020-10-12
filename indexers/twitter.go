package indexers

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const twitterSource string = "twitter"

func startTwitterIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	if opts.TwitterConsumerKey == "" || opts.TwitterAccessToken == "" {
		log.Print("No twitter login provided")
		return nil
	}
	if len(opts.TwitterNames) == 0 || opts.TwitterNames[0] == "" {
		log.Print("No twitter names provided")
		return nil
	}
	log.Println("Listening to tweets from", opts.TwitterNames)

	config := oauth1.NewConfig(opts.TwitterConsumerKey, opts.TwitterConsumerSecret)
	token := oauth1.NewToken(opts.TwitterAccessToken, opts.TwitterAccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	followIDs := make([]string, 0)
	for _, name := range opts.TwitterNames {
		user, _, err := client.Users.Show(&twitter.UserShowParams{
			ScreenName: name,
		})
		if err != nil {
			log.Print(err)
			return err
		}
		followIDs = append(followIDs, user.IDStr)
	}
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		isOG := false
		for _, id := range followIDs {
			if tweet.User.IDStr == id {
				isOG = true
				break
			}
		}
		if !isOG {
			return
		}
		var tweetMap map[string]interface{}
		jsonTweet, _ := json.Marshal(tweet)
		json.Unmarshal(jsonTweet, &tweetMap)
		evt := &events.Event{
			Source:       twitterSource,
			Type:         "tweet",
			Author:       tweet.User.ScreenName,
			TimeReported: tweet.CreatedAt,
			Content:      cleanTweet(tweetToText(tweet)),
			Extras:       tweetMap,
			CacheHash:    events.HashKey(tweet.IDStr),
		}
		es.OnEvent(evt)
	}
	filterParams := &twitter.StreamFilterParams{
		Follow:        followIDs,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Print(err)
		return err
	}
	demux.HandleChan(stream.Messages)
	return nil
}

func cleanTweet(text string) string {
	cleanText := strings.ReplaceAll(text, "…", "...")
	cleanText = scraping.RegexReplace(cleanText, "https:\\/\\/t.co\\/\\w+", "")
	return strings.TrimSpace(cleanText)
}

func tweetToText(tweet *twitter.Tweet) string {
	if tweet.RetweetedStatus != nil && tweet.RetweetedStatus.ExtendedTweet != nil && tweet.RetweetedStatus.ExtendedTweet.FullText != "" {
		return tweet.RetweetedStatus.ExtendedTweet.FullText
	} else if tweet.RetweetedStatus != nil && tweet.RetweetedStatus.Text != "" {
		return tweet.RetweetedStatus.Text
	} else if tweet.FullText != "" {
		return tweet.FullText
	}
	return tweet.Text
}
