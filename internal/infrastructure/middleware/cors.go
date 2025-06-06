
package middleware

import (


	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
    config := cors.Config{
        AllowOrigins:     []string{"*"}, 
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS","HEAD"},
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"*"},
    }
    return cors.New(config)
}
