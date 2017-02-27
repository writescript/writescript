package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/paulvollmer/go-verbose"
	"github.com/writescript/writescript"
	"io/ioutil"
	"os"
)

// .............................................................................
// . write .....................................................................
// ... script ..................................................................
// ..... v0.3.1 ................................................................
// .............................................................................

var (
	debug verbose.Verbose
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
			Name:  "header, H",
			Usage: "enables default header output",
		},
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose mode",
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
		flagHeader := c.Bool("header")
		flagVerbose := c.Bool("verbose")

		// initialize verbose logger
		debug = *verbose.New(os.Stdout, flagVerbose)
		debug.Println("==> debug active")

		// read plugin
		debug.Println("==> load plugins")
		pluginBytes, err := writescript.LoadPlugin(flagPlugin, debug)
		if err != nil {
			fmt.Println("writescript load plugin error!\n", err)
			os.Exit(1)
		}

		// read data
		debug.Println("==> load data")
		data := writescript.Data{}
		data.Init(flagData, debug)
		// fmt.Println(data)

		// run the generator
		debug.Println("==> execute script")
		ws := writescript.WriteScript{}
		err = ws.Process(string(pluginBytes), string(data.JSON), flagHeader, debug)
		if err != nil {
			fmt.Println("writescript plugin error!\n", err)
			os.Exit(1)
		}

		// write output
		if flagLinebreak == "\\n" {
			flagLinebreak = "\n"
		}
		if flagWhitespace == "\\t" {
			flagWhitespace = "\t"
		}
		if flagOutput == "" {
			fmt.Println(string(ws.Content.Get(flagLinebreak, flagWhitespace)))
		} else {
			debug.Println("==> write file", flagOutput)
			fileBytes := ws.Content.Get(flagLinebreak, flagWhitespace)
			err := ioutil.WriteFile(flagOutput, fileBytes, 0644)
			if err != nil {
				fmt.Println("writescript output error!\n", err)
			}
		}
		return nil
	}

	app.Run(os.Args)
}
