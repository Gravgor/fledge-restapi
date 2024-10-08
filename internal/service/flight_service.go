package service

import (
	"context"
	"fledge-restapi/internal/domain/entity"
	"fledge-restapi/internal/domain/repository"
	"fledge-restapi/pkg/errors"
	"time"

	"github.com/google/uuid"
)

type FlightService interface {
	SearchFlights(ctx context.Context, req *entity.FlightSearchRequest) ([]entity.Flight, error)
	GetFlightByID(ctx context.Context, id uint) (*entity.Flight, error)
	BookFlight(ctx context.Context, userID uuid.UUID, bookingReq *entity.BookingRequest) (*entity.Booking, error)
	ListAllFlights(ctx context.Context) ([]entity.Flight, error)
	ListFlightsByOrigin(ctx context.Context, origin string) ([]entity.Flight, error)
}

type flightService struct {
	flightRepo  repository.FlightRepository
	bookingRepo repository.BookingRepository
}

func NewFlightService(flightRepo repository.FlightRepository, bookingRepo repository.BookingRepository) FlightService {
	return &flightService{
		flightRepo:  flightRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *flightService) SearchFlights(ctx context.Context, req *entity.FlightSearchRequest) ([]entity.Flight, error) {
	// Validate search criteria
	if req.DepartureDate.Before(time.Now()) {
		return nil, errors.ErrInvalidDepartureDate
	}

	if req.ReturnDate != nil && req.ReturnDate.Before(req.DepartureDate) {
		return nil, errors.ErrInvalidReturnDate
	}

	// Search for flights
	flights, err := s.flightRepo.Search(ctx, repository.FlightSearchParams{
		DepartureCity: req.DepartureCity,
		ArrivalCity:   req.ArrivalCity,
		DepartureDate: req.DepartureDate,
		Passengers:    req.Passengers,
		Class:         req.Class,
	})

	if err != nil {
		return nil, err
	}

	return flights, nil
}

func (s *flightService) GetFlightByID(ctx context.Context, id uint) (*entity.Flight, error) {
	return s.flightRepo.FindByID(ctx, id)
}

func (s *flightService) ListAllFlights(ctx context.Context) ([]entity.Flight, error) {
	return s.flightRepo.FindAll(ctx)
}

func (s *flightService) ListFlightsByOrigin(ctx context.Context, origin string) ([]entity.Flight, error) {
	// Get all flights filtered by origin
	flights, err := s.flightRepo.FindByOrigin(ctx, origin)
	if err != nil {
		return nil, err
	}

	return flights, nil
}

func (s *flightService) BookFlight(ctx context.Context, userID uuid.UUID, bookingReq *entity.BookingRequest) (*entity.Booking, error) {
	// Validate flight exists and has availability
	flight, err := s.flightRepo.FindByID(ctx, *bookingReq.FlightID)
	if err != nil {
		return nil, err
	}

	if flight.AvailableSeats < bookingReq.NumGuests {
		return nil, errors.ErrInsufficientSeats
	}

	// Create booking
	booking := &entity.Booking{
		UserID:          userID,
		BookingType:     "flight",
		FlightID:        bookingReq.FlightID,
		Status:          "confirmed",
		BookingDate:     time.Now(),
		TotalPrice:      float64(bookingReq.NumGuests) * flight.Price,
		PaymentStatus:   "pending",
		NumGuests:       bookingReq.NumGuests,
		SpecialRequests: bookingReq.SpecialRequests,
	}

	// Save booking
	err = s.bookingRepo.Create(ctx, booking)
	if err != nil {
		return nil, err
	}

	// Update flight availability
	flight.AvailableSeats -= bookingReq.NumGuests
	err = s.flightRepo.Update(ctx, flight)
	if err != nil {
		return nil, err
	}

	return booking, nil
}
