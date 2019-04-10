package cli

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/dergoegge/go-functions-sdk/internal/pkg/build"
	"github.com/dergoegge/go-functions-sdk/internal/pkg/deploy"
	"github.com/dergoegge/go-functions-sdk/internal/pkg/parse"
	"github.com/google/logger"
	"github.com/urfave/cli"
)

func deployFuncs(ctx *cli.Context) error {
	funcNames := strings.Split(ctx.String("only"), ",")

	pkgs := parse.GetPackages()
	availableFunctions := pkgs.Functions()
	stagedFunctions := availableFunctions

	if strings.Compare(funcNames[0], "all") != 0 {
		stagedFunctions = []types.Object{}
		for _, fn := range availableFunctions {
			keep := false
			for _, funcName := range funcNames {
				if strings.Compare(fn.Name(), funcName) == 0 {
					keep = true
					break
				}
			}

			if keep {
				stagedFunctions = append(stagedFunctions, fn)
			}
		}
	}

	if len(stagedFunctions) == 0 {
		return fmt.Errorf("No functions staged for deployment")
	}

	funcNames = []string{}
	for _, fnObj := range stagedFunctions {
		funcNames = append(funcNames, fnObj.Name())
	}
	logger.Info("Preparing deployment of ", funcNames)

	err := build.Plugins(pkgs)
	if err != nil {
		return err
	}

	deployCmds, err := deploy.Prepare(stagedFunctions)
	if err != nil {
		return err
	}

	return deploy.Functions(stagedFunctions, deployCmds)
}

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
		// {
		// 	Name:    "list",
		// 	Aliases: []string{"ls"},
		// 	Usage:   "List all deployed cloud functions",
		// 	Action: func(c *cli.Context) error {
		// 		fmt.Println("added task: ", c.String("only"))
		// 		return nil
		// 	},
		// },
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
