package pathin

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type chanResp struct {
	path chan string
	err  chan error
}

func TestNewFS(t *testing.T) {
	myFs := New("bucket-name")
	assert.Equal(t, myFs.Name(), "bucket-name")

	inBucketDest := myFs.AddDestGroup("companyBucket", groupHandler)
	inBucketDest.AddDest("cad-files", rawHandler)

	inUserDest := inBucketDest.AddDestGroup("userBucket", userHandler)
	inUserDest.AddDest("profile-picture", rawHandler)
}

func groupHandler(typeName string, values interface{}) (string, error) {
	info, ok := values.(*bucketInfo)
	if ok == false || info == nil {
		return "", errors.New("No readable data")
	}

	if info.bucketId > 0 {
		return "buckets/bucket_" + strconv.Itoa(info.bucketId), nil
	}

	return "", errors.New("No bucket Id defined")
}

func userHandler(typeName string, values interface{}) (string, error) {
	info, ok := values.(*bucketInfo)
	if ok == false || info == nil {
		return "", errors.New("No readable data")
	}

	if info.userId > 0 {
		return "users/" + string(typeName) + strconv.Itoa(info.userId), nil
	}

	return "", errors.New("No user Id defined")
}
