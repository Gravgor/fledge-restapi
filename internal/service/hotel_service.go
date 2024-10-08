package service

import (
	"context"
	"fledge-restapi/internal/domain/entity"
	"fledge-restapi/internal/domain/repository"
	"fledge-restapi/pkg/errors"
	"time"

	"github.com/google/uuid"
)

type HotelService interface {
	SearchHotels(ctx context.Context, req *entity.HotelSearchRequest) ([]entity.Hotel, error)
	GetHotelByID(ctx context.Context, id uint) (*entity.Hotel, error)
	BookHotel(ctx context.Context, userID uuid.UUID, bookingReq *entity.BookingRequest) (*entity.Booking, error)
}

type hotelService struct {
	hotelRepo   repository.HotelRepository
	bookingRepo repository.BookingRepository
}

func NewHotelService(hotelRepo repository.HotelRepository, bookingRepo repository.BookingRepository) HotelService {
	return &hotelService{
		hotelRepo:   hotelRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *hotelService) SearchHotels(ctx context.Context, req *entity.HotelSearchRequest) ([]entity.Hotel, error) {
	// Validate dates
	if req.CheckIn.Before(time.Now()) {
		return nil, errors.ErrInvalidCheckInDate
	}

	if req.CheckOut.Before(req.CheckIn) {
		return nil, errors.ErrInvalidCheckOutDate
	}

	// Search for hotels
	hotels, err := s.hotelRepo.Search(ctx, repository.HotelSearchParams{
		City:      req.City,
		CheckIn:   req.CheckIn,
		CheckOut:  req.CheckOut,
		Guests:    req.Guests,
		RoomType:  req.RoomType,
		MaxPrice:  req.MaxPrice,
		MinRating: req.MinRating,
	})

	if err != nil {
		return nil, err
	}

	return hotels, nil
}

func (s *hotelService) GetHotelByID(ctx context.Context, id uint) (*entity.Hotel, error) {
	return s.hotelRepo.FindByID(ctx, id)
}

func (s *hotelService) BookHotel(ctx context.Context, userID uuid.UUID, bookingReq *entity.BookingRequest) (*entity.Booking, error) {
	// Validate hotel exists and has availability
	hotel, err := s.hotelRepo.FindByID(ctx, *bookingReq.HotelID)
	if err != nil {
		return nil, err
	}

	if hotel.AvailableRooms < 1 {
		return nil, errors.ErrNoRoomsAvailable
	}

	// Calculate total nights
	nights := int(bookingReq.CheckOutDate.Sub(*bookingReq.CheckInDate).Hours() / 24)
	if nights < 1 {
		return nil, errors.ErrInvalidStayDuration
	}

	// Create booking
	booking := &entity.Booking{
		UserID:          userID,
		BookingType:     "hotel",
		HotelID:         bookingReq.HotelID,
		Status:          "confirmed",
		BookingDate:     time.Now(),
		TotalPrice:      float64(nights) * hotel.Price,
		PaymentStatus:   "pending",
		CheckInDate:     *bookingReq.CheckInDate,
		CheckOutDate:    *bookingReq.CheckOutDate,
		NumGuests:       bookingReq.NumGuests,
		SpecialRequests: bookingReq.SpecialRequests,
	}

	// Save booking
	err = s.bookingRepo.Create(ctx, booking)
	if err != nil {
		return nil, err
	}

	// Update hotel availability
	hotel.AvailableRooms--
	err = s.hotelRepo.Update(ctx, hotel)
	if err != nil {
		return nil, err
	}

	return booking, nil
}
