package repository

import (
	"context"
	"fledge-restapi/internal/domain/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Generic repository interface
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*T, error)
}

// Base repository implementation
type baseRepository[T any] struct {
	db *gorm.DB
}

func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *baseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *baseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

func (r *baseRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// Flight Repository
type FlightRepository interface {
	Repository[entity.Flight]
	Search(ctx context.Context, params FlightSearchParams) ([]entity.Flight, error)
	FindAll(ctx context.Context) ([]entity.Flight, error)
	FindByOrigin(ctx context.Context, origin string) ([]entity.Flight, error)
}

type FlightSearchParams struct {
	DepartureCity string
	ArrivalCity   string
	DepartureDate time.Time
	ReturnDate    *time.Time
	Passengers    int
	Class         string
}

type flightRepository struct {
	baseRepository[entity.Flight]
}

func NewFlightRepository(db *gorm.DB) FlightRepository {
	return &flightRepository{baseRepository[entity.Flight]{db: db}}
}

func (r *flightRepository) Search(ctx context.Context, params FlightSearchParams) ([]entity.Flight, error) {
	var flights []entity.Flight
	query := r.db.WithContext(ctx).
		Where("departure_city = ? AND arrival_city = ?", params.DepartureCity, params.ArrivalCity).
		Where("departure_time >= ? AND departure_time <= ?",
			params.DepartureDate, params.DepartureDate.Add(24*time.Hour)).
		Where("available_seats >= ?", params.Passengers).
		Where("class = ?", params.Class)

	if err := query.Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}

func (r *flightRepository) FindAll(ctx context.Context) ([]entity.Flight, error) {
	var flights []entity.Flight
	if err := r.db.WithContext(ctx).Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}

func (r *flightRepository) FindByOrigin(ctx context.Context, origin string) ([]entity.Flight, error) {
	var flights []entity.Flight
	if err := r.db.WithContext(ctx).
		Where("departure_city = ?", origin).
		Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}

// Hotel Repository
type HotelRepository interface {
	Repository[entity.Hotel]
	Search(ctx context.Context, params HotelSearchParams) ([]entity.Hotel, error)
}

type HotelSearchParams struct {
	City      string
	CheckIn   time.Time
	CheckOut  time.Time
	Guests    int
	RoomType  string
	MaxPrice  *float64
	MinRating *float32
}

type hotelRepository struct {
	baseRepository[entity.Hotel]
}

func NewHotelRepository(db *gorm.DB) HotelRepository {
	return &hotelRepository{baseRepository[entity.Hotel]{db: db}}
}

func (r *hotelRepository) Search(ctx context.Context, params HotelSearchParams) ([]entity.Hotel, error) {
	var hotels []entity.Hotel
	query := r.db.WithContext(ctx).
		Where("city = ?", params.City).
		Where("available_rooms > 0")

	if params.MaxPrice != nil {
		query = query.Where("price_per_night <= ?", *params.MaxPrice)
	}
	if params.MinRating != nil {
		query = query.Where("rating >= ?", *params.MinRating)
	}

	if err := query.Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}

// Booking Repository
type BookingRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Booking, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error)
	Create(ctx context.Context, booking *entity.Booking) error
	Update(ctx context.Context, id uint, updates map[string]interface{}) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) FindByID(ctx context.Context, id uint) (*entity.Booking, error) {
	var booking entity.Booking
	if err := r.db.WithContext(ctx).First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error) {
	var bookings []entity.Booking
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) Create(ctx context.Context, booking *entity.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&entity.Booking{}).Where("id = ?", id).Updates(updates).Error
}
