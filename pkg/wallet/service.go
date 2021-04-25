package wallet

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
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



func (s *Service) Export(dir string) error{
    

	if len(s.accounts)>0 {

		file,err := os.OpenFile(dir + "/accounts.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

if err != nil {
	log.Print(err)
	return err
}
defer func ()  {
		
	if cerr := file.Close(); cerr!= nil {
		if err == nil {
			cerr=err
		}
	}
}()
accstr:=""

	

	for _, account := range s.accounts {

		accstr+= strconv.Itoa(int(account.ID))+";"
		accstr+=string(account.Phone)+";"
		accstr+=strconv.Itoa(int(account.Balance))+ "\n"
		
	}

	file.WriteString(accstr)
	}
	
	if len(s.payments)>0 {
	fil,err := os.OpenFile(dir + "/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	if err != nil {
		log.Print(err)
		return err
	}
	defer func ()  {
			
		if cerr := fil.Close(); cerr!= nil {
			if err == nil {
				cerr=err
			}
		}
	}()

	paystr:=""

	

	for _, payment := range s.payments {

		 paystr+=string(payment.ID) + ";"

 		 paystr+=strconv.Itoa(int(payment.AccountID))+ ";"

		 paystr+= strconv.Itoa(int(payment.Amount))+ ";"

		 paystr+= string(payment.Category)+ ";"

         paystr+=string(payment.Status)+ "\n"
	}
	fil.WriteString(paystr)
}

if len (s.favorites) > 0 {

	files,err := os.OpenFile(dir + "/favorites.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	if err != nil {
		log.Print(err)
		return err
	}
	defer func ()  {
			
		if cerr := files.Close(); cerr!= nil {
			if err == nil {
				cerr=err
			}
		}
	}()

	favstr := ""

	

	
	for _, favorite := range s.favorites {

		favstr+= favorite.ID + ";"
		favstr+= strconv.Itoa(int(favorite.AccountID))+";"
		favstr+= favorite.Name + ";"
		favstr+= strconv.Itoa(int(favorite.Amount)) + ";"
		favstr+= string(favorite.Category) + "\n"

		
	}
	files.WriteString(favstr)
}

	

	return nil
}

//Import ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (s *Service)  Import(dir string) error{

	_,err:=os.Stat(dir + "/accounts.dump")
	if err == nil {
     
		file,err := os.ReadFile(dir + "/accounts.dump")

		if err != nil{
			log.Print(err)
			return err
		}

		accstr := string(file)

		accounts:=strings.Split(accstr,"\n")

		if len(accounts) > 0{
  
			accounts = accounts[:len(accounts)-1]
		 }
  
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
	}

// Payments==========================================

_,err= os.Stat(dir + "/payments.dump")

if err == nil {

	file,err := os.ReadFile(dir + "/payments.dump")

	if err != nil {
	log.Print(err)
	return err
}
paystr := string(file)

payments:= strings.Split(paystr, "\n")

if len(payments) > 0{
	payments= payments[:len(payments)-1]
} 
for _, payment := range payments {

	splits:= strings.Split(payment, ";")
	id := splits[0]

	accountid,err := strconv.Atoi(splits[1]) 
	if err != nil {
		log.Print(err)
		return err
	}
  
	amount,err := strconv.Atoi(splits[2])
	if err != nil{
		log.Print(err)
		return err
	}  
  
	category := splits[3]
  
	status := splits[4]
  
	  s.payments=append(s.payments, &types.Payment{
		  ID: id,
		  AccountID:int64(accountid) ,
		  Amount: types.Money(amount),
		  Category: types.PaymentCategory(category),
		  Status: types.PaymentStatus(status),
	  })
	
}

}
//Favorites =======================================================
_,err = os.Stat(dir + "/favorites.dump") 

if err == nil {
	
	file,err := os.ReadFile(dir + "/favorites.dump")
	if err != nil {
		log.Print(err)
		return err
	}

	favstr :=string(file)
	
	favorites := strings.Split(favstr, "\n")
	
	if len(favorites) > 0 {
		favorites = favorites[:len(favorites)-1]
	}
	for _, favorite := range favorites {

		splits:= strings.Split(favorite, ";")
	 
		id := splits[0]
		accountid,err := strconv.Atoi(splits[1])
		if err != nil {
			log.Print(err)
			return err
		}
		name := splits[2]
		amount,err := strconv.Atoi(splits[3])
		if err != nil {
		 log.Print(err)
		 return err
		}
		category:= types.PaymentCategory(splits[4])
	
	
	
		s.favorites=append(s.favorites, &types.Favorite{
			ID: id,
			AccountID:int64(accountid) ,
			Name: name,
			Amount: types.Money(amount),
			Category: types.PaymentCategory(category),
			
		})	  
	  }				

}


return nil
}


