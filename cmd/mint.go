package cmd

import (
	hedera "github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/spf13/cobra"
)

var tokenID string
var amount uint64

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint tokens",
	Example: `
Mint tokens:
    token-creator create \
        --operatorKey="<operatorKey>" \
        --operatorID="<operatorID>" \
        --tokenID="<tokenID>" \
        --amount="<amount>" \
        --supplyKey="<used for signing; required if supplyKey of token is not operatorKey>"
    `,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var privateKey hedera.PrivateKey

		initClient()

		transaction := hedera.NewTokenMintTransaction()

		token, err := hedera.TokenIDFromString(tokenID)
		if err != nil {
			panic("failed to parse token ID")
		}

		transaction.SetTokenID(token)
		transaction.SetAmount(amount)

		if supplyKey != "" {
			privateKey, err = hedera.PrivateKeyFromString(supplyKey)
			_, err = transaction.FreezeWith(client)
			if err != nil {
				panic(err)
			}

			transaction.Sign(privateKey)
		}

		response, err := transaction.Execute(client)
		if err != nil {
			panic(err)
		}

		receipt, err := response.GetReceipt(client)
		if err != nil {
			panic(err)
		}

		println("New Total Supply:", receipt.TotalSupply)
	},
}

func init() {
	mintCmd.PersistentFlags().StringVar(&tokenID, "tokenID", "", "")
	mintCmd.PersistentFlags().StringVar(&supplyKey, "suppyKey", "", "")
	mintCmd.PersistentFlags().Uint64Var(&amount, "amount", 0, "")

	mintCmd.MarkPersistentFlagRequired("tokenID")
	mintCmd.MarkPersistentFlagRequired("amount")
}
