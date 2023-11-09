package handlers

import (
	"encoding/json"
	"errors"
	requestparams "go-sample/api/handlers/request-params"
	"go-sample/api/handlers/services"
	"go-sample/api/utils"
	accountrepo "go-sample/storage/account-repo"
	transactionrepo "go-sample/storage/transaction-repo"
	"go-sample/types"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type AccountHandler struct {
	accountRepo     accountrepo.IAccountRepo
	transactionRepo transactionrepo.ITransactionRepo
	accountSrv      services.Account
}

func NewAccountRepoHandler(
	accountRepo accountrepo.AccountRepo,
	transactionRepo transactionrepo.TransactionRepo,
	accountSrv services.Account,
) *AccountHandler {
	return &AccountHandler{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		accountSrv:      accountSrv,
	}
}

func (ah *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	var account types.Account
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	entity := types.NewAccount(account.Owner, account.Balance)
	isValidated := entity.ValidateAccount()

	if !isValidated {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accountExists, err := ah.accountSrv.IsThereAlreadyAccountWithThisOwner(r.Context(), account.Owner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if accountExists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "This owner already holds an account",
		})
		return
	}

	err = ah.accountSrv.CreateAccount(r.Context(), entity)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (ah *AccountHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "id")
	balance, err := ah.accountSrv.GetAccountBalance(r.Context(), accountId)

	if errors.Is(err, services.ErrInexistentAccount) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "There's no any bank account associated with this id",
		})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{
		"balance": balance,
	})
}

func (ah *AccountHandler) DepositMoney(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "id")
	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exists := ah.accountSrv.IsAccountExistent(r.Context(), accountId)
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Account doest not exist"})
		return
	}
	err = ah.accountSrv.DepositMoney(r.Context(), accountId, amount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ah *AccountHandler) WithdrawMoney(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "id")
	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exists := ah.accountSrv.IsAccountExistent(r.Context(), accountId)
	if !exists {
		json.NewEncoder(w).Encode(map[string]string{"error": "Account doest not exist"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ah.accountSrv.HasInsufficientFunds(r.Context(), accountId, amount) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Insufficient funds to transfer",
		})
		return
	}
	err = ah.accountSrv.WithdrawMoney(r.Context(), accountId, amount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ah *AccountHandler) TransferMoney(w http.ResponseWriter, r *http.Request) {
	var request requestparams.TransferMoneyRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if !request.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request params",
		})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ah.accountSrv.HasInsufficientFunds(r.Context(), request.From, request.Amount) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Insufficient funds to transfer",
		})
		return
	}

	err = ah.accountSrv.TransferMoney(r.Context(), request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ah *AccountHandler) GetTransactionsHistory(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "id")
	data, err := utils.ExtracteQueryParams(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filter := new(requestparams.GetTransactionsHistoryRequest)
	if err = json.Unmarshal(data, filter); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if filter.Page == "0" {
		filter.Page = "1"
	}
	transactions, err := ah.accountSrv.GetTransactionsHistory(r.Context(), accountId, *filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(transactions) == 0 {
		json.NewEncoder(w).Encode([]map[string]string{})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

func (ah *AccountHandler) RefundMoney(w http.ResponseWriter, r *http.Request) {

	var request requestparams.RefundMoneyRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if request.IsMultibenificiary {
		err := ah.accountSrv.RefundMultibeneficiaryTransfer(r.Context(), request.MultiBeneficiaryId)

		if err != nil {

			if err == services.ErrUnableToRefundARefund {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Unable to refund a refund transaction",
				})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	err = ah.accountSrv.RefundMoneyNormalTransfer(r.Context(), request.TransactionId)

	if err != nil {

		if err == services.ErrUnableToRefundARefund {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Unable to refund a refund transaction",
			})
			return
		}

		if err == services.ErrInsufficientFunds {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Insufficient funds to refund",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
