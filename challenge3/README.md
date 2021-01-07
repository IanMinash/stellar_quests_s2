## Challenge 3: Create and submit a fee bump transaction

Fee channels are a common best practice in Stellar development. Their goal is to delegate fee payments away from user accounts for an improved UX.

In this challenge your task is to create and execute a fee bump transaction which consumes the sequence number from your account but the transaction fee from some other account.

## Requirements

Create a `.env` file in the current directory and add the following:

```
SPONSOR_SIGN=<private key of the account to as the fee source> // You can use the private key generated in challenge 1
SENDER_SIGN=<private key of the account being sponsored>
```
