# Simple Nostr relay server

A simple Nostr relay server that stores events in a SQLite database.

Based on [fiatjaf/relayer](https://github.com/fiatjaf/relayer/tree/master/basic) basic relay implementation.

## Running

```
docker run -p 7447:7447 -v $PWD/data:/data ghcr.io/rubiojr/snr:latest
```

## Building

```
go build
```
