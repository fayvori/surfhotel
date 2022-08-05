package main

import (
	"github.com/gin-gonic/gin"
	"surf-hotel-engine/handlers"
)

func main() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	v1 := r.Group("/hotels/api/v1")
	{
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"service": "Surf Hotel Engine Service",
			})
		})
		v1.POST("/search_id", handlers.GetSearchId)
		v1.POST("/hotel_search", handlers.SearchHotel)
		v1.GET("/hotels_db", handlers.HotelsDb)
		v1.GET("/single_search_id", handlers.SingleHotelSearchId)
		v1.GET("/update_currencies", handlers.UpdateCurrencies)
	}

	r.Run(":4000")
}
