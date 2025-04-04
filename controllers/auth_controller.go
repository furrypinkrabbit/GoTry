package controllers

import(
	"MyApp/models"
	"MyApp/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"MyApp/global"
	"log"
	
)


func Register(ctx *gin.Context){
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg" : "Invalid request",
		})
		return
	}


	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg" : "Error hashing password",
		})
		return
	}

	user.Password = hashedPwd

	token,err := utils.GenerateJWT(user.Username)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg" : "Error generating token , got error",
		})
		return
	}

	if err := global.DB.AutoMigrate(&user); err != nil{
		log.Fatalf("failed to migrate user model , got error %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg" : "Error migrating user model",
	} )
	return
		}

	if err := global.DB.Create(&user).Error; err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg" : "Error creating user",
		})
		return
	}
	



	ctx.JSON(http.StatusOK, gin.H{"token " : token})

}


func Login(ctx *gin.Context){
	var input struct{
		Username string "json:username"
		Password string "json:password"

	}


	if err := ctx.ShouldBindJSON(&input); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg" : "Invalid request",
		})
		return
	}

	var user models.User

	if err := global.DB.Where("username = ?", input.Username).First(&user).Error; err != nil{
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg" : "User not found",
		})
		return
	}

	if !utils.CheckPassword(input.Password, user.Password){
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error" : "wrong credenttials",
		})
		return	
	}

	token,err := utils.GenerateJWT(user.Username)

	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg" : "Error generating token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token" :token})
}
