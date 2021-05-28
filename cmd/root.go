package cmd

import (
	hedera "github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/spf13/cobra"
)

var operatorKey string
var operatorID string
var network string

var rootCmd = &cobra.Command{
	Use:   "third-act",
	Short: "Create and mint tokens from CLI",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
var client *hedera.Client

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&operatorKey, "operatorKey", "", "")
	rootCmd.PersistentFlags().StringVar(&operatorID, "operatorID", "", "")
	rootCmd.PersistentFlags().StringVar(&network, "network", "testnet", "")

	rootCmd.MarkPersistentFlagRequired("operatorKey")
	rootCmd.MarkPersistentFlagRequired("operatorID")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(mintCmd)
}

func initClient() {
	var err error
	operatorPrivateKey, err := hedera.PrivateKeyFromString(operatorKey)
	if err != nil {
		panic("failed to parse operator key")
	}

	operatorAccountID, err := hedera.AccountIDFromString(operatorID)
	if err != nil {
		panic("failed to parse operator ID")
	}

	client, err = hedera.ClientForName(network)
	if err != nil {
		panic("Invalid network provided. Must be one of: `mainnet`, `testnet`, or `previewnet`")
	}

	client.SetOperator(operatorAccountID, operatorPrivateKey)
}
