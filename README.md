
# eWallet Application Testing Guide

Welcome to the testing guide for the eWallet application. This README provides instructions on how to download the Docker image for eWallet and includes a list of commands to interact with the app.

## Download the Application

To get started with testing the eWallet app, you'll need to download the Docker image. Execute the following command in your terminal to pull the latest version of the eWallet Docker image:

Bash

docker pull artforme/ewallet:latest

## Testing Commands

Below you will find a set of HTTP endpoints that you can use to test various functionalities of the eWallet application. Ensure your Docker container is running and exposing the appropriate port (default is 8082) before performing these operations.

### Create Wallet

To create a new wallet, send a POST request to the following endpoint:

http://0.0.0.0:8082/api/v1/wallet

### Transfer Funds

To transfer funds from one wallet to another, send a POST request with the payload containing the wallet ID and the amount to transfer:

http://0.0.0.0:8082/api/v1/wallet/{yourWalletID}/send

Payload:
{
    "walletId": "{yourWalletID}",
    "amount": "9.23"
}

### Show Transaction History

To view the transaction history of a wallet, send a GET request to:

http://0.0.0.0:8082/api/v1/wallet/{yourWalletID}/history

### Show Wallet Details

To retrieve the details of a specific wallet, send a GET request to:
```bash
http://0.0.0.0:8082/api/v1/wallet/{yourWalletID}
```
