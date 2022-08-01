package cmd

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/spf13/cobra"
)

var (
	subscriptionCmd = &cobra.Command{
		Use:   "sub",
		Short: "Manage subscriptions",
	}
	subscriptionListCmd = &cobra.Command{
		Use:   "alias-list",
		Short: "List alias subscriptions created by the user",
		RunE: func(cmd *cobra.Command, args []string) error {
			aliasClient, err := getAliasClient(cmd)
			if err != nil {
				return err
			}
			result := make([]*armsubscription.AliasResponse, 0)
			subs, err := aliasClient.List(context.Background(), nil)
			if err != nil {
				return err
			}
			result = append(result, subs.Value...)
			s, _ := json.MarshalIndent(result, "", "\t")
			cmd.Print(string(s), "\n")
			return nil
		},
	}
	subscriptionCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a subscription",
		RunE: func(cmd *cobra.Command, args []string) error {
			aliasClient, err := getAliasClient(cmd)
			if err != nil {
				return err
			}
			displayName := cmd.Flag("display-name").Value.String()
			subscriptionOwner, _ := cmd.Flags().GetString("subscription-owner")
			tenantID, _ := cmd.Flags().GetString("tenant-id")
			workload := getWorkload(cmd)
			subs, err := aliasClient.BeginCreate(context.Background(), "test", armsubscription.PutAliasRequest{
				Properties: &armsubscription.PutAliasRequestProperties{
					DisplayName: &displayName,
					AdditionalProperties: &armsubscription.PutAliasRequestAdditionalProperties{
						SubscriptionTenantID: &tenantID,
						SubscriptionOwnerID:  &subscriptionOwner,
					},
					Workload: &workload,
				}}, nil)
			if err != nil {
				return err
			}
			resp, err := subs.PollUntilDone(context.Background(), nil)
			if err != nil {
				return err
			}
			cmd.Println(resp.Properties.SubscriptionID)
			return nil
		},
	}

	subscriptionAcceptCmd = &cobra.Command{
		Use:   "accept",
		Short: "Accepts a subscription ownership",
		RunE: func(cmd *cobra.Command, args []string) error {
			subscriptionClient, err := getSubscriptionClient(cmd)
			if err != nil {
				return err
			}
			displayName, err := cmd.Flags().GetString("display-name")
			if err != nil {
				return err
			}
			managementGroupID, err := cmd.Flags().GetString("management-group-id")
			if err != nil {
				return err
			}
			subscriptionClient.BeginAcceptOwnership(context.Background(), "", armsubscription.AcceptOwnershipRequest{
				Properties: &armsubscription.AcceptOwnershipRequestProperties{
					DisplayName:       &displayName,
					ManagementGroupID: &managementGroupID,
				},
			}, nil)
			return nil
		},
	}
)

func getWorkload(cmd *cobra.Command) armsubscription.Workload {
	workload := cmd.Flag("workload").Value.String()
	if workload == "DevTest" {
		return armsubscription.WorkloadDevTest
	}
	return armsubscription.WorkloadProduction
}

func initSubscription(rootCmd *cobra.Command) {

	subscriptionCmd.AddCommand(subscriptionListCmd)

	subscriptionCreateCmd.Flags().StringP("display-name", "d", "", "Display name of the subscription")
	subscriptionCreateCmd.Flags().StringP("tenant-id", "t", "", "Tenant ID of the subscription")
	subscriptionCreateCmd.Flags().StringP("workload", "w", "Production", "Workload of the subscription")
	subscriptionCreateCmd.Flags().StringP("subscription-owner", "o", "", "Subscription owner")
	subscriptionCreateCmd.MarkFlagRequired("display-name")
	subscriptionCreateCmd.MarkFlagRequired("tenant-id")
	subscriptionCreateCmd.MarkFlagRequired("workload")
	subscriptionCmd.AddCommand(subscriptionCreateCmd)

	subscriptionAcceptCmd.Flags().StringP("display-name", "d", "", "Display name of the subscription")
	subscriptionAcceptCmd.Flags().StringP("management-group-id", "m", "", "Management group ID of the subscription")
	subscriptionAcceptCmd.MarkFlagRequired("display-name")
	subscriptionCmd.AddCommand(subscriptionAcceptCmd)

	rootCmd.AddCommand(subscriptionCmd)
}
