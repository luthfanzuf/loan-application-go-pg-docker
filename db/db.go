package db

import (
	"context"
	"fmt"
	"net/http"

	x "vanilla-loan-application/handlers/cust"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	z "vanilla-loan-application/ctx"
)

// membuat key dbcontext untuk pass context value ke handler

func DBMiddleware(next http.Handler, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), z.DBContext, db)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DBConnect(dsn string) (db *gorm.DB, e error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	sqlDB, _ := db.DB()

	err = sqlDB.Ping()
	if err != nil {
		fmt.Println("Cannot ping the database")
		panic(err)
	}
	fmt.Println("Connected to DB successfully!")

	// automigrate struct
	db.AutoMigrate(&x.LoanApplication{})
	db.AutoMigrate(&x.CustomerAccountProfile{})

	return db, e
}
