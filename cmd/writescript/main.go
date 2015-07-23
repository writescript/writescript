package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/writescript/writescript"
	"os"
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
	app.Email = "paul.vollmer@fh-potsdam.de"

	//
	// flags
	//
	app.Flags = []cli.Flag{
		// javascript code
		cli.StringFlag{
			Name:  "plugin, p",
			Value: "",
			Usage: "the generator plugin as file",
		},
		// json data
		cli.StringFlag{
			Name:  "data, d",
			Value: "",
			Usage: "the data formatted as json or yaml, as file or string",
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
	app.Action = func(c *cli.Context) {
		//
		// cli flags
		//
		flagPlugin := c.String("plugin")
		flagData := c.String("data")
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
		data := Data{}
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
		fmt.Println(ws.Content.GetString(flagLinebreak, flagWhitespace))
	}

	app.Run(os.Args)
}
