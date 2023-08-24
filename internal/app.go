package internal

import (
	"fmt"

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
	for _, link := range app.Links {
		w := NewWebsite(link)

		err := w.Processing()
		if err != nil {
			logrus.Error(err)
			continue
		}

		if !meta {
			err = w.SaveFile()
			if err != nil {
				logrus.Error(err)
			}
		} else {
			fmt.Println(w.Meta)
		}
	}
}
