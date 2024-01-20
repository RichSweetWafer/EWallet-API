package api

import (
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/RichSweetWafer/EWallet-API/wallets"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

/*----------WALLET RESPONSE----------*/
type WalletResponse struct {
	ID      string `json:"id"`
	Balance int64  `json:"balance"`
}

func (response WalletResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewWalletResponse(w wallets.Wallet) WalletResponse {
	bin, _ := w.ID.MarshalBinary()
	return WalletResponse{
		ID:      hex.EncodeToString(bin),
		Balance: w.Balance,
	}
}

/*--------TRANSACTION RESPONSE--------*/
type TransactionResponse struct {
	Time   time.Time `json:"time"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount int64     `json:"amount"`
}

func (response TransactionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewTransactionResponse(t wallets.Transaction) TransactionResponse {
	return TransactionResponse{
		Time:   t.Time,
		From:   t.From,
		To:     t.To,
		Amount: t.Amount,
	}
}

func NewHistoryResponse(transactions []wallets.Transaction) []render.Renderer {
	list := []render.Renderer{}
	for _, transaction := range transactions {
		response := NewTransactionResponse(transaction)
		list = append(list, response)
	}
	return list
}

/*--------TRANSACTION REQUEST--------*/

type TransactionRequest struct {
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

func (request *TransactionRequest) Bind(r *http.Request) error {
	return nil
}

/*--------HANDLERS--------*/

func (s *Server) handleStatusCheck(w http.ResponseWriter, r *http.Request) {
	resp := []byte{'o', 'k', 'a', 'y'}
	w.Write(resp)
}

func (s *Server) handleCreateWallet(w http.ResponseWriter, r *http.Request) {

	wallet, err := s.wallets.CreateWallet(r.Context())
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	render.Render(w, r, NewWalletResponse(wallet))
}

func (s *Server) handleTransaction(w http.ResponseWriter, r *http.Request) {
	walletIdParam := chi.URLParam(r, "walletId")
	walletId, err := uuid.Parse(walletIdParam)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	request := &TransactionRequest{}

	if err := render.Bind(r, request); err != nil {
		render.Render(w, r, ErrBadRequest)
	}

	params := wallets.CreateTransactionParams{
		From:   walletId,
		To:     uuid.MustParse(request.To),
		Amount: request.Amount,
	}

	err = s.wallets.CreateTransaction(r.Context(), params)
	if err != nil {
		var dupKeyErr *wallets.WalletNotFoundError
		if errors.As(err, &dupKeyErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrBadRequest)
		}
		return
	}

	w.WriteHeader(200)
	w.Write(nil)

}

func (s *Server) handleWalletHistory(w http.ResponseWriter, r *http.Request) {
	walletIdParam := chi.URLParam(r, "walletId")
	walletId, err := uuid.Parse(walletIdParam)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	history, err := s.wallets.GetHistory(r.Context(), walletId)
	if err != nil {
		var rfnErr *wallets.WalletNotFoundError
		if errors.As(err, &rfnErr) {
			render.Render(w, r, ErrNotFound)
		}
		// } else {
		// 	render.Render(w, r, ErrInternalServerError)
		// }
		return
	}

	render.RenderList(w, r, NewHistoryResponse(history))
}

func (s *Server) handleWalletState(w http.ResponseWriter, r *http.Request) {
	walletIdParam := chi.URLParam(r, "walletId")
	walletId, err := uuid.Parse(walletIdParam)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	wallet, err := s.wallets.GetWallet(r.Context(), walletId)
	if err != nil {
		var rfnErr *wallets.WalletNotFoundError
		if errors.As(err, &rfnErr) {
			render.Render(w, r, ErrNotFound)
		}
		// } else {
		// 	render.Render(w, r, ErrInternalServerError)
		// }
		return
	}

	render.Render(w, r, NewWalletResponse(wallet))
}
