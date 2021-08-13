package cmd

import (
	"time"

	hedera "github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/spf13/cobra"
)

var tokenName string
var tokenSymbol string
var tokenDecimals uint32
var initialSupply uint64
var treasuryAccountID string
var adminKey string
var kycKey string
var freezeKey string
var wipeKey string
var supplyKey string
var freezeDefault bool
var expirationTime string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create tokens",
	Example: `
Create a mutable token (set an adminKey):
    token-creator create \
        --operatorKey="<operatorKey>" \
        --operatorID="<operatorID>" \
        --tokenName="<tokenName>" \
        --tokenSymbol="<tokenSymbol>" \
        --treasuryAccountID="<treasuryAccountID>" \
        --adminKey="<the key required to sign token update transactions>"


Create a token which requires KYC (set a kycKey):
    token-creator create \
        --operatorKey="<operatorKey>" \
        --operatorID="<operatorID>" \
        --tokenName="<tokenName>" \
        --tokenSymbol="<tokenSymbol>" \
        --treasuryAccountID="<treasuryAccountID>" \
        --kycKey="<the key required to add KYC to token holders>"


Create token with limited supply (Do not set a supplyKey):
    token-creator create \
        --operatorKey="<operatorKey>" \
        --operatorID="<operatorID>" \
        --tokenName="<tokenName>" \
        --tokenSymbol="<tokenSymbol>" \
        --treasuryAccountID="<treasuryAccountID>" \
        --initialSupply="<limited supply amount>"
    `,
	Run: func(cmd *cobra.Command, args []string) {
		var expiration time.Time
		var err error

		var privateKey hedera.PrivateKey
		var publicKey hedera.PublicKey

		initClient()

		if expirationTime != "" {
			expiration, err = time.Parse(time.RFC3339, expirationTime)
			if err != nil {
				panic("expiration time provided in invalid format")
			}
		}

		transaction := hedera.NewTokenCreateTransaction()
		transaction.SetTokenName(tokenName)
		transaction.SetTokenSymbol(tokenSymbol)
		transaction.SetDecimals(uint(tokenDecimals))
		transaction.SetInitialSupply(initialSupply)

		treasury, err := hedera.AccountIDFromString(treasuryAccountID)
		if err != nil {
			panic("failed to parse treasury account ID")
		}

		transaction.SetTreasuryAccountID(treasury)

		if adminKey != "" {
			privateKey, err = hedera.PrivateKeyFromString(adminKey)

			if err != nil {
				publicKey, err = hedera.PublicKeyFromString(adminKey)
			} else {
				publicKey = privateKey.PublicKey()
			}

			transaction.SetAdminKey(publicKey)
		}

		if kycKey != "" {
			privateKey, err = hedera.PrivateKeyFromString(kycKey)

			if err != nil {
				publicKey, err = hedera.PublicKeyFromString(kycKey)
			} else {
				publicKey = privateKey.PublicKey()
			}

			transaction.SetKycKey(publicKey)
		}

		if freezeKey != "" {
			privateKey, err = hedera.PrivateKeyFromString(freezeKey)

			if err != nil {
				publicKey, err = hedera.PublicKeyFromString(freezeKey)
			} else {
				publicKey = privateKey.PublicKey()
			}

			transaction.SetFreezeKey(publicKey)
		}

		if wipeKey != "" {
			privateKey, err = hedera.PrivateKeyFromString(wipeKey)

			if err != nil {
				publicKey, err = hedera.PublicKeyFromString(wipeKey)
			} else {
				publicKey = privateKey.PublicKey()
			}

			transaction.SetWipeKey(publicKey)
		}

		if supplyKey != "" {
			privateKey, err = hedera.PrivateKeyFromString(supplyKey)

			if err != nil {
				publicKey, err = hedera.PublicKeyFromString(supplyKey)
			} else {
				publicKey = privateKey.PublicKey()
			}

			transaction.SetSupplyKey(publicKey)
		}

		transaction.SetFreezeDefault(freezeDefault)

		if expirationTime != "" {
			transaction.SetExpirationTime(expiration)
		}

		response, err := transaction.Execute(client)
		if err != nil {
			panic(err)
		}

		receipt, err := response.GetReceipt(client)
		if err != nil {
			panic(err)
		}

		println(receipt.TokenID.String())
	},
}

func init() {
	createCmd.PersistentFlags().StringVar(&tokenName, "tokenName", "", "")
	createCmd.PersistentFlags().StringVar(&tokenSymbol, "tokenSymbol", "", "")
	createCmd.PersistentFlags().Uint32Var(&tokenDecimals, "decimals", 0, "")
	createCmd.PersistentFlags().Uint64Var(&initialSupply, "initialSupply", 0, "")
	createCmd.PersistentFlags().StringVar(&treasuryAccountID, "treasuryAccountID", "", "")
	createCmd.PersistentFlags().StringVar(&adminKey, "adminKey", "", "")
	createCmd.PersistentFlags().StringVar(&kycKey, "kycKey", "", "")
	createCmd.PersistentFlags().StringVar(&freezeKey, "freezeKey", "", "")
	createCmd.PersistentFlags().StringVar(&wipeKey, "wipeKey", "", "")
	createCmd.PersistentFlags().StringVar(&supplyKey, "supplyKey", "", "")
	createCmd.PersistentFlags().BoolVar(&freezeDefault, "freezeDefault", false, "")
	createCmd.PersistentFlags().StringVar(&expirationTime, "expirationTime", "", "RFC339 format")

	createCmd.MarkPersistentFlagRequired("tokenName")
	createCmd.MarkPersistentFlagRequired("tokenSymbol")
	createCmd.MarkPersistentFlagRequired("treasuryAccountID")
}
