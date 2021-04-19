package wallet

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

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


var ErrAccountNotFound = errors.New("account not found")
var ErrAccountMustBePositive = errors.New("Amount must be greater than zero")
var ErrPhoneRegistered = errors.New("phone already registered")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrCannotRepeat = errors.New("can nor Repeat")
var ErrFavoriteNotFound = errors.New("can not find favorite")

type Service struct{
	nextAccountID int64
	accounts [] *types.Account
	payments [] *types.Payment
	favorites [] *types.Favorite

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


 func (s *Service) Repeat(paymentID string) (*types.Payment,error){

payment,err := s.FindPaymentByID(paymentID)
if err != nil{
	return nil, fmt.Errorf("can not find payment, %v",ErrPaymentNotFound)
}

pay,err:=s.Pay(payment.AccountID,payment.Amount,payment.Category)	
if err != nil{
	return nil, fmt.Errorf("can do new payment, %v",ErrPaymentNotFound)
}
 pay.ID=uuid.New().String()

return pay,nil

 }

 func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error){

	payment,err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil,ErrPaymentNotFound
	}
 
	favorite := &types.Favorite{
	  ID: uuid.New().String(),
	  AccountID: payment.AccountID,
	  Name: name,
	  Amount: payment.Amount,
	  Category: payment.Category,
	  
  }
  s.favorites=append(s.favorites, favorite)
return favorite,nil
 } 


func(s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error){


	findpay,err := s.FindFavoriteByID(favoriteID)

	if err != nil {
		return nil,ErrFavoriteNotFound
	}
 pay,err := s.Pay(findpay.AccountID,findpay.Amount,findpay.Category)

 if err != nil {
	 return nil,ErrPaymentNotFound
 }

 return pay,nil

}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {

	var favorite *types.Favorite

	for _, fav := range s.favorites{
		
		if (fav.ID==favoriteID){
			favorite=fav
			break

		}
	}
	if favorite == nil {
		return nil, ErrFavoriteNotFound
	}
	return favorite,nil
}

func (s *Service) ExportToFile(path string) error{

 str:=""

 file,err := os.Create(path)
	
 
	if err != nil{
		log.Print(err)
		return err
	}

	defer func ()  {
		
		err = file.Close()
		if err != nil {
			log.Print(err)
		}
	}()

for _, account := range s.accounts {

	

 str+= strconv.Itoa(int(account.ID))+";"
 str+=string(account.Phone)+";"
 str+=strconv.Itoa(int(account.Balance))+"|"

	 
}
_,err = file.Write([]byte (str))
if err != nil {
    log.Print(err)
	return err
}

return nil
}


func (s *Service) ImportFromFile(path string) error{

	file,err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}

content := make([] byte,0)
buf := make([]byte,4)

for{

	read,err := file.Read(buf)

	if err == io.EOF{
    
		content=append(content, buf[:read]...)
		break

	}

	if err != nil {
		log.Print(err)
		return err
	}

	content = append(content, buf[:read]...)

}
data := string(content)


  accounts:=strings.Split(data,"|")
  accounts = accounts[:len(accounts)-1]

  for _, account := range accounts {

	splits := strings.Split(account, ";")

	id,err := strconv.Atoi(splits[0])
	if err != nil {
		log.Print(err)
		return err
	}

	phone := splits[1]

	balance,err := strconv.Atoi(splits[2])
	if err != nil{
		log.Print(err)
		return err
	}  

	  s.accounts=append(s.accounts, &types.Account{
		  ID: int64(id),
		  Phone: types.Phone(phone),
		  Balance: types.Money(balance),
	  })
  }
 return nil
  
  }







