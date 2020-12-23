package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/multiverse-os/cli"
	framework "github.com/multiverse-os/starshipyard/framework"
)

// TODO: Review the additional functionality provided by Rails binary, so that
// features like `rails notes` can be added (which scans files for TODO and for
// our purposes NOTE, and likely define what it looks for via YAML configuration
// for a generally useful system; then take these TODOs and others and build a
// notes file to help guide development).
func main() {
	cmd := cli.New(&cli.CLI{
		Name:        "Starshipyard",
		Version:     cli.Version{Major: 0, Minor: 1, Patch: 1},
		Description: "A command-line tool for controling the starshipyard server, scaffolding boilerplate code, and executing developer defined commands",
		GlobalFlags: []cli.Flag{
			cli.Flag{
				Name:        "env",
				Alias:       "e",
				Default:     "development",
				Description: "Specify the server environment",
			},
			cli.Flag{
				Name:        "address",
				Alias:       "a",
				Default:     "0.0.0.0",
				Description: "Specify the address for the HTTP server to listen",
			},
			cli.Flag{
				Name:        "port",
				Alias:       "p",
				Default:     "3000",
				Description: "Specify the listening port for the HTTP server",
			},
		},
		Commands: []cli.Command{
			{
				Name:        "server",
				Alias:       "s",
				Description: "Options for controlling starshipyard HTTP server",
				Subcommands: []cli.Command{
					{
						Name:        "start",
						Description: "Start the starship yard http server",
						Flags: []cli.Flag{
							cli.Flag{
								Name:        "daemonize",
								Alias:       "d",
								Description: "Daemonize the http server",
							},
						},
						Action: func(c *cli.Context) error {
							fmt.Println("Starting the starship yard http server...")
							config, err := framework.LoadConfig("config/app.yaml")
							if err != nil {
								fmt.Println("[starship] missing 'config/app.yaml' starship app configuration")
								config = framework.DefaultConfig()
							}
							config.Address = c.Flag("address").String()
							config.Port = c.Flag("port").Int()

							s := framework.Init(config)

							daemonize := c.Flag("daemonize").Bool()
							if daemonize {
								fmt.Println("[starship] launching in daemon mode...")
								fmt.Println("[starship] not currently implemented, work on this functionality is in progress")
								//s.StartAsDaemon()
							} else {
								fmt.Println("[starship] launching with terminal attached to server")
								s.Start()
							}
							return nil
						},
					},
				},
			},
			{
				Name:        "generate",
				Alias:       "g",
				Description: "Generate new go source code for models, controllers, and views",
				Subcommands: []cli.Command{
					{
						Name:        "model",
						Description: "Build a model template with the specified object data",
						Action: func(c *cli.Context) error {
							fmt.Println("[starship] code generation functionality is not implemented yet")
							fmt.Println("[starship] test code has been built, it just needs to be migrated into the base and will be available shortly")
							return nil
						},
					},
					{
						Name:        "controller",
						Description: "Build a controller template with the specified object data",
						Action: func(c *cli.Context) error {
							fmt.Println("[starship] code generation functionality is not implemented yet")
							fmt.Println("[starship] test code has been built, it just needs to be migrated into the base and will be available shortly")
							return nil
						},
					},
					{
						Name:        "view",
						Description: "Build a view template with the specified object data",
						Action: func(c *cli.Context) error {
							fmt.Println("[starship] code generation functionality is not implemented yet")
							fmt.Println("[starship] test code has been built, it just needs to be migrated into the base and will be available shortly")
							return nil
						},
					},
					{
						Name:        "job",
						Description: "Build a job template with the specified object data",
						Action: func(c *cli.Context) error {
							fmt.Println("[starship] code generation functionality is not implemented yet")
							fmt.Println("[starship] test code has been built, it just needs to be migrated into the base and will be available shortly")
							return nil
						},
					},
					{
						Name:        "model",
						Description: "Build a model template with the specified object data",
						Action: func(c *cli.Context) error {
							fmt.Println("[starship] code generation functionality is not implemented yet")
							fmt.Println("[starship] test code has been built, it just needs to be migrated into the base and will be available shortly")
							return nil
						},
					},
				},
			},
			{
				Name:        "new",
				Alias:       "n",
				Description: "Create a new starship project",
				Action: func(c *cli.Context) error {
					fmt.Println("Building a new starship project directory:")

					fmt.Println("[CREATE] README.md")
					fmt.Println("[CREATE] Taskfile")
					fmt.Println("[CREATE] Dependencies")
					fmt.Println("[CREATE] app")
					fmt.Println("[CREATE] app/assets")
					fmt.Println("[CREATE] app/assets/stylesheets")
					fmt.Println("[CREATE] app/controllers")
					fmt.Println("[CREATE] app/models")
					fmt.Println("[CREATE] app/views")
					fmt.Println("[CREATE] bin")
					fmt.Println("[CREATE] bin/task")
					fmt.Println("[CREATE] bin/starship")
					fmt.Println("[CREATE] bin/dep")
					fmt.Println("[CREATE] config")
					fmt.Println("[CREATE] config/database.yml")
					fmt.Println("[CREATE] config/application.yml")
					fmt.Println("[CREATE] config/environments")
					fmt.Println("[CREATE] config/environments/development.yml")
					fmt.Println("[CREATE] config/environments/production.yml")
					fmt.Println("[CREATE] config/environments/test.yml")
					fmt.Println("[CREATE] config/initializers")
					fmt.Println("[CREATE] config/initializers/mime_types.go")
					fmt.Println("[CREATE] config/initializers/inflections.go")
					fmt.Println("[CREATE] config/initializers/cors.go")
					fmt.Println("[CREATE] config/initializers/cookie_serialization.go")
					fmt.Println("[CREATE] config/initializers/content_security_policy.go")
					fmt.Println("[CREATE] config/initializers/backtrace_silencers.go")
					fmt.Println("[CREATE] config/initializers/assets.go")
					fmt.Println("[CREATE] config/locales")
					fmt.Println("[CREATE] config/boot.go")
					fmt.Println("[CREATE] db")
					fmt.Println("[CREATE] db/seed.go")
					fmt.Println("[CREATE] log")
					fmt.Println("[CREATE] public")
					fmt.Println("[CREATE] public/404.html")
					fmt.Println("[CREATE] public/422.html")
					fmt.Println("[CREATE] public/500.html")
					fmt.Println("[CREATE] public/apple-touch-icon-precomposed.png")
					fmt.Println("[CREATE] public/apple-touch-icon.png")
					fmt.Println("[CREATE] public/favicon.ico")
					fmt.Println("[CREATE] tmp")
					fmt.Println("[CREATE] tmp/cache")
					fmt.Println("[CREATE] tmp/cache/assets")
					fmt.Println("[CREATE] test")
					fmt.Println("[CREATE] test/fixtures")
					fmt.Println("[CREATE] .gitignore")

					return nil
				},
			},
			{
				Name:        "console",
				Alias:       "c",
				Description: "Start the starship yard console interface",
				Action: func(c *cli.Context) error {
					fmt.Println("[starship][CONSOLE] console interface is not implemented yes")
					return nil
				},
			},
		},
	})

	_, err := cmd.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
