package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	ctx "vanilla-loan-application/ctx"

	"gorm.io/gorm"
)

type LoanApplication struct {
	Application_id       int       `json:"application_id" gorm:"primaryKey;autoIncrement"`
	Loan_type            string    `json:"loan_type"`
	Customer_id          int64     `json:"customer_id"`
	Cust_name            string    `json:"cust_name"`
	Cust_phonenumber     string    `json:"cust_phonenumber"`
	Cust_address         string    `json:"cust_address"`
	Monthly_income       int64     `json:"monthly_income"`
	Submit_date          time.Time `json:"submit_date"`
	Application_status   string    `json:"application_status"`
	Nominal_borrowed     int64     `json:"nominal_borrowed"`
	Loan_period_m        int64     `json:"loan_period_m"`
	Annual_interest_rate float32   `json:"annual_interest_rate"`
	Interest_tobe_paid   int64     `json:"interest_tobe_paid"`
	Total_tobe_paid      int64     `json:"total_tobe_paid"`
	Period_type          string    `json:"period_type"`
	Payment_per_period   int64     `json:"payment_per_period"`
}

var arrLoan []LoanApplication

func CreateLoanApplication(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Origin", "*")

		var applicationObj LoanApplication

		json.NewDecoder(r.Body).Decode(&applicationObj)

		// kalkulasi data berdasarkan data di request body
		totalInterest := int64(float32(applicationObj.Nominal_borrowed) * applicationObj.Annual_interest_rate)
		totalToPay := applicationObj.Nominal_borrowed + totalInterest
		paymentPerMonth := totalToPay / applicationObj.Loan_period_m

		// membuat object baru untuk di insert
		loanToInsert := LoanApplication{
			Loan_type:            applicationObj.Loan_type,
			Customer_id:          applicationObj.Customer_id,
			Cust_name:            applicationObj.Cust_name,
			Cust_phonenumber:     applicationObj.Cust_phonenumber,
			Cust_address:         applicationObj.Cust_address,
			Monthly_income:       applicationObj.Monthly_income,
			Submit_date:          time.Now(),       //
			Application_status:   "To Be Reviewed", //
			Nominal_borrowed:     applicationObj.Nominal_borrowed,
			Loan_period_m:        applicationObj.Loan_period_m,
			Annual_interest_rate: applicationObj.Annual_interest_rate,
			Interest_tobe_paid:   totalInterest, //
			Total_tobe_paid:      totalToPay,    //
			Period_type:          "perbulan",
			Payment_per_period:   paymentPerMonth, //
		}

		// access context db ###########################################
		db, ok := r.Context().Value(ctx.DBContext).(*gorm.DB)
		if !ok {
			fmt.Println("something is broke with http dbcontext")
			panic("failed passing context to create handler")
		}

		// GORM create ################################################

		result := db.Create(&loanToInsert)
		if err := result.Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println("something is broke with gorm create handler")
				panic("create handler failed to create row to database")
			}
			return
		}

		// ##############################################################

		// show response
		json.NewEncoder(rw).Encode(loanToInsert)

	}

}

func GetAllLoanApplication(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Origin", "*")

		// access context db ###########################################
		db, ok := r.Context().Value(ctx.DBContext).(*gorm.DB)
		if !ok {
			fmt.Println("something is broke with http dbcontext")
			panic("failed passing context to create handler")
		}

		// array of object
		var arrLoanApplication []LoanApplication

		// gorm get all
		result := db.Find(&arrLoanApplication)

		// error handling
		if err := result.Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				json.NewEncoder(rw).Encode(arrLoanApplication)
				fmt.Println("something is broke with gorm getAll handler")
				panic("getAll handler failed to find row in database")
			}
			if len(arrLoanApplication) == 0 {
				json.NewEncoder(rw).Encode("data is empty, please insert data first!")
			}
			return
		}

		// show response
		json.NewEncoder(rw).Encode(arrLoanApplication)

	}

}

func UpdateLoanApplication(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Origin", "*")

		// mendapatkan id yang mau diedit
		query := r.URL.Query()
		idToEdit, _ := strconv.Atoi(query.Get("id"))

		// access context db ###########################################
		db, ok := r.Context().Value(ctx.DBContext).(*gorm.DB)
		if !ok {
			fmt.Println("something is broke with http dbcontext")
			panic("failed passing context to create handler")
		}

		// buat loan baru yang mengambil data dari request.Body
		var newLoanObj LoanApplication
		json.NewDecoder(r.Body).Decode(&newLoanObj)

		//cari Id dan dapatkan loan yang mau diedit
		var oldLoanObj LoanApplication
		db.Where("application_id = ?", idToEdit).Find(&oldLoanObj)

		//melakukan perubahan data
		newMonthlyIncome := newLoanObj.Monthly_income
		newNominalBorrowed := newLoanObj.Nominal_borrowed
		newLoanPeriod := newLoanObj.Loan_period_m

		//kalkulasi perubahan data
		totalInterest := int64(float32(newLoanObj.Nominal_borrowed) * oldLoanObj.Annual_interest_rate)
		totalToPay := newNominalBorrowed + totalInterest
		paymentPerMonth := totalToPay / newLoanPeriod

		// gorm update
		db.Model(oldLoanObj).Updates(LoanApplication{
			Monthly_income:     newMonthlyIncome,
			Nominal_borrowed:   newNominalBorrowed,
			Loan_period_m:      newLoanPeriod,
			Interest_tobe_paid: totalInterest,   //
			Total_tobe_paid:    totalToPay,      //
			Payment_per_period: paymentPerMonth, //
		})

		// show response
		json.NewEncoder(rw).Encode(&newLoanObj)

	}

}

func DeleteLoanApplication(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Origin", "*")

		// mendapatkan id object yg akan di delete dengan query strings
		query := r.URL.Query()
		idToDelete, _ := strconv.Atoi(query.Get("id"))

		loanToDelete := LoanApplication{ // need to improve
			Application_id: idToDelete,
		}

		// console log id yg akan didelete
		log.Printf("Delete loan handler is triggered for id :%s", strconv.Itoa(idToDelete))

		// access context db ###########################################
		db, ok := r.Context().Value(ctx.DBContext).(*gorm.DB)
		if !ok {
			fmt.Println("something is broke with http dbcontext")
			panic("failed passing context to create handler")
		}

		// gorm delete
		db.Delete(&LoanApplication{}, idToDelete)

		// show response
		json.NewEncoder(rw).Encode(loanToDelete)

	}

}

func GetLoanById(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Origin", "*")

		if len(arrLoan) == 0 {
			json.NewEncoder(rw).Encode("data is empty, please insert data first!")
		}
		if len(arrLoan) > 0 {
			json.NewEncoder(rw).Encode(arrLoan)
		}

	}
}
