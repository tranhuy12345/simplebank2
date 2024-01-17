package api

import (
	"db/token"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "payloadKey"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Kiểm tra xem authorization có trong Header hay không?
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provider")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			// Hủy bỏ yêu cầu từ request không next tới middleware tiếp theo
			return
		}
		//Kiểm tra xem authorization này có đúng định dang hay không (Bear :fdfdfdfdf)
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header is not format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		//Kiểm tra xem loại authorization có phải là Bearer hay không?
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := errors.New("authorization header is not bearer")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		//Kiểm tra lại accesstoken
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		//Nếu ok hết thì chuyển đến middleware tiếp theo
		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}
