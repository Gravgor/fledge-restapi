package main

import (
	"fledge-restapi/internal/config"
	"fledge-restapi/internal/domain/repository"
	"fledge-restapi/internal/handler"
	"fledge-restapi/internal/middleware"
	"fledge-restapi/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db := config.InitDB()

	userRepo := repository.NewUserRepository(db)
	flightRepo := repository.NewFlightRepository(db)
	hotelRepo := repository.NewHotelRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	flightService := service.NewFlightService(flightRepo, bookingRepo)
	hotelService := service.NewHotelService(hotelRepo, bookingRepo)
	bookingService := service.NewBookingService(bookingRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	flightHandler := handler.NewFlightHandler(flightService)
	hotelHandler := handler.NewHotelHandler(hotelService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	// Setup router
	r := gin.Default()

	// Middleware
	r.Use(middleware.RateLimiter())
	r.Use(middleware.Cors())

	// Public routes
	r.POST("/auth/signup", userHandler.Signup)
	r.POST("/auth/login", userHandler.Login)
	//r.POST("/auth/refresh", userHandler.RefreshToken)

	// API routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Flight routes
		api.GET("/flights/search", flightHandler.SearchFlights)
		api.GET("/flights/:id", flightHandler.GetFlight)
		api.POST("/flights/:id/book", flightHandler.BookFlight)

		// Hotel routes
		api.GET("/hotels/search", hotelHandler.SearchHotels)
		api.GET("/hotels/:id", hotelHandler.GetHotel)
		api.POST("/hotels/:id/book", hotelHandler.BookHotel)

		// Booking routes
		api.GET("/bookings", bookingHandler.ListBookings)
		api.GET("/bookings/:id", bookingHandler.GetBooking)
		api.PATCH("/bookings/:id", bookingHandler.UpdateBooking)
		api.DELETE("/bookings/:id", bookingHandler.CancelBooking)

		// Profile routes
		api.GET("/profile", userHandler.GetProfile)
		api.PUT("/profile", userHandler.UpdateProfile)
		api.GET("/profile/preferences", userHandler.GetPreferences)
		api.PUT("/profile/preferences", userHandler.UpdatePreferences)
	}

	// Start server
	r.Run(":8080")
}
