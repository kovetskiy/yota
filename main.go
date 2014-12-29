package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/go-yota"
)

func main() {
	rawArgs := mergeArgsWithConfig(os.Getenv("HOME") + "/.config/yotarc")
	args, err := parseArgs(rawArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if args["--user"] == nil || args["--pass"] == nil {
		fmt.Println("--user and --pass should be specified.")
		os.Exit(1)
	}

	username := args["--user"].(string)
	password := args["--pass"].(string)

	yotaCli := yota.NewClient(username, password, nil)

	err = yotaCli.Login()
	if err != nil {
		fmt.Printf("Could not login: %s", err.Error())
		os.Exit(1)
	}

	switch {
	case args["list"]:
		listMode(yotaCli)

	case args["switch"]:
		query := map[string]string{
			"Code":  "",
			"Speed": "",
		}

		if args["--code"] != nil {
			query["Code"] = args["--code"].(string)
		} else if args["--speed"] != nil {
			query["Speed"] = args["--speed"].(string)
		}

		switchMode(yotaCli, query)

	case args["balance"]:
		balanceMode(yotaCli)
	}
}

func listMode(yotaCli *yota.Client) {
	tariffs, err := yotaCli.GetTariffs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	format := "%-12s %-5s %s\n"
	fmt.Printf(format, "Code", "Speed", "Name")
	for _, tariff := range tariffs {
		name := ""
		if tariff.Active {
			name = "[active] "
		}
		name = name + tariff.Name

		fmt.Printf(format, tariff.Code, tariff.Speed, name)
	}
}

func balanceMode(yotaCli *yota.Client) {
	balance, err := yotaCli.GetBalance()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%d", balance)
}

func switchMode(yotaCli *yota.Client, query map[string]string) {
	tariffs, err := yotaCli.GetTariffs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, tariff := range tariffs {
		found := tariff.Code == query["Code"] || tariff.Speed == query["--speed"]

		if found {
			err := yotaCli.ChangeTariff(tariff)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Tariff successfully switched")
			os.Exit(0)
		}
	}

	fmt.Printf("Something went wrong")
	os.Exit(1)
}

func parseArgs(cmd []string) (map[string]interface{}, error) {
	help := `Yota Cli.

Usage:
  yota-cli [options] list --user USERNAME --pass PASSWORD
  yota-cli [options] switch (--code CODE|--speed SPEED) --user USERNAME --pass PASSWORD
  yota-cli [options] balance --user USERNAME --pass PASSWORD
  yota-cli -h | --help
  yota-cli -v | --version

Options:
  -h --help           Show this help.
  -v --version        Show version
  --code=<code>       Set tariff by code like 'POS-1234-567'
  --speed=<speed>     Set tariff by speed (float like '1.0' or just 'max')
  --user=<username>   Your username in system
  --pass=<password>   Your password in system
`

	args, err := docopt.Parse(help, cmd, true, "1.0", false)

	return args, err
}

func mergeArgsWithConfig(path string) []string {
	args := make([]string, 0)

	conf, err := ioutil.ReadFile(path)
	if err == nil {
		confLines := strings.Split(string(conf), "\n")
		for _, line := range confLines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			args = append(args, line)
		}
	}

	args = append(args, os.Args[1:]...)

	return args
}
