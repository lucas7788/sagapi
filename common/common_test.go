package common

import (
	"encoding/json"
	"fmt"
	"github.com/ontio/sagapi/models/tables"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildApiBasicInfo(t *testing.T) {
	type A struct {
		AB string
	}
	type B struct {
		ABC string
	}
	type C struct {
		*A
		*B
	}
	c := C {
		&A{
			AB:"ab",
		},
		&B{
			ABC:"abc",
		},
	}
	bs, _ := json.Marshal(c)
	fmt.Println(string(bs))
}

func TestGenerateOrderId(t *testing.T) {
	orderId := make([]string, 0)
	for i := 0; i < 10; i++ {
		orderId = append(orderId, GenerateUUId())
		fmt.Println(orderId)
	}

	for i := 0; i < 10; i++ {
		for j := i + 1; j < 10; j++ {
			assert.NotEqual(t, orderId[i], orderId[j])
		}
	}
}

func TestBuildQrCodeResult(t *testing.T) {
	res := ApiDetailResponse{
		ApiId: 1,
		RequestParams: []*tables.RequestParam{
			&tables.RequestParam{
				ApiDetailInfoId: 1,
				ParamName:       "aa",
			},
		},
		ErrorCodes: []*tables.ErrorCode{
			&tables.ErrorCode{
				ApiDetailInfoId: 1,
				ErrorCode:       2,
			},
		},
	}

	bs, _ := json.Marshal(res)
	fmt.Println(string(bs))
}
