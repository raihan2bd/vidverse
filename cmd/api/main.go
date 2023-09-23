package main

import "github.com/raihan2bd/vidverse/initializers"

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.ConnectToCloudinary()
	initializers.SyncDatabase()
}

func main() {
	r := NewRouter()

	r.Run()
}