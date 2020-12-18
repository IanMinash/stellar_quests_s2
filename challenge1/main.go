package main

import (
	"encoding/hex"
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

	newAccountID, newAccountIDPresent := os.LookupEnv("ACC_ID")
	if !newAccountIDPresent {
		log.Fatalln("ACC_ID is not defined in the environment file. Please define it and try again")
	}

	client := horizonclient.DefaultTestNetClient

	// Create keypair for issuing account
	creatorKp := keypair.MustRandom()
	fmt.Printf("Creator keypair: \n\tAddress:\t%s\n\tSecret Key:\t%s\nPlease store these keys if you wish to use the account later.\n", creatorKp.Address(), creatorKp.Seed())

	// Create the issuing account by funding it with the testnet's friendbot
	client.Fund(creatorKp.Address())

	creatorReq := horizonclient.AccountRequest{AccountID: creatorKp.Address()}
	creatorAccount, err := client.AccountDetail(creatorReq)
	if err != nil {
		log.Fatalln(err)
	}

	var txnMemo [32]byte
	hexBytes, err := hex.DecodeString("e3366fcb087bdb2381b7069a19405b748da831c18145eba25654d1092e93ef37")
	if err != nil {
		log.Fatalln(err)
	}
	copy(txnMemo[:], hexBytes)

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &creatorAccount,
			Operations: []txnbuild.Operation{&txnbuild.CreateAccount{
				Destination: newAccountID,
				Amount:      "5000",
			}},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
			IncrementSequenceNum: true,
			Memo:                 txnbuild.MemoHash(txnMemo),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, creatorKp)
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
		fmt.Printf("Account %s successfully created.\n", newAccountID)
	}

}
