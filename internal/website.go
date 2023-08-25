package internal

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type Website struct {
	URLString string
	Content   []byte

	URL  *url.URL
	Data map[string]Metadata
}

func NewWebsite(u *url.URL, data map[string]Metadata) *Website {
	return &Website{
		URL:       u,
		URLString: u.Scheme + "://" + u.Host + u.Path,
		Data:      data,
	}
}

func (w *Website) Processing() error {
	err := w.Fetch()
	if err != nil {
		return err
	}

	err = w.SaveMetadata()
	if err != nil {
		return err
	}

	err = w.SaveFile()
	if err != nil {
		return err
	}

	return nil
}

func (w *Website) Fetch() error {
	response, err := http.Get(w.URLString)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("wrong status: " + response.Status)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	w.Content = bodyBytes
	return nil
}

func (w *Website) SaveMetadata() error {
	bodyString := string(w.Content)
	urlRegex, err := regexp.Compile("<a[^>]*href=[\"']([^\"']*)[\"'][^>]*>")
	if err != nil {
		return err
	}
	urls := urlRegex.FindAllString(bodyString, -1)

	imgRegex, err := regexp.Compile("<img[^>]*src=[\"']([^\"']*)[\"'][^>]*>")
	if err != nil {
		return err
	}
	imgs := imgRegex.FindAllString(bodyString, -1)

	w.Data[w.URL.Host+w.URL.Path] = Metadata{
		Site:      w.URL.Host + w.URL.Path,
		NumLinks:  len(urls),
		Images:    len(imgs),
		LastFetch: time.Now().UTC(),
	}

	return SaveMetadata(w.Data)
}

func (w *Website) SaveFile() error {
	fileName := w.URL.Host + w.URL.Path + ".html"
	fileName = strings.Replace(fileName, "/", "_", -1)

	return os.WriteFile(fileName, w.Content, 0644)
}
