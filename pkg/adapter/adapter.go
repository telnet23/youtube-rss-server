package adapter

import (
	"bytes"
	"encoding/xml"
	"strconv"

	"google.golang.org/api/youtube/v3"
)

const (
	XMLNSYT    = "http://www.youtube.com/xml/schemas/2015"
	XMLNSMedia = "http://search.yahoo.com/mrss/"
	XMLNS      = "http://www.w3.org/2005/Atom"
)

type Feed struct {
	XMLName      xml.Name `xml:"feed"`
	XMLNSYT      string   `xml:"xmlns:yt,attr"`
	XMLNSMedia   string   `xml:"xmlns:media,attr"`
	XMLNS        string   `xml:"xmlns,attr"`
	ID           string   `xml:"id"`
	YTPlaylistID string   `xml:"yt:playlistId,omitempty"`
	YTChannelID  string   `xml:"yt:channelId"`
	Title        string   `xml:"title"`
	Links        []Link   `xml:"link"`
	Author       Author   `xml:"author"`
	Published    string   `xml:"published"`
	Entries      []*Entry `xml:"entry"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type Author struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

type Entry struct {
	ID          string     `xml:"id"`
	YTVideoID   string     `xml:"yt:videoId"`
	YTChannelID string     `xml:"yt:channelId"`
	Title       string     `xml:"title"`
	Link        Link       `xml:"link"`
	Author      Author     `xml:"author"`
	Published   string     `xml:"published"`
	Updated     string     `xml:"updated"`
	MediaGroup  MediaGroup `xml:"media:group"`
}

type MediaGroup struct {
	MediaTitle       string         `xml:"media:title"`
	MediaContent     MediaContent   `xml:"media:content"`
	MediaThumbnail   MediaThumbnail `xml:"media:thumbnail"`
	MediaDescription string         `xml:"media:description"`
	// MediaCommunity   MediaCommunity `xml:"media:community"`
}

type MediaContent struct {
	URL string `xml:"url,attr"`
}

type MediaThumbnail struct {
	URL    string `xml:"url,attr"`
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
}

// type MediaCommunity struct {
// 	MediaStarRating MediaStarRating `xml:"media:starRating"`
// 	MediaStatistics MediaStatistics `xml:"media:statistics"`
// }

// type MediaStarRating struct {
// 	Count   string `xml:"count,attr"`
// 	Average string `xml:"average,attr"`
// 	Min     string `xml:"min,attr"`
// 	Max     string `xml:"max,attr"`
// }

// type MediaStatistics struct {
// 	Views string `xml:"views,attr"`
// }

func channelURL(channelID string) string {
	return "https://www.youtube.com/channel/" + channelID
}

func playlistURL(playlistID string) string {
	return "https://www.youtube.com/playlist?list=" + playlistID
}

func videoURL(videoID string) string {
	return "https://www.youtube.com/watch?v=" + videoID
}

func FeedFromChannel(channel *youtube.Channel, username string) *Feed {
	selfHref := "http://www.youtube.com/feeds/videos.xml?channel_id=" + channel.Id
	if username != "" {
		selfHref = "http://www.youtube.com/feeds/videos.xml?user=" + username
	}

	return &Feed{
		XMLNSYT:     XMLNSYT,
		XMLNSMedia:  XMLNSMedia,
		XMLNS:       XMLNS,
		ID:          "yt:channel:" + channel.Id,
		YTChannelID: channel.Id,
		Title:       channel.Snippet.Title,
		Links: []Link{
			{
				Rel:  "self",
				Href: selfHref,
			},
			{
				Rel:  "alternate",
				Href: channelURL(channel.Id),
			},
		},
		Author: Author{
			Name: channel.Snippet.Title,
			URI:  channelURL(channel.Id),
		},
		Published: channel.Snippet.PublishedAt,
	}
}

func FeedFromPlaylist(playlist *youtube.Playlist) *Feed {
	return &Feed{
		XMLNSYT:      XMLNSYT,
		XMLNSMedia:   XMLNSMedia,
		XMLNS:        XMLNS,
		ID:           "yt:playlist:" + playlist.Id,
		YTChannelID:  playlist.Snippet.ChannelId,
		YTPlaylistID: playlist.Id,
		Title:        playlist.Snippet.Title,
		Links: []Link{
			{
				Rel:  "self",
				Href: "http://www.youtube.com/feeds/videos.xml?playlist_id=" + playlist.Id,
			},
			{
				Rel:  "alternate",
				Href: playlistURL(playlist.Id),
			},
		},
		Author: Author{
			Name: playlist.Snippet.ChannelTitle,
			URI:  channelURL(playlist.Snippet.ChannelId),
		},
		Published: playlist.Snippet.PublishedAt,
	}
}

func EntryFromPlaylistItem(playlistItem *youtube.PlaylistItem) *Entry {
	// This happens when the video is private or deleted
	if playlistItem.Snippet.Thumbnails.High == nil {
		return nil
	}

	return &Entry{
		ID:          "yt:video:" + playlistItem.Snippet.ResourceId.VideoId,
		YTVideoID:   playlistItem.Snippet.ResourceId.VideoId,
		YTChannelID: playlistItem.Snippet.VideoOwnerChannelId,
		Title:       playlistItem.Snippet.Title,
		Link: Link{
			Rel:  "alternate",
			Href: videoURL(playlistItem.Snippet.ResourceId.VideoId),
		},
		Author: Author{
			Name: playlistItem.Snippet.VideoOwnerChannelTitle,
			URI:  channelURL(playlistItem.Snippet.VideoOwnerChannelId),
		},
		Published: playlistItem.ContentDetails.VideoPublishedAt,
		Updated:   playlistItem.Snippet.PublishedAt,
		MediaGroup: MediaGroup{
			MediaTitle: playlistItem.Snippet.Title,
			MediaContent: MediaContent{
				URL: videoURL(playlistItem.Snippet.ResourceId.VideoId),
			},
			MediaThumbnail: MediaThumbnail{
				URL:    playlistItem.Snippet.Thumbnails.High.Url,
				Width:  strconv.FormatInt(playlistItem.Snippet.Thumbnails.High.Width, 10),
				Height: strconv.FormatInt(playlistItem.Snippet.Thumbnails.High.Height, 10),
			},
			MediaDescription: playlistItem.Snippet.Description,
		},
	}
}

func EncodeFeedXML(feed *Feed) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	enc := xml.NewEncoder(&buf)
	enc.Indent("", " ")
	if err := enc.Encode(feed); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
