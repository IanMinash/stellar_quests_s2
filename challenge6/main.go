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

	// Get private key of your account
	signKey, signKeyPresent := os.LookupEnv("SIGN_KEY")
	if !signKeyPresent {
		log.Fatalln("SIGN_KEY is not defined in the environment file. Please define it and try again")
	}

	kp := keypair.MustParse(signKey)

	// Keypair of account we are creating
	createdKp := keypair.MustRandom()

	client := horizonclient.DefaultTestNetClient

	request := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(request)
	if err != nil {
		log.Fatalln(err)
	}

	// Create account transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
			Operations: []txnbuild.Operation{
				&txnbuild.BeginSponsoringFutureReserves{
					SponsoredID:   createdKp.Address(),
					SourceAccount: &sourceAccount,
				},
				&txnbuild.CreateAccount{
					Destination:   createdKp.Address(),
					Amount:        "0",
					SourceAccount: &sourceAccount,
				},
				&txnbuild.EndSponsoringFutureReserves{
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: createdKp.Address(),
					},
				},
			},
		},
	)

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full), createdKp)
	if err != nil {
		log.Fatalln(err)
	}

	txe, err := tx.Base64()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(txe)

	resp, err := client.SubmitTransactionXDR(txe)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Create Account Transaction ID:\t%s\n", resp.ID)
}
