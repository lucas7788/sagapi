package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateOrderId(t *testing.T) {
	orderId := make([]string, 0)
	for i := 0; i < 10; i++ {
		orderId = append(orderId, GenerateOrderId())
		fmt.Printf("orderId: %s, i:%d", orderId, i)
	}

	for i := 0; i < 10; i++ {
		for j := i + 1; j < 10; j++ {
			assert.NotEqual(t, orderId[i], orderId[j])
		}
	}
}
