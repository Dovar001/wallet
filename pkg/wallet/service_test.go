package wallet

import (
	"fmt"
	"testing"

	"github.com/Dovar001/wallet/pkg/types"
	

	"reflect"
)

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