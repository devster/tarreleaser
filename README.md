# TarReleaser

[![Release](https://img.shields.io/github/release/devster/tarreleaser.svg?style=for-the-badge)](https://github.com/devster/tarreleaser/releases/latest)
[![Travis](https://img.shields.io/travis/devster/tarreleaser/master.svg?style=for-the-badge)](https://travis-ci.org/devster/tarreleaser)

**Under development**

Heavily inspired by [GoReleaser](https://github.com/goreleaser/goreleaser)

## Install

### Downloader script

The downloader script will take care of the checksum validation

```bash
# Download, extract and install the latest version in ./bin dir
curl -sSL https://git.io/tarreleaser | bash

# Specify a tag and install dir
curl -sSL https://git.io/tarreleaser | sudo bash -s -- -b /usr/local/bin 0.1.0-alpha
```

### Go install

	go install github.com/devster/tarreleaser

### Binaries

Download the already compiled binary you want -> [releases](https://github.com/devster/tarreleaser/releases)

## CI

Oneliner to install and execute latest version of tarreleaser

Example with travis-ci:

```yaml
deploy:
  - provider: script
    skip_cleanup: true
    script: curl -sSL https://git.io/tarreleaser | bash && bin/tarreleaser
    on:
      tags: true
```

## Dev

Build a dev version

	make 

Run tests

	make test
