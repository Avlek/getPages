package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/Avlek/getPages/internal"
)

func main() {
	if len(os.Args) < 2 {
		logrus.Fatal("you need urls")
	}

	links := os.Args[1:]
	app := internal.NewApp(links)
	app.Exec()
}
