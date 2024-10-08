package entity

import (
	"time"

	"gorm.io/gorm"
)

// User represents the application user
type User struct {
	gorm.Model
	Email       string          `json:"email" gorm:"unique;not null"`
	Password    string          `json:"-" gorm:"not null"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	PhoneNumber string          `json:"phone_number"`
	Preferences UserPreferences `json:"preferences" gorm:"foreignKey:UserID"`
	Bookings    []Booking       `json:"bookings,omitempty" gorm:"foreignKey:UserID"`
	Role        string          `json:"role" gorm:"default:'user'"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SignupRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserPreferences stores user's travel preferences
type UserPreferences struct {
	gorm.Model
	UserID            uint   `json:"user_id"`
	PreferredSeat     string `json:"preferred_seat"` // window, aisle
	MealPreference    string `json:"meal_preference"`
	PreferredAirlines string `json:"preferred_airlines"`
	PreferredHotels   string `json:"preferred_hotels"`
}

// Flight represents a flight offering
type Flight struct {
	gorm.Model
	FlightNumber   string    `json:"flight_number"`
	Airline        string    `json:"airline"`
	DepartureCity  string    `json:"departure_city"`
	ArrivalCity    string    `json:"arrival_city"`
	DepartureTime  time.Time `json:"departure_time"`
	ArrivalTime    time.Time `json:"arrival_time"`
	AvailableSeats int       `json:"available_seats"`
	Price          float64   `json:"price"`
	Class          string    `json:"class"` // economy, business, first
	Status         string    `json:"status"`
}

// Hotel represents a hotel offering
type Hotel struct {
	gorm.Model
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	City           string    `json:"city"`
	Country        string    `json:"country"`
	Rating         float32   `json:"rating"`
	Price          float64   `json:"price_per_night"`
	AvailableRooms int       `json:"available_rooms"`
	Amenities      []Amenity `json:"amenities" gorm:"many2many:hotel_amenities;"`
}

// Amenity represents hotel amenities
type Amenity struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Hotels      []Hotel `json:"hotels" gorm:"many2many:hotel_amenities;"`
}

// VacationPackage represents a pre-built vacation package
type VacationPackage struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Destination string    `json:"destination"`
	Duration    int       `json:"duration_days"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Price       float64   `json:"price"`
	Includes    string    `json:"includes"`
	MaxPeople   int       `json:"max_people"`
	Available   bool      `json:"available"`
}

// Booking represents a user booking
type Booking struct {
	gorm.Model
	UserID            uint      `json:"user_id"`
	BookingType       string    `json:"booking_type"` // flight, hotel, package
	FlightID          *uint     `json:"flight_id,omitempty"`
	HotelID           *uint     `json:"hotel_id,omitempty"`
	VacationPackageID *uint     `json:"vacation_package_id,omitempty"`
	Status            string    `json:"status"` // confirmed, cancelled, pending
	BookingDate       time.Time `json:"booking_date"`
	TotalPrice        float64   `json:"total_price"`
	PaymentStatus     string    `json:"payment_status"`
	CheckInDate       time.Time `json:"check_in_date,omitempty"`
	CheckOutDate      time.Time `json:"check_out_date,omitempty"`
	NumGuests         int       `json:"num_guests"`
	SpecialRequests   string    `json:"special_requests"`
}

// Search request structs
type FlightSearchRequest struct {
	DepartureCity string     `json:"departure_city" binding:"required"`
	ArrivalCity   string     `json:"arrival_city" binding:"required"`
	DepartureDate time.Time  `json:"departure_date" binding:"required"`
	ReturnDate    *time.Time `json:"return_date"`
	Passengers    int        `json:"passengers" binding:"required,min=1"`
	Class         string     `json:"class" binding:"required"`
}

type HotelSearchRequest struct {
	City      string    `json:"city" binding:"required"`
	CheckIn   time.Time `json:"check_in" binding:"required"`
	CheckOut  time.Time `json:"check_out" binding:"required"`
	Guests    int       `json:"guests" binding:"required,min=1"`
	RoomType  string    `json:"room_type"`
	MaxPrice  *float64  `json:"max_price"`
	MinRating *float32  `json:"min_rating"`
}

type BookingRequest struct {
	BookingType       string     `json:"booking_type" binding:"required"`
	FlightID          *uint      `json:"flight_id"`
	HotelID           *uint      `json:"hotel_id"`
	VacationPackageID *uint      `json:"vacation_package_id"`
	CheckInDate       *time.Time `json:"check_in_date"`
	CheckOutDate      *time.Time `json:"check_out_date"`
	NumGuests         int        `json:"num_guests" binding:"required,min=1"`
	SpecialRequests   string     `json:"special_requests"`
}
