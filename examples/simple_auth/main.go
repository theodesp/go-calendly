package main

import (
	"context"
	"fmt"
	"go-calendly/calendly"
	"os"
)

var (
	apiKey = os.Getenv("CALENDLY_API_KEY")
)

func main() {
	ctx := context.Background()
	authClient := calendly.NewTokenAuthClient(&calendly.Config{ApiKey: apiKey})
	client := calendly.NewClient(authClient)
	resp, _, _ := client.Echo(ctx)

	fmt.Println(resp)
}
