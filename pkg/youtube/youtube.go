package youtube

import (
	"context"
	"strings"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	service            *youtube.Service
	channelCache       *ttlcache.Cache[string, *youtube.Channel]
	playlistCache      *ttlcache.Cache[string, *youtube.Playlist]
	playlistItemsCache *ttlcache.Cache[string, []*youtube.PlaylistItem]
	maxResults         int64
	timeout            time.Duration
}

func NewClient(apiKey string, metadataCacheTTL, itemsCacheTTL time.Duration, maxResults int64, timeout time.Duration) (*Client, error) {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	c := &Client{
		service:    service,
		maxResults: maxResults,
		timeout:    timeout,
	}

	if metadataCacheTTL > 0 {
		c.channelCache = ttlcache.New(
			ttlcache.WithTTL[string, *youtube.Channel](metadataCacheTTL),
		)
		c.playlistCache = ttlcache.New(
			ttlcache.WithTTL[string, *youtube.Playlist](metadataCacheTTL),
		)
	}
	if itemsCacheTTL > 0 {
		c.playlistItemsCache = ttlcache.New(
			ttlcache.WithTTL[string, []*youtube.PlaylistItem](itemsCacheTTL),
		)
	}

	return c, nil
}

func (c *Client) GetChannelForUsername(username string) (*youtube.Channel, error) {
	username = strings.ToLower(username)

	if c.channelCache != nil && c.channelCache.Has(username) {
		return c.channelCache.Get(username).Value(), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	resp, err := c.service.Channels.
		List([]string{"contentDetails", "snippet"}).
		Context(ctx).
		ForUsername(username).
		Do()
	if err != nil {
		return nil, err
	}
	if len(resp.Items) < 1 {
		return nil, nil
	}

	if c.channelCache != nil {
		c.channelCache.Set(username, resp.Items[0], ttlcache.DefaultTTL)
	}
	return resp.Items[0], nil
}

func (c *Client) GetChannel(channelID string) (*youtube.Channel, error) {
	if c.channelCache != nil && c.channelCache.Has(channelID) {
		return c.channelCache.Get(channelID).Value(), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	resp, err := c.service.Channels.
		List([]string{"contentDetails", "snippet"}).
		Context(ctx).
		Id(channelID).
		Do()
	if err != nil {
		return nil, err
	}
	if len(resp.Items) < 1 {
		return nil, nil
	}

	if c.channelCache != nil {
		c.channelCache.Set(channelID, resp.Items[0], ttlcache.DefaultTTL)
	}
	return resp.Items[0], nil
}

func (c *Client) GetPlaylist(playlistID string) (*youtube.Playlist, error) {
	if c.playlistCache != nil && c.playlistCache.Has(playlistID) {
		return c.playlistCache.Get(playlistID).Value(), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	resp, err := c.service.Playlists.
		List([]string{"snippet"}).
		Context(ctx).
		Id(playlistID).
		Do()
	if err != nil {
		return nil, err
	}
	if len(resp.Items) < 1 {
		return nil, nil
	}

	if c.playlistCache != nil {
		c.playlistCache.Set(playlistID, resp.Items[0], ttlcache.DefaultTTL)
	}
	return resp.Items[0], nil
}

func (c *Client) GetPlaylistItems(playlistID string) ([]*youtube.PlaylistItem, error) {
	if c.playlistItemsCache != nil && c.playlistItemsCache.Has(playlistID) {
		return c.playlistItemsCache.Get(playlistID).Value(), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	resp, err := c.service.PlaylistItems.
		List([]string{"contentDetails", "snippet"}).
		Context(ctx).
		PlaylistId(playlistID).
		MaxResults(c.maxResults).
		Do()
	if err != nil {
		return nil, err
	}

	if c.playlistItemsCache != nil {
		c.playlistItemsCache.Set(playlistID, resp.Items, ttlcache.DefaultTTL)
	}
	return resp.Items, nil
}
