package main

import (
	"github.com/codegangsta/cli"
	"github.com/thoas/picfit/application"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "picfit"
	app.Author = "thoas"
	app.Email = "florent.messa@gmail.com"
	app.Usage = "Display, manipulate, transform and cache your images"
	app.Version = application.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Config file path",
			EnvVar: "PICFIT_CONFIG_PATH",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "signature",
			ShortName: "s",
			Usage:     "Verify that your client application is generating correct signatures",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "key",
					Usage: "The signing key",
				},
			},
			Action: func(c *cli.Context) {
				key := c.String("key")
				if key == "" {
					application.App.Logger.Info.Printf("You must provide a key")
					os.Exit(1)
				}

				if len(c.Args()) < 1 {
					application.App.Logger.Info.Printf("You must provide a Query String")
					os.Exit(1)
				}

				qs := c.Args()[0]

				signature := application.Sign(key, qs)

				appended := application.AppendSign(key, qs)

				application.App.Logger.Info.Printf("Query String: %s", qs)
				application.App.Logger.Info.Printf("Signature: %s", signature)
				application.App.Logger.Info.Printf("Signed Query String: %s", appended)
			},
		},
	}
	app.Action = func(c *cli.Context) {
		config := c.String("config")

		if config != "" {
			if _, err := os.Stat(config); err != nil {
				application.App.Logger.Error.Printf("Can't find config file `%s`\n", config)
				os.Exit(1)
			}
		} else {
			application.App.Logger.Error.Print("Can't find config file\n")
			os.Exit(1)
		}

		err := application.Run(config)

		if err != nil {
			application.App.Logger.Error.Print(err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}