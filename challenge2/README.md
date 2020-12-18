## Challenge 2: Construct and execute a multi-operational transaction

Each Stellar transaction can include as many as 100 unique operations. This is an incredible feature as each transaction is atomic meaning either the whole group of operations succeeds or fails together.

In this challenge your task is to create a multi-operational transaction which creates a custom asset trustline on your account and pays that asset to your account from the issuing account all in the same transaction.

### Requirements

Create a `.env` file in the current directory and add the following:

```
ISSUER_SIGN=<private key of the account to issue the custom token> // You can use the private key generated in challenge 1
RECEIVER_SIGN=<private key of the account to receive the custom token>
```
