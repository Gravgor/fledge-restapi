package handler

import (
	"fledge-restapi/internal/domain/entity"
	"fledge-restapi/internal/service"
	"fledge-restapi/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HotelHandler struct {
	hotelService service.HotelService
}

func NewHotelHandler(hotelService service.HotelService) *HotelHandler {
	return &HotelHandler{
		hotelService: hotelService,
	}
}

// SearchHotels godoc
// @Summary Search for hotels
// @Description Search for hotels based on criteria
// @Tags hotels
// @Accept json
// @Produce json
// @Param search body entity.HotelSearchRequest true "Hotel search criteria"
// @Success 200 {array} entity.Hotel
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/hotels/search [get]
func (h *HotelHandler) SearchHotels(c *gin.Context) {
	var req entity.HotelSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotels, err := h.hotelService.SearchHotels(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search hotels"})
		return
	}

	c.JSON(http.StatusOK, hotels)
}

// GetHotel godoc
// @Summary Get hotel details
// @Description Get detailed information about a specific hotel
// @Tags hotels
// @Accept json
// @Produce json
// @Param id path int true "Hotel ID"
// @Success 200 {object} entity.Hotel
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/hotels/{id} [get]
func (h *HotelHandler) GetHotel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	hotel, err := h.hotelService.GetHotelByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

// BookHotel godoc
// @Summary Book a hotel
// @Description Create a new hotel booking
// @Tags hotels
// @Accept json
// @Produce json
// @Param id path int true "Hotel ID"
// @Param booking body entity.BookingRequest true "Booking details"
// @Success 201 {object} entity.Booking
// @Failure 400 {object} errors.ErrorResponse
// @Security Bearer
// @Router /api/hotels/{id}/book [post]
func (h *HotelHandler) BookHotel(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var req entity.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	req.HotelID = new(uint)
	*req.HotelID = uint(hotelID)
	req.BookingType = "hotel"

	booking, err := h.hotelService.BookHotel(c.Request.Context(), userID, &req)
	if err != nil {
		switch err {
		case errors.ErrNoRoomsAvailable:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient rooms available"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		}
		return
	}

	c.JSON(http.StatusCreated, booking)
}
