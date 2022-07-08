package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	z "vanilla-loan-application/ctx"

	"gorm.io/gorm"
)

type CustomerAccountProfile struct {
	Cust_id               int       `json:"cust_id" gorm:"primaryKey;autoIncrement"`
	Cust_email            string    `json:"cust_email"`
	Cust_country_code     string    `json:"cust_country_code"`
	Cust_phone            string    `json:"cust_phone"`
	Email_verified        bool      `json:"email_verified"`
	Phone_verified        bool      `json:"phone_verified"`
	Roles                 string    `json:"roles"`
	Time_created          time.Time `json:"time_created"`
	Last_edited           time.Time `json:"last_edited"`
	Salt                  string    `json:"salt"`
	Hashpass              string    `json:"hashpass"`
	Access_token          string    `json:"access_token"`
	Cust_firstname        string    `json:"cust_firstname"`
	Cust_lastname         string    `json:"cust_lastname"`
	Cust_address          string    `json:"cust_address"`
	Date_of_birth         string    `json:"date_of_birth"`
	Age                   int       `json:"age"`
	Identity_type         string    `json:"identity_type"`
	Identity_number       string    `json:"identity_number"`
	Identity_scan_url     string    `json:"identity_scan_url"`
	Bank_name             string    `json:"bank_name"`
	Bank_account_number   string    `json:"bank_account_number"`
	Account_book_scan_url string    `json:"account_book_scan_url"`
	Profile_verified      bool      `json:"profile_verified"`
}

func CreateCustomerAccount(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Origin", "*")

		var applicationObj CustomerAccountProfile

		json.NewDecoder(r.Body).Decode(&applicationObj)

		// create to local Array

		t1 := time.Now().Year()
		t2, err := time.Parse("2006-01-02", applicationObj.Date_of_birth)
		if err != nil {
			fmt.Println(err)
			return
		}
		ageCalc := (t1 - t2.Year())

		objToInsert := CustomerAccountProfile{
			Cust_email:            applicationObj.Cust_email,
			Cust_country_code:     applicationObj.Cust_country_code,
			Cust_phone:            applicationObj.Cust_phone,
			Email_verified:        false,
			Phone_verified:        false,
			Roles:                 "customer",
			Time_created:          time.Now(),
			Last_edited:           time.Now(),
			Salt:                  "",
			Hashpass:              "",
			Access_token:          "",
			Cust_firstname:        applicationObj.Cust_firstname,
			Cust_lastname:         applicationObj.Cust_lastname,
			Cust_address:          applicationObj.Cust_address,
			Date_of_birth:         t2.String(),
			Age:                   ageCalc,
			Identity_type:         applicationObj.Identity_type,
			Identity_number:       applicationObj.Identity_number,
			Identity_scan_url:     applicationObj.Identity_scan_url,
			Bank_name:             applicationObj.Bank_name,
			Bank_account_number:   applicationObj.Bank_account_number,
			Account_book_scan_url: applicationObj.Account_book_scan_url,
			Profile_verified:      false,
		}

		// access context db ###########################################
		db, ok := r.Context().Value(z.DBContext).(*gorm.DB)
		if !ok {
			fmt.Println("something is broke with http dbcontext")
			panic("failed passing context to create handler")
		}
		// GORM create ################################################

		result := db.Create(&objToInsert)
		if err := result.Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println("something is broke with gorm create handler")
				panic("create handler failed to create row to database")
			}
			return
		}
		// ##############################################################

		json.NewEncoder(rw).Encode(objToInsert)

	}

}

func GetAllCustomerAccount(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Origin", "*")

		// access context db ###########################################
		db, ok := r.Context().Value(z.DBContext).(*gorm.DB)
		if !ok {
			fmt.Println("something is broke with http dbcontext")
			panic("failed passing context to create handler")
		}

		var arrAccountProfile []CustomerAccountProfile

		result := db.Find(&arrAccountProfile)

		if err := result.Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				json.NewEncoder(rw).Encode(arrAccountProfile)
				fmt.Println("something is broke with gorm getAll handler")
				panic("getAll handler failed to find row in database")
			}
			if len(arrAccountProfile) == 0 {
				json.NewEncoder(rw).Encode("data is empty, please insert data first!")
			}
			return
		}

		json.NewEncoder(rw).Encode(arrAccountProfile)

	}

}

func UpdateCustomerAccount(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Origin", "*")

		query := r.URL.Query()
		idToEdit, _ := strconv.Atoi(query.Get("id"))

		// access context db ###########################################
		db, ok := r.Context().Value(z.DBContext).(*gorm.DB)
		if !ok {
			fmt.Println("something is broke with http dbcontext")
			panic("failed passing context to create handler")
		}

		// buat account baru yang mengambil data dari request.Body
		var newAccountProfile CustomerAccountProfile
		json.NewDecoder(r.Body).Decode(&newAccountProfile)

		//cari Id dan dapatkan account yang mau diedit
		var oldAccountProfile CustomerAccountProfile
		db.Where("cust_id = ?", idToEdit).Find(&oldAccountProfile)

		//melakukan perubahan data
		newFN := newAccountProfile.Cust_firstname
		newLN := newAccountProfile.Cust_lastname
		newADDR := newAccountProfile.Cust_address

		db.Model(oldAccountProfile).Updates(CustomerAccountProfile{
			Cust_firstname: newFN,
			Cust_lastname:  newLN,
			Cust_address:   newADDR,
		})

		// show response
		json.NewEncoder(rw).Encode(&newAccountProfile)

	}
}

func DeleteCustomerAccount(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Origin", "*")

		query := r.URL.Query()
		idToDelete, _ := strconv.Atoi(query.Get("id"))

		accProfDeleted := CustomerAccountProfile{ // need to improve
			Cust_id: idToDelete,
		}

		log.Printf("Delete Account Handler is triggered for id :%s", strconv.Itoa(idToDelete))

		// access context db ###########################################
		db, ok := r.Context().Value(z.DBContext).(*gorm.DB)
		if !ok {
			fmt.Println("something is broke with http dbcontext")
			panic("failed passing context to create handler")
		}

		db.Delete(&CustomerAccountProfile{}, idToDelete)

		json.NewEncoder(rw).Encode(accProfDeleted)

	}

}
