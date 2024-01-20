# EWallet-API
A simple EWallet HTTP application that implements a transaction processing system.

### Project structure

#### api
Server, routes, handlers and errors for HTTP Server using Chi.

#### config
Configuration parsing from environment variables.

#### wallets
Some general structures for data representation, errors, an interface to the database and it's implementation using MongoDB.

#### Database "schema"
- Collection "wallets" holds documents representing wallets:
    - "_id" - wallet id;
    - "balance" - wallet's money amount;
- Collection "history_{walletId}" holds transaction history for the given wallet. The format is:
    - "time": transaction time;
    - "from": wallet id to send money from;
    - "to": wallet id to send money to;
    - "amount": how much money to send;

### Build
    sudo docker compose up --build (-d)
