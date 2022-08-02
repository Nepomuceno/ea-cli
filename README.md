# EA CLI 

EA Cli helps you to manage the creation of your subscriptions and helps you to utilize them.



## Installation

### Prerequisites

There are no prereq if you want to use the relase binaries.

If you want to get it from source you will need go 1.8 or higher.

    ```bash
    go install github.com/nepomuceno/ea-cli
    ```

## Usage

To use ea-cli you will need to hava a calid azure user.

### Authentication

The ea-cli respect the same methods of authentication as you can use [here](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication?tabs=bash#2-authenticate-with-azure).

You can also specify the parameter `--service-principal`, `--username`, `--password` and `--tenant` on any call to use a service principal authentication.

The order of precedence when aythenticating is:
- command line arguments
- environment variables
- manged identity
- az cli

### Available commands 

You can see a list of commands by just typing `ea-cli help`. If you want to see the help of a specific command you can type `ea-cli help <command>`.

## Sample operations

### Give a service principal subscription creation permissions

In order to do that you will need to know your billing account number ( not the GUID but the actual number) and your enrrolment account number ( not the GUID but the account number).

```bash
ea-cli account give-creator-permission --billing-account-number <billing-account-number> --enrollment-account-number <enrollment-account-number> --principal-id <principal-id> --principal-tenant-id <principal-tenant-id>
```

## Contributing

In order to develop you will need to have go 1.18 installed.
It is recommended that you also have the following tools installed:
- golangci-lint   - Installation: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- goreleaser - Installation: `go install github.com/goreleaser/goreleaser@latest`

Before pushing your code you will need to run the following command:

```bash
golangci-lint run
```