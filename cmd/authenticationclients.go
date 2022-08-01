package cmd

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/spf13/cobra"
)

func getAccountClient(cmd *cobra.Command) (*armbilling.AccountsClient, error) {
	cred, err := getCredentials(cmd)
	if err != nil {
		return nil, err
	}
	client, err := armbilling.NewAccountsClient(cred, nil)
	return client, err
}

func getEnrollmentAccountClient(cmd *cobra.Command) (*armbilling.EnrollmentAccountsClient, error) {
	cred, err := getCredentials(cmd)
	if err != nil {
		return nil, err
	}
	client, err := armbilling.NewEnrollmentAccountsClient(cred, nil)
	return client, err
}

func getRoleDefinitionsClient(cmd *cobra.Command) (*armbilling.RoleDefinitionsClient, error) {
	cred, err := getCredentials(cmd)
	if err != nil {
		return nil, err
	}
	client, err := armbilling.NewRoleDefinitionsClient(cred, nil)
	return client, err
}

func getAliasClient(cmd *cobra.Command) (*armsubscription.AliasClient, error) {
	cred, err := getCredentials(cmd)
	if err != nil {
		return nil, err
	}
	client, err := armsubscription.NewAliasClient(cred, nil)
	return client, err
}
func getCredentials(cmd *cobra.Command) (azcore.TokenCredential, error) {
	useServicePrincipal, err := cmd.Flags().GetBool("service-principal")
	if err != nil {
		return nil, err
	}
	if useServicePrincipal {
		tenantID, err := cmd.Flags().GetString("login-tenant")
		if err != nil {
			return nil, err
		}
		azidentity.NewClientSecretCredential(tenantID, "", "", nil)
	}

	opt := azidentity.DefaultAzureCredentialOptions{
		TenantID: cmd.Flag("login-tenant").Value.String(),
	}
	credDefault, err := azidentity.NewDefaultAzureCredential(&opt)

	return credDefault, err
}
