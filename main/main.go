package main

import (
	"log"
	"net/http"

	"github.com/thedevsaddam/renderer"
)

import "time"
import "strconv"

type PaymentType string

const (
	Credit PaymentType = "Credit"
	Debit              = "Debit"
)

type Account struct {
	Transactions []*Transaction
	UserId       int
	Balance      int
}

type Transaction struct {
	Id         int
	Amount     int
	PaymenType PaymentType
	Date       time.Time
}

var rnd *renderer.Render
var account *Account

func init() {
	opts := renderer.Options{
		ParseGlobPattern: "../front/*.html",
	}

	rnd = renderer.New(opts)

	account = NewAccount(123, 10000)
}

func handlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/history", about)
	mux.HandleFunc("/debit", debit)
	mux.HandleFunc("/credit", credit)
	mux.HandleFunc("/pay", pay)
	mux.HandleFunc("/successful", successful)
	mux.HandleFunc("/transactions", transactions)
	return mux
}

func main() {
	handlers := handlers()
	port := ":9000"
	log.Println("Listening on port ", port)
	http.ListenAndServe(port, handlers)
}

func home(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "home", nil)
}

func about(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "about", nil)
}

func debit(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "debit", nil)
}

func credit(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "credit", nil)
}

func pay(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	paymentType := r.Form.Get("type")
	amount, _ := strconv.Atoi(r.Form.Get("amount"))

	account.AddPayment(paymentType, amount)

	http.Redirect(w, r, "/successful", 301)
}

func successful(w http.ResponseWriter, r *http.Request) {
	account := map[string]interface{}{"Account": account}
	rnd.HTML(w, http.StatusOK, "successful", account)
}

func transactions(w http.ResponseWriter, r *http.Request) {
	account := map[string]interface{}{"Account": account}
	rnd.HTML(w, http.StatusOK, "transactions", account)
}

/* Could not make imports work */

func NewAccount(userId int, balance int) *Account {
	return &Account{
		Transactions: make([]*Transaction, 0),
		UserId:       userId,
		Balance:      balance,
	}
}

func NewTransaction(id int, amount int, paymenType PaymentType, date time.Time) *Transaction {
	return &Transaction{
		Id:         id,
		Amount:     amount,
		PaymenType: paymenType,
		Date:       date,
	}
}

func (a *Account) AddPayment(paymentType string, amount int) error {
	newID := a.CreateID()
	now := time.Now()
	paymentTypeEnum := getPaymentTypeEnum(paymentType)

	NewPayment := NewTransaction(newID, amount, paymentTypeEnum, now)

	a.Transactions = append(a.Transactions, NewPayment)

	a.GetNewBalance(amount, paymentTypeEnum)

	return nil
}

func (a *Account) CreateID() int {
	transactionSize := len(a.Transactions)
	if transactionSize == 0 {
		return 1
	}
	return transactionSize + 1
}

func (a *Account) GetNewBalance(amount int, paymentType PaymentType) {
	switch paymentType {
	case Credit:
		a.Balance = a.Balance + amount
		return
	default:
		a.Balance = a.Balance - amount
	}
}

func getPaymentTypeEnum(paymentType string) PaymentType {
	switch paymentType {
	case "credit":
		return Credit
	default:
		return Debit
	}
}
