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

	// Get public key of the account you're revoking sponsorship from
	sponsoredID, sponsoredIDPresent := os.LookupEnv("SPONSORED_ID")
	if !sponsoredIDPresent {
		log.Fatalln("SPONSORED_ID is not defined in the environment file. Please define it and try again")
	}

	kp := keypair.MustParse(signKey)

	client := horizonclient.DefaultTestNetClient

	request := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(request)
	if err != nil {
		log.Fatalln(err)
	}

	// Because we are sponsoring the accounts minimum balance, revoking the sponsorship would cause the account to be in violation
	// of the minimum balance requirement if it has a nil balance. Therefore, we need to send XLM to the account to cover this
	// fee requirement before revoking the sponsorship.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Asset:       txnbuild.NativeAsset{},
					Amount:      "1",
					Destination: sponsoredID,
				},
				&txnbuild.RevokeSponsorship{
					Account:         &sponsoredID,
					SponsorshipType: txnbuild.RevokeSponsorshipTypeAccount,
				},
			},
		},
	)

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	txe, err := tx.Base64()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(txe)

	resp, err := client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Revoke Sponsorship Transaction ID:\t%s\n", resp.ID)
}
