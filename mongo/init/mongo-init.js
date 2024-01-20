db.createUser(
    {
        user: "GolangBanker",
        pwd: "Golang",
        roles: [
            {
                role: "readWrite",
                db: "EWallet"
            }
        ]
    }
);