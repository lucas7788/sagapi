package nasa

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApod(t *testing.T) {
	bs, err := Apod()
	assert.Nil(t, err)
	fmt.Println("res:", string(bs))
}

func TestFeed(t *testing.T) {
	res, err := Feed("2015-09-07", "2015-09-08")
	assert.Nil(t, err)
	fmt.Println("res:", string(res))
}
