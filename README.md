# Gatter (WIP)

A lightweight and performant ActivityPub server written in Go.
- Supports multiple users, each with his own domain (for easy cohosting)
- Postgres JSONB as storage backend
- Async task queue implemented with 
- Compatible with most existing Mastodon apps
- Blogging support (with ActivityPub implementation)
- Easy Plugin support

## Building

```sh
git clone https://github.com/ruhrscholz/gatter.git && cd gatter
go build .
```

## Running (Development)

## Running (Production)