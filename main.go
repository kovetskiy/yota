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
	usage = `yota-cli 2.2

Usage:
  yota-cli [options] -C -c <code>
  yota-cli [options] -C -s <speed>
  yota-cli [options] -L
  yota-cli [options] -B
  yota-cli [options] -R
  yota-cli -h | --help
  yota-cli -v | --version

Options:
  -C            Change tariff.
    -c <code>   Specify tariff by code like 'POS-1234-567'
    -s <speed>  Specify tariff by speed (float like '1.0' or just 'max')
  -L            List all tariffs.
  -B            Show balance.
  -R            Show remains.
  -f <config>   Specify configuratin file [default: ~/.config/yotarc]
  -h --help     Show this screen.
  -v --version  Show version.
`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, "2.2", true)
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

	case args["-R"].(bool):
		err = showRemains(yotaClient)
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

	format := "%s %-12s %-5s %s\n"

	fmt.Printf(format, " ", "code", "speed", "name")
	for _, tariff := range tariffs {
		label := " "
		if tariff.Active {
			label = "*"
		}

		fmt.Printf(format, label, tariff.Code, tariff.Speed, tariff.Name)
	}

	return err
}

func showBalance(yotaClient *yota.Client) error {
	balance, currency, err := yotaClient.GetBalance()
	if err != nil {
		return err
	}

	fmt.Printf("%.2f %s\n", balance, currency)
	return nil
}

func showRemains(yotaClient *yota.Client) error {
	remains, err := yotaClient.GetRemains()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", remains)
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
