# Gatter (WIP)

A lightweight and performant ActivityPub server written in Go.
- Supports multiple users, each with his own domain (for easy cohosting)
- Native Postgres for performance
- Compatible with most existing Mastodon apps
- Theming support

## Building

```sh
git clone https://github.com/ruhrscholz/Gatter.git && cd Gatter
go build cmd/gatter/main.go
```

## Running (Development)

## Running (Production)

## Administration

Compile `gatterctl`:

```sh
go build cmd/gatterctl/main.go
```

## Scripts