package routes

import (
	"html"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/gorilla/feeds"
)

var (
	atomFeed, jsonFeed, rssFeed []byte
	feed                        feeds.Feed
)

// TZ=UTC git log --date='format-local:%Y-%m-%d %X' --format='%h %cd %s'
func init() {
	feed = feeds.Feed{
		Link:  &feeds.Link{Href: "https://code.golf/"},
		Title: "Code Golf",
	}

	for _, i := range []struct {
		sha, created, id string
		hole             bool
	}{
		{"15bc065", "2017-07-06 21:37:43", "99-bottles-of-beer", true},
	} {
		var name, link string
		if i.hole {
			hole := config.HoleByID[i.id]
			name = hole.Name
			link = "https://code.golf/" + i.id
		} else {
			name = config.LangByID[i.id].Name
			link = "https://code.golf/rankings/holes/all/" + i.id + "/bytes"
		}

		item := feeds.Item{
			Description: "Added the <a href=" + link + ">“" + html.EscapeString(name) + "”</a> ",
			Id:          link,
			Link:        &feeds.Link{Href: link},
			Title:       "Added “" + name + "” ",
		}

		if i.hole {
			item.Title += "Hole"
			item.Description += "hole"
		} else {
			item.Title += "Language"
			item.Description += "language"
		}

		item.Description += " via <a href=https://github.com/code-golf/code-golf/commit/" +
			i.sha + ">" + i.sha + "</a>."

		var err error
		if item.Created, err = time.Parse(time.DateTime, i.created); err != nil {
			panic(err)
		}

		feed.Items = append(feed.Items, &item)

		if feed.Created.IsZero() {
			feed.Created = item.Created
		}
	}

	feed.Title = "Code Golf (Atom Feed)"

	if data, err := feed.ToAtom(); err != nil {
		panic(err)
	} else {
		atomFeed = []byte(data)
	}

	feed.Title = "Code Golf (JSON Feed)"

	if data, err := feed.ToJSON(); err != nil {
		panic(err)
	} else {
		jsonFeed = []byte(data)
	}

	feed.Title = "Code Golf (RSS Feed)"

	if data, err := feed.ToRss(); err != nil {
		panic(err)
	} else {
		rssFeed = []byte(data)
	}
}

// GET /feeds
func feedsGET(w http.ResponseWriter, r *http.Request) {
	render(w, r, "feeds", feed, "Feeds")
}

// GET /feeds/{feed}
func feedGET(w http.ResponseWriter, r *http.Request) {
	switch param(r, "feed") {
	case "atom":
		w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
		w.Write(atomFeed)
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonFeed)
	case "rss":
		w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		w.Write(rssFeed)
	}
}
