package cmd

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling"
	"github.com/spf13/cobra"
)

var (
	billingCmd = &cobra.Command{
		Use:   "billing",
		Short: "Manage EA billing",
	}
	billingListCmd = &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			acClient, err := getAccountClient(cmd)
			if err != nil {
				return err
			}
			acPager := acClient.NewListPager(nil)
			result := make([]*armbilling.Account, 0)
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

func initBilling(rootCmd *cobra.Command) {
	billingCmd.AddCommand(billingListCmd)
	rootCmd.AddCommand(billingCmd)
}
