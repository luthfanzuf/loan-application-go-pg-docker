package main

import (
	"net/http"
	cust "vanilla-loan-application/handlers/cust"
)

func (app *Application) routes() *http.ServeMux {

	router := http.NewServeMux()

	// home route
	router.HandleFunc("/", cust.HomeHandler)

	//customer route
	router.HandleFunc("/user/account/create", cust.CreateCustomerAccount)  //v customer bisa membuat akun/profil baru
	router.HandleFunc("/user/account/listall", cust.GetAllCustomerAccount) //v employee bisa melihat semua list akun yang terdaftar
	router.HandleFunc("/user/account/update", cust.UpdateCustomerAccount)  //v customer bisa memperbarui akun/profil
	router.HandleFunc("/user/account/delete", cust.DeleteCustomerAccount)  //v customer bisa menghapus akun/profil

	router.HandleFunc("/user/loan/create", cust.CreateLoanApplication)  //v customer bisa membuat pengajuan pinjaman
	router.HandleFunc("/user/loan/listall", cust.GetAllLoanApplication) //v employee bisa melihat pengajuan yang pernah dibuat
	router.HandleFunc("/user/loan/update", cust.UpdateLoanApplication)  //v customer bisa mengedit pengjuan pinjaman yang pernah dibuat dan belum disetujui
	router.HandleFunc("/user/loan/delete", cust.DeleteLoanApplication)  //v customer bisa menghapus pengajuan pinjaman yang pernah dibuat dan belum disetujui

	return router
}
