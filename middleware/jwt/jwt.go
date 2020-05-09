package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology-crypto/keypair"
	common2 "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/ontology/core/signature"
	"github.com/ontio/sagapi/restful/api/common"
	"github.com/ontio/sagapi/sagaconfig"
	"net/http"
	"strings"
	"time"
)

const (
	ROLE_ADMIN string = "ADMIN"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		header := map[string][]string(c.Request.Header)
		token := header["Authorization"]
		var ontid string
		if token == nil || token[0] == "" {
			err = fmt.Errorf("token is nil")
		} else {
			ontid, err = validateToken(token[0], false)
		}
		if err != nil {
			log.Errorf("token error:%s", err)
			c.JSON(http.StatusUnauthorized, common.ResponseFailed(common.VERIFY_TOKEN_ERROR, err))
			c.Abort()
			return
		}
		c.Set(sagaconfig.Key_OntId, ontid)
		c.Next()
	}
}

func JWTAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		header := map[string][]string(c.Request.Header)
		token := header["Authorization"]
		var ontid string
		if token == nil || token[0] == "" {
			err = fmt.Errorf("token is nil")
		} else {
			ontid, err = validateToken(token[0], true)
		}
		if err != nil {
			log.Errorf("token error:%s", err)
			c.JSON(http.StatusUnauthorized, common.ResponseFailed(common.VERIFY_TOKEN_ERROR, err))
			c.Abort()
			return
		}
		c.Set(sagaconfig.Key_OntId, ontid)
		c.Next()
	}
}

func validateToken(token string, admin bool) (string, error) {
	//header.payloadBs.sig
	arr := strings.Split(token, ".")
	if len(arr) != 3 {
		return "", fmt.Errorf("wrong token: %s", token)
	}
	sig, err := base64.RawURLEncoding.DecodeString(arr[2])
	if err != nil {
		return "", err
	}
	pubKeyStr, _ := common2.HexToBytes(sagaconfig.DefSagaConfig.OperatorPublicKey)
	pubKey, err := keypair.DeserializePublicKey(pubKeyStr)
	if err != nil {
		return "", err
	}
	data := arr[0] + "." + arr[1]
	sig, err = common2.HexToBytes(string(sig))
	if err != nil {
		return "", err
	}
	err = signature.Verify(pubKey, []byte(data), sig)
	if err != nil {
		return "", err
	}
	payloadBs, err := base64.RawURLEncoding.DecodeString(arr[1])
	if err != nil {
		return "", err
	}
	pl := &Payload{}
	err = json.Unmarshal(payloadBs, pl)
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()
	if pl.Exp > int(now) {
		return "", fmt.Errorf("jwt token expired")
	}

	if admin && pl.Content.Role != ROLE_ADMIN {
		return "", fmt.Errorf("jwt token sould be admin")
	}

	return pl.Content.OntId, nil
}
