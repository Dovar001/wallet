package main

import (
	"log"

	"fmt"

	"github.com/Dovar001/wallet/pkg/wallet"
)


func main() {

	
 svc:= &wallet.Service{}


 svc.RegisterAccount("+992909796600")
 svc.RegisterAccount("+992927592402")
svc.RegisterAccount("+992000000000")
svc.RegisterAccount("+9920000001200")
svc.RegisterAccount("+9920000001220")

 //svc.ExportToFile("../data/export.txt")
 //svc.ImportFromFile(".../data/export.txt")
 svc.Deposit(1,10_000)
 svc.Deposit(2, 234_000)
 svc.Deposit(3, 234_000)

pay,err:= svc.Pay(1, 500, "Машина")


if err != nil {
	log.Print(err)
	return
}


 svc.FavoritePayment(pay.ID, "Довар")
 
//svc.Export("../data")
 fmt.Print(svc.Import("../data"))

 }
	

	
