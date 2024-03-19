package api

import (
	"rent-car/api/handler"
	"rent-car/storage"

	"github.com/gin-gonic/gin"
)

func NewCarApi(store storage.IStorage) *gin.Engine {
	h := handler.NewStrg(store)
	r := gin.Default()

	r.POST("/car", h.CreateCar)
	// r.GET("/car/:id", h.GetAllCars)
	r.GET("/car", h.GetAllCars)
	r.PUT("/car/:id", h.UpdateCar)
	r.DELETE("/car/:id", h.DeleteCar)
	// r.PATCH("/car/:id", h.UpdateUserPassword)
	return r
}
