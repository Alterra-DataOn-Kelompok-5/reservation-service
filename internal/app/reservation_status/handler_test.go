package locations

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
	testEmail                 = "vincentlhubbard@superrito.com"
	testEmployeeID            = uint(1)
	testAdminRoleID           = uint(enum.Admin)
	testUserRoleID            = uint(enum.User)
	testDivisionID            = uint(enum.Finance)
	adminClaims               = util.CreateJWTClaims(testEmail, testEmployeeID, testAdminRoleID, testDivisionID)
	userClaims                = util.CreateJWTClaims(testEmail, testEmployeeID, testUserRoleID, testDivisionID)
	db                        = database.GetConnection()
	reservationStatusHandler  = NewHandler(&f)
	f                         = factory.Factory{ReservationStatusRepository: repository.NewStatusRepository(db)}
	testCreatePayload         = dto.CreateReservationStatusRequestBody{ReservationStatusName: &testReservationStatusName}
	echoMock                  = mocks.EchoMock{E: echo.New()}
	testReservationStatusID   = uint(enum.Reserved)
	testReservationStatusName = enum.Status(testReservationStatusID).String()
	testUpdatePayload         = dto.UpdateReservationStatusRequestBody{ID: &testReservationStatusID, ReservationStatusName: &testReservationStatusName}
)

func TestReservationStatusHandlerGetInvalidPayload(t *testing.T) {
	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/api/v1/status")
	c.QueryParams().Add("page", "a")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.Get(c)) {
		asserts.Equal(400, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Bad Request")
	}
}
func TestReservationStatusHandlerGetUnauthorized(t *testing.T) {
	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	c.SetPath("/api/v1/status")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.Get(c)) {
		asserts.Equal(401, rec.Code)
		body := rec.Body.String()
		asserts.Contains(body, "unauthorized")
	}
}

func TestReservationStatusHandlerGetSuccess(t *testing.T) {
	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	c.SetPath("/api/v1/status")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.Get(c)) {
		asserts.Equal(200, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "id")
		asserts.Contains(body, "reservation_status_name")
	}
}

func TestReservationStatusHandlerGetByIdInvalidPayload(t *testing.T) {
	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	reservatoinStatusID := "a"

	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/api/v1/status")
	c.SetParamNames("id")
	c.SetParamValues(reservatoinStatusID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.GetById(c)) {
		asserts.Equal(400, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Bad Request")
	}
}

func TestReservationStatusHandlerGetByIdNotFound(t *testing.T) {
	seeder.NewSeeder().DeleteAll()

	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	reservatoinStatusID := strconv.Itoa(int(testReservationStatusID))
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/api/v1/status")
	c.SetParamNames("id")
	c.SetParamValues(reservatoinStatusID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.GetById(c)) {
		asserts.Equal(404, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Data not found")
	}
}

func TestReservationStatusHandlerGetByIdUnauthorized(t *testing.T) {
	seeder.NewSeeder().DeleteAll()

	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	reservatoinStatusID := strconv.Itoa(int(testReservationStatusID))

	c.SetPath("/api/v1/status")
	c.SetParamNames("id")
	c.SetParamValues(reservatoinStatusID)

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.GetById(c)) {
		asserts.Equal(401, rec.Code)
		body := rec.Body.String()
		asserts.Contains(body, "unauthorized")
	}
}

func TestReservationStatusHandlerGetByIdSuccess(t *testing.T) {
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	c, rec := echoMock.RequestMock(http.MethodGet, "/", nil)
	reservationStatusID := strconv.Itoa(int(testReservationStatusID))
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/api/v1/status")
	c.SetParamNames("id")
	c.SetParamValues(reservationStatusID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.GetById(c)) {
		asserts.Equal(200, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "id")
		asserts.Contains(body, "name")
	}
}

func TestReservationStatusHandlerUpdateByIdInvalidPayload(t *testing.T) {
	c, rec := echoMock.RequestMock(http.MethodPut, "/", nil)
	reservationStatusID := "a"
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/api/v1/status")
	c.SetParamNames("id")
	c.SetParamValues(reservationStatusID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.UpdateById(c)) {
		asserts.Equal(400, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Bad Request")
	}
}

func TestReservationStatusHandlerUpdateByIdNotFound(t *testing.T) {
	seeder.NewSeeder().DeleteAll()

	payload, err := json.Marshal(testUpdatePayload)
	if err != nil {
		t.Fatal(err)
	}
	c, rec := echoMock.RequestMock(http.MethodPut, "/", bytes.NewBuffer(payload))
	reservationStatusID := strconv.Itoa(int(testReservationStatusID))
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/status")
	c.SetParamNames("id")
	c.SetParamValues(reservationStatusID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Request().Header.Set("Content-Type", "application/json")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.UpdateById(c)) {
		asserts.Equal(404, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Data not found")
	}
}
func TestReservationStatusHandlerUpdateByIdUnauthorized(t *testing.T) {
	seeder.NewSeeder().DeleteAll()

	c, rec := echoMock.RequestMock(http.MethodPut, "/", nil)
	reservationStatusID := strconv.Itoa(int(testReservationStatusID))

	c.SetPath("/api/v1/status")
	c.SetParamNames("id")
	c.SetParamValues(reservationStatusID)

	token, err := util.CreateJWTToken(userClaims)
	if err != nil {
		t.Fatal(err)
	}
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.UpdateById(c)) {
		asserts.Equal(401, rec.Code)
		body := rec.Body.String()
		asserts.Contains(body, "unauthorized")
	}
}

func TestReservationStatusHandlerUpdateByIdSuccess(t *testing.T) {
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	payload, err := json.Marshal(testUpdatePayload)
	if err != nil {
		t.Fatal(err)
	}
	c, rec := echoMock.RequestMock(http.MethodPut, "/", bytes.NewBuffer(payload))
	reservationStatusID := strconv.Itoa(int(testReservationStatusID))
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}

	c.SetPath("/status")
	c.SetParamNames("id")
	c.SetParamValues(reservationStatusID)
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Request().Header.Set("Content-Type", "application/json")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.UpdateById(c)) {
		asserts.Equal(200, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "id")
		asserts.Contains(body, "name")
	}
}

