package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dergoegge/go-functions-sdk/internal/pkg/config"
	"github.com/urfave/cli"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/option"
)

func printTable(funcs []*cloudfunctions.CloudFunction) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t%s\t%s\t", "Name", "Region", "Timeout", "Memory", "Status", "Trigger")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t%s\t%s\t", "----", "------", "-------", "------", "------", "-------")

	for _, fn := range funcs {
		split := strings.Split(fn.Name, "/")
		name := split[len(split)-1]
		region := split[len(split)-3]
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%dMB\t%s\t", name, region, fn.Timeout, fn.AvailableMemoryMb, fn.Status)

		if fn.EventTrigger != nil {
			fmt.Fprintf(w, "%s\t", fn.EventTrigger.EventType)
		}

		if fn.HttpsTrigger != nil {
			fmt.Fprintf(w, "%s\t", fn.HttpsTrigger.Url)
		}
	}

	w.Flush()
}

func list(c *cli.Context) error {
	conf, err := config.NewSDKConfig()
	if err != nil {
		return err
	}

	ctx := context.Background()
	cloudfunctionsService, err := cloudfunctions.NewService(
		ctx,
		option.WithTokenSource(conf),
		option.WithScopes("https://www.googleapis.com/auth/cloudfunctions"),
	)
	if err != nil {
		return err
	}

	projectID, err := conf.ProjectID()
	if err != nil {
		return err
	}

	res, err := cloudfunctionsService.Projects.
		Locations.Functions.List(fmt.Sprintf("projects/%s/locations/-", projectID)).Do()
	if err != nil {
		fmt.Println(res)
		return err
	}

	printTable(res.Functions)

	return nil
}
