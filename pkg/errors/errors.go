package errors

import "errors"

var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")

	// Flight errors
	ErrFlightNotFound       = errors.New("flight not found")
	ErrInvalidDepartureDate = errors.New("departure date must be in the future")
	ErrInvalidReturnDate    = errors.New("return date must be after departure date")
	ErrInsufficientSeats    = errors.New("insufficient seats available")

	// Hotel errors
	ErrHotelNotFound       = errors.New("hotel not found")
	ErrInvalidCheckInDate  = errors.New("check-in date must be in the future")
	ErrInvalidCheckOutDate = errors.New("check-out date must be after check-in date")
	ErrNoRoomsAvailable    = errors.New("no rooms available")
	ErrInvalidStayDuration = errors.New("invalid stay duration")

	// Booking errors
	ErrBookingNotFound      = errors.New("booking not found")
	ErrInvalidBookingType   = errors.New("invalid booking type")
	ErrBookingCancelled     = errors.New("booking already cancelled")
	ErrInvalidBookingStatus = errors.New("invalid booking status")

	// Payment errors
	ErrPaymentFailed        = errors.New("payment failed")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrPaymentDeclined      = errors.New("payment declined")

	// Validation errors
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidPassengers = errors.New("invalid number of passengers")
	ErrInvalidPrice      = errors.New("invalid price")
)

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Error       string `json:"error"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
