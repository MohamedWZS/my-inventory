package main

// Our Program Entry Point
func main() {
	app := App{}
	app.Initialize()
	app.Run("localhost:10000")
}
