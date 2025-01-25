package main

// func exportJSFunctions() {
// 	js.Global().Set("setSession", js.FuncOf(setSession))
// }

// //go:export setSession
// func setSession(this js.Value, args []js.Value) interface{} {
// 	session := args[0]
// 	username := args[1]

// 	if !session.Truthy() || !username.Truthy() {
// 		panic("ERROR: session and username are required")
// 	}

// 	// TODO: store session in a global variable

// 	return js.ValueOf(session)
// }
