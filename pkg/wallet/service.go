package wallet

import (
	"errors"
	"fmt"

	"github.com/Dovar001/wallet/pkg/types"
	"github.com/google/uuid"
)

func New(text string ) error{
	return &errorString{text}
}

type errorString struct{
	s string
}
func (e *errorString) Error() string {
	return e.s
}
var defaultTestAccount = testAccount{
	phone: "+992000000001",
	balance: 10_000_00,
	payments: []struct {
    amount types.Money
	category types.PaymentCategory
}{
	{amount: 1_000_00, category: "auto"},
	},
}

var ErrAccountNotFound = errors.New("account not found")
var ErrAccountMustBePositive = errors.New("Amount must be greater than zero")
var ErrPhoneRegistered = errors.New("phone already registered")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrCannotRepeat = errors.New("can nor Repeat")

type Service struct{
	nextAccountID int64
	accounts [] *types.Account
	payments [] *types.Payment
}

type testService struct{
	*Service
}
type testAccount struct{

	phone types.Phone
	balance types.Money
	payments []struct{
		amount types.Money
		category types.PaymentCategory
	}
}




func newTestService() *testService{
	return &testService{Service: &Service{}}
}


func (s *Service) RegisterAccount(phone types.Phone) (*types.Account,error){

	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil,ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID: s.nextAccountID,
		Phone: phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account,nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory)(*types.Payment,error){

	if amount <= 0 {
		return nil, ErrAccountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID{
			account =acc
			break
		}
		
	}
	if account == nil{
		return nil, ErrAccountNotFound
	}
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}
	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID: paymentID,
		AccountID: accountID,
		Amount: amount,
		Category: category,
		Status: types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}
func (s *Service) Deposit(accountID int64 , amount types.Money)error{
	if amount <= 0 {
		return ErrAccountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		
		if acc.ID == accountID{
			account = acc
			break
		}

	}
	if account == nil{
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}



func  (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {

var account *types.Account

for _, acc := range s.accounts {

	if (acc.ID==accountID){
		account = acc
		break
	}
	
}
if account == nil {
	return  nil, ErrAccountNotFound
}

return account,nil

}


func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {

	var payment *types.Payment

	for _, pay := range s.payments{
		
		if (pay.ID==paymentID){
			payment = pay
			break

		}
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment,nil
}
	
	
 func (s *Service) Reject(paymentID string) error{
	 payment,err := s.FindPaymentByID(paymentID)
	
	 if payment.ID == paymentID{
	
	 payment.Status = types.PaymentStatusFail 
	 }else if (payment.ID==err.Error()) {
		 return ErrPaymentNotFound
	 }
	 
	 account,err := s.FindAccountByID(payment.AccountID)

	 if payment.AccountID==account.ID{
		 account.Balance+=payment.Amount
	 }
	 return nil
 }

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error){

	//Регистрируем там пользователья
	account, err := s.RegisterAccount(data.phone)

	if err!= nil {
		return nil, nil, fmt.Errorf("can not register account, error = %v", err)

	}
	//Пополняем его счёт 
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can not deposit account, error = %v", err)

	}
	//Выполняем платежи
	//можем создать слайс сразу нужной длины, поскольку знаем размер
	payments:= make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments{
		//тогда здесь работаем через index, а не через append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err!= nil {
			return nil, nil, fmt.Errorf("can not make payment, error = %v", err)

		}
	}
	return account, payments, nil
}






 
 func (s *Service) Repeat(paymentID string) (*types.Payment,error){

	

payment,err := s.FindPaymentByID(paymentID)
if err != nil{
	return nil, fmt.Errorf("can not find payment, %v",ErrPaymentNotFound)
}

repayment := &types.Payment {
  ID : uuid.New().String(),
  AccountID: payment.AccountID,
  Amount: payment.Amount,
  Category: payment.Category,
  Status: payment.Status,
}
return repayment,nil

 }






 