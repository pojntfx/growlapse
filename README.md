# Growlapse

Visualize plant growth over time with Go, WebDAV and WASM; @pojntfx's entry for #growlab.

## Installation

### Containerized

You can get the Docker container like so:

```shell
$ docker pull pojntfx/growlapse-agent
```

### Natively

If you prefer a native installation, static binaries are also available on [GitHub releases](https://github.com/pojntfx/growlapse/releases).

You can install them like so:

```shell
$ curl -L -o /tmp/growlapse-agent https://github.com/pojntfx/growlapse/releases/download/latest/growlapse-agent.linux-$(uname -m)
$ sudo install /tmp/growlapse-agent /usr/local/bin
```

### About the Frontend

The frontend is also available on [GitHub releases](https://github.com/pojntfx/growlapse/releases) in the form of a static `.tar.gz` archive; to deploy it, simply upload it to a CDN or copy it to a web server. For most users, this shouldn't be necessary though; thanks to [@maxence-charriere](https://github.com/maxence-charriere)'s [go-app package](https://go-app.dev/), Growlapse is a progressive web app. By simply visiting the [public deployment](https://pojntfx.github.io/growlapse/) once, it will be available for offline use whenever you need it.

## Usage

🚧 This project is a work-in-progress! Instructions will be added as soon as it is usable. 🚧

## License

growlapse (c) 2021 Felicitas Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
