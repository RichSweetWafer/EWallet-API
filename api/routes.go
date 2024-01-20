package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (s *Server) routes() {
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	s.router.Get("/", s.handleStatusCheck)

	s.router.Route("/api/v1/wallet", func(r chi.Router) {
		r.Post("/", s.handleCreateWallet)
		r.Post("/{walletId}/send", s.handleTransaction)

		r.Get("/{walletId}/history", s.handleWalletHistory)
		r.Get("/{walletId}", s.handleWalletState)

	})
}
