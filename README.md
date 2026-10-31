<picture>
  <source media="(prefers-color-scheme: dark)" srcset=".github/banner.svg">
  <source media="(prefers-color-scheme: light)" srcset=".github/banner-light.svg">
  <img alt="gitmon - GitHub issue and pull-request paths forming an RSS feed symbol" src=".github/banner-light.svg">
</picture>

A small server that turns GitHub issue and pull request searches into JSON, RSS and Atom feeds.

## Usage

1. Download the binary for your platform from [Releases](https://github.com/coalaura/gitmon/releases).
2. Copy [example.config.yml](example.config.yml) to `config.yml` in your working directory and edit the targets.
3. Run the binary (`chmod +x` first on Linux and macOS).

The server listens on port `9229` by default. Each target uses a GitHub search `filter` and a `type` of `issues` or `pull-requests`. Set `server.github_token` to authenticate requests and `server.base_url` to your public URL when hosting behind a proxy. Results are cached for `server.fetch_interval` (default `30m`).

For a target named `issues`:

- JSON: `http://localhost:9229/issues.json` (or `/issues`)
- RSS: `http://localhost:9229/issues.rss`
- Atom: `http://localhost:9229/issues.atom`

## Glance

Use the included [widget configuration](glance/example.config.yml) and [template](glance/gitmon.html) to display feeds in Glance. Update the URLs to match your server and configured targets.

## Releases

Push a `v*` tag to build signed Linux, macOS and Windows binaries for amd64 and arm64, with SHA-256 checksums and build provenance. The workflow requires the `COALAURA_PFX` and `COALAURA_PFX_PASSWORD` repository secrets for signing.
