package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var (
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "Manage EA accounts",
	}
	accountListCmd = &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			acClient, err := getEnrollmentAccountClient(cmd)
			if err != nil {
				return err
			}
			acPager := acClient.NewListPager(nil)
			for acPager.More() {
				page, err := acPager.NextPage(context.Background())
				if err != nil {
					return err
				}
				for _, v := range page.Value {
					cmd.Println("Name:", *v.Name, "Principal:", *v.Properties.PrincipalName, "ID:", *v.ID)
				}
			}
			return nil
		},
	}
)

func initAccount(rootCmd *cobra.Command) {
	accountCmd.AddCommand(accountListCmd)
	rootCmd.AddCommand(accountCmd)
}
