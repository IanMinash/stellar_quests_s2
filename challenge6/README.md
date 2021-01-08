## Challenge 6: Sponsor the absolute minimum balance for a new account

All Stellar accounts must maintain a minimum balance of lumens. The mnimum balance is calculated by the `base reserve` which is currently 0.5 XLM using the following expression:

    Minimum Balance = (2 + # of entries + # of sponsoring entries - # of sponsored entries) * base reserve

Until protocol 15 the minimum balance had to be paid for by the account itself. However there are instances where it would be more convenient or even essential for these fees to be staked by some other account, a "sponsor" account.

In this challenge your task is to create a brand new 0 XLM balance account with the absolute minimum balance sponsored by your account.

## Requirements

Create a `.env` file in the current directory and add the following:

```
SIGN_KEY=<private key of your account>
```
