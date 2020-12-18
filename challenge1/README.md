## Challenge 1: Create and fund a Stellar account

In this challenge your task is to create and fund a brand new Stellar account with 5000 XLM on the testnet.

Include the SHA256 hash of the string `Stellar Quest Series 2` as the `MEMO_HASH` in the transaction memo field.

You will be required to use the `createAccount` operation.

### Requirements

For this implementation, you'll need to have an initial account with some Stellar Lumens. You can create one from the [Stellar Laboratory](https://laboratory.stellar.org/#account-creator?network=test) and use the Friendbot to fund it with the initial balance. This account will be used to issue the `CreateAccount` operation.

Create a `.env` file in the current directory and add the following:

```
ACC_ID=<public key of the account to be created>
SIGN_KEY=<secret key of the account created on the Stellar Laboratory>
```