func TestReservationStatusHandlerCreateInvalidPayload(t *testing.T) {
	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(dto.CreateReservationStatusRequestBody{})
	if err != nil {
		t.Fatal(err)
	}

	c, rec := echoMock.RequestMock(http.MethodPost, "/", bytes.NewBuffer(payload))

	c.SetPath("/api/v1/status")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Request().Header.Set("Content-Type", "application/json")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.Create(c)) {
		asserts.Equal(400, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "Invalid parameters or payload")
	}
}

func TestReservationStatusHandlerCreateDivisionAlreadyExist(t *testing.T) {
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()

	token, err := util.CreateJWTToken(adminClaims)
	if err != nil {
		t.Fatal(err)
	}
	name := enum.Status(testReservationStatusID).String()
	payload, err := json.Marshal(dto.CreateReservationStatusRequestBody{ReservationStatusName: &name})
	if err != nil {
		t.Fatal(err)
	}

	c, rec := echoMock.RequestMock(http.MethodPost, "/", bytes.NewBuffer(payload))

	c.SetPath("/api/v1/status")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Request().Header.Set("Content-Type", "application/json")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.Create(c)) {
		asserts.Equal(409, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "duplicate")
	}
}

func TestReservationStatusHandlerCreateUnauthorized(t *testing.T) {
	payload, err := json.Marshal(testCreatePayload)
	if err != nil {
		t.Fatal(err)
	}

	c, rec := echoMock.RequestMock(http.MethodPost, "/", bytes.NewBuffer(payload))

	c.SetPath("/api/v1/status")
	c.Request().Header.Set("Content-Type", "application/json")

	token, err := util.CreateJWTToken(userClaims)
	if err != nil {
		t.Fatal(err)
	}
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.Create(c)) {
		asserts.Equal(401, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "unauthorized")
	}
}

func TestReservationStatusHandlerCreateSuccess(t *testing.T) {
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

	c.SetPath("/api/v1/status")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Request().Header.Set("Content-Type", "application/json")

	// testing
	asserts := assert.New(t)
	if asserts.NoError(reservationStatusHandler.Create(c)) {
		asserts.Equal(200, rec.Code)

		body := rec.Body.String()
		asserts.Contains(body, "id")
		asserts.Contains(body, "name")
	}
}
