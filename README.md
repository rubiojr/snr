# Simple Nostr relay server

A simple Nostr relay server that stores events in a SQLite database.

## Running

```
docker run -p 7447:7447 -v $PWD/data:/data snr
```

## Building

```
go build
```
