package cmd

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling"
	"github.com/spf13/cobra"
)

var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Manage EA roles",
}

var roleListCmd = &cobra.Command{
	Use: "list",
	RunE: func(cmd *cobra.Command, args []string) error {
		acClient, err := getAccountClient(cmd)
		if err != nil {
			return err
		}
		roleClients, err := getRoleDefinitionsClient(cmd)
		if err != nil {
			return err
		}
		acPager := acClient.NewListPager(nil)
		result := make([]*armbilling.RoleDefinition, 0)
		for acPager.More() {
			page, err := acPager.NextPage(context.Background())
			if err != nil {
				return err
			}

			for _, v := range page.Value {
				roleResult := roleClients.NewListByBillingAccountPager(*v.Name, nil)
				for roleResult.More() {
					pageresult, err := roleResult.NextPage(context.Background())
					if err != nil {
						return err
					}
					result = append(result, pageresult.Value...)

				}
			}
		}
		s, _ := json.MarshalIndent(result, "", "\t")
		cmd.Print(string(s), "\n")
		return nil
	},
}

func init() {
	roleCmd.AddCommand(roleListCmd)
	rootCmd.AddCommand(roleCmd)
}
