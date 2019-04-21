package cli

import (
	"github.com/urfave/cli"
)

// Init intialises the cli's commands and flags
func Init() *cli.App {
	app := cli.NewApp()
	app.Name = "gocf"
	app.Description = "Manage golang cloud functions"
	app.Usage = "Manage golang cloud functions"

	app.Commands = []cli.Command{
		{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "Deploy cloud functions",
			Action:  deployFuncs,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "only",
					Usage: "Specify cloud functions to deploy",
					Value: "all",
				},
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all deployed cloud functions",
			Action:  list,
		},
		// {
		// 	Name:    "login",
		// 	Aliases: []string{"l"},
		// 	Usage:   "Login with your google account",
		// 	Action: func(c *cli.Context) error {
		// 		fmt.Println("completed task: ", c.Args().First())
		// 		return nil
		// 	},
		// },
	}

	return app
}
