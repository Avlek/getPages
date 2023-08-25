package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

func GetMetadata() (map[string]Metadata, error) {
	filename := "./metadata.json"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		f, err2 := os.Create(filename)
		if err2 != nil {
			return nil, err2
		}
		f.WriteString("{}")
		defer f.Close()
		return nil, nil
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data map[string]Metadata
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	return data, nil
}

func SaveMetadata(data map[string]Metadata) error {
	filename := "./metadata.json"
	jsonString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonString, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
