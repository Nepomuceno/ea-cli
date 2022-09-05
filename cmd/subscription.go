package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/spf13/cobra"
)

var subscriptionCmd = &cobra.Command{
	Use:   "sub",
	Short: "Manage subscriptions",
}
var subscriptionListCmd = &cobra.Command{
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
var subscriptionCreateCmd = &cobra.Command{
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
		managementGroupID, _ := cmd.Flags().GetString("management-group-id")
		workload := getWorkload(cmd)
		subs, err := aliasClient.BeginCreate(context.Background(), name, armsubscription.PutAliasRequest{
			Properties: &armsubscription.PutAliasRequestProperties{
				BillingScope: &eaAccount,
				DisplayName:  &displayName,
				AdditionalProperties: &armsubscription.PutAliasRequestAdditionalProperties{
					SubscriptionTenantID: &tenantID,
					SubscriptionOwnerID:  &subscriptionOwner,
					ManagementGroupID:    &managementGroupID,
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

var subscriptionAcceptCmd = &cobra.Command{
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
		if managementGroupID != "" {
			managementGroupID = fmt.Sprintf("/providers/Microsoft.Management/managementGroups/%s", managementGroupID)
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
		managementGroupID = fmt.Sprintf("/providers/Microsoft.Management/managementGroups/%s", managementGroupID)
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
		// print json body
		s, err := json.MarshalIndent(body, "", "\t")
		if err != nil {
			return err
		}
		cmd.Print(string(s), "\n")
		// send request
		err = runtime.MarshalAsJSON(req, body)
		if err != nil {
			return err
		}
		client := &http.Client{}
		response, err := client.Do(req.Raw())
		if err != nil {
			return err
		}
		s, err = json.MarshalIndent(response.Body, "", "\t")
		if err != nil {
			return err
		}
		cmd.Print(string(s), "\n Create Response \n")
		if response.StatusCode > 300 {
			return fmt.Errorf("failed to accept ownership of subscription %s code %d", subscriptionID, response.StatusCode)
		}
		cmd.Println("Subscription ownership accepted", response.StatusCode)
		return nil
	},
}

func getWorkload(cmd *cobra.Command) armsubscription.Workload {
	workload := cmd.Flag("workload").Value.String()
	if workload == "DevTest" {
		return armsubscription.WorkloadDevTest
	}
	return armsubscription.WorkloadProduction
}

func init() {

	subscriptionCmd.AddCommand(subscriptionListCmd)
	subscriptionCreateCmd.Flags().StringP("name", "n", "", "Name of the subscription")
	subscriptionCreateCmd.Flags().StringP("display-name", "d", "", "Display name of the subscription")
	subscriptionCreateCmd.Flags().StringP("sub-tenant", "t", "", "Tenant ID of the subscription")
	subscriptionCreateCmd.Flags().StringP("workload", "w", "Production", "Workload of the subscription")
	subscriptionCreateCmd.Flags().StringP("subscription-owner", "o", "", "Subscription owner")
	subscriptionCreateCmd.Flags().StringP("enrollment-account", "e", "", "Enrollment account")
	subscriptionCreateCmd.Flags().StringP("management-group-id", "g", "", "Management group ID")
	err := subscriptionCreateCmd.MarkFlagRequired("name")
	PanicIfErr(err)
	err = subscriptionCreateCmd.MarkFlagRequired("sub-tenant")
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
