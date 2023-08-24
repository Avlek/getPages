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
	"time"
)

type Metadata struct {
	Site      string
	NumLinks  int
	Images    int
	LastFetch time.Time
}

func (meta Metadata) String() string {
	return fmt.Sprintf("site: %s\nnum_links: %d\nimages: %d\nlast_fetch: %s",
		meta.Site, meta.NumLinks, meta.Images, meta.LastFetch.Format("Mon Jan 02 2006 15:04 MST"))
}

type Website struct {
	URLString string
	Content   []byte

	URL  *url.URL
	Meta *Metadata
}

func NewWebsite(u string) *Website {
	return &Website{
		URLString: u,
	}
}

func (w *Website) Processing() error {
	err := w.CheckURL()
	if err != nil {
		return err
	}

	err = w.Fetch()
	if err != nil {
		return err
	}

	err = w.SaveMetadata()
	if err != nil {
		return err
	}

	//err = w.SaveFile()
	//if err != nil {
	//	return err
	//}

	return nil
}

func (w *Website) CheckURL() (err error) {
	if !strings.HasPrefix(w.URLString, "http") {
		w.URLString = "https://" + w.URLString
	}

	w.URL, err = url.Parse(w.URLString)
	return
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

	w.Meta = &Metadata{
		Site:      w.URL.Host + w.URL.Path,
		NumLinks:  len(urls),
		Images:    len(imgs),
		LastFetch: time.Now().UTC(),
	}

	return nil
}

func (w *Website) SaveFile() error {
	fileName := w.URL.Host + w.URL.Path + ".html"
	fileName = strings.TrimRight(fileName, "/")
	fileName = strings.Replace(fileName, "/", "_", -1)

	return os.WriteFile(fileName, w.Content, 0644)
}
