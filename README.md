# EWallet-API
A simple EWallet HTTP application that implements a transaction processing system.

### Project structure

#### api
Server, routes, handlers and errors for HTTP Server using Chi.

#### config
Configuration parsing from environment variables.

#### wallets
Some general structures for data representation, errors, an interface to the database and it's implementation using MongoDB.

### Build
    sudo docker compose up --build (-d)