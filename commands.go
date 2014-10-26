package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mackerelio/gomkr/utils"
	mkr "github.com/mackerelio/mackerel-client-go"
)

var Commands = []cli.Command{
	commandStatus,
	commandHosts,
	commandCreate,
	commandUpdate,
	commandThrow,
	commandFetch,
	commandRetire,
}

var commandStatus = cli.Command{
	Name:  "status",
	Usage: "Show host status",
	Description: `
`,
	Action: doStatus,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
	},
}

var commandHosts = cli.Command{
	Name:  "hosts",
	Usage: "Show hosts",
	Description: `
`,
	Action: doHosts,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "", Usage: "Show hosts only matched with <name>"},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "Show hosts only belongs to <service>"},
		cli.StringSliceFlag{Name: "role, r", Value: &cli.StringSlice{}, Usage: "Show hosts only belongs to <role>. Multiple choice allow. Required --service"},
		cli.StringSliceFlag{Name: "status, st", Value: &cli.StringSlice{}, Usage: "Show hosts only matched <status>. Multiple choice allow."},
		cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
	},
}

var commandCreate = cli.Command{
	Name:  "create",
	Usage: "",
	Description: `
`,
	Action: doCreate,
}

var commandUpdate = cli.Command{
	Name:  "update",
	Usage: "",
	Description: `
`,
	Action: doUpdate,
}

var commandThrow = cli.Command{
	Name:  "throw",
	Usage: "",
	Description: `
`,
	Action: doThrow,
}

var commandFetch = cli.Command{
	Name:  "fetch",
	Usage: "",
	Description: `
`,
	Action: doFetch,
}

var commandRetire = cli.Command{
	Name:  "retire",
	Usage: "",
	Description: `
`,
	Action: doRetire,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func newMackerel() *mkr.Client {
	apiKey := os.Getenv("MACKEREL_APIKEY")
	if apiKey == "" {
		utils.Log("error", `
Not set MACKEREL_APIKEY environment variable. (Try "export MACKEREL_APIKEY='<Your apikey>'")
`)
		os.Exit(1)
	}

	if os.Getenv("DEBUG") != "" {
		mackerel, err := mkr.NewClientForTest(apiKey, "https://mackerel.io/api/v0", true)
		utils.DieIf(err)

		return mackerel
	} else {
		return mkr.NewClient(apiKey)
	}
}

func doStatus(c *cli.Context) {
	argHostId := c.Args().Get(0)
	isVerbose := c.Bool("verbose")

	if argHostId == "" {
		cli.ShowCommandHelp(c, "status")
		os.Exit(1)
	}

	mackerel := newMackerel()
	host, err := mackerel.FindHost(argHostId)
	utils.DieIf(err)

	if isVerbose {
		data, err := json.MarshalIndent(host, "", "    ")
		utils.DieIf(err)

		fmt.Fprintln(os.Stdout, string(data))
	} else {
		format := &HostFormat{
			Id:            host.Id,
			Name:          host.Name,
			Status:        host.Status,
			RoleFullnames: host.GetRoleFullnames(),
			IsRetired:     host.IsRetired,
			CreatedAt:     host.DateStringFromCreatedAt(),
		}

		data, err := json.MarshalIndent(format, "", "    ")
		utils.DieIf(err)

		fmt.Fprintln(os.Stdout, string(data))
	}
}

func doHosts(c *cli.Context) {
	isVerbose := c.Bool("verbose")
	argName := c.String("name")
	argService := c.String("service")
	argRoles := c.StringSlice("role")
	argStatuses := c.StringSlice("status")

	mackerel := newMackerel()
	hosts, err := mackerel.FindHosts(&mkr.FindHostsParam{
		Name:     argName,
		Service:  argService,
		Roles:    argRoles,
		Statuses: argStatuses,
	})
	utils.DieIf(err)

	if isVerbose {
		data, err := json.MarshalIndent(hosts, "", "    ")
		utils.DieIf(err)

		fmt.Fprintln(os.Stdout, string(data))
	} else {
		var hosts_format []*HostFormat
		for _, host := range hosts {
			format := &HostFormat{
				Id:            host.Id,
				Name:          host.Name,
				Status:        host.Status,
				RoleFullnames: host.GetRoleFullnames(),
				IsRetired:     host.IsRetired,
				CreatedAt:     host.DateStringFromCreatedAt(),
			}
			hosts_format = append(hosts_format, format)
		}

		data, err := json.MarshalIndent(hosts_format, "", "    ")
		utils.DieIf(err)

		fmt.Fprintln(os.Stdout, string(data))
	}
}

func doCreate(c *cli.Context) {
}

func doUpdate(c *cli.Context) {
}

func doThrow(c *cli.Context) {
}

func doFetch(c *cli.Context) {
}

func doRetire(c *cli.Context) {
}
