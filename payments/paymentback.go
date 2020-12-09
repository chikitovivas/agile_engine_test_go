package payments

import "time"

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

func NewAccount(userId int, balance int) *Account {
	return &Account{
		Transactions: make([]*Transaction, 0),
		UserId:       userId,
		Balance:      balance,
	}
}

type Transaction struct {
	Id         int
	Amount     int
	PaymenType PaymentType
	Date       time.Time
}

func NewTransaction(id int, amount int, paymenType PaymentType, date time.Time) *Transaction {
	return &Transaction{
		Id:         id,
		Amount:     amount,
		PaymenType: paymenType,
		Date:       date,
	}
}

func (a *Account) addPayment(paymentType string, amount int) error {
	newID := a.createID()
	now := time.Now()
	paymentTypeEnum := getPaymentTypeEnum(paymentType)

	NewPayment := NewTransaction(newID, amount, paymentTypeEnum, now)

	a.Transactions = append(a.Transactions, NewPayment)

	return nil
}

func (a *Account) createID() int {
	transactionSize := len(a.Transactions)
	if transactionSize == 0 {
		return 1
	}
	return transactionSize + 1
}

func getPaymentTypeEnum(paymentType string) PaymentType {
	switch paymentType {
	case "Credit":
		return Credit
	default:
		return Debit
	}
}
