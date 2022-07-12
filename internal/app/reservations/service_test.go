package reservations

import (
	"context"
	"fmt"
	"testing"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database/seeder"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/dto"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/factory"
	pkgdto "github.com/Alterra-DataOn-Kelompok-5/reservation-service/pkg/dto"
	"github.com/stretchr/testify/assert"
)

var (
	ctx                = context.Background()
	reservationService = NewService(factory.NewFactory())
)

func TestReservationServiceFindAllSuccess(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	var (
		asserts = assert.New(t)
		payload = pkgdto.SearchGetRequest{}
	)

	res, err := reservationService.Find(ctx, &payload)
	if err != nil {
		t.Fatal(err)
	}

	asserts.Len(res.Data, 3)
	for _, val := range res.Data {
		asserts.NotEmpty(val.ReservationCode)
		asserts.NotEmpty(val.ID)
	}
}
func TestReservationServiceFindByIdSuccess(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	var (
		asserts = assert.New(t)
		payload = pkgdto.ByIDRequest{ID: 1}
	)

	res, err := reservationService.FindByID(ctx, &payload)
	if err != nil {
		t.Fatal(err)
	}

	asserts.Equal(uint(1), res.ID)
}

func TestReservationServiceFindByIdRecordNotFound(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()

	var (
		asserts = assert.New(t)
		payload = pkgdto.ByIDRequest{ID: 1}
	)

	_, err := reservationService.FindByID(ctx, &payload)
	if err != nil {
		asserts.Equal(err.Error(), "error code 404")
	}
}

func TestReservationServiceUpdateByIdSuccess(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	var (
		asserts = assert.New(t)
	)
	res, err := reservationService.UpdateById(ctx, &testUpdatePayload)
	fmt.Println("res stat id", res.ReservationStatusID)
	fmt.Println("res", res)
	if err != nil {
		t.Fatal(err)
	}
	asserts.Equal(testReservationStatusID, res.ReservationStatusID)
}

func TestReservationServiceUpdateByIdRecordNotFound(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()

	var (
		asserts = assert.New(t)
		id      = uint(10)
		status  = uint(2)
		payload = dto.UpdateReservationRequestBody{
			ID:                  &id,
			ReservationStatusID: &status,
		}
	)

	_, err := reservationService.UpdateById(ctx, &payload)
	if err != nil {
		asserts.Equal(err.Error(), "error code 404")
	}
}

func TestReservationServiceCreateReservationSuccess(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedReservationStatus()

	asserts := assert.New(t)
	res, err := reservationService.Store(ctx, &testCreatePayload)
	fmt.Println("res:", res)
	if err != nil {
		t.Fatal(err)
	}
	asserts.NotEmpty(res.ID)
	asserts.Equal(testCreatePayload.ReservationCode, res.ReservationCode)
}

func TestReservationServiceCreateReservationAlreadyExist(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	var (
		asserts = assert.New(t)
		code    = "RSVN/20220710/001"
		payload = dto.CreateReservationRequestBody{
			ReservationCode: code,
		}
	)

	_, err := reservationService.Store(ctx, &payload)
	if asserts.Error(err) {
		asserts.Equal(err.Error(), "error code 409")
	}
}
