package main

import (
	"flag"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/Avlek/getPages/internal"
)

func main() {
	metaFlag := flag.Bool("metadata", false, "")
	flag.Parse()

	var links []string
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			links = append(links, arg)
		}
	}
	if len(links) == 0 {
		logrus.Fatal("you need urls")
	}
	app := internal.NewApp(links)
	app.Exec(*metaFlag)
}
