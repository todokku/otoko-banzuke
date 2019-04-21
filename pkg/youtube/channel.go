package youtube

import (
	"log"
	"time"
)

// Channel fa
type Channel struct {
	ChannelID       string    `gorm:"column:channel_id"`
	Name            string    `gorm:"column:name"`
	Description     string    `gorm:"column:description"`
	ThumbnailURL    string    `gorm:"column:thumbnail_url"`
	PlaylistID      string    `gorm:"column:playlist_id"`
	ViewCount       int64     `gorm:"column:view_count"`
	VideoCount      int32     `gorm:"column:video_count"`
	SubscriberCount int32     `gorm:"column:subscriber_count"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

// Select channel
func (c *Channel) Select() Channel {
	db := newGormConnect()
	defer db.Close()

	db.First(&c, "channel_id=?", c.ChannelID)

	return *c
}

// Insert Channel
func (c *Channel) Insert() error {
	db := newGormConnect()
	defer db.Close()

	r := db.Create(&c)
	log.Printf("Insert channel: %v\n", r)

	return r.Error
}

func (c *Channel) update() {

}

func (c *Channel) delete() {

}

func (c *Channel) selectVideos() []Video {
	db := newGormConnect()
	defer db.Close()

	videos := []Video{}
	db.Find(&videos, "channel_id=?", c.ChannelID)

	return videos
}

func (c *Channel) deleteVideos() {

}

// SetDetailInfo PlaylistID, ViewCount, SubscriberCount, VideoCount
func (c *Channel) SetDetailInfo() {
	service := newYoutubeService()
	call := service.Channels.List("snippet,contentDetails,statistics").
		Id(c.ChannelID).
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}
	item := response.Items[0]

	c.Name = item.Snippet.Title
	c.Description = item.Snippet.Description
	c.ThumbnailURL = item.Snippet.Thumbnails.High.Url
	c.PlaylistID = item.ContentDetails.RelatedPlaylists.Uploads
	c.ViewCount = int64(item.Statistics.ViewCount)
	c.SubscriberCount = int32(item.Statistics.SubscriberCount)
	c.VideoCount = int32(item.Statistics.VideoCount)
}

// GetNewVideos assume to run once a day
func (c *Channel) GetNewVideos() []Video {
	// put 1 day period afer video published
	beginAt := time.Now().Add(-time.Duration(24*2) * time.Hour).Format(time.RFC3339)
	endAt := time.Now().Add(-time.Duration(24) * time.Hour).Format(time.RFC3339)

	service := newYoutubeService()
	call := service.Search.List("id,snippet").
		Type("video").
		ChannelId(c.ChannelID).
		PublishedAfter(beginAt).
		PublishedBefore(endAt).
		Order("date").
		MaxResults(50)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	videos := []Video{}
	for _, item := range response.Items {
		videoID := item.Id.VideoId
		title := item.Snippet.Title
		description := item.Snippet.Description
		thumbnailURL := item.Snippet.Thumbnails.High.Url
		channelID := item.Snippet.ChannelId
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			log.Fatalf("%v", err)
		}

		video := Video{
			VideoID:      videoID,
			Title:        title,
			Description:  description,
			ThumbnailURL: thumbnailURL,
			ChannelID:    channelID,
			PublishedAt:  publishedAt,
		}
		videos = append(videos, video)
	}
	log.Printf("Get %v videos\n", len(videos))

	return videos
}

// GetAllVideos hoge
func (c *Channel) GetAllVideos(pageToken string) []Video {
	service := newYoutubeService()
	call := service.PlaylistItems.List("id,snippet,contentDetails").
		PlaylistId(c.PlaylistID).
		PageToken(pageToken).
		MaxResults(50)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	videos := []Video{}
	for _, item := range response.Items {
		videoID := item.ContentDetails.VideoId
		title := item.Snippet.Title
		description := item.Snippet.Description
		thumbnailURL := item.Snippet.Thumbnails.High.Url
		channelID := item.Snippet.ChannelId
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			log.Fatalf("%v", err)
		}

		video := Video{
			VideoID:      videoID,
			Title:        title,
			Description:  description,
			ThumbnailURL: thumbnailURL,
			ChannelID:    channelID,
			PublishedAt:  publishedAt,
		}

		videos = append(videos, video)
	}

	pageToken = response.NextPageToken
	if pageToken != "" {
		videos = append(videos, c.GetAllVideos(pageToken)...)
	}
	log.Printf("Get %v videos\n", len(videos))

	return videos
}

func getHighRatedVideos(playlistID string, pageToken string) {

}

// SearchChannels hoge
func SearchChannels(q string) []Channel {
	service := newYoutubeService()
	call := service.Search.List("id,snippet").
		Type("channel").
		Q(q).
		Order("relevance").
		MaxResults(10)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	channels := []Channel{}
	for _, item := range response.Items {
		channelID := item.Id.ChannelId
		title := item.Snippet.Title
		description := item.Snippet.Description
		thumbnailURL := item.Snippet.Thumbnails.High.Url

		channel := Channel{
			ChannelID:    channelID,
			Name:         title,
			Description:  description,
			ThumbnailURL: thumbnailURL,
		}
		channels = append(channels, channel)
	}
	log.Printf("Get %v channels\n", len(channels))

	return channels
}

// SelectAllChannels hoge
func SelectAllChannels() []Channel {
	db := newGormConnect()
	defer db.Close()

	channels := []Channel{}
	db.Find(&channels)

	return channels
}
