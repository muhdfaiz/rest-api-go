package main

import "bitbucket.org/shoppermate/application"

func main() {
	application.Bootstrap().Run("127.0.0.1:8080")
}
