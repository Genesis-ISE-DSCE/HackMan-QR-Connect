package main

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}



func main() {

	type Product struct {
		gorm.Model
		Code  string
		Price uint
	}

	http.HandleFunc("/", getHello)
	fmt.Println("Server is running on localhost:7500...")

	db, err := gorm.Open(postgres.Open("postgresql://hackman_owner:6kRxDQUd2ATg@ep-empty-fire-a5vgftnm.us-east-2.aws.neon.tech/hackman?sslmode=require"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to db %v", err)
	}
	fmt.Println(db)

	db.AutoMigrate(&Product{})
	db.Create(&Product{Code: "D42", Price: 100})

    if err := http.ListenAndServe(":7500", nil); err!= nil {
        log.Fatalf("Failed to start server: %v", err)
    }
	// dbConnect("postgresql://hackman_owner:6kRxDQUd2ATg@ep-empty-fire-a5vgftnm.us-east-2.aws.neon.tech/hackman?sslmode=require")
}