package project

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/zscole/cli/templates"
)

var Licenses = make(map[string]License)

type License struct {
	Name            string
	PossibleMatches []string
	Text            string
	Header          string
}

func init() {
	Licenses["none"] = License{"None", []string{"none", "false"}, "", ""}

	initializeLicense("agpl", "GNU Affero General Public License", []string{"agpl", "affero gpl", "gnu agpl"})
	initializeLicense("apache", "Apache 2.0", []string{"apache", "apache20", "apache 2.0", "apache2.0", "apache-2.0"})
	initializeLicense("freebsd", "Simplified BSD License", []string{"freebsd", "simpbsd", "simple bsd", "2-clause bsd", "2 clause bsd", "simplified bsd license"})
	initializeLicense("bsd", "NewBSD", []string{"bsd", "newbsd", "3 clause bsd", "3-clause bsd"})
	initializeLicense("gpl2", "GNU General Public License 2.0", []string{"gpl2", "gnu gpl2", "gplv2"})
	initializeLicense("gpl3", "GNU General Public License 3.0", []string{"gpl3", "gplv3", "gpl", "gnu gpl3", "gnu gpl"})
	initializeLicense("lgpl", "GNU Lesser General Public License", []string{"lgpl", "lesser gpl", "gnu lgpl"})
	initializeLicense("mit", "MIT License", []string{"mit"})
}

func initializeLicense(asset, name string, possibleMatches []string) error {
	data := map[string]string{
		"copyright": copyrightLine(),
	}

	buf, err := templates.ExecuteTemplate(fmt.Sprintf("licenses/%s/header.tpl", asset), data)
	if err != nil {
		return err
	}
	header := buf.String()

	buf, err = templates.ExecuteTemplate(fmt.Sprintf("licenses/%s/text.tpl", asset), data)
	if err != nil {
		return err
	}
	text := buf.String()

	Licenses[asset] = License{name, possibleMatches, header, text}
	return nil
}

func getLicense() License {
	if viper.IsSet("license.header") || viper.IsSet("license.text") {
		return License{Header: viper.GetString("license.header"),
			Text: "license.text"}
	}

	if viper.IsSet("license") {
		return findLicense(viper.GetString("license"))
	}

	return Licenses["apache"]
}

func copyrightLine() string {
	author := viper.GetString("author")

	year := viper.GetString("year")
	if year == "" {
		year = time.Now().Format("2018")
	}

	return "Copyright © " + year + " " + author
}

func findLicense(name string) License {
	found := matchLicense(name)
	if found == "" {
		fmt.Println("unknown license: " + name)
		os.Exit(1)
	}
	return Licenses[found]
}

func matchLicense(name string) string {
	if name == "" {
		return ""
	}

	for key, lic := range Licenses {
		for _, match := range lic.PossibleMatches {
			if strings.EqualFold(name, match) {
				return key
			}
		}
	}

	return ""
}
