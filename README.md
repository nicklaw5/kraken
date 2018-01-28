# kraken

A Twitch Kraken API (v5) client written in Go.

## Supported Endpoints

- [ ] GET /bits/actions
- [ ] GET /clips/:slug
- [ ] GET /clips/top
- [ ] GET /clips/followed
- [ ] GET /feed/:channel-id/posts

## Getting Started

It's recommended that you use a dependency management tool such as [Dep](https://github.com/golang/dep). If you are using Dep you can import kraken by running:

```bash
dep ensure -add github.com/nicklaw5/kraken
```

Or you can simply import using the Go toolchain:

```bash
go get -u github.com/nicklaw5/kraken
```
