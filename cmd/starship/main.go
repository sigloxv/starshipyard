package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/multiverse-os/cli"
	color "github.com/multiverse-os/cli/text/ansi/color"
	framework "github.com/multiverse-os/starshipyard/framework"
)

// TODO: Review the additional functionality provided by Rails binary, so that
// features like `rails notes` can be added (which scans files for TODO and for
// our purposes NOTE, and likely define what it looks for via YAML configuration
// for a generally useful system; then take these TODOs and others and build a
// notes file to help guide development).
func main() {
	cmd := cli.New(&cli.CLI{
		Name:    "Starshipyard",
		Version: cli.Version{Major: 0, Minor: 1, Patch: 1},
		Usage:   "A command-line tool for controling the starshipyard server, scaffolding boilerplate code, and executing developer defined commands",
		Flags: []cli.Flag{
			cli.Flag{
				Name:  "env, e",
				Value: "development",
				Usage: "Specify the server environment",
			},
			cli.Flag{
				Name:  "address, a",
				Usage: "Specify the address for the HTTP server to listen",
			},
			cli.Flag{
				Name:  "port, p",
				Usage: "Specify the listening port for the HTTP server",
			},
		},
		Commands: []cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Options for controlling starshipyard HTTP server",
				Subcommands: []cli.Command{
					{
						Name:  "start",
						Usage: "Start the starship yard http server",
						Flags: []cli.Flag{
							cli.Flag{
								Name:  "daemonize, d",
								Usage: "Daemonize the http server",
							},
						},
						Action: func(c *cli.Context) error {
							// TODO: Use flags to get port and host address and environment to
							// start the server in or use envirnonemtnal variables. We take
							// these and overwrite the values in the config object in this
							// function
							config, err := framework.LoadConfig("config/app.yaml")
							if err != nil {
								fmt.Println("[starship] missing 'config/app.yaml' starship app configuration")
								// TODO: Should write this default config to config/app.yaml
								config = framework.DefaultConfig()
							}
							// TODO: Should validate address is sane
							address, ok := c.Flags["address"].String()
							if ok {
								config.Address = address
							}
							port, ok := c.Flags["port"].Int()
							if ok {
								// TODO: Should validate for sane value, as in must be between
								// valid range of ports, for 80 and 443 will need to add
								// capabilities to the binary or run as root and drop priviledges
								// liked done by nginx
								config.Port = port
							}

							s := framework.Init(config)

							daemonize, ok := c.Flags["daemonize"]
							if ok {
								fmt.Println("[starship] launching in daemon mode...")
								fmt.Println("[starship] not currently implemented, work on this functionality is in progress")
								//s.StartAsDaemon()
							} else {
								s.Start()
							}
							return nil
						},
					},
				},
			},
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "Generate new go source code for models, controllers, and views",
				Subcommands: []cli.Command{
					{
						Name:  "model",
						Usage: "Build a model template with the specified object data",
						Action: func(c *cli.Context) error {
							fmt.Println("[starship] code generation functionality is not implemented yet")
							fmt.Println("[starship] test code has been built, it just needs to be migrated into the base and will be available shortly")
							return nil
						},
					},
					{
						Name:  "controller",
						Usage: "Build a controller template with the specified object data",
						Action: func(c *cli.Context) error {
							fmt.Println("[starship] code generation functionality is not implemented yet")
							fmt.Println("[starship] test code has been built, it just needs to be migrated into the base and will be available shortly")
							return nil
						},
					},
					{
						Name:  "view",
						Usage: "Build a view template with the specified object data",
						Action: func(c *cli.Context) error {
							fmt.Println("[starship] code generation functionality is not implemented yet")
							fmt.Println("[starship] test code has been built, it just needs to be migrated into the base and will be available shortly")
							return nil
						},
					},
					{
						Name:  "job",
						Usage: "Build a job template with the specified object data",
						Action: func(c *cli.Context) error {
							fmt.Println("[starship] code generation functionality is not implemented yet")
							fmt.Println("[starship] test code has been built, it just needs to be migrated into the base and will be available shortly")
							return nil
						},
					},
					{
						Name:  "model",
						Usage: "Build a model template with the specified object data",
						Action: func(c *cli.Context) error {
							fmt.Println("[starship] code generation functionality is not implemented yet")
							fmt.Println("[starship] test code has been built, it just needs to be migrated into the base and will be available shortly")
							return nil
						},
					},
				},
			},
			{
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "Create a new starship project",
				Action: func(c *cli.Context) error {
					fmt.Println("Building a new starship project directory:")

					fmt.Println("  " + color.Green("CREATE") + " README.md")
					fmt.Println("  " + color.Green("CREATE") + " Taskfile")
					fmt.Println("  " + color.Green("CREATE") + " Dependencies")
					fmt.Println("  " + color.Green("CREATE") + " app")
					fmt.Println("  " + color.Green("CREATE") + " app/assets")
					fmt.Println("  " + color.Green("CREATE") + " app/assets/stylesheets")
					fmt.Println("  " + color.Green("CREATE") + " app/controllers")
					fmt.Println("  " + color.Green("CREATE") + " app/models")
					fmt.Println("  " + color.Green("CREATE") + " app/views")
					fmt.Println("  " + color.Green("CREATE") + " bin")
					fmt.Println("  " + color.Green("CREATE") + " bin/task")
					fmt.Println("  " + color.Green("CREATE") + " bin/starship")
					fmt.Println("  " + color.Green("CREATE") + " bin/dep")
					fmt.Println("  " + color.Green("CREATE") + " config")
					fmt.Println("  " + color.Green("CREATE") + " config/database.yml")
					fmt.Println("  " + color.Green("CREATE") + " config/application.yml")
					fmt.Println("  " + color.Green("CREATE") + " config/environments")
					fmt.Println("  " + color.Green("CREATE") + " config/environments/development.yml")
					fmt.Println("  " + color.Green("CREATE") + " config/environments/production.yml")
					fmt.Println("  " + color.Green("CREATE") + " config/environments/test.yml")
					fmt.Println("  " + color.Green("CREATE") + " config/initializers")
					fmt.Println("  " + color.Green("CREATE") + " config/initializers/mime_types.go")
					fmt.Println("  " + color.Green("CREATE") + " config/initializers/inflections.go")
					fmt.Println("  " + color.Green("CREATE") + " config/initializers/cors.go")
					fmt.Println("  " + color.Green("CREATE") + " config/initializers/cookie_serialization.go")
					fmt.Println("  " + color.Green("CREATE") + " config/initializers/content_security_policy.go")
					fmt.Println("  " + color.Green("CREATE") + " config/initializers/backtrace_silencers.go")
					fmt.Println("  " + color.Green("CREATE") + " config/initializers/assets.go")
					fmt.Println("  " + color.Green("CREATE") + " config/locales")
					fmt.Println("  " + color.Green("CREATE") + " config/boot.go")
					fmt.Println("  " + color.Green("CREATE") + " db")
					fmt.Println("  " + color.Green("CREATE") + " db/seed.go")
					fmt.Println("  " + color.Green("CREATE") + " log")
					fmt.Println("  " + color.Green("CREATE") + " public")
					fmt.Println("  " + color.Green("CREATE") + " public/404.html")
					fmt.Println("  " + color.Green("CREATE") + " public/422.html")
					fmt.Println("  " + color.Green("CREATE") + " public/500.html")
					fmt.Println("  " + color.Green("CREATE") + " public/apple-touch-icon-precomposed.png")
					fmt.Println("  " + color.Green("CREATE") + " public/apple-touch-icon.png")
					fmt.Println("  " + color.Green("CREATE") + " public/favicon.ico")
					fmt.Println("  " + color.Green("CREATE") + " tmp")
					fmt.Println("  " + color.Green("CREATE") + " tmp/cache")
					fmt.Println("  " + color.Green("CREATE") + " tmp/cache/assets")
					fmt.Println("  " + color.Green("CREATE") + " test")
					fmt.Println("  " + color.Green("CREATE") + " test/fixtures")
					fmt.Println("  " + color.Green("CREATE") + " .gitignore")

					return nil
				},
			},
			{
				Name:    "console",
				Aliases: []string{"c"},
				Usage:   "Start the starship yard console interface",
				Action: func(c *cli.Context) error {
					fmt.Println("[starship][CONSOLE] console interface is not implemented yes")
					return nil
				},
			},
		},
	})

	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
