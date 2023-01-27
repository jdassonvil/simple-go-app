package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

/*
   rdb := redis.NewClient(&redis.Options{
       Addr:     "localhost:6379",
       Password: "", // no password set
       DB:       0,  // use default DB
   })

   err := rdb.Set(ctx, "key", "value", 0).Err()
   if err != nil {
       panic(err)
   }

   val, err := rdb.Get(ctx, "key").Result()
   if err != nil {
       panic(err)
   }
*/

/*func getOrGeneratToken(username string) string {

}*/

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	token := generateToken(username)
	fmt.Fprintf(w, token)
	// Generate token
}
