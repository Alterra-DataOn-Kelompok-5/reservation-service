package reservations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database/seeder"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/dto"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/factory"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/mocks"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/pkg/enum"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/pkg/util"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	adminClaims             = util.CreateJWTClaims(testEmail, testEmployeeID, testAdminRoleID, testDivisionID)
	userClaims              = util.CreateJWTClaims(testEmail, testEmployeeID, testUserRoleID, testDivisionID)
	testEmail               = "vincentlhubbard@superrito.com"
	testEmployeeID          = uint(1)
	testAdminRoleID         = uint(enum.Admin)
	testUserRoleID          = uint(enum.User)
	testDivisionID          = uint(enum.Finance)
	db                      = database.GetConnection()
	echoMock                = mocks.EchoMock{E: echo.New()}
	f                       = factory.Factory{ReservationsRepository: repository.NewReservationsRepository(db)}
	reservationHandler      = NewHandler(&f)
	testReservationID       = uint(3)
	testReservationCode     = "RSVN/20220709/005"
	testReservationDesc     = "reservation test"
	testReservationStatusID = uint(1)
	testRoomID              = uint(2)
	testTimeStart           = "2022-05-10 07:00:00"
	testTimeEnd             = "2022-05-10 08:00:00"
	testCreatePayload       = dto.CreateReservationRequestBody{
		ReservationCode:      &testReservationCode,
		ReservationDesc:      &testReservationDesc,
		EmployeeID:           &testEmployeeID,
		RoomID:               &testRoomID,
		ReservationTimeStart: &testTimeStart,
		ReservationTimeEnd:   &testTimeEnd,
	}
	testUpdatePayload = dto.UpdateReservationRequestBody{ID: &testReservationID, AdminID: &testEmployeeID, ReservationStatusID: &testReservationStatusID}
)

func TestReservationsHandlerGetInvalidPayload(t *testing.T) {
	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/api/v1/reservations")
	c.QueryParams().Add("page", "a")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.Get(c)) {
		asserts.Equal(400, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Bad Request")
	}
}

func TestReservationsHandlerGetUnauthorized(t *testing.T) {
	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	c.SetPath("/api/v1/reservations")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.Get(c)) {
		asserts.Equal(401, rec.Code)
		body := rec.Body.String()
		asserts.Contains(body, "unauthorized")
	}
}

func TestReservationsHandlerGetSuccess(t *testing.T) {
	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	c.SetPath("/api/v1/reservations")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.Get(c)) {
		asserts.Equal(200, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "id")
		asserts.Contains(body, "reservation_code")
	}
}

