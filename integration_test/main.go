package main

import "context"

func main() {
	var url = "http://localhost:8080"
	var ctx = context.Background()
	marvelTest(ctx, url)
}
