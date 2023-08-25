package internal

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type App struct {
	Links []string
}

func NewApp(links []string) *App {
	return &App{
		Links: links,
	}
}

func (app *App) Exec(meta bool) {
	data, err := GetMetadata()
	if err != nil {
		logrus.Error(err)
		return
	}
	if data == nil {
		data = make(map[string]Metadata)
	}

	for _, link := range app.Links {
		u, err := checkURL(link)
		if err != nil {
			logrus.Errorf("%s: %s", link, err)
			continue
		}

		if !meta {
			w := NewWebsite(u, data)
			err = w.Processing()
			if err != nil {
				logrus.Error(err)
				continue
			}
		} else {
			if m, ok := data[u.Host+u.Path]; ok {
				fmt.Println(m)
			}
		}
	}
}

func checkURL(u string) (*url.URL, error) {
	if !strings.HasPrefix(u, "http") {
		u = "https://" + u
	}

	l, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return l, nil
}
