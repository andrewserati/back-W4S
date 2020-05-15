package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"w4s/DB"
	"w4s/controllers"
	"w4s/middleware"
)

func main() {
	//creating connection with database
	r := gin.Default()     //starting the gin. //Iniciando o gin
	db := DB.SetupModels() //Connection database //Conexão banco de dados
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.Static("/css", "tela_alterar_senha/css")
	r.Static("/images", "tela_alterar_senha/images")
	r.LoadHTMLFiles("tela_alterar_senha/index.html")

	authorized := r.Group("/v1")
	r.POST("/login", controllers.Login)
	r.POST("/user/create", controllers.CreateUser)
	r.GET("/user/confirm", controllers.ConfirmUser)

	r.POST("user/password/recovery", controllers.RecoveryPasswordUser)

	recoveryPassword := r.Group("/user/password/recovery")
	recoveryPassword.Use(middleware.AuthRequired2)
	{
		recoveryPassword.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		})
		recoveryPassword.PUT("", controllers.ChangeExternalPassword)
	}
	authorized.Use(middleware.AuthRequired)
	{
		authorized.GET("/searchall", controllers.FindUser)
		authorized.GET("/search", controllers.FindUserByNick)
		authorized.PATCH("/update/user/createprofile", controllers.CreateProfile)

		authorized.PATCH("/update/user", controllers.UpdateUser)
		authorized.PATCH("/logoff", controllers.Logoff)
		authorized.DELETE("/delete/user", controllers.SoftDeletedUserByNick)

	}

	err := r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080") // listando e escutando no localhost:8080
	if err != nil {
		panic("NOT POSSIBLE RUN")
	}
}