package utils

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHTTPPostRawData(t *testing.T) {
	resp, err := HTTPPostRawData("http://127.0.0.1:8888/v1/chain/get_info", "")
	assert.Nil(t, err)

	fmt.Println(string(resp))
}