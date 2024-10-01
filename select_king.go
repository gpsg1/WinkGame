package main

import (
	"fmt"
	"math/rand"
	"time"
)

func Get(user_list []string) string {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random number between 0 and 100
	randomNumber := r.Intn(len(user_list)) // Intn(n) returns a random integer from 0 to n-1, so 101 gives 0 to 100

	fmt.Println(randomNumber) // Print the random number
	return user_list[randomNumber]
}
