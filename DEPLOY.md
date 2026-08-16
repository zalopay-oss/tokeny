# Deploying

**Tokeny** uses [GoReleaser](https://goreleaser.com/) to build and deploy binaries to GitHub Releases.

To run a manual release with the same GoReleaser major version as the GitHub Actions workflow:

```bash
docker run --rm \
  -e GITHUB_TOKEN="$GITHUB_TOKEN" \
  -v "$(pwd):/go/src/github.com/zalopay-oss/tokeny" \
  -w /go/src/github.com/zalopay-oss/tokeny \
  goreleaser/goreleaser:v2.17.1 release --clean
```

* `$GITHUB_TOKEN` represents your GitHub's personal access token
