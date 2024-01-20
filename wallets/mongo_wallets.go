package wallets

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/RichSweetWafer/EWallet-API/config"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoWallets struct {
	database   config.Database
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoWallets(config config.Database) *MongoWallets {
	return &MongoWallets{
		database: config,
	}
}

func (m *MongoWallets) connect(ctx context.Context) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(m.database.DatabaseURL).SetServerAPIOptions(serverAPI),
	)

	if err != nil {
		return err
	}

	m.client = client
	m.collection = m.client.Database(m.database.DatabaseName).Collection(m.database.CollectionName)

	return nil
}

func (m *MongoWallets) close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *MongoWallets) CreateWallet(ctx context.Context) (Wallet, error) {

	err := m.connect(ctx)
	if err != nil {
		return Wallet{}, err
	}
	defer m.close(ctx)

	// create wallet
	id := uuid.New()
	wallet := Wallet{
		ID:      id,
		Balance: defaultBalance,
	}

	if _, err := m.collection.InsertOne(ctx, wallet); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return Wallet{}, &DuplicateKeyError{ID: id}
		}
		return Wallet{}, err
	}

	// create wallet history collection
	idBin, err := id.MarshalBinary()
	if err != nil {
		return wallet, err
	}
	history := TransactionHistoryPrefix + hex.EncodeToString(idBin)
	if err := m.client.Database(m.database.DatabaseName).CreateCollection(ctx, history); err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (m *MongoWallets) CreateTransaction(ctx context.Context, params CreateTransactionParams) error {
	err := m.connect(ctx)
	if err != nil {
		return err
	}
	defer m.close(ctx)

	// find if both wallets exist

	var fromWallet Wallet

	if err := m.collection.FindOne(ctx, bson.M{"_id": params.From}).Decode(&fromWallet); err != nil {
		if err == mongo.ErrNoDocuments {
			return &WalletNotFoundError{}
		}
		return err
	}

	var toWallet Wallet

	if err := m.collection.FindOne(ctx, bson.M{"_id": params.To}).Decode(&toWallet); err != nil {
		if err == mongo.ErrNoDocuments {
			return &WalletNotFoundError{}
		}
		return err
	}

	// find if 'from' balance is large enough
	if fromWallet.Balance < params.Amount {
		return &WalletBalanceLow{}
	}

	// update wallets
	fromWallet.Balance -= params.Amount
	toWallet.Balance += params.Amount

	update := bson.M{
		"$set": bson.M{
			"balance": fromWallet.Balance,
		},
	}

	if _, err := m.collection.UpdateOne(ctx, bson.M{"_id": params.From}, update); err != nil {
		return err
	}

	update = bson.M{
		"$set": bson.M{
			"balance": toWallet.Balance,
		},
	}

	if _, err := m.collection.UpdateOne(ctx, bson.M{"_id": params.To}, update); err != nil {
		return err
	}

	// add history
	fromBin, err := params.From.MarshalBinary()
	if err != nil {
		return err
	}
	toBin, err := params.To.MarshalBinary()
	if err != nil {
		return err
	}

	transaction := Transaction{
		Time:   time.Now(),
		From:   hex.EncodeToString(fromBin),
		To:     hex.EncodeToString(toBin),
		Amount: params.Amount,
	}

	history := TransactionHistoryPrefix + transaction.From
	collection := m.client.Database(m.database.DatabaseName).Collection(history)

	if _, err := collection.InsertOne(ctx, transaction); err != nil {
		return err
	}

	history = TransactionHistoryPrefix + transaction.To
	collection = m.client.Database(m.database.DatabaseName).Collection(history)

	if _, err := collection.InsertOne(ctx, transaction); err != nil {
		return err
	}

	return nil
}

func (m *MongoWallets) GetHistory(ctx context.Context, id uuid.UUID) ([]Transaction, error) {
	err := m.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer m.close(ctx)

	idBin, err := id.MarshalBinary()
	if err != nil {
		return []Transaction{}, err
	}

	history := TransactionHistoryPrefix + hex.EncodeToString(idBin)
	collection := m.client.Database(m.database.DatabaseName).Collection(history)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var transactions []Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (m *MongoWallets) GetWallet(ctx context.Context, id uuid.UUID) (Wallet, error) {
	err := m.connect(ctx)
	if err != nil {
		return Wallet{}, err
	}
	defer m.close(ctx)

	var wallet Wallet

	if err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&wallet); err != nil {
		if err == mongo.ErrNoDocuments {
			return Wallet{}, &WalletNotFoundError{}
		}
		return Wallet{}, err
	}

	return wallet, nil
}
