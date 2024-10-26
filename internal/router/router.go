package router

import (
	"os"
	"wallet/internal/config"
	"wallet/internal/controller"
	"wallet/internal/gin/middleware"
	"wallet/internal/gin/middleware/jwt"
	"wallet/internal/service"

	_ "wallet/docs"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/gin-gonic/gin"
)

// Init router
func Init(conf *config.Config) *gin.Engine {
	setMode()

	r := gin.New()
	r.Use(middleware.Base()...)
	jwt := jwt.Jwt(&conf.Jwt)
	svcCtx := service.NewServiceContext(conf)

	{
		user := controller.NewUserController(svcCtx)
		r := r.Group("/api/v1/user")

		r.POST("/login", user.Login)
	}
	{
		wallet := controller.NewWalletController(svcCtx)
		r := r.Group("/api/v1/wallet")

		r.POST("/deposit", jwt, wallet.Deposit)
		r.POST("/withdraw", jwt, wallet.Withdraw)
		r.POST("/transfer", jwt, wallet.Transfer)
		r.GET("/balance", jwt, wallet.Balance)
		r.GET("/record", jwt, wallet.Record)
	}

	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER"))

	return r
}

func setMode() {
	runMode := "release"
	if env := os.Getenv("GO_MODE"); env == "debug" {
		runMode = "debug"
	}
	gin.SetMode(runMode)
}
