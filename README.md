
# eWallet Application Testing Guide

Welcome to the testing guide for the eWallet application. This README provides instructions on how to download the Docker image for eWallet and includes a list of commands to interact with the app.

## Download the Application

To get started with testing the eWallet app, you'll need to download the Docker image. Execute the following command in your terminal to pull the latest version of the eWallet Docker image:

```cmd
docker pull artforme/ewallet:latest
```

## Testing Commands

Below you will find a set of HTTP endpoints that you can use to test various functionalities of the eWallet application. Ensure your Docker container is running and exposing the appropriate port (default is 8082) before performing these operations.

### Create Wallet

```cmd
http://0.0.0.0:8082/api/v1/wallet
```
### Transfer

```cmd
http://0.0.0.0:8082/api/v1/wallet/{yourWalletID}/send
```
```JSON
{
    "walletId": "{yourWalletID}",
    "amount": "9.23"
}
```
### Show Transaction History

```cmd
http://0.0.0.0:8082/api/v1/wallet/{yourWalletID}/history
```
### Show Wallet Details

```cmd
http://0.0.0.0:8082/api/v1/wallet/{yourWalletID}
```
