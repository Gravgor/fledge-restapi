package service

import (
	"context"
	"errors"
	"time"

	"fledge-restapi/internal/domain/entity"
	"fledge-restapi/internal/domain/repository"
)

type BookingService struct {
	bookingRepo repository.BookingRepository
}

func NewBookingService(bookingRepo repository.BookingRepository) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
	}
}

func (s *BookingService) ListBookings(ctx context.Context, userID uint) ([]entity.Booking, error) {
	return s.bookingRepo.FindByUserID(ctx, userID)
}

func (s *BookingService) GetBooking(ctx context.Context, id uint, userID uint) (*entity.Booking, error) {
	booking, err := s.bookingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if booking.UserID != userID {
		return nil, errors.New("unauthorized access to booking")
	}

	return booking, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, id uint, userID uint, updates map[string]interface{}) error {
	booking, err := s.GetBooking(ctx, id, userID)
	if err != nil {
		return err
	}

	// Validate that the booking can be updated
	if booking.Status == "cancelled" {
		return errors.New("cannot update cancelled booking")
	}

	return s.bookingRepo.Update(ctx, id, updates)
}

func (s *BookingService) CancelBooking(ctx context.Context, id uint, userID uint) error {
	booking, err := s.GetBooking(ctx, id, userID)
	if err != nil {
		return err
	}

	// Check if booking is within cancellation period
	if time.Until(booking.CheckInDate) < 24*time.Hour {
		return errors.New("booking cannot be cancelled within 24 hours of start date")
	}

	updates := map[string]interface{}{
		"status": "cancelled",
	}

	return s.bookingRepo.Update(ctx, id, updates)
}
