package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/telnet23/youtube-rss-server/pkg/adapter"
	"github.com/telnet23/youtube-rss-server/pkg/youtube"
)

var youtubeClient *youtube.Client

func handler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling request", "url", r.URL.String())

	username := r.URL.Query().Get("user")
	channelID := r.URL.Query().Get("channel_id")
	playlistID := r.URL.Query().Get("playlist_id")

	var feed *adapter.Feed
	if username != "" && channelID == "" && playlistID == "" {
		channel, err := youtubeClient.GetChannelForUsername(username)
		if err != nil {
			slog.Error("Failed to get channel", "username", username, "error", err)
			http.Error(w, "Failed to get channel", http.StatusInternalServerError)
			return
		}
		if channel == nil {
			http.Error(w, "Channel not found", http.StatusNotFound)
			return
		}

		feed = adapter.FeedFromChannel(channel, username)
		playlistID = channel.ContentDetails.RelatedPlaylists.Uploads
	} else if username == "" && channelID != "" && playlistID == "" {
		channel, err := youtubeClient.GetChannel(channelID)
		if err != nil {
			slog.Error("Failed to get channel", "channelID", channelID, "error", err)
			http.Error(w, "Failed to get channel", http.StatusInternalServerError)
			return
		}
		if channel == nil {
			http.Error(w, "Channel not found", http.StatusNotFound)
			return
		}

		feed = adapter.FeedFromChannel(channel, "")
		playlistID = channel.ContentDetails.RelatedPlaylists.Uploads
	} else if username == "" && channelID == "" && playlistID != "" {
		playlist, err := youtubeClient.GetPlaylist(playlistID)
		if err != nil {
			slog.Error("Failed to get playlist", "playlistID", playlistID, "error", err)
			http.Error(w, "Failed to get playlist", http.StatusInternalServerError)
			return
		}
		if playlist == nil {
			http.Error(w, "Playlist not found", http.StatusNotFound)
			return
		}

		feed = adapter.FeedFromPlaylist(playlist)
	} else {
		http.Error(w, "Exactly one parameter is required: channel_id, playlist_id, user", http.StatusBadRequest)
		return
	}

	playlistItems, err := youtubeClient.GetPlaylistItems(playlistID)
	if err != nil {
		slog.Error("Failed to get playlist items", "playlistID", playlistID, "error", err)
		http.Error(w, "Failed to get playlist items", http.StatusInternalServerError)
		return
	}

	for _, playlistItem := range playlistItems {
		if entry := adapter.EntryFromPlaylistItem(playlistItem); entry != nil {
			feed.Entries = append(feed.Entries, entry)
		}
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	content, err := adapter.EncodeFeedXML(feed)
	if err != nil {
		slog.Error("Failed to encode feed XML", "error", err)
		http.Error(w, "Failed to encode feed XML", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(content)
	if err != nil {
		slog.Error("Failed to write content", "error", err)
	}
}

func parseDurationEnv(key string, defaultValue time.Duration) time.Duration {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(rawValue)
	if err != nil {
		slog.Error("Failed to parse duration", "error", err)
		os.Exit(1)
	}
	return value
}

func parseIntEnv(key string, defaultValue int64) int64 {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		return defaultValue
	}

	value, err := strconv.ParseInt(rawValue, 10, 0)
	if err != nil {
		slog.Error("Failed to parse integer", "error", err)
		os.Exit(1)
	}
	return value
}

func main() {
	var err error
	youtubeClient, err = youtube.NewClient(
		os.Getenv("API_KEY"),
		parseDurationEnv("METADATA_CACHE_TTL", 12*time.Hour),
		parseDurationEnv("ITEMS_CACHE_TTL", 15*time.Minute),
		parseIntEnv("MAX_RESULTS", 15),
		parseDurationEnv("TIMEOUT", 30*time.Second),
	)
	if err != nil {
		slog.Error("Failed to create YouTube client", "error", err)
		os.Exit(1)
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	slog.Info("Listening on " + listenAddr)
	http.HandleFunc("/feeds/videos.xml", handler)
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		slog.Error("Failed to listen and serve", "error", err)
		os.Exit(1)
	}
}
