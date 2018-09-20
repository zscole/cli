package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"bitbucket.com/wb/project"
)

var RootCmd = &cobra.Command{
	Use:   "wb",
	Short: "A golang development environment for Ethereum",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	prj, err := project.FindProject()
	if prj != nil && err == nil {
		viper.SetConfigFile(filepath.Join(prj.AbsPath(), project.ProjectConfigFilename))
		if err := viper.ReadInConfig(); err != nil {
			Fatal(err)
		}
	}
}
