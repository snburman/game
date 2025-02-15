package config

// Env returns Vars struct of environment variables
func Env() Vars {
	return LocalVars()
}
