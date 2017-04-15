package pathin

import (
	"errors"
	"strconv"
)

type bucketInfo struct {
	bucketId int
	userId   int
}

func rawHandler(typeName handlerName, values interface{}) (string, error) {
	return string(typeName), nil
}

func groupHandler(typeName handlerName, values interface{}) (string, error) {
	info, ok := values.(*bucketInfo)
	if ok == false || info == nil {
		return "", errors.New("No readable data")
	}

	if info.bucketId > 0 {
		return "buckets/bucket_" + strconv.Itoa(info.bucketId), nil
	}

	return "", errors.New("No bucket Id defined")
}

func userHandler(typeName handlerName, values interface{}) (string, error) {
	info, ok := values.(*bucketInfo)
	if ok == false || info == nil {
		return "", errors.New("No readable data")
	}

	if info.userId > 0 {
		return "users/" + string(typeName) + strconv.Itoa(info.userId), nil
	}

	return "", errors.New("No user Id defined")
}
