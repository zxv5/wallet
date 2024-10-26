package middleware

import (
	"wallet/internal/gin/middleware/cors"
	"wallet/internal/gin/middleware/ignorerr"
	"wallet/internal/gin/middleware/recovery"
	"wallet/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Base() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Logger(),
		recovery.Recovery(logger.Logger, true),
		ignorerr.Ignorerr(),
		cors.Cors(),
	}
}
