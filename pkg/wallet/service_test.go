package wallet

import (
	"fmt"
	"testing"

	"github.com/Dovar001/wallet/pkg/types"
	

	"reflect"
)
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


func TestService_FindAccountByID_success(t *testing.T) {
service := &Service{ 
	
	accounts: []*types.Account{

	{
		ID: 11111111,
		Phone:"90-979-6600" ,
		Balance: 10_000 ,
	},
	{
		ID: 22222222,
		Phone:"92-890-4443" ,
		Balance: 20_000 ,
	},
	{
		ID: 33333333,
		Phone:"90-000-00-00" ,
		Balance: 30_000 ,
	},
	{
		ID: 44444444,
		Phone:"90-111-11-11" ,
		Balance: 40_000 ,
	},
	{
		ID: 55555555,
		Phone:"90-444-44-44" ,
		Balance: 50_000 ,
	},
	},
	
}

expect := types.Account{
		ID: 55555555,
		Phone:"90-444-44-44" ,
		Balance: 50_000 ,		
}

result,err := service.FindAccountByID(55555555)
if err == nil {
	fmt.Println(err)
	return
}


if !reflect.DeepEqual(expect,result){
	t.Errorf("invalid result, expected: %v, actual:%v", expect , result)
}


}



func TestService_FindAccountByID_notfound(t *testing.T) {



	service := &Service{ 
		
		accounts: []*types.Account{
	
		{
			ID: 11111111,
			Phone:"90-979-6600" ,
			Balance: 10_000 ,
		},
		{
			ID: 22222222,
			Phone:"92-890-4443" ,
			Balance: 20_000 ,
		},
		{
			ID: 33333333,
			Phone:"90-000-00-00" ,
			Balance: 30_000 ,
		},
		{
			ID: 44444444,
			Phone:"90-111-11-11" ,
			Balance: 40_000 ,
		},
		{
			ID: 55555555,
			Phone:"90-444-44-44" ,
			Balance: 50_000 ,
		},
		},
		
	}
	
	expect := ErrAccountNotFound
	
	result,err := service.FindAccountByID(555555)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	
	if !reflect.DeepEqual(expect,result){
		t.Errorf("invalid result, expected: %v, actual:%v", expect , result)
	}
	
	
	}

	func TestService_FindPaymentByID_success(t *testing.T) {
service := &Service{ 
	
	payments: []*types.Payment{

	{
		ID: "1111",
		Amount: 1000 ,
		
	},
	{
		ID: "2222",
		Amount: 2000 ,
	},
	{
		ID: "3333",
		Amount: 3000 ,
	},
	{
		ID: "4444",
		Amount: 4000 ,
	},
	{
		ID: "5555",
		Amount: 5000 ,
	},
	},
	
}

expect := types.Payment{
	ID: "1111",
	Amount: 1000 ,	
}

result,err := service.FindPaymentByID("1111")
if err == nil {
	fmt.Println(err)
	return
}


if !reflect.DeepEqual(expect,result){
	t.Errorf("invalid result, expected: %v, actual:%v", expect , result)
}


}



func TestService_FindPaymentByID_fail(t *testing.T) {
	service := &Service{ 
		
		payments: []*types.Payment{
	
		{
			ID: "1111",
			Amount: 1000 ,
			
		},
		{
			ID: "2222",
			Amount: 2000 ,
		},
		{
			ID: "3333",
			Amount: 3000 ,
		},
		{
			ID: "4444",
			Amount: 4000 ,
		},
		{
			ID: "5555",
			Amount: 5000 ,
		},
		},
		
	}
	
	expect := ErrPaymentNotFound
	
	result,err := service.FindPaymentByID("1112")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	
	if !reflect.DeepEqual(expect,result){
		t.Errorf("invalid result, expected: %v, actual:%v", expect , result)
	}
	
}
	

func TestService_Reject_success(t *testing.T) {

	//создаём сервис
	s := Service{}

	//Регистрируем там пользователья 
	phone := types.Phone("+992000000001")
	account, err := s.RegisterAccount(phone)
	
	if err != nil {
     t.Errorf("Reject(): can not register account, errror = %v", err)
  return
	}
	//пополняем его счёт 

	err = s.Deposit(account.ID, 10_000_00)
	
	if err != nil {
		t.Errorf("Reject(): can not deposit account, error = %v", err)
		return
	}
	//осуществляем платёж на его счёт

	payment, err := s.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		t.Errorf("Reject(): can not creat payment, error = %v", err)
		return
	}
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}
}	

func TestService_Reject_notfound(t *testing.T) {

	//создаём сервис
	s := Service{}

	//Регистрируем там пользователья 
	phone := types.Phone("+992000000001")
	account, err := s.RegisterAccount(phone)
	
	if err != nil {
     t.Errorf("Reject(): can not register account, errror = %v", err)
  return
	}
	//пополняем его счёт 

	err = s.Deposit(account.ID, 10_000_00)
	
	if err != nil {
		t.Errorf("Reject(): can not deposit account, error = %v", err)
		return
	}
	//осуществляем платёж на его счёт

	payment, err := s.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		t.Errorf("Reject(): can not creat payment, error = %v", err)
		return
	}
	err = s.Reject(payment.ID)
	if err == ErrAccountNotFound {
		t.Errorf("Reject(): error = %v", err)
		return
	}
}

