package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling"
	"github.com/spf13/cobra"
)

// Magic string the defines a subscription creator role it is constant for all azure
const SUBSCRIPTION_CREATOR string = "cfff8e42-45ec-463a-9ae9-276083fcf6a9"

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
	accountGivePermissionCmd = &cobra.Command{
		Use: "give-permission",
		RunE: func(cmd *cobra.Command, args []string) error {
			azcred, err := getCredentials(cmd)
			if err != nil {
				return err
			}
			token, err := azcred.GetToken(context.Background(), policy.TokenRequestOptions{
				Scopes: []string{"https://management.azure.com/.default"},
			})
			if err != nil {
				return err
			}
			billingAccountNumber, _ := cmd.Flags().GetString("billing-account-number")
			enrrolementAccountNumber, _ := cmd.Flags().GetString("enrollment-account-number")
			roleDefinitionID := fmt.Sprintf("billingAccounts/%s/enrollmentAccounts/%s/billingRoleAssignments/%s", billingAccountNumber, enrrolementAccountNumber, SUBSCRIPTION_CREATOR)
			urlPath := "https://management.azure.com/providers/Microsoft.Billing/{roleDefinitionID}?api-version=2019-10-01-preview"
			urlPath = strings.ReplaceAll(urlPath, "{roleDefinitionID}", roleDefinitionID)
			req, err := runtime.NewRequest(context.Background(), http.MethodPost, urlPath)
			if err != nil {
				return err
			}
			req.Raw().Header["Accept"] = []string{"application/json"}
			req.Raw().Header["Authorization"] = []string{"Bearer " + token.Token}
			principalID, _ := cmd.Flags().GetString("principal-id")
			principalTenantID, _ := cmd.Flags().GetString("principal-tenant-id")
			body := armbilling.RoleAssignment{
				Properties: &armbilling.RoleAssignmentProperties{
					PrincipalID:       &principalID,
					PrincipalTenantID: &principalTenantID,
					RoleDefinitionID:  &roleDefinitionID,
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
			cmd.Println("Role applied", response.StatusCode)
			return nil
		}}
)

func initAccount(rootCmd *cobra.Command) {
	accountCmd.AddCommand(accountListCmd)
	accountGivePermissionCmd.Flags().String("principal-id", "u", "User ID to assign the permission to")
	accountGivePermissionCmd.Flags().String("principal-tenant-id", "t", "Tenant ID to assign the permission to")
	accountGivePermissionCmd.Flags().String("billing-account-number", "b", "Billing account number to assign the permission to")
	accountGivePermissionCmd.Flags().String("enrollment-account-number", "e", "Enrollment account number to assign the permission to")

	err := accountGivePermissionCmd.MarkFlagRequired("principal-id")
	PanicIfErr(err)
	err = accountGivePermissionCmd.MarkFlagRequired("principal-tenant-id")
	PanicIfErr(err)
	err = accountGivePermissionCmd.MarkFlagRequired("billing-account-number")
	PanicIfErr(err)
	err = accountGivePermissionCmd.MarkFlagRequired("enrollment-account-number")
	PanicIfErr(err)

	accountCmd.AddCommand(accountGivePermissionCmd)
	rootCmd.AddCommand(accountCmd)
}
