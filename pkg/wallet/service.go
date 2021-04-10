package wallet

import (
	"errors"

	"github.com/Dovar001/wallet/pkg/types"
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

type Service struct{
	nextAccountID int64
	accounts [] *types.Account
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

