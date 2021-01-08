package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	timestampStr, timestampPresent := os.LookupEnv("TIMESTAMP")
	if !timestampPresent {
		log.Fatalln("TIMESTAMP is not defined in the environment file. Please define it and try again")
	}
	timestamp, _ := strconv.ParseInt(timestampStr, 10, 64)

	kp := keypair.MustParse(signKey)

	// Rules for when the claimable balance can be claimed
	claimPredicate := txnbuild.NotPredicate(txnbuild.BeforeAbsoluteTimePredicate(timestamp))

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
				&txnbuild.CreateClaimableBalance{
					Destinations: []txnbuild.Claimant{
						txnbuild.NewClaimant(kp.Address(), &claimPredicate),
					},
					Asset:  txnbuild.NativeAsset{},
					Amount: "100",
				},
			},
		},
	)

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	claimableBalanceID, err := tx.ClaimableBalanceID(0)
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

	fmt.Printf("Create Claimable Balance Transaction ID:\t%s\n", resp.ID)

	fmt.Printf("Claimable Balance ID:\t%s\nStore this ID to use when claiming\n", claimableBalanceID)
}
