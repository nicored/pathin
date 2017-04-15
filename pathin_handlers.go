package pathin

type bucketInfo struct {
	bucketId int
	userId   int
}

func rawHandler(typeName string, values interface{}) (string, error) {
	return string(typeName), nil
}
