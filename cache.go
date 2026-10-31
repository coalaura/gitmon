package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/coalaura/lock"
)

type Cache struct {
	cache   *lock.SingleFlightMap[string, []Entry]
	targets map[string]*Target
}

type Entry struct {
	User User `json:"user"`

	ID        string     `json:"id"`
	Type      string     `json:"type"`
	URL       string     `json:"url"`
	Title     string     `json:"title"`
	State     string     `json:"state"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Draft     bool       `json:"draft"`
}

type User struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	URL    string `json:"url"`
}

var errNotFound = errors.New("target not found")

func (c *Cache) Get(name string) ([]Entry, error) {
	target, ok := c.targets[name]
	if !ok {
		return nil, errNotFound
	}

	return c.cache.GetOrFillWithTTL(name, config.Server.interval, func() ([]Entry, error) {
		data, err := target.Fetch()
		if err != nil {
			return nil, err
		}

		entries := make([]Entry, 0, len(data))

		for _, entry := range data {
			var (
				id  string
				typ string
			)

			switch target.Type {
			case "issues":
				id = fmt.Sprintf("#%d", entry.Number)
				typ = "issue"
			case "pull-requests":
				id = fmt.Sprintf("PR#%d", entry.Number)
				typ = "pr"
			}

			entries = append(entries, Entry{
				User: User{
					Name:   entry.User.Login,
					Avatar: entry.User.AvatarURL,
					URL:    entry.User.HtmlURL,
				},

				ID:        id,
				Type:      typ,
				URL:       entry.HtmlURL,
				Title:     id + " " + entry.Title,
				State:     entry.State,
				CreatedAt: entry.CreatedAt,
				UpdatedAt: entry.UpdatedAt,
				Draft:     entry.Draft,
			})
		}

		return entries, nil
	})
}

func NewCache() *Cache {
	targets := make(map[string]*Target, len(config.Targets))

	for i, target := range config.Targets {
		targets[target.Name] = &config.Targets[i]
	}

	return &Cache{
		cache:   lock.NewSingleFlightMap[string, []Entry](),
		targets: targets,
	}
}
