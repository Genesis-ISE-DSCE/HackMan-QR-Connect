package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/skip2/go-qrcode"

	"bufio"
	"io"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func generateQR(code string, filename string) {
    qrCode, err := qrcode.Encode(code, qrcode.Medium, 256)
    if err!= nil {
        fmt.Printf("Error encoding QR code: %s\n", err)
        return
    }

    err = os.WriteFile(filename, qrCode, 0644)
    if err!= nil {
        fmt.Printf("Error writing QR code to file: %s\n", err)
        return
    }

    fmt.Printf("QR code generated and saved as %s\n", filename)
}

func getUserDetails(w http.ResponseWriter, r *http.Request) {
	/*
		1. DB calls to gather user information with team
		2. Get Dinner, Breakfast and snacks status of the user from DB
		3. Send all the details as in JSON format to the dashboard
	*/ 
}

func postFoodUpdate(w http.ResponseWriter, r *http.Request) {
	/* 
		1. Get the body in request
		2. Update the status for the food
		3. Return confirmation
	*/
}

func LoadEnv(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())

        // Skip empty lines and comments
        if len(line) == 0 || strings.HasPrefix(line, "#") {
            continue
        }

        // Split line into key and value
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            return fmt.Errorf("invalid line: %s", line)
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])

        // Set the environment variable
        err = os.Setenv(key, value)
        if err != nil {
            return err
        }
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    return nil
}

func auth(w http.ResponseWriter, r *http.Request) {
	type TokenResponse struct {
		Token string `json:"token"`
	}

	err := LoadEnv(".env")
    if err != nil {
        fmt.Printf("Error loading .env file: %v\n", err)
        os.Exit(1)
    }

	if r.Method != "POST" {
		http.Error(w,"Method is not supported", http.StatusNotFound)
		return
	}
	bodybytes,err := io.ReadAll(r.Body)
	if err!=nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err = json.Unmarshal(bodybytes, &creds)
	if err!= nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	envUsername := os.Getenv("Username")
	envPassword := os.Getenv("Password")

	if creds.Username != envUsername || creds.Password != envPassword {
		http.Error(w, "Unauhorised access", http.StatusUnauthorized)
		return;
	}

	
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = creds.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err:= token.SignedString([]byte("yoursec"))
	if err !=nil {
		http.Error(w,"failed token creation", http.StatusInternalServerError)
		return
	}
	tokenS := TokenResponse{
        Token: tokenString,
    }

	jsonResponse, err := json.Marshal(tokenS)
    if err != nil {
        http.Error(w, "Error generating JSON", http.StatusInternalServerError)
        return
    }

    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    
    w.Write(jsonResponse)
}

func main() {

	err := LoadEnv(".env")
    if err != nil {
        fmt.Printf("Error loading .env file: %v\n", err)
        os.Exit(1)
    }

	type Product struct {
		gorm.Model
		Code  string
		Price uint
	}

	generateQR("http://your-dashboard-url/123", "adi.png")

	http.HandleFunc("/", getHello)
	http.HandleFunc("/auth", auth)
	http.HandleFunc("/user/details/:id", getUserDetails)
	http.HandleFunc("/user/update/:id", postFoodUpdate)

	

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to db %v", err)
	}
	fmt.Println(db)

	db.AutoMigrate(&Product{})
	db.Create(&Product{Code: "D42", Price: 100})

    if err := http.ListenAndServe(":7500", nil); err!= nil {
        log.Fatalf("Failed to start server: %v", err)
    }else {
		fmt.Println("Server is running on localhost:7500...")
	}
}
