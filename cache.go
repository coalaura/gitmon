package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/coalaura/lock"
)

type Cache struct {
	cache *lock.SingleFlightMap[string, []Entry]
}

type Entry struct {
	ID        string
	Type      string
	URL       string
	Title     string
	Summary   string
	User      string
	State     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Draft     bool
}

func (c *Cache) Get(config *Config, target *Target, name string) ([]Entry, error) {
	return c.cache.GetOrFillWithTTL(name, config.Server.interval, func() ([]Entry, error) {
		data, err := target.Fetch(config)
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
				ID:        id,
				Type:      typ,
				URL:       entry.HtmlURL,
				Title:     id + " " + entry.Title,
				Summary:   truncate(entry.Body, 128),
				User:      entry.User.Login,
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
	return &Cache{
		cache: lock.NewSingleFlightMap[string, []Entry](),
	}
}

func truncate(s string, n int) string {
	if n < 0 {
		n = 0
	}

	var (
		end     int
		size    int
		changed bool
	)

	for end < len(s) && size < n {
		start := end
		char, width := utf8.DecodeRuneInString(s[end:])

		if !unicode.IsSpace(char) {
			if width > n-size {
				break
			}

			end += width
			size += width

			continue
		}

		end += width

		for end < len(s) {
			char, width = utf8.DecodeRuneInString(s[end:])
			if !unicode.IsSpace(char) {
				break
			}

			end += width
		}

		size++

		if end-start != 1 || s[start] != ' ' {
			changed = true
		}
	}

	truncated := end < len(s)

	if !changed {
		if truncated {
			return s[:end] + "..."
		}

		return s
	}

	var buf strings.Builder

	if truncated {
		buf.Grow(size + len("..."))
	} else {
		buf.Grow(size)
	}

	var (
		last    int
		inSpace bool
	)

	for i, char := range s[:end] {
		if unicode.IsSpace(char) {
			if !inSpace {
				buf.WriteString(s[last:i])
				buf.WriteByte(' ')

				inSpace = true
			}
		} else if inSpace {
			last = i
			inSpace = false
		}
	}

	if !inSpace {
		buf.WriteString(s[last:end])
	}

	if truncated {
		buf.WriteString("...")
	}

	return buf.String()
}
