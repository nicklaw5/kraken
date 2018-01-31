# kraken

A Twitch Kraken API (v5) client written in Go. If you are looking for a client for Twitch's Helix API, see [helix](https://github.com/nicklaw5/helix).

[![Build Status](https://travis-ci.org/nicklaw5/kraken.svg?branch=master)](https://travis-ci.org/nicklaw5/kraken)
[![Coverage Status](https://coveralls.io/repos/github/nicklaw5/kraken/badge.svg)](https://coveralls.io/github/nicklaw5/kraken)

## Package Status

This project is a work in progress. Below is a list of currently supported endpoints. Happy for others to contribute.

## Supported Endpoints

- [ ] GET /bits/actions
- [x] GET /clips/:slug
- [x] GET /clips/top
- [x] GET /clips/followed
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

## License

This package is distributed under the terms of the [MIT](License) License.
