package handler

import (
	"fledge-restapi/internal/domain/entity"
	"fledge-restapi/internal/service"
	"fledge-restapi/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FlightHandler struct {
	flightService service.FlightService
}

func NewFlightHandler(flightService service.FlightService) *FlightHandler {
	return &FlightHandler{
		flightService: flightService,
	}
}

// SearchFlights godoc
// @Summary Search for flights
// @Description Search for flights based on criteria
// @Tags flights
// @Accept json
// @Produce json
// @Param search body entity.FlightSearchRequest true "Flight search criteria"
// @Success 200 {array} entity.Flight
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/flights/search [get]
func (h *FlightHandler) SearchFlights(c *gin.Context) {
	var req entity.FlightSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flights, err := h.flightService.SearchFlights(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search flights"})
		return
	}

	c.JSON(http.StatusOK, flights)
}

// GetFlight godoc
// @Summary Get flight details
// @Description Get detailed information about a specific flight
// @Tags flights
// @Accept json
// @Produce json
// @Param id path int true "Flight ID"
// @Success 200 {object} entity.Flight
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/flights/{id} [get]
func (h *FlightHandler) GetFlight(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	flight, err := h.flightService.GetFlightByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(http.StatusOK, flight)
}

// BookFlight godoc
// @Summary Book a flight
// @Description Create a new flight booking
// @Tags flights
// @Accept json
// @Produce json
// @Param id path int true "Flight ID"
// @Param booking body entity.BookingRequest true "Booking details"
// @Success 201 {object} entity.Booking
// @Failure 400 {object} errors.ErrorResponse
// @Security Bearer
// @Router /api/flights/{id}/book [post]
func (h *FlightHandler) BookFlight(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req entity.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flightID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	req.FlightID = new(uint)
	*req.FlightID = uint(flightID)
	req.BookingType = "flight"

	booking, err := h.flightService.BookFlight(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		switch err {
		case errors.ErrInsufficientSeats:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient seats available"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		}
		return
	}

	c.JSON(http.StatusCreated, booking)
}
