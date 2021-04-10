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