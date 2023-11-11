package handlers

import (
	"context"
	"fmt"
	requestparams "go-sample/api/handlers/request-params"
	"go-sample/api/handlers/services"
	accountrepo "go-sample/storage/account-repo"
	transactionrepo "go-sample/storage/transaction-repo"
	"go-sample/types"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

func createTables(db *sqlx.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.transaction_ (
		id VARCHAR(36) PRIMARY KEY,
		from_account VARCHAR(36) NOT NULL,
		to_account VARCHAR(36) NOT NULL,
		tr_status TEXT NOT NULL,
		operation VARCHAR(10) NOT NULL,
		amount MONEY NOT NULL,
		multiBeneficiaryId VARCHAR(36) NOT NULL,
		is_Refund BOOLEAN NOT NULL, 
		refunded_transaction_id VARCHAR(36) NOT NULL,
		createdAt TIMESTAMP NOT NULL
	  );`)

	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
	  CREATE TABLE IF NOT EXISTS public.accounts (
		  id VARCHAR(36) PRIMARY KEY,
		  owner_id VARCHAR(36) NOT NULL,
		  balance MONEY NOT NULL,
		  created_at TIMESTAMP NOT NULL,
		  updated_at TIMESTAMP NOT NULL,
		  deletedAt TIMESTAMP NULL
		);`)
	if err != nil {
		panic(err)
	}
}

