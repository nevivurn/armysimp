package scrape

import (
	"context"
	"time"

	"cloud.google.com/go/translate"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/language"
	"google.golang.org/api/youtube/v3"

	"github.com/nevivurn/armysimp/report"
)

// Client is a scraping client.
type Client struct {
	tl *translate.Client
	yt *youtube.Service
}

// New creates a new Client.
func New(tl *translate.Client, yt *youtube.Service) *Client {
	return &Client{tl, yt}
}

// Scrape youtube and related data.
func (c *Client) Scrape(ctx context.Context) (report.Data, error) {
	data := report.Data{
		Timestamp:   time.Now(),
		Generations: make([]report.Generation, len(seedData)),
	}

	// Load seed data
	mapYoutubeToChan := make(map[string]*report.Channel)
	for i, gen := range seedData {
		data.Generations[i] = report.Generation{
			Title:    gen.title,
			Channels: make([]report.Channel, len(gen.channels)),
		}
		for j, ch := range gen.channels {
			data.Generations[i].Channels[j].Title = ch.title
			data.Generations[i].Channels[j].Favorite = ch.favorite
			mapYoutubeToChan[ch.youtubeID] = &data.Generations[i].Channels[j]
		}
	}

	ytPlayIDs, err := c.fetchChannels(ctx, mapYoutubeToChan)
	if err != nil {
		return report.Data{}, err
	}

	titles, err := c.fetchPlaylists(ctx, mapYoutubeToChan, ytPlayIDs)
	if err != nil {
		return report.Data{}, err
	}

	err = c.fetchTranslations(ctx, titles)
	if err != nil {
		return report.Data{}, err
	}

	return data, nil
}

// Fetch channels, fill in subscriber counts, and return playlist IDs.
func (c *Client) fetchChannels(ctx context.Context, mapYoutubeToChan map[string]*report.Channel) ([]string, error) {
	ytChanIDs := make([]string, 0, len(mapYoutubeToChan))
	for ch := range mapYoutubeToChan {
		ytChanIDs = append(ytChanIDs, ch)
	}

	resp, err := c.yt.Channels.
		List([]string{"contentDetails", "id", "statistics"}).
		Context(ctx).
		Id(ytChanIDs...).
		Do()
	if err != nil {
		return nil, err
	}

	ytPlayIDs := make([]string, len(resp.Items))
	for i, ch := range resp.Items {
		mapYoutubeToChan[ch.Id].Subscribers = ch.Statistics.SubscriberCount
		ytPlayIDs[i] = ch.ContentDetails.RelatedPlaylists.Uploads
	}

	return ytPlayIDs, nil
}

type translationPair struct {
	src string
	dst *string
}

// Fetch most recent playlist items and set video titles, returns list of titles.
func (c *Client) fetchPlaylists(ctx context.Context, mapYoutubeToChan map[string]*report.Channel, ytPlayIDs []string) ([]translationPair, error) {
	errg, ctx := errgroup.WithContext(ctx)

	titlesCh := make(chan translationPair)
	done := make(chan struct{})
	var titles []translationPair
	go func() {
		for title := range titlesCh {
			titles = append(titles, title)
		}
		close(done)
	}()

	f := func(play string) error {
		resp, err := c.yt.PlaylistItems.
			List([]string{"snippet"}).
			Context(ctx).
			PlaylistId(play).
			MaxResults(3).
			Do()
		if err != nil {
			return err
		}
		if len(resp.Items) == 0 {
			return nil
		}

		ch := resp.Items[0].Snippet.ChannelId
		mapYoutubeToChan[ch].Videos = make([]report.Video, len(resp.Items))

		for i, vid := range resp.Items {
			mapYoutubeToChan[ch].Videos[i].Title = vid.Snippet.Title
			titlesCh <- translationPair{
				src: vid.Snippet.Title,
				dst: &mapYoutubeToChan[ch].Videos[i].TitleTranslated,
			}
		}

		return nil
	}

	for _, play := range ytPlayIDs {
		play := play
		errg.Go(func() error { return f(play) })
	}

	err := errg.Wait()
	close(titlesCh)
	<-done

	return titles, err
}

// Fetch translations for titles
func (c *Client) fetchTranslations(ctx context.Context, titles []translationPair) error {
	errg, ctx := errgroup.WithContext(ctx)

	f := func(tls []translationPair) error {
		srcs := make([]string, len(tls))
		for i, tl := range tls {
			srcs[i] = tl.src
		}

		resp, err := c.tl.Translate(ctx, srcs, language.English, &translate.Options{
			Source: language.Japanese,
			Format: translate.Text,
		})
		if err != nil {
			return err
		}

		for i, tl := range resp {
			*tls[i].dst = tl.Text
		}

		return nil
	}

	for len(titles) > 0 {
		end := 5
		if len(titles) < end {
			end = len(titles)
		}

		cur := titles[:end]
		titles = titles[end:]
		errg.Go(func() error { return f(cur) })
	}

	return errg.Wait()
}
