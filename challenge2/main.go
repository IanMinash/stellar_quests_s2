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

	issuerSignKey, issuerSignKeyPresent := os.LookupEnv("ISSUER_SIGN")
	if !issuerSignKeyPresent {
		log.Fatalln("ISSUER_SIGN is not defined in the environment file. Please define it and try again")
	}

	receiverSignKey, receiverSignKeyPresent := os.LookupEnv("RECEIVER_SIGN")
	if !receiverSignKeyPresent {
		log.Fatalln("RECEIVER_SIGN is not defined in the environment file. Please define it and try again")
	}

	issuerKp := keypair.MustParse(issuerSignKey)
	receiverKp := keypair.MustParse(receiverSignKey)

	client := horizonclient.DefaultTestNetClient

	receiverAccount, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: receiverKp.Address()})
	if err != nil {
		log.Fatalln(err)
	}

	issuerAcount, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: issuerKp.Address()})
	if err != nil {
		log.Fatalln(err)
	}

	mbogiToken := txnbuild.CreditAsset{
		Code:   "MBOGI",
		Issuer: issuerKp.Address(),
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &receiverAccount,
			Operations: []txnbuild.Operation{&txnbuild.ChangeTrust{
				Line: mbogiToken,
			},
				&txnbuild.Payment{
					Destination:   receiverKp.Address(),
					Amount:        "1",
					Asset:         mbogiToken,
					SourceAccount: &issuerAcount,
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

	tx, err = tx.Sign(network.TestNetworkPassphrase, receiverKp.(*keypair.Full), issuerKp.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	txEnvelope, err := tx.Base64()
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SubmitTransactionXDR(txEnvelope)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.Successful {
		fmt.Printf("Asset transferred successfully.\n")
	}

}