func dropTables(db *sqlx.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS public.transaction_;`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS public.accounts;`)
	if err != nil {
		panic(err)
	}
}
func TestAccountHandler(t *testing.T) {

	TEST_DATABASE_NAME := "nell_challenge_test"
	TEST_DATABASE_URL := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", TEST_DATABASE_NAME)
	db, _ := sqlx.Connect("postgres", TEST_DATABASE_URL)

	createTables(db)

	transactionRepo := transactionrepo.NewATransactionRepo(db)
	accountRepo := accountrepo.NewAccountRepo(db)

	accountSrv := services.NewAccount(accountRepo, transactionRepo)

	accountHandler := &AccountHandler{
		accountSrv: accountSrv,
	}

	t.Run("POST /accounts", func(t *testing.T) {

		t.Run("POST /accounts should respond with 201 on success", func(t *testing.T) {
			expected := 201

			owner_id := "1e55e90d-427e-4f0b-b71f-5e38e57c2ae8"
			jsonPayload := fmt.Sprintf(`
				{
					"owner_id": "%v",
					"balance": 100.0
				}
			`, owner_id)

			body := strings.NewReader(jsonPayload)
			req := httptest.NewRequest(http.MethodPost, "/accounts", body)
			w := httptest.NewRecorder()

			accountHandler.CreateAccount(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}
			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts WHERE owner_id=$1", owner_id)
			})
		})

		t.Run("POST /accounts should respond with 400 if owner_id absent in body payload ", func(t *testing.T) {
			expected := 400

			jsonPayload := `
				{
					"balance": 100.0
				}
			`

			body := strings.NewReader(jsonPayload)
			req := httptest.NewRequest(http.MethodPost, "/accounts", body)
			w := httptest.NewRecorder()

			accountHandler.CreateAccount(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}
			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts WHERE owner_id=$1", "1e55e90d-427e-4f0b-b71f-5e38e57c2ae8")
			})
		})

		t.Run("POST /accounts should respond with 400 if there's already an account the specified owner_id", func(t *testing.T) {
			expected := 400

			ctx := context.Background()
			owner_id := "1e55e90d-427e-4f0b-b71f-5e38e57c2ae8"
			jsonPayload := fmt.Sprintf(`
				{
					"owner_id": "%v",
					"balance": 100.0
				}
			`, owner_id)

			account := types.NewAccount(owner_id, 0)
			accountSrv.CreateAccount(ctx, account)

			body := strings.NewReader(jsonPayload)
			req := httptest.NewRequest(http.MethodPost, "/accounts", body)
			w := httptest.NewRecorder()

			accountHandler.CreateAccount(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}
			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts WHERE owner_id=$1", owner_id)
			})
		})
	})

	t.Run("POST /accounts/{id}/deposit", func(t *testing.T) {

		t.Run("POST /accounts/{id}/deposit should respond with 200 on success", func(t *testing.T) {

			expected := 200

			ctx := context.Background()
			owner_id := "1e55e903-427e-4f0b-b71f-5e38e57c2ae8"

			account := types.NewAccount(owner_id, 0)
			accountSrv.CreateAccount(ctx, account)

			amount := 100
			endpoint := fmt.Sprintf("/accounts/%v/deposit?amount=%v", account.ID, amount)

			req := httptest.NewRequest(http.MethodPost, endpoint, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", account.ID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.DepositMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}
			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts WHERE owner_id=$1", owner_id)
			})
		})

		t.Run("POST /accounts/{id}/deposit should respond with 400 if amount negative", func(t *testing.T) {

			expected := 400

			ctx := context.Background()
			owner_id := "1e55e903-427e-4f0b-b71f-5e38e57c2ae8"

			account := types.NewAccount(owner_id, 0)
			accountSrv.CreateAccount(ctx, account)

			amount := -100
			endpoint := fmt.Sprintf("/accounts/%v/deposit?amount=%v", account.ID, amount)

			req := httptest.NewRequest(http.MethodPost, endpoint, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", account.ID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.DepositMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}
			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts WHERE owner_id=$1", owner_id)
			})
		})

		t.Run("POST /accounts/{id}/deposit should respond with 400 if account doest not exist", func(t *testing.T) {

			expected := 400

			inexistentAccountId := "1e55e903-427e-4f0b-b71f-5e38e57c2ae8"
			amount := -100
			endpoint := fmt.Sprintf("/accounts/%v/deposit?amount=%v", inexistentAccountId, amount)

			req := httptest.NewRequest(http.MethodPost, endpoint, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", inexistentAccountId)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.DepositMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}
			t.Cleanup(func() {
				//db.Exec("DELETE FROM accounts WHERE owner_id=$1", owner_id)
			})
		})

	})

	t.Run("POST /accounts/{id}/withdraw", func(t *testing.T) {
		t.Run("POST /accounts/{id}/withdraw should respond with 200 on success", func(t *testing.T) {

			expected := 200

			ctx := context.Background()
			owner_id := "1e55e903-427e-4f0b-b71f-5e38e57c2ae8"

			account := types.NewAccount(owner_id, 200)
			accountSrv.CreateAccount(ctx, account)

			amount := 100
			endpoint := fmt.Sprintf("/accounts/%v/withdraw?amount=%v", account.ID, amount)

			req := httptest.NewRequest(http.MethodPost, endpoint, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", account.ID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.WithdrawMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}
			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts")
			})
		})

		t.Run("POST /accounts/{id}/withdraw should respond with 400 if amount negative", func(t *testing.T) {

			expected := 400

			ctx := context.Background()
			owner_id := "1e55e903-427e-4f0b-b71f-5e38e57c2ae8"

			account := types.NewAccount(owner_id, 0)
			accountSrv.CreateAccount(ctx, account)

			amount := -100
			endpoint := fmt.Sprintf("/accounts/%v/withdraw?amount=%v", account.ID, amount)

			req := httptest.NewRequest(http.MethodPost, endpoint, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", account.ID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.WithdrawMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}
			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts")
			})
		})

		t.Run("POST /accounts/{id}/withdraw should respond with 400 if account doest not exist", func(t *testing.T) {

			expected := 400

			inexistentAccountId := "1e55e903-427e-4f0b-b71f-5e38e57c2ae8"
			amount := -100
			endpoint := fmt.Sprintf("/accounts/%v/withdraw?amount=%v", inexistentAccountId, amount)

			req := httptest.NewRequest(http.MethodPost, endpoint, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", inexistentAccountId)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.WithdrawMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}
			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts")
			})
		})

	})

	t.Run("POST /accounts/{id}/transfer_money", func(t *testing.T) {
		t.Run("POST /accounts/{id}/transfer_money should respond with 200 on success on normal transfer", func(t *testing.T) {

			expected := 200
			amount := 50

			ctx := context.Background()

			accounts := []*types.Account{
				{
					Owner:   "1e55e903-427e-4f0b-b71f-5e38e57c2ae8",
					Balance: 100,
				},
				{
					Owner:   "1e55e903-427e-4f0b-b74f-5e38e57c2ae8",
					Balance: 0,
				},
			}

			for i, account := range accounts {
				entity := types.NewAccount(account.Owner, account.Balance)
				accountSrv.CreateAccount(ctx, entity)
				accounts[i].ID = entity.ID
			}

			endpoint := fmt.Sprintf("/accounts/%v/transfer_money", accounts[0].ID)

			jsonPayload := fmt.Sprintf(`
				{	
					"from": "%v",
					"amount": %v,
					"subject": "uma transferencia",
					"recipients": [
						{
							"accountId": "%v",
							"amount": %v
						}
					]
				}
			`, accounts[0].ID, amount, accounts[1].ID, amount)

			body := strings.NewReader(jsonPayload)
			req := httptest.NewRequest(http.MethodPost, endpoint, body)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", accounts[0].ID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.TransferMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}

			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts")
				db.Exec("DELETE FROM transaction_")
			})
		})

		t.Run("POST /accounts/{id}/transfer_money should respond with 400 if the sum recepients amount not equal the transfer amount", func(t *testing.T) {

			expected := 400
			from := "1e55e903-427e-4f0b-b71f-5e38e57c2ae8"
			to := "1e55e903-427e-4f0b-b74f-5e38e57c2ae2"
			amount := 50

			endpoint := fmt.Sprintf("/accounts/%v/transfer_money", from)

			jsonPayload := fmt.Sprintf(`
				{	
					"from": "%v",
					"amount": %v,
					"subject": "uma transferencia",
					"recipients": [
						{
							"accountId": "%v",
							"amount": %v
						},
						{
							"accountId": "%v",
							"amount": %v
						}
					]
				}
			`, from, amount, to, amount, to, amount)

			body := strings.NewReader(jsonPayload)
			req := httptest.NewRequest(http.MethodPost, endpoint, body)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", from)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.TransferMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}

			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts")
				db.Exec("DELETE FROM transaction_")
			})
		})

		t.Run("POST /accounts/{id}/transfer_money should respond with 400 if the amount less or equal 0", func(t *testing.T) {

			expected := 400
			from := "1e55e903-427e-4f0b-b71f-5e38e57c2ae8"
			to := "1e55e903-427e-4f0b-b74f-5e38e57c2ae2"
			amount := 50

			endpoint := fmt.Sprintf("/accounts/%v/transfer_money", from)

			jsonPayload := fmt.Sprintf(`
				{	
					"from": "%v",
					"amount": %v,
					"subject": "uma transferencia",
					"recipients": [
						{
							"accountId": "%v",
							"amount": %v
						},
						{
							"accountId": "%v",
							"amount": %v
						}
					]
				}
			`, from, -1, to, amount, to, amount)

			body := strings.NewReader(jsonPayload)
			req := httptest.NewRequest(http.MethodPost, endpoint, body)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", from)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.TransferMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}

			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts")
				db.Exec("DELETE FROM transaction_")
			})
		})

		t.Run("POST /accounts/{id}/transfer_money should respond with 400 if user transfering does not exists", func(t *testing.T) {
			expected := 400
			from := "1e55e903-427e-4f0b-b71f-5e38e57c2ae8"
			to := "1e55e903-427e-4f0b-b74f-5e38e57c2ae2"
			fromAmount := 100
			amount := 50

			endpoint := fmt.Sprintf("/accounts/%v/transfer_money", from)

			jsonPayload := fmt.Sprintf(`
				{	
					"from": "%v",
					"amount": %v,
					"subject": "uma transferencia",
					"recipients": [
						{
							"accountId": "%v",
							"amount": %v
						},
						{
							"accountId": "%v",
							"amount": %v
						}
					]
				}
			`, from, fromAmount, to, amount, to, amount)

			body := strings.NewReader(jsonPayload)
			req := httptest.NewRequest(http.MethodPost, endpoint, body)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", from)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.TransferMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}

			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts")
				db.Exec("DELETE FROM transaction_")
			})
		})

		t.Run("POST /accounts/{id}/transfer_money should respond with 400 if user transfering hasnt sufficient funds", func(t *testing.T) {
			expected := 400
			amount := 150

			ctx := context.Background()

			accounts := []*types.Account{
				{
					Owner:   "1e55e903-427e-4f0b-b71f-5e38e57c2ae8",
					Balance: 100,
				},
				{
					Owner:   "1e55e903-427e-4f0b-b74f-5e38e57c2ae8",
					Balance: 0,
				},
			}

			for i, account := range accounts {
				entity := types.NewAccount(account.Owner, account.Balance)
				accountSrv.CreateAccount(ctx, entity)
				accounts[i].ID = entity.ID
			}

			endpoint := fmt.Sprintf("/accounts/%v/transfer_money", accounts[0].ID)

			jsonPayload := fmt.Sprintf(`
				{	
					"from": "%v",
					"amount": %v,
					"subject": "uma transferencia",
					"recipients": [
						{
							"accountId": "%v",
							"amount": %v
						}
					]
				}
			`, accounts[0].ID, amount, accounts[1].ID, amount)

			body := strings.NewReader(jsonPayload)
			req := httptest.NewRequest(http.MethodPost, endpoint, body)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", accounts[0].ID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.TransferMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}

			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts")
				db.Exec("DELETE FROM transaction_")
			})
		})

		t.Run("POST /accounts/{id}/transfer_money should respond with 400 if at least one of recipients doest not exists", func(t *testing.T) {

			expected := 400
			amount := 50

			ctx := context.Background()

			accounts := []*types.Account{
				{
					Owner:   "1e55e903-427e-4f0b-b71f-5e38e57c2ae8",
					Balance: 100,
				},
				{
					Owner:   "1e55e903-427e-4f0b-b74f-5e38e57c2ae8",
					Balance: 0,
				},
			}

			for i, account := range accounts {
				entity := types.NewAccount(account.Owner, account.Balance)
				accountSrv.CreateAccount(ctx, entity)
				accounts[i].ID = entity.ID
			}

			endpoint := fmt.Sprintf("/accounts/%v/transfer_money", accounts[0].ID)
			inexistentAccountId := "1e45e903-427e-d4f0b-b71f-5e38e57c2ae8"
			jsonPayload := fmt.Sprintf(`
				{	
					"from": "%v",
					"amount": %v,
					"subject": "uma transferencia",
					"recipients": [
						{
							"accountId": "%v",
							"amount": %v
						},
						{
							"accountId": "%v",
                            "amount": %v
						}
					]
				}
			`, accounts[0].ID, amount, accounts[1].ID, 25, inexistentAccountId, 25)

			body := strings.NewReader(jsonPayload)
			req := httptest.NewRequest(http.MethodPost, endpoint, body)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", accounts[0].ID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.TransferMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}

			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts")
				db.Exec("DELETE FROM transaction_")
			})
		})

	})

	t.Run("POST /accounts/{id}/refund_money", func(t *testing.T) {

		t.Run("POST /accounts/{id}/refund_money should respond 400 if multibeneficiary_id is empty", func(t *testing.T) {

			expected := 400
			ctx := context.Background()

			accounts := []*types.Account{
				{
					Owner:   "1e55e903-427e-4f0b-b71f-5e38e57c2ae8",
					Balance: 100,
				},
				{
					Owner:   "1e55e903-427e-4f0b-b74f-5e38e57c2ae8",
					Balance: 0,
				},
			}

			for i, account := range accounts {
				entity := types.NewAccount(account.Owner, account.Balance)
				accountSrv.CreateAccount(ctx, entity)
				accounts[i].ID = entity.ID
			}

			amount := 100.0
			transferParams := requestparams.TransferMoneyRequest{
				From:    accounts[0].ID,
				Amount:  amount,
				Subject: "uma transferencia",
				Repcipients: []requestparams.Recipient{
					{
						AccountId: accounts[1].ID,
						Amount:    amount,
					},
				},
			}
			accountSrv.TransferMoney(ctx, transferParams)

			jsonPayload := `
			{
				"is_multi_benificiary": true,
				"multi_beneficiary_id": "",
				"transaction_id": ""
			}`

			endpoint := fmt.Sprintf("/accounts/%v/refund_money", accounts[1].ID)

			body := strings.NewReader(jsonPayload)
			req := httptest.NewRequest(http.MethodPost, endpoint, body)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", accounts[0].ID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			accountHandler.RefundMoney(w, req)
			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if res.StatusCode != expected {
				t.Errorf("Expected Hello john but got %v", string(data))
			}

			t.Cleanup(func() {
				db.Exec("DELETE FROM accounts")
				db.Exec("DELETE FROM transaction_")
			})
		})

	})
	t.Cleanup(func() {
		dropTables(db)
	})
}
