package app

func StartApplicaton() {
	app := GinRouter{}

	app.Init()
	app.MapRoutes()
	app.Run(":8080")
}
