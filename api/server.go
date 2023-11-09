package api

import (
	"go-sample/api/handlers"
	"go-sample/api/handlers/services"
	accountrepo "go-sample/storage/account-repo"
	transactionrepo "go-sample/storage/transaction-repo"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	listenAddr string
	db         *sqlx.DB
}

func NewServer(listenAddr string, db *sqlx.DB) *Server {
	return &Server{
		listenAddr: listenAddr,
		db:         db,
	}
}
func (s *Server) Start() error {
	accountRepo := accountrepo.NewAccountRepo(s.db)
	transactionRepo := transactionrepo.NewATransactionRepo(s.db)

	accountSrv := services.NewAccount(accountRepo, transactionRepo)

	accountHandler := handlers.NewAccountRepoHandler(accountSrv)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post("/accounts", accountHandler.CreateAccount)
	r.Get("/accounts/{id}/balance", accountHandler.GetAccountBalance)
	r.Get("/accounts/{id}/transactions", accountHandler.GetTransactionsHistory)
	r.Post("/accounts/{id}/deposit", accountHandler.DepositMoney)
	r.Post("/accounts/{id}/withdraw", accountHandler.WithdrawMoney)
	r.Post("/accounts/{id}/transfer_money", accountHandler.TransferMoney)
	r.Post("/accounts/{id}/refund_money", accountHandler.RefundMoney)
	r.Get("/accounts/{id}/transactions", accountHandler.GetTransactionsHistory)

	return http.ListenAndServe(":"+s.listenAddr, r)
}
