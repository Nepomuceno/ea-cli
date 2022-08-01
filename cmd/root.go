package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ea-cli",
		Short: "A CLI to help you create subscriptions",
		Long:  `Creating subscriptions can sometimes be a confusing process. This CLI helps you create subscriptions on Azure.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().Bool("service-principal", false, "Use service principal to authenticate")
	rootCmd.PersistentFlags().StringP("username", "u", "", "Username to authenticate")
	rootCmd.PersistentFlags().StringP("password", "p", "", "Password to authenticate")
	rootCmd.PersistentFlags().String("tenant", "", "Tenant ID to authenticate")
	rootCmd.MarkFlagsRequiredTogether("username", "password", "tenant")
	initAccount(rootCmd)
	initRole(rootCmd)
	initBilling(rootCmd)
	initSubscription(rootCmd)
}
