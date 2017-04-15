package pathin

// Returns the value of the handler name to be added to generated path
// For instance:
//
//     p.AddDest("profile-picture", pathin.RawHandler)
//
// will add /profile-picture to the path
func RawHandler(handlerName string, values interface{}) (string, error) {
	return "/" + string(handlerName), nil
}
