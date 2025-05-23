package main

import(
"MyApp/config"
"MyApp/router"
"log"
"context"
"net/http"
"os"
"os/signal"
"time"
)

func main(){
config.InitConfig()

r := router.SetUpRouter()


 port := config.AppConfig.App.Port


	if port ==""{
		port = "127.0.0.1:8080"
	}

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	   }
	   
		go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	   }
	   }()
	   
	   quit := make(chan os.Signal, 1)
	   signal.Notify(quit, os.Interrupt)
	   <-quit
	   log.Println("Shutdown Server ...")
	   
	   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	   defer cancel()
	   if err := srv.Shutdown(ctx); err != nil {
	   log.Fatal("Server Shutdown:", err)
	   }
	   log.Println("Server exiting")

	



}