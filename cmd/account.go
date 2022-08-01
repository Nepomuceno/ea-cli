package cmd

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling"
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
			result := make([]*armbilling.EnrollmentAccountSummary, 0)
			for acPager.More() {
				page, err := acPager.NextPage(context.Background())
				if err != nil {
					return err
				}
				result = append(result, page.Value...)
			}
			s, _ := json.MarshalIndent(result, "", "\t")
			cmd.Print(string(s), "\n")
			return nil
		},
	}
)

func initAccount(rootCmd *cobra.Command) {
	accountCmd.AddCommand(accountListCmd)
	rootCmd.AddCommand(accountCmd)
}
