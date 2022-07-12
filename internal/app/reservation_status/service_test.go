package locations

import (
	"context"
	"testing"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database/seeder"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/factory"
	pkgdto "github.com/Alterra-DataOn-Kelompok-5/reservation-service/pkg/dto"
	"github.com/stretchr/testify/assert"
)

var (
	ctx                      = context.Background()
	reservationStatusService = NewService(factory.NewFactory())
	testFindAllPayload       = pkgdto.SearchGetRequest{}
	testFindByIdPayload      = pkgdto.ByIDRequest{ID: 1}
)

func TestReservationStatusServiceFindAllSuccess(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	asserts := assert.New(t)
	res, err := reservationStatusService.Find(ctx, &testFindAllPayload)
	if err != nil {
		t.Fatal(err)
	}

	asserts.Len(res.Data, 4)
	for _, val := range res.Data {
		asserts.NotEmpty(val.ReservationStatusName)
		asserts.NotEmpty(val.ID)
	}
}
func TestReservationStatusServiceFindByIdSuccess(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	asserts := assert.New(t)
	res, err := reservationStatusService.FindByID(ctx, &testFindByIdPayload)
	if err != nil {
		t.Fatal(err)
	}

	asserts.Equal(uint(1), res.ID)
}

func TestReservationStatusServiceFindByIdRecordNotFound(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()

	asserts := assert.New(t)
	_, err := reservationStatusService.FindByID(ctx, &testFindByIdPayload)
	if err != nil {
		asserts.Equal(err.Error(), "error code 404")
	}
}

func TestReservationStatusServiceUpdataByIdSuccess(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	asserts := assert.New(t)
	res, err := reservationStatusService.UpdateById(ctx, &testUpdatePayload)
	if err != nil {
		t.Fatal(err)
	}
	asserts.Equal(testReservationStatusName, res.ReservationStatusName)
}

func TestReservationStatusServiceUpdateByIdRecordNotFound(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()

	asserts := assert.New(t)
	_, err := reservationStatusService.UpdateById(ctx, &testUpdatePayload)
	if err != nil {
		asserts.Equal(err.Error(), "error code 404")
	}
}

func TestReservationStatusServiceCreateStatusSuccess(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()

	asserts := assert.New(t)
	res, err := reservationStatusService.Store(ctx, &testCreatePayload)
	if err != nil {
		t.Fatal(err)
	}
	asserts.NotEmpty(res.ID)
	asserts.Equal(*testCreatePayload.ReservationStatusName, res.ReservationStatusName)
}

func TestReservationStatusServiceCreateStatusAlreadyExist(t *testing.T) {
	database.GetConnection()
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	asserts := assert.New(t)
	_, err := reservationStatusService.Store(ctx, &testCreatePayload)
	if asserts.Error(err) {
		asserts.Equal(err.Error(), "error code 409")
	}
}
