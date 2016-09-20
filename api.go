package main

import "bitbucket.org/shoppermate-api/application"

func main() {
	application.Bootstrap().Run(":8080")
}
