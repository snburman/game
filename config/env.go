package config

// Env returns Vars struct of environment variables
func Env() Vars {
	// if flag.Lookup("test.v") == nil {
	// 	err := godotenv.Load()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	return Vars{
	// 		SERVER_URL:    os.Getenv("SERVER_URL"),
	// 		CLIENT_ID:     os.Getenv("CLIENT_ID"),
	// 		CLIENT_SECRET: os.Getenv("CLIENT_SECRET"),
	// 		MONGO_URI:     os.Getenv("MONGO_URI"),
	// 	}
	// }
	return LocalVars()
}
