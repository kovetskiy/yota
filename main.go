package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/go-yota"
	"github.com/zazab/zhash"
)

const (
	usage = `yota-cli 2.0

Usage:
  yota-cli [options] -C -c <code>
  yota-cli [options] -C -s <speed>
  yota-cli [options] -L
  yota-cli [options] -B
  yota-cli -h | --help
  yota-cli -v | --version

Options:
  -h --help      Show this help.
  -v --version   Show version
  -C             Change tariff.
	-c  <code>   Specify tariff by code like 'POS-1234-567'
	-s <speed>   Specify tariff by speed (float like '1.0' or just 'max')
  -L             List all tariffs.
  -B             Show balance.
  -f <config>    Specify configuratin file [default: ~/.config/yotarc]
`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, "2.0", true)
	if err != nil {
		panic(err)
	}

	config, err := getConfig(args["-f"].(string))
	if err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	username, err := config.GetString("username")
	if err != nil {
		log.Fatal(err)
	}

	password, err := config.GetString("password")
	if err != nil {
		log.Fatal(err)
	}

	yotaClient := yota.NewClient(username, password, nil)

	err = yotaClient.Login()
	if err != nil {
		log.Fatalf("can't login: %s", err)
	}

	switch {
	case args["-L"].(bool):
		err = listTariffs(yotaClient)

	case args["-C"].(bool):
		var (
			code, _  = args["-c"].(string)
			speed, _ = args["-s"].(string)
		)

		err = searchAndChangeTariff(yotaClient, code, speed)

	case args["-B"].(bool):
		err = showBalance(yotaClient)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func listTariffs(yotaClient *yota.Client) error {
	tariffs, err := yotaClient.GetTariffs()
	if err != nil {
		return err
	}

	format := "%-12s %-5s %s\n"

	fmt.Printf(format, "code", "speed", "name")
	for _, tariff := range tariffs {
		name := ""
		if tariff.Active {
			name = "[active] "
		}
		name = name + tariff.Name

		fmt.Printf(format, tariff.Code, tariff.Speed, name)
	}

	return err
}

func showBalance(yotaClient *yota.Client) error {
	balance, err := yotaClient.GetBalance()
	if err != nil {
		return err
	}

	fmt.Printf("%d\n", balance)
	return nil
}

func searchAndChangeTariff(
	yotaClient *yota.Client, code, speed string,
) error {
	tariffs, err := yotaClient.GetTariffs()
	if err != nil {
		return err
	}

	found := false
	newTariff := yota.Tariff{}
	for _, tariff := range tariffs {
		if tariff.Code == code || tariff.Speed == speed {
			found = true
			newTariff = tariff
			break
		}
	}

	if found {
		err = yotaClient.ChangeTariff(newTariff)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("tariff successfully changed")
		return nil
	}

	return fmt.Errorf("can't find specified tariff")
}

func getConfig(path string) (zhash.Hash, error) {
	var configData map[string]interface{}

	if strings.HasPrefix(path, "~/") {
		path = strings.Replace(path, "~/", os.Getenv("HOME")+"/", 1)
	}

	_, err := toml.DecodeFile(path, &configData)
	if err != nil {
		return zhash.Hash{}, err
	}

	return zhash.HashFromMap(configData), nil
}
