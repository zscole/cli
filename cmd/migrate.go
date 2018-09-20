package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run migrations to deploy contracts",
	Run: func(cmd *cobra.Command, args []string) {
		stub_args := []string{"--network", viper.GetString("default_network")}
		if viper.GetBool("reset") {
			stub_args = append(stub_args, "--reset")
		}

		if err := runStub("migrate", stub_args...); err != nil {
			Fatal(err)
		}
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "(alias for migrate)",
	Run:   migrateCmd.Run,
}

func init() {
	RootCmd.AddCommand(migrateCmd)
	RootCmd.AddCommand(deployCmd)

	migrateCmd.PersistentFlags().StringP("network", "n", "NAME", "network to run migrations on")
	migrateCmd.PersistentFlags().Bool("reset", false, "redeploy all migrations")
	deployCmd.PersistentFlags().StringP("network", "n", "NAME", "network to run migrations on")
	deployCmd.PersistentFlags().Bool("reset", false, "redeploy all migrations")

	viper.BindPFlag("reset", migrateCmd.PersistentFlags().Lookup("reset"))
	viper.BindPFlag("default_network", migrateCmd.PersistentFlags().Lookup("network"))
	viper.SetDefault("default_network", "dev")
}
