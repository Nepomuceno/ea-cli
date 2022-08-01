package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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
			name, _ := cmd.Flags().GetString("name")
			displayName, _ := cmd.Flags().GetString("display-name")
			subscriptionOwner, _ := cmd.Flags().GetString("subscription-owner")
			tenantID, _ := cmd.Flags().GetString("sub-tenant")
			eaAccount, _ := cmd.Flags().GetString("enrollment-account")
			workload := getWorkload(cmd)
			subs, err := aliasClient.BeginCreate(context.Background(), name, armsubscription.PutAliasRequest{
				Properties: &armsubscription.PutAliasRequestProperties{
					BillingScope: &eaAccount,
					DisplayName:  &displayName,
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
			s, _ := json.MarshalIndent(resp, "", "\t")
			cmd.Print(string(s), "\n")
			return nil
		},
	}

	subscriptionAcceptCmd = &cobra.Command{
		Use:   "accept",
		Short: "Accepts a subscription ownership",
		RunE: func(cmd *cobra.Command, args []string) error {
			azcred, err := getCredentials(cmd)
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
			subscriptionID, err := cmd.Flags().GetString("subscription")
			if err != nil {
				return err
			}
			token, err := azcred.GetToken(context.Background(), policy.TokenRequestOptions{
				Scopes: []string{"https://management.azure.com/.default"},
			})
			if err != nil {
				return err
			}
			urlPath := "https://management.azure.com/providers/Microsoft.Subscription/subscriptions/{subscriptionId}/acceptOwnership?api-version=2021-10-01"
			urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
			req, err := runtime.NewRequest(context.Background(), http.MethodPost, urlPath)
			if err != nil {
				return err
			}
			req.Raw().Header["Accept"] = []string{"application/json"}
			req.Raw().Header["Authorization"] = []string{"Bearer " + token.Token}
			body := armsubscription.AcceptOwnershipRequest{
				Properties: &armsubscription.AcceptOwnershipRequestProperties{
					DisplayName:       &displayName,
					ManagementGroupID: &managementGroupID,
				},
			}
			err = runtime.MarshalAsJSON(req, body)
			if err != nil {
				return err
			}
			client := &http.Client{}
			response, err := client.Do(req.Raw())
			if err != nil {
				return err
			}
			s, _ := json.MarshalIndent(response, "", "\t")
			cmd.Print(string(s), "\n")
			cmd.Println("Subscription ownership accepted", response.StatusCode)
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
	subscriptionCreateCmd.Flags().StringP("name", "n", "", "Name of the subscription")
	subscriptionCreateCmd.Flags().StringP("display-name", "d", "", "Display name of the subscription")
	subscriptionCreateCmd.Flags().StringP("sub-tenant", "t", "", "Tenant ID of the subscription")
	subscriptionCreateCmd.Flags().StringP("workload", "w", "Production", "Workload of the subscription")
	subscriptionCreateCmd.Flags().StringP("subscription-owner", "o", "", "Subscription owner")
	subscriptionCreateCmd.Flags().StringP("enrollment-account", "e", "", "Enrollment account")
	err := subscriptionCreateCmd.MarkFlagRequired("name")
	PanicIfErr(err)
	err = subscriptionCreateCmd.MarkFlagRequired("tenant-id")
	PanicIfErr(err)
	err = subscriptionCreateCmd.MarkFlagRequired("subscription-owner")
	PanicIfErr(err)
	err = subscriptionCreateCmd.MarkFlagRequired("enrollment-account")
	PanicIfErr(err)
	subscriptionCmd.AddCommand(subscriptionCreateCmd)

	subscriptionAcceptCmd.Flags().StringP("subscription", "s", "", "Id of the subscription that you are accepting ownership for")
	subscriptionAcceptCmd.Flags().StringP("display-name", "d", "", "Display name of the subscription")
	subscriptionAcceptCmd.Flags().StringP("management-group-id", "m", "", "Management group ID of the subscription")
	err = subscriptionAcceptCmd.MarkFlagRequired("display-name")
	PanicIfErr(err)
	err = subscriptionAcceptCmd.MarkFlagRequired("subscription")
	PanicIfErr(err)
	subscriptionCmd.AddCommand(subscriptionAcceptCmd)

	rootCmd.AddCommand(subscriptionCmd)
}