func TestReservationsHandlerGetByIdInvalidPayload(t *testing.T) {
	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	roleID := "a"

	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/api/v1/reservations")
	c.SetParamNames("id")
	c.SetParamValues(roleID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.GetById(c)) {
		asserts.Equal(400, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Bad Request")
	}
}

func TestReservationsHandlerGetByIdNotFound(t *testing.T) {
	seeder.NewSeeder().DeleteAll()

	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	roleID := strconv.Itoa(int(testAdminRoleID))
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/api/v1/reservations")
	c.SetParamNames("id")
	c.SetParamValues(roleID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.GetById(c)) {
		asserts.Equal(404, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Data not found")
	}
}

func TestReservationsHandlerGetByIdUnauthorized(t *testing.T) {
	seeder.NewSeeder().DeleteAll()

	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	reservationID := strconv.Itoa(int(testReservationID))

	c.SetPath("/api/v1/reservations")
	c.SetParamNames("id")
	c.SetParamValues(reservationID)

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.GetById(c)) {
		asserts.Equal(401, rec.Code)
		body := rec.Body.String()
		asserts.Contains(body, "unauthorized")
	}
}

func TestReservationsHandlerGetByIdSuccess(t *testing.T) {
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	roleID := strconv.Itoa(int(testAdminRoleID))
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/api/v1/reservations")
	c.SetParamNames("id")
	c.SetParamValues(roleID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.GetById(c)) {
		asserts.Equal(200, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "id")
		asserts.Contains(body, "name")
	}
}

func TestReservationsHandlerUpdateByIdInvalidPayload(t *testing.T) {
	c, rec := echoMock.RequestMock(http.MethodPut, "/", nil)
	roleID := "a"
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/api/v1/reservations")
	c.SetParamNames("id")
	c.SetParamValues(roleID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.UpdateById(c)) {
		asserts.Equal(400, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Bad Request")
	}
}

func TestReservationsHandlerUpdateByIdNotFound(t *testing.T) {
	seeder.NewSeeder().DeleteAll()

	payload, err := json.Marshal(testUpdatePayload)
	if err != nil {
		t.Fatal(err)
	}
	c, rec := echoMock.RequestMock(http.MethodPut, "/", bytes.NewBuffer(payload))
	roleID := strconv.Itoa(int(testAdminRoleID))
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/reservations")
	c.SetParamNames("id")
	c.SetParamValues(roleID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Request().Header.Set("Content-Type", "application/json")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.UpdateById(c)) {
		asserts.Equal(404, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Data not found")
	}
}
func TestReservationsHandlerUpdateByIdUnauthorized(t *testing.T) {
	token, err := util.CreateJWTToken(userClaims)
	if err != nil {
		t.Fatal(err)
	}
	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	c.SetPath("/api/v1/reservations")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	c.SetPath("/api/v1/reservations")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(testReservationID)))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.UpdateById(c)) {
		asserts.Equal(401, rec.Code)
		body := rec.Body.String()
		asserts.Contains(body, "unauthorized")
	}
}

func TestReservationsHandlerUpdateByIdSuccess(t *testing.T) {
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	payload, err := json.Marshal(testUpdatePayload)
	if err != nil {
		t.Fatal(err)
	}
	c, rec := echoMock.RequestMock(http.MethodPut, "/", bytes.NewBuffer(payload))
	roleID := strconv.Itoa(int(testAdminRoleID))
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/reservations")
	c.SetParamNames("id")
	c.SetParamValues(roleID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Request().Header.Set("Content-Type", "application/json")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.UpdateById(c)) {
		asserts.Equal(200, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "id")
		asserts.Contains(body, "reservation_code")
	}
}

func TestReservationsHandlerCreateInvalidPayload(t *testing.T) {
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(dto.CreateReservationRequestBody{})
	if err != nil {
		t.Fatal(err)
	}

	c, rec := echoMock.RequestMock(http.MethodPost, "/", bytes.NewBuffer(payload))

	c.SetPath("/api/v1/reservations")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Request().Header.Set("Content-Type", "application/json")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.Create(c)) {
		asserts.Equal(400, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Invalid parameters or payload")
	}
}

func TestReservationsHandlerCreateReservationAlreadyExist(t *testing.T) {
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}
	code := "RSVN/20220710/003"
	payload, err := json.Marshal(dto.CreateReservationRequestBody{ReservationCode: &code})
	if err != nil {
		t.Fatal(err)
	}

	c, rec := echoMock.RequestMock(http.MethodPost, "/", bytes.NewBuffer(payload))

	c.SetPath("/api/v1/reservations")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Request().Header.Set("Content-Type", "application/json")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.Create(c)) {
		asserts.Equal(409, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "duplicate")
	}
}

func TestReservationsHandlerCreateSuccess(t *testing.T) {
	seeder.NewSeeder().DeleteAll()

	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(testCreatePayload)
	if err != nil {
		t.Fatal(err)
	}

	c, rec := echoMock.RequestMock(http.MethodPost, "/", bytes.NewBuffer(payload))

	c.SetPath("/api/v1/reservations")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Request().Header.Set("Content-Type", "application/json")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationHandler.Create(c)) {
		asserts.Equal(200, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "id")
		asserts.Contains(body, "reservation_code")
	}
}
