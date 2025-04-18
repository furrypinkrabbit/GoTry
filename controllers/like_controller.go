package controllers

import (
	"MyApp/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)


func LikeArticle(ctx *gin.Context){
	articleID := ctx.Param("id")

	likeKey := "article:" + articleID + ":likes"

	if err := global.RedisDB.Incr(ctx, likeKey).Err(); err !=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully liked the article"})
}

func GetArticleLikes(ctx *gin.Context){
	articleID := ctx.Param("id")

	likeKey := "article:" + articleID + ":likes"

	likes, err := global.RedisDB.Get(ctx, likeKey).Result()

	if err == redis.Nil{
		likes = "0"
	}else if err !=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"likes": likes})
}