package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/youtube/v3"
)

var testChannel = &youtube.Channel{
	ContentDetails: &youtube.ChannelContentDetails{
		RelatedPlaylists: &youtube.ChannelContentDetailsRelatedPlaylists{
			Uploads: "UUAAAAAAAaaaaaaa00000000",
		},
	},
	Id:   "UCAAAAAAAaaaaaaa00000000",
	Kind: "youtube#channel",
	Snippet: &youtube.ChannelSnippet{
		PublishedAt: "2024-11-11T11:11:11Z",
		Title:       "Test Channel",
	},
}

var testPlaylist = &youtube.Playlist{
	Id:   "PLCCCCCCCCCCCccccccccccc2222222222",
	Kind: "youtube#playlist",
	Snippet: &youtube.PlaylistSnippet{
		ChannelId:    "UCBBBBBBBbbbbbbb11111111",
		ChannelTitle: "Other Test Channel",
		PublishedAt:  "2026-11-11T11:11:11Z",
		Title:        "Test Playlist",
	},
}

var testPlaylistItem = &youtube.PlaylistItem{
	ContentDetails: &youtube.PlaylistItemContentDetails{
		VideoPublishedAt: "2025-11-11T11:11:11Z",
	},
	Id:   "UEDDDDDDDDDDDDDDDdddddddddddddddd3333333333333333",
	Kind: "youtube#playlistItem",
	Snippet: &youtube.PlaylistItemSnippet{
		ChannelId:    "UCBBBBBBBbbbbbbb11111111",
		ChannelTitle: "Other Test Channel",
		Description:  "Test Description",
		PlaylistId:   "PLBBBBBBBBBBBbbbbbbbbbbb1111111111",
		PublishedAt:  "2026-11-11T11:11:11Z",
		ResourceId: &youtube.ResourceId{
			Kind:    "youtube#video",
			VideoId: "EEEEeeee444",
		},
		Thumbnails: &youtube.ThumbnailDetails{
			High: &youtube.Thumbnail{
				Height: 360,
				Url:    "https://i.ytimg.com/vi/EEEEeeee444/hqdefault.jpg",
				Width:  480,
			},
		},
		Title:                  "Test Title",
		VideoOwnerChannelId:    "UCAAAAAAAaaaaaaa00000000",
		VideoOwnerChannelTitle: "Test Channel",
	},
}

var testPlaylistItemDeleted = &youtube.PlaylistItem{
	Id:   "UEDDDDDDDDDDDDDDDdddddddddddddddd3333333333333333",
	Kind: "youtube#playlistItem",
	Snippet: &youtube.PlaylistItemSnippet{
		ChannelId:    "UCBBBBBBBbbbbbbb11111111",
		ChannelTitle: "Other Test Channel",
		Description:  "This video is unavailable.",
		PlaylistId:   "PLBBBBBBBBBBBbbbbbbbbbbb1111111111",
		PublishedAt:  "2026-11-11T11:11:11Z",
		ResourceId: &youtube.ResourceId{
			Kind:    "youtube#video",
			VideoId: "EEEEeeee444",
		},
		Thumbnails: &youtube.ThumbnailDetails{},
		Title:      "Deleted video",
	},
}

func TestFeedFromChannel(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected *Feed
	}{
		{
			name:     "from channel_id",
			username: "",
			expected: &Feed{
				XMLNSYT:     "http://www.youtube.com/xml/schemas/2015",
				XMLNSMedia:  "http://search.yahoo.com/mrss/",
				XMLNS:       "http://www.w3.org/2005/Atom",
				ID:          "yt:channel:UCAAAAAAAaaaaaaa00000000",
				YTChannelID: "UCAAAAAAAaaaaaaa00000000",
				Title:       "Test Channel",
				Links: []Link{
					{
						Rel:  "self",
						Href: "http://www.youtube.com/feeds/videos.xml?channel_id=UCAAAAAAAaaaaaaa00000000",
					},
					{
						Rel:  "alternate",
						Href: "https://www.youtube.com/channel/UCAAAAAAAaaaaaaa00000000",
					},
				},
				Author: Author{
					Name: "Test Channel",
					URI:  "https://www.youtube.com/channel/UCAAAAAAAaaaaaaa00000000",
				},
				Published: "2024-11-11T11:11:11Z",
				Entries:   nil,
			},
		},
		{
			name:     "from username",
			username: "TestUsername",
			expected: &Feed{
				XMLNSYT:     "http://www.youtube.com/xml/schemas/2015",
				XMLNSMedia:  "http://search.yahoo.com/mrss/",
				XMLNS:       "http://www.w3.org/2005/Atom",
				ID:          "yt:channel:UCAAAAAAAaaaaaaa00000000",
				YTChannelID: "UCAAAAAAAaaaaaaa00000000",
				Title:       "Test Channel",
				Links: []Link{
					{
						Rel:  "self",
						Href: "http://www.youtube.com/feeds/videos.xml?user=TestUsername",
					},
					{
						Rel:  "alternate",
						Href: "https://www.youtube.com/channel/UCAAAAAAAaaaaaaa00000000",
					},
				},
				Author: Author{
					Name: "Test Channel",
					URI:  "https://www.youtube.com/channel/UCAAAAAAAaaaaaaa00000000",
				},
				Published: "2024-11-11T11:11:11Z",
				Entries:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feed := FeedFromChannel(testChannel, tt.username)
			assert.Equal(t, tt.expected, feed)
		})
	}
}

func TestFeedFromPlaylist(t *testing.T) {
	expected := &Feed{
		XMLNSYT:      "http://www.youtube.com/xml/schemas/2015",
		XMLNSMedia:   "http://search.yahoo.com/mrss/",
		XMLNS:        "http://www.w3.org/2005/Atom",
		ID:           "yt:playlist:PLCCCCCCCCCCCccccccccccc2222222222",
		YTChannelID:  "UCBBBBBBBbbbbbbb11111111",
		YTPlaylistID: "PLCCCCCCCCCCCccccccccccc2222222222",
		Title:        "Test Playlist",
		Links: []Link{
			{
				Rel:  "self",
				Href: "http://www.youtube.com/feeds/videos.xml?playlist_id=PLCCCCCCCCCCCccccccccccc2222222222",
			},
			{
				Rel:  "alternate",
				Href: "https://www.youtube.com/playlist?list=PLCCCCCCCCCCCccccccccccc2222222222",
			},
		},
		Author: Author{
			Name: "Other Test Channel",
			URI:  "https://www.youtube.com/channel/UCBBBBBBBbbbbbbb11111111",
		},
		Published: "2026-11-11T11:11:11Z",
		Entries:   nil,
	}

	feed := FeedFromPlaylist(testPlaylist)
	assert.Equal(t, expected, feed)
}

func TestEntryFromPlaylistItem(t *testing.T) {
	expected := &Entry{
		ID:          "yt:video:EEEEeeee444",
		YTVideoID:   "EEEEeeee444",
		YTChannelID: "UCAAAAAAAaaaaaaa00000000",
		Title:       "Test Title",
		Link: Link{
			Rel:  "alternate",
			Href: "https://www.youtube.com/watch?v=EEEEeeee444",
		},
		Author: Author{
			Name: "Test Channel",
			URI:  "https://www.youtube.com/channel/UCAAAAAAAaaaaaaa00000000",
		},
		Published: "2025-11-11T11:11:11Z",
		Updated:   "2026-11-11T11:11:11Z",
		MediaGroup: MediaGroup{
			MediaTitle: "Test Title",
			MediaContent: MediaContent{
				URL: "https://www.youtube.com/watch?v=EEEEeeee444",
			},
			MediaThumbnail: MediaThumbnail{
				URL:    "https://i.ytimg.com/vi/EEEEeeee444/hqdefault.jpg",
				Width:  "480",
				Height: "360",
			},
			MediaDescription: "Test Description",
		},
	}

	entry := EntryFromPlaylistItem(testPlaylistItem)
	assert.Equal(t, expected, entry)
}

func TestEntryFromPlaylistItemDeleted(t *testing.T) {
	entry := EntryFromPlaylistItem(testPlaylistItemDeleted)
	assert.Nil(t, entry)
}

