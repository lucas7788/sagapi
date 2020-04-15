package jwt

import (
	"encoding/base64"
	"fmt"
	"github.com/candybox-sig/log"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology-crypto/keypair"
	common2 "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/signature"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/restful/api/common"
	"net/http"
	"strings"
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
			log.Errorf("jwt error:%s", err)
			c.JSON(http.StatusUnauthorized, common.ResponseFailed(common.VERIFY_TOKEN_ERROR, err))
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
		return err
	}
	pubKeyStr, _ := common2.HexToBytes(config.DefConfig.OperatorPublicKey)
	pubKey, err := keypair.DeserializePublicKey(pubKeyStr)
	if err != nil {
		return err
	}
	data := arr[0] + "." + arr[1]
	sig, err = common2.HexToBytes(string(sig))
	if err != nil {
		return err
	}
	return signature.Verify(pubKey, []byte(data), sig)
}
