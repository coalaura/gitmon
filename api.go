package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/feeds"
)

func HandleRequest(config *Config, cache *Cache) http.HandlerFunc {
	targets := make(map[string]*Target, len(config.Targets))

	for i, target := range config.Targets {
		targets[target.Name] = &config.Targets[i]
	}

	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		target, ok := targets[name]
		if !ok {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		data, err := cache.Get(config, target, name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Warnln(err)

			return
		}

		format := chi.URLParam(r, "format")

		switch format {
		case "atom":
			format = "atom"

			w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
		case "rss":
			format = "rss"

			w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		default:
			format = "json"

			w.Header().Set("Content-Type", "application/feed+json; charset=utf-8")
		}

		uri := config.Server.baseUri.ResolveReference(&url.URL{
			Path: "/" + name + "." + format,
		})

		feed := buildFeed(target, data, uri.String())

		switch format {
		case "atom":
			err = feed.WriteAtom(w)
		case "rss":
			err = feed.WriteRss(w)
		default:
			err = feed.WriteJSON(w)
		}

		if err != nil {
			log.Warnf("%s: %v\n", name, err)
		}
	}
}

func buildFeed(target *Target, entries []Entry, siteURL string) *feeds.Feed {
	out := &feeds.Feed{
		Title:       fmt.Sprintf("%s %s", target.Name, target.Type),
		Link:        &feeds.Link{Href: siteURL},
		Description: fmt.Sprintf("%s for %s", target.Type, target.Name),
		Id:          siteURL,
		Items:       make([]*feeds.Item, 0, len(entries)),
	}

	for _, entry := range entries {
		id := entry.URL
		if id == "" {
			id = fmt.Sprintf("%s:%s:%s", target.Name, target.Type, entry.ID)
		}

		description := entry.Summary
		if description == "" {
			description = fmt.Sprintf("%s: %s", entry.Type, entry.State)
		}

		if entry.Draft {
			description += " (draft)"
		}

		item := &feeds.Item{
			Id:          id,
			IsPermaLink: "false",
			Title:       entry.Title,
			Link:        &feeds.Link{Href: entry.URL},
			Description: description,
		}

		if entry.User != "" {
			item.Author = &feeds.Author{Name: entry.User}
		}

		if entry.CreatedAt != nil {
			item.Created = *entry.CreatedAt

			if out.Created.IsZero() || item.Created.Before(out.Created) {
				out.Created = item.Created
			}

			if item.Created.After(out.Updated) {
				out.Updated = item.Created
			}
		}

		if entry.UpdatedAt != nil {
			item.Updated = *entry.UpdatedAt

			if item.Updated.After(out.Updated) {
				out.Updated = item.Updated
			}
		}

		out.Items = append(out.Items, item)
	}

	if out.Created.IsZero() {
		out.Created = out.Updated
	}

	if out.Created.IsZero() {
		out.Created = time.Now()
	}

	if out.Updated.IsZero() {
		out.Updated = out.Created
	}

	return out
}
