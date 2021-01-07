package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stellar/go/txnbuild"

	"github.com/joho/godotenv"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("An error occured while trying to read from .env")
	}

	sponsorSignKey, sponsorSignKeyPresent := os.LookupEnv("SPONSOR_SIGN")
	if !sponsorSignKeyPresent {
		log.Fatalln("SPONSOR_SIGN is not defined in the environment file. Please define it and try again")
	}

	senderSignKey, senderSignKeyPresent := os.LookupEnv("SENDER_SIGN")
	if !senderSignKeyPresent {
		log.Fatalln("SENDER_SIGN is not defined in the environment file. Please define it and try again")
	}

	sponsorKp := keypair.MustParse(sponsorSignKey)
	senderKp := keypair.MustParse(senderSignKey)

	client := horizonclient.DefaultTestNetClient

	senderAccount, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: senderKp.Address()})
	if err != nil {
		log.Fatalln(err)
	}

	sponsorAccount, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: sponsorKp.Address()})
	if err != nil {
		log.Fatalln(err)
	}

	paymentTx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &senderAccount,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: sponsorKp.Address(),
					Amount:      "15",
					Asset:       txnbuild.NativeAsset{},
				},
			},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
			IncrementSequenceNum: true,
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	paymentTx, err = paymentTx.Sign(network.TestNetworkPassphrase, senderKp.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	feebumpedTx, err := txnbuild.NewFeeBumpTransaction(
		txnbuild.FeeBumpTransactionParams{
			BaseFee:    txnbuild.MinBaseFee,
			Inner:      paymentTx,
			FeeAccount: sponsorAccount.AccountID,
		},
	)

	feebumpedTx, err = feebumpedTx.Sign(network.TestNetworkPassphrase, sponsorKp.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	feebumpedTxEnvelope, err := feebumpedTx.Base64()
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SubmitTransactionXDR(feebumpedTxEnvelope)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.Successful {
		fmt.Printf("Fee bumped transaction submitted successfully. %s.\n", resp.ID)
	}

}
