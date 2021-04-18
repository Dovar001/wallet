package main

import (
	

	
	
	"github.com/Dovar001/wallet/pkg/wallet"
)


func main() {

	
 svc:= &wallet.Service{}

 svc.RegisterAccount("+992909796600")
 svc.RegisterAccount("+992927592402")
 svc.RegisterAccount("+992000000000")
 svc.RegisterAccount("+992111111111")
 svc.RegisterAccount("+992222222222")
 svc.ExportToFile("../data/export.txt")

 }
	

	
