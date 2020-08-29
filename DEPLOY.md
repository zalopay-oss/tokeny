# Deploying

**Tokeny** is using [GoReleaser](https://goreleaser.com/) to build and deploy binaries to GitHub Releases.

In the meantime, GoReleaser is having [problem](https://github.com/goreleaser/goreleaser/issues/708) with cross compiling when `CGO_ENABLED` is required.

The solution is running GoReleaser inside a Docker container as proposed by @robdefeo at [goreleaser-xcgo](https://github.com/mailchain/goreleaser-xcgo). 

```bash
docker run --rm --privileged -e GITHUB_TOKEN=$GITHUB_TOKEN -v $(pwd):/go/src/github.com/zalopay-oss/tokeny -v /var/run/docker.sock:/var/run/docker.sock -w /go/src/github.com/zalopay-oss/tokeny mailchain/goreleaser-xcgo --rm-dist
```

* `$GITHUB_TOKEN` represents your GitHub's personal access token

