package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/writescript/writescript"
	"os"
	"io/ioutil"
)

// main cli tool
func main() {

	//
	// meta
	//
	app := cli.NewApp()
	app.Name = "writescript"
	app.Version = writescript.Version
	app.Usage = "plugin based generator tool backpacked with data"
	app.Author = "Paul Vollmer"
	app.Email = "github.com/writescript"

	//
	// flags
	//
	app.Flags = []cli.Flag{
		// input
		cli.StringFlag{
			Name:  "plugin, p",
			Value: "",
			Usage: "the generator plugin as file",
		},
		cli.StringFlag{
			Name:  "data, d",
			Value: "",
			Usage: "the data formatted as json or yaml, as file or string",
		},
		// output
		cli.StringFlag{
			Name:  "output, o",
			Value: "",
			Usage: "the output filepath",
		},
		// settings
		cli.StringFlag{
			Name:  "linebreak, l",
			Value: "\\n",
			Usage: "the linebreak for each row",
		},
		cli.StringFlag{
			Name:  "whitespace, w",
			Value: "\\t",
			Usage: "the level whitespace",
		},
		cli.BoolFlag{
			Name:  "header-off, H",
			Usage: "disables header output",
		},
	}

	//
	// process
	//
	app.Action = func(c *cli.Context) error {
		//
		// cli flags
		//
		flagPlugin := c.String("plugin")
		flagData := c.String("data")
		flagOutput := c.String("output")
		flagLinebreak := c.String("linebreak")
		flagWhitespace := c.String("whitespace")
		flagHeaderOff := c.Bool("header-off")

		// read plugin
		pluginBytes, err := writescript.LoadPlugin(flagPlugin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// read data
		data := writescript.Data{}
		data.Init(flagData)

		// run the generator
		ws := writescript.WriteScript{}
		err = ws.Process(string(pluginBytes), data.JSON, !flagHeaderOff)
		if err != nil {
			fmt.Println("writescript plugin error!\n", err)
			os.Exit(1)
		}

		if flagLinebreak == "\\n" {
			flagLinebreak = "\n"
		}
		if flagWhitespace == "\\t" {
			flagWhitespace = "\t"
		}

		// write output
		if flagOutput == "" {
			fmt.Println(string(ws.Content.Get(flagLinebreak, flagWhitespace)))
		} else {
			fileBytes := ws.Content.Get(flagLinebreak, flagWhitespace)
			ioutil.WriteFile(flagOutput, fileBytes, 0644)
		}
		return nil
	}

	app.Run(os.Args)
}
