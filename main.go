package main

import (
	"fmt"
	"log"
	"os"

	"github.com/binarydud/archiver/archive"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "archiver"
	app.Usage = "make an explosive entrance"
	app.Action = func(c *cli.Context) error {
		zipfile := c.Args().Get(0)
		files := c.Args()[1:]
		archiver := archive.GetZipArchiver(zipfile)
		archiver.Start()
		defer archiver.CloseWriter()
		//fmt.Printf("package %s, files %q", zipfile, files)
		for _, filename := range files {
			if err := archiver.AddItem(filename); err != nil {
				log.Fatal(err)
			}

		}
		hash, err := archiver.GetHash()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("{\"hash\":\"%s\",\"filename\":\"%s\"}", hash, zipfile)
		// archiver.Close()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
