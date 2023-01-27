package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", helloWorldHandler).Methods("GET")
	r.HandleFunc("/token/{username}", tokenHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World !")
}

func generateToken(username string) string {
	timestamp := time.Now().Unix()
	h := sha256.New()
	randomized_username := username + strconv.FormatInt(timestamp, 10)
	h.Write([]byte(randomized_username))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return sha
}

var ctx = context.Background()

func getOrGeneratToken(username string) string {
	redisEndpoint := fmt.Sprintf("%s:%s", os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PORT"))
	log.Printf("Redis address %s", redisEndpoint)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisEndpoint,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	token, err := rdb.Get(username).Result()
	if err != nil {
		log.Printf("No token in cache for %s", username)
		newToken := generateToken(username)
		err := rdb.Set(username, newToken, 60*time.Second).Err()
		if err != nil {
			panic(err)
		}
		return newToken
	}

	log.Printf("Using cache for %s", username)
	return token
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	token := getOrGeneratToken(username)
	fmt.Fprintf(w, token)
	// Generate token
}
