package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/restful/api/common"
	"strings"
	"fmt"
	"encoding/base64"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/ontology-crypto/keypair"
	common2 "github.com/ontio/ontology/common"
	"net/http"
	"github.com/ontio/ontology/core/signature"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		header := map[string][]string(c.Request.Header)
		token := header["Authorization"]
		if token[0] == "" {
			err = fmt.Errorf("token is nil")
		} else {
			err = validateToken(token[0])
		}
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":common.VERIFY_TOKEN_ERROR,
				"msg":err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func validateToken(token string) error {
	//header.payloadBs.sig
	arr := strings.Split(token, ".")
	if len(arr) != 3 {
		return fmt.Errorf("wrong token: %s", token)
	}
	sig, err := base64.RawURLEncoding.DecodeString(arr[2])
	if err != nil {
		fmt.Println(err)
		return err
	}
	pubKeyStr,_ := common2.HexToBytes(config.DefConfig.OperatorPublicKey)
	pubKey,err := keypair.DeserializePublicKey(pubKeyStr)
	if err != nil {
		return err
	}
	data := arr[0]+"."+arr[1]
	fmt.Println("sig:",string(sig))
	sig,err = common2.HexToBytes(string(sig))
	if err != nil {
		return err
	}
	return signature.Verify(pubKey, []byte(data), sig)
}