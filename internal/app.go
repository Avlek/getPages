package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

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

func (app *App) Exec() {
	for _, link := range app.Links {
		b, err := Fetch(link)
		if err != nil {
			logrus.Error(err)
			continue
		}

		err = os.WriteFile(link+".html", b, 0644)
		if err != nil {
			fmt.Println("Error saving to disk:", err)
			return
		}
	}
}

func Fetch(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("wrong status: " + response.Status)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
