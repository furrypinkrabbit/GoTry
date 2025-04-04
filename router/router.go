package router

import (
	"MyApp/controllers"
	"MyApp/middlewares"
	"time"
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 1 << 30 // 1 GB
	
	 // 自定义静态文件处理函数
	 r.GET("/videos/*filepath", func(c *gin.Context) {
        // 添加 CORS 头
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

        http.ServeFile(c.Writer, c.Request, "D:/NewBack/scripts/4C2025/video_action_cls/video_test_output/shear_test"+c.Param("filepath"))
    })


	// CORS 配置
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "*","http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization","accpet"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))

	// 认证路由
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}

	// 需要认证的 API 路由
	api := r.Group("/api")
	api.GET("/getVideo", controllers.GetMatchedFiles)
	api.GET("/getVideoTest", controllers.GetProcessedVideo)
	api.POST("/log",controllers.ReceiveFrontendLog)
	api.Use(middlewares.AuthMiddleWare())
	{
		api.POST("/articles", controllers.CreateArticle)
		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/:id", controllers.GetArticleByID)
		api.POST("/articles/:id/like", controllers.LikeArticle)
		api.GET("/articles/:id/like", controllers.GetArticleLikes)
		api.POST("/uploadVideo", controllers.UploadVideoL)

}
	return r
}
