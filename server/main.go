package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"io"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func GenerateQR(code string, filename string) {
    qrCode, err := qrcode.Encode(fmt.Sprintf("http://localhost:5173/%v", code), qrcode.Medium, 256)
    if err != nil {
        fmt.Printf("Error encoding QR code: %s\n", err)
        return
    }

    err = os.WriteFile(fmt.Sprintf("./qrs/%v", filename), qrCode, 0644)
    if err != nil {
        fmt.Printf("Error writing QR code to file: %s\n", err)
        return
    }

    fmt.Printf("QR code generated and saved as %s\n", filename)
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

func dbStart() (*mongo.Client, context.Context) {
 
    var uri = os.Getenv("DB_URL")
    
    client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect(ctx)

    err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
    return client, ctx
}


func init() {
    err := LoadEnv(".env")
    if err != nil {
        fmt.Printf("Error loading .env file: %v\n", err)
        os.Exit(1)
    }
}

type Participants struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // MongoDB document ID
	TeamName  string             `bson:"teamName"`      // Required field
	Email     string             `bson:"email"`         // Required field
	PhoneNum  string             `bson:"phoneNum"`      // Required field
	Name      string             `bson:"name"`          // Required field
	Breakfast bool               `bson:"breakfast"`     // Default: false
	Lunch     bool               `bson:"lunch"`         // Default: false
	Dinner    bool               `bson:"dinner"`        // Default: false
	Snack1    bool               `bson:"snack1"`        // Default: false
	Snack2    bool               `bson:"snack2"`        // Default: false
	CreatedAt time.Time          `bson:"createdAt,omitempty"` // Timestamps
	UpdatedAt time.Time          `bson:"updatedAt,omitempty"` // Timestamps
}

func getUserDetails(w http.ResponseWriter, r *http.Request) {
    // Extract ID from the URL path
    type UserDetailsResponse struct {
        Participant Participants `json:"participant"`
    }

    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 4 {
        http.Error(w, "Invalid path", http.StatusBadRequest)
        return
    }
    idStr := pathParts[3]

    // Convert the string ID to a MongoDB ObjectID
    objID, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        http.Error(w, "Invalid ObjectID", http.StatusBadRequest)
        return
    }

    // Create a timeout context for the MongoDB operation
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Assuming you have a MongoDB client instance named 'client'
    database := client.Database("hackman-qr")
    participantCollection := database.Collection("participants")

    // Query the collection using the ObjectID
    var participant Participants
    err = participantCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&participant)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "No participant found", http.StatusNotFound)
        } else {
            log.Printf("Error finding participant: %v", err)
            http.Error(w, "Error fetching participant", http.StatusInternalServerError)
        }
        return
    }

    // Prepare the response
    participantS := UserDetailsResponse{
        Participant: participant,
    }

    jsonResponse, err := json.Marshal(participantS)
    if err != nil {
        log.Printf("Error marshaling JSON: %v", err)
        http.Error(w, "Error generating JSON", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func postFoodUpdate(w http.ResponseWriter, r *http.Request) {
    // Extract ID from the URL path
    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 4 {
        http.Error(w, "Invalid path", http.StatusBadRequest)
        return
    }
    idStr := pathParts[3]

    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    var meal struct {
        Meal string `json:"meal"`
    }

    err = json.Unmarshal(bodyBytes, &meal)
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    objID, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        http.Error(w, "Invalid ObjectID", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    database := client.Database("hackman-qr")
    participantCollection := database.Collection("participants")

    var participant Participants
    err = participantCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&participant)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "No participant found", http.StatusNotFound)
        } else {
            log.Printf("Error finding participant: %v", err)
            http.Error(w, "Error fetching participant", http.StatusInternalServerError)
        }
        return
    }

    var res bool
    switch meal.Meal {
    case "breakfast":
        res = !participant.Breakfast
        participant.Breakfast = res
    case "lunch":
        res = !participant.Lunch
        participant.Lunch = res
    case "dinner":
        res = !participant.Dinner
        participant.Dinner = res
    case "snack1":
        res = !participant.Snack1
        participant.Snack1 = res
    case "snack2":
        res = !participant.Snack2
        participant.Snack2 = res
    default:
        http.Error(w, "Invalid meal type", http.StatusBadRequest)
        return
    }

    _, err = participantCollection.UpdateByID(ctx, objID, bson.M{"$set": bson.M{meal.Meal: res}})
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "No participant found", http.StatusNotFound)
        } else {
            log.Printf("Error updating participant: %v", err)
            http.Error(w, "Error updating participant", http.StatusInternalServerError)
        }
        return
    }

    jsonResponse, err := json.Marshal(participant)
    if err != nil {
        log.Printf("Error marshaling JSON: %v", err)
        http.Error(w, "Error generating JSON", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}


func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("SECRET"), nil
			})
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "props", claims)
				// Access context values in handlers like this
				// props, _ := r.Context().Value("props").(jwt.MapClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}
		}
	})
}

func auth(w http.ResponseWriter, r *http.Request) {
    type TokenResponse struct {
        Token string `json:"token"`
    }

    if r.Method != "POST" {
        http.Error(w, "Method is not supported", http.StatusNotFound)
        return
    }
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    var creds struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    err = json.Unmarshal(bodyBytes, &creds)
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    envUsername := os.Getenv("Username")
    envPassword := os.Getenv("Password")

    if creds.Username != envUsername || creds.Password != envPassword {
        http.Error(w, "Unauthorized access", http.StatusUnauthorized)
        return
    }

    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["username"] = creds.Username

    tokenString, err := token.SignedString([]byte("SECRET"))
    if err != nil {
        http.Error(w, "Failed token creation", http.StatusInternalServerError)
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

    client, _ = dbStart()
    
    if err != nil {
        fmt.Printf("Error loading .env file: %v\n", err)
        os.Exit(1)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Assuming you have a MongoDB client instance named 'client'
    // database := client.Database("hackman-qr")
    // participantsCollection := database.Collection("participants")

    // Query the collection using the ObjectID

    // cur, err := participantsCollection.Find(ctx, bson.M{})
    // if err != nil {
    //     log.Printf("Error finding participants: %v", err)
    //     return
    // }
    // defer cur.Close(ctx)

    // var participants []Participants
    // if err = cur.All(ctx, &participants); err != nil {
    //     log.Printf("Error decoding participants: %v", err)
    //     return
    // }

    // fmt.Println(participants)

    // for i := 0; i < len(participants); i++ {
    //     idText, err := participants[i].ID.MarshalText()
    //     if err != nil {
    //         log.Printf("Error marshaling ID: %v", err)
    //         continue
    //     }
    //     GenerateQR(string(idText), fmt.Sprintf("%v.png", participants[i].Name))
    // }

    databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

    fmt.Println(databases)
    http.HandleFunc("/auth", auth)
    http.Handle("/user/details/", middleware(http.HandlerFunc(getUserDetails)))
    http.Handle("/user/update/", middleware(http.HandlerFunc(postFoodUpdate)))

    if err := http.ListenAndServe(":7500", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    } else {
        fmt.Println("Server is running on localhost:7500...")
    }
}