func TestEncodeFeedFromChannel(t *testing.T) {
	feed := FeedFromChannel(testChannel, "")
	entry := EntryFromPlaylistItem(testPlaylistItem)
	feed.Entries = append(feed.Entries, entry)
	content, err := EncodeFeedXML(feed)
	assert.NoError(t, err)
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns:yt="http://www.youtube.com/xml/schemas/2015" xmlns:media="http://search.yahoo.com/mrss/" xmlns="http://www.w3.org/2005/Atom">
 <id>yt:channel:UCAAAAAAAaaaaaaa00000000</id>
 <yt:channelId>UCAAAAAAAaaaaaaa00000000</yt:channelId>
 <title>Test Channel</title>
 <link rel="self" href="http://www.youtube.com/feeds/videos.xml?channel_id=UCAAAAAAAaaaaaaa00000000"></link>
 <link rel="alternate" href="https://www.youtube.com/channel/UCAAAAAAAaaaaaaa00000000"></link>
 <author>
  <name>Test Channel</name>
  <uri>https://www.youtube.com/channel/UCAAAAAAAaaaaaaa00000000</uri>
 </author>
 <published>2024-11-11T11:11:11Z</published>
 <entry>
  <id>yt:video:EEEEeeee444</id>
  <yt:videoId>EEEEeeee444</yt:videoId>
  <yt:channelId>UCAAAAAAAaaaaaaa00000000</yt:channelId>
  <title>Test Title</title>
  <link rel="alternate" href="https://www.youtube.com/watch?v=EEEEeeee444"></link>
  <author>
   <name>Test Channel</name>
   <uri>https://www.youtube.com/channel/UCAAAAAAAaaaaaaa00000000</uri>
  </author>
  <published>2025-11-11T11:11:11Z</published>
  <updated>2026-11-11T11:11:11Z</updated>
  <media:group>
   <media:title>Test Title</media:title>
   <media:content url="https://www.youtube.com/watch?v=EEEEeeee444"></media:content>
   <media:thumbnail url="https://i.ytimg.com/vi/EEEEeeee444/hqdefault.jpg" width="480" height="360"></media:thumbnail>
   <media:description>Test Description</media:description>
  </media:group>
 </entry>
</feed>`, string(content))
}

func TestEncodeFeedFromPlaylist(t *testing.T) {
	feed := FeedFromPlaylist(testPlaylist)
	entry := EntryFromPlaylistItem(testPlaylistItem)
	feed.Entries = append(feed.Entries, entry)
	content, err := EncodeFeedXML(feed)
	assert.NoError(t, err)
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns:yt="http://www.youtube.com/xml/schemas/2015" xmlns:media="http://search.yahoo.com/mrss/" xmlns="http://www.w3.org/2005/Atom">
 <id>yt:playlist:PLCCCCCCCCCCCccccccccccc2222222222</id>
 <yt:playlistId>PLCCCCCCCCCCCccccccccccc2222222222</yt:playlistId>
 <yt:channelId>UCBBBBBBBbbbbbbb11111111</yt:channelId>
 <title>Test Playlist</title>
 <link rel="self" href="http://www.youtube.com/feeds/videos.xml?playlist_id=PLCCCCCCCCCCCccccccccccc2222222222"></link>
 <link rel="alternate" href="https://www.youtube.com/playlist?list=PLCCCCCCCCCCCccccccccccc2222222222"></link>
 <author>
  <name>Other Test Channel</name>
  <uri>https://www.youtube.com/channel/UCBBBBBBBbbbbbbb11111111</uri>
 </author>
 <published>2026-11-11T11:11:11Z</published>
 <entry>
  <id>yt:video:EEEEeeee444</id>
  <yt:videoId>EEEEeeee444</yt:videoId>
  <yt:channelId>UCAAAAAAAaaaaaaa00000000</yt:channelId>
  <title>Test Title</title>
  <link rel="alternate" href="https://www.youtube.com/watch?v=EEEEeeee444"></link>
  <author>
   <name>Test Channel</name>
   <uri>https://www.youtube.com/channel/UCAAAAAAAaaaaaaa00000000</uri>
  </author>
  <published>2025-11-11T11:11:11Z</published>
  <updated>2026-11-11T11:11:11Z</updated>
  <media:group>
   <media:title>Test Title</media:title>
   <media:content url="https://www.youtube.com/watch?v=EEEEeeee444"></media:content>
   <media:thumbnail url="https://i.ytimg.com/vi/EEEEeeee444/hqdefault.jpg" width="480" height="360"></media:thumbnail>
   <media:description>Test Description</media:description>
  </media:group>
 </entry>
</feed>`, string(content))
}
