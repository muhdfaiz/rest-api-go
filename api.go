package main

import "bitbucket.org/shoppermate-api/application"

func main() {
	application.Bootstrap().Run("127.0.0.1:8080")
}
