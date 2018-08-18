#!/bin/sh

set -ex

# ~/.config/goreleaser/github_token must contain your github token
# for testing add --skip-publish
goreleaser release --rm-dist
