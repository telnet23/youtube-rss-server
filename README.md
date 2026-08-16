# youtube-rss-server

Provides RSS feeds for YouTube channels and playlists using the YouTube API.

## Supported URLs

- `http://localhost:8080/feeds/videos.xml?channel_id=CHANNEL_ID`
- `http://localhost:8080/feeds/videos.xml?playlist_id=PLAYLIST_ID`
- `http://localhost:8080/feeds/videos.xml?user=USER`

## Why?

YouTube provides equivalent RSS feeds at the following URLs:

- `https://www.youtube.com/feeds/videos.xml?channel_id=CHANNEL_ID`
- `https://www.youtube.com/feeds/videos.xml?playlist_id=PLAYLIST_ID`
- `https://www.youtube.com/feeds/videos.xml?user=USER`

However, these URLs have been returning 404 and 500 errors intermittently for several months, often for several hours at a time:

- https://community.n8n.io/t/youtube-rss-feed-endpoint-returns-404-errors/241692
- https://discuss.ai.google.dev/t/youtube-rss-feed-endpoint-returns-404-errors/113379
- https://www.reddit.com/r/rss/comments/1sowhwl/youtube_rss_again/
- https://www.reddit.com/r/youtube/comments/1r61jpo/all_youtube_channel_rss_feeds_are_down_return_404/

Running `youtube-rss-server` locally provides a drop-in replacement for these RSS feeds without the reliability issues.

## Configuration

| Environment Variable | Description                                                   | Default |
| -------------------- | ------------------------------------------------------------- | ------- |
| `API_KEY`            | API key for YouTube API requests                              | `""`    |
| `ITEMS_CACHE_TTL`    | Cache TTL for channel/playlist items (e.g. videos)            | `15m`   |
| `LISTEN_ADDR`        | Address and port to listen on for RSS requests                | `:8080` |
| `MAX_RESULTS`        | Maximum number of items to return                             | `15`    |
| `METADATA_CACHE_TTL` | Cache TTL for channel/playlist metadata (e.g. channel titles) | `12h`   |
| `TIMEOUT`            | Timeout for YouTube API requests                              | `30s`   |

If `API_KEY` is not provided, then Application Default Credentials must be provided:

- https://developers.google.com/youtube/v3/getting-started
- https://docs.cloud.google.com/docs/authentication/application-default-credentials

## Running Using Go

```
go install github.com/telnet23/youtube-rss-server@latest
API_KEY=... youtube-rss-server
```

## Running Using Pre-Built Binaries

Binaries can be downloaded from GitHub Releases: https://github.com/telnet23/youtube-rss-server/releases

```
tar xzf youtube-rss-server_*.tar.gz
API_KEY=... ./youtube-rss-server
```

## Running Using Pre-Built Docker Images

Docker images can be pulled from GitHub Container Registry:

```
docker run -e API_KEY=... ghcr.io/telnet23/youtube-rss-server
```
