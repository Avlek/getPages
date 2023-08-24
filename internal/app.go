package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
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

func (app *App) Exec() {
	for _, link := range app.Links {
		if !strings.HasPrefix(link, "http") {
			link = "https://" + link
		}

		url2, err := url.Parse(link)
		if err != nil {
			logrus.Error(err)
			continue
		}

		b, err := Fetch(link)
		if err != nil {
			logrus.Error(err)
			continue
		}

		fileName := url2.Host + url2.Path + ".html"
		fileName = strings.TrimRight(fileName, "/")
		fileName = strings.Replace(fileName, "/", "_", -1)
		err = os.WriteFile(fileName, b, 0644)
		if err != nil {
			fmt.Println("Error saving to disk:", err)
			continue
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

	bodyString := string(bodyBytes)
	urlRegex, _ := regexp.Compile("<a[^>]*href=[\"']([^\"']*)[\"'][^>]*>")
	urls := urlRegex.FindAllString(bodyString, -1)
	imgRegex, _ := regexp.Compile("<img[^>]*src=[\"']([^\"']*)[\"'][^>]*>")
	imgs := imgRegex.FindAllString(bodyString, -1)
	logrus.Info(len(urls), len(imgs))

	return bodyBytes, nil
}