func (s *Service) ExportAccountHistory(accountID int64) ([]*types.Payment, error){

	var account  *types.Account
	var paymentsss [] *types.Payment

	for _, acc := range s.accounts {
		
		if ( accountID == acc.ID){
            account= acc
			break
		}
		
	}
	if account == nil {
		return nil , ErrAccountNotFound
	}

	for _, payment := range s.payments {
  
  if (payment.AccountID == account.ID){
 
	paymentsss=append(paymentsss, payment)
  } 
		
	}

	if paymentsss == nil {
		return nil, ErrPaymentNotFound
	}

	return paymentsss,nil

} 

 func (s *Service) HistoryToFiles(payments []*types.Payment, dir string, records int) error{

	if (len(payments)>0 && len(payments)<=records) {
		file,err := os.OpenFile(dir + "/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	
		if err != nil {
			log.Print(err)
			return err
		}
		defer func ()  {
				
			if cerr := file.Close(); cerr!= nil {
				if err == nil {
					cerr=err
				}
			}
		}()
		
		paystr:=""

		for _, payment := range payments {
	
			 paystr+=string(payment.ID) + ";"
	
			  paystr+=strconv.Itoa(int(payment.AccountID))+ ";"
	
			 paystr+= strconv.Itoa(int(payment.Amount))+ ";"
	
			 paystr+= string(payment.Category)+ ";"
	
			 paystr+=string(payment.Status)+ "\n"
		}
		file.WriteString(paystr)
	} else if records >0{
		k:=0
		t:=1
 
		s:= int(math.Ceil(float64(len(payments)/records)))
		
         for i:=0; i<=s; i++ {
			 file,err:= os.OpenFile(dir+"/payments" + fmt.Sprint(t) + ".dump",os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			 if err!= nil {
				 log.Print(err)
				 return err
			 }
			 defer func ()  {
				
				if cerr := file.Close(); cerr!= nil {
					if err == nil {
						cerr=err
					}
				}
			}()
			paystr:= ""

			for i := k; i <records+k ; i++{
   
          
				paystr+=string(payments[i].ID) + ";"
	   
				paystr+=strconv.Itoa(int(payments[i].AccountID))+ ";"
	  
			   paystr+= strconv.Itoa(int(payments[i].Amount))+ ";"
	  
			   paystr+= string(payments[i].Category)+ ";"
	  
			   paystr+=string(payments[i].Status)+ "\n"   
			   
			   
		   }
		   file.WriteString(paystr)	
		   
		   paystr=""
		   t++
		   k+=records
		 }

		}
		return nil
	}




		////////////////////////////////////////////////////////////////////////////////////////////////////

		/* file,err := os.OpenFile(dir + "/payments1.dump",os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		 if err != nil {
			 log.Print(err)
			 return err
		 }
		 defer func ()  {
				
			if cerr := file.Close(); cerr!= nil {
				if err == nil {
					cerr=err
				}
			}
		}()
		 
     paystr:= ""
		
	 if records < len(payments){
	 for i := 0; i <records ; i++{
   
          
			 paystr+=string(payments[i].ID) + ";"
	
			 paystr+=strconv.Itoa(int(payments[i].AccountID))+ ";"
   
			paystr+= strconv.Itoa(int(payments[i].Amount))+ ";"
   
			paystr+= string(payments[i].Category)+ ";"
   
			paystr+=string(payments[i].Status)+ "\n"    
		}
		file.WriteString(paystr)	
		
		paystr=""

		fil,err := os.OpenFile(dir + "/payments2.dump",os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Print(err)
			return err
		}
		defer func ()  {
			   
		   if cerr := fil.Close(); cerr!= nil {
			   if err == nil {
				   cerr=err
			   }
		   }
	   }()

	   for i := records; i < len(payments); i++{
		   
			 paystr+=string(payments[i].ID) + ";"
	
			 paystr+=strconv.Itoa(int(payments[i].AccountID))+ ";"
   
			paystr+= strconv.Itoa(int(payments[i].Amount))+ ";"
   
			paystr+= string(payments[i].Category)+ ";"
   
			paystr+=string(payments[i].Status)+ "\n"    
		}
		file.WriteString(paystr)	

	   }
	}
return nil
}
*/