func TestService_Repeat_success(t *testing.T){

	s:= newTestService()
		//Регистрируем там пользователья 
		phone := types.Phone("+992000000001")
		account, err := s.RegisterAccount(phone)
		
		if err != nil {
		 t.Errorf(" can not register account, errror = %v", err)
	  return
		}
		//пополняем его счёт 
	
		err = s.Deposit(account.ID, 10_000_00)
		
		if err != nil {
			t.Errorf(" can not deposit account, error = %v", err)
			return
		}
		//осуществляем платёж на его счёт
	
		payment, err := s.Pay(account.ID, 1000_00, "auto")
		if err != nil {
			t.Errorf(" can not creat payment, error = %v", err)
			return
		}
       
		got,err := s.FindPaymentByID(payment.ID)
		if err != nil {
			t.Errorf("can not find payment =%v",err)
		return
		}

		repeat,err := s.Repeat(payment.ID)
		if err != nil {
			t.Errorf("can not repeat payment, error = %v",err)
			return
		}
		
		if reflect.DeepEqual(got,repeat){
			t.Errorf("wrong repeat of payment = %v",err)
			return

		}
}

func TestService_Pay_succes(t *testing.T){
//создаём сервис
s := Service{}

//Регистрируем там пользователья 
phone := types.Phone("+992000000001")
account, err := s.RegisterAccount(phone)

if err != nil {
 t.Errorf(" can not register account, errror = %v", err)
return
}
//пополняем его счёт 

err = s.Deposit(account.ID, 10_000_00)

if err != nil {
	t.Errorf("can not deposit account, error = %v", err)
	return
}
//осуществляем платёж на его счёт

payment, err := s.Pay(account.ID, 1000_00, "auto")
if err != nil {
	t.Errorf(" can not creat payment, error = %v", err)
	return
}


if payment==nil{
	t.Errorf("wrong repeat of payment = %v",err)
	return

}
}



func TestService_Deposit_succes(t *testing.T){
	//создаём сервис
	s := Service{}
	
	//Регистрируем там пользователья 
	phone := types.Phone("+992000000001")
	account, err := s.RegisterAccount(phone)
	
	if err != nil {
	 t.Errorf(" can not register account, errror = %v", err)
	return
	}
	//пополняем его счёт 
	
	err = s.Deposit(account.ID, 10_000_00)
	
	if err != nil {
		t.Errorf("can not deposit account, error = %v", err)
		return
	}
	
	}


	func TestService_addAccount_success(t *testing.T){

		s:=newTestService()
    
		account,_,err:=s.addAccount(defaultTestAccount)
		if err != nil{
			t.Error(err)
		}

		expect:= defaultTestAccount

			
		if reflect.DeepEqual(account,expect){
			t.Errorf("account not right = %v",err)
			return

		}

	}

func TestFavoritePayment_success(t *testing.T){

//создаём сервис
s := newTestService()

//Регистрируем там пользователья 
phone := types.Phone("+992000000001")
account, err := s.RegisterAccount(phone)

if err != nil {
 t.Errorf(" can not register account, errror = %v", err)
return
}
//пополняем его счёт 

err = s.Deposit(account.ID, 10_000_00)

if err != nil {
	t.Errorf("can not deposit account, error = %v", err)
	return
}
//осуществляем платёж на его счёт

payment, err := s.Pay(account.ID, 1000_00, "auto")

if err != nil {
	t.Errorf(" can not creat payment, error = %v", err)
	return
}

pay,err:= s.FindPaymentByID(payment.ID)

if err != nil {
	t.Errorf("can not find payment =%v", err)
	return
}

fav,err := s.FavoritePayment(payment.ID, "Довар")

if err != nil {
	t.Errorf("favorite can not found= %v", err)
	return
}

if reflect.DeepEqual(fav,pay){
	
	t.Errorf("favorite equal to payment = %v",err)
			return

}

}

func TestPayFromFavorite_success(t *testing.T){
//создаём сервис
s := newTestService()

//Регистрируем там пользователья 
phone := types.Phone("+992000000001")
account, err := s.RegisterAccount(phone)

if err != nil {
 t.Errorf(" can not register account, errror = %v", err)
return
}
//пополняем его счёт 

err = s.Deposit(account.ID, 10_000_00)

if err != nil {
	t.Errorf("can not deposit account, error = %v", err)
	return
}
//осуществляем платёж на его счёт

payment, err := s.Pay(account.ID, 1000_00, "auto")

if err != nil {
	t.Errorf(" can not creat payment, error = %v", err)
	return
}

pay,err:= s.FindPaymentByID(payment.ID)

if err != nil {
	t.Errorf("can not find payment =%v", err)
	return
}

fav,err := s.FavoritePayment(pay.ID, "Довар")

if err != nil {
	t.Errorf("favorite can not found= %v", err)
	return
}
findfav,err:= s.FindFavoriteByID(fav.ID)

if err != nil {
	t.Errorf("can not find payment =%v", err)
	return
}

payfromfav,err := s.PayFromFavorite(findfav.ID) 

if err != nil {
	t.Errorf("can not pay from favorite = %v",err)
	return
}
if reflect.DeepEqual(payfromfav,findfav){

	t.Errorf("can not pay from favorite = %v",err)
	return
}

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


	



	
