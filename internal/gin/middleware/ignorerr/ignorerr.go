package ignorerr

import (
	"wallet/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Ignorerr Ignore panic return
func Ignorerr() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if c.GetBool("ignorePainc") {
				if err := recover(); err != nil {
					logger.Warnf("Ignore Panic-> %+v", err)
				}
			}
		}()

		c.Next()
	}
}
