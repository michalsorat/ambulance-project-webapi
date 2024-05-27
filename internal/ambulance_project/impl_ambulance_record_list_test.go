package ambulance_project

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/michalsorat/ambulance-project-webapi/internal/db_service"
)

type AmbulanceRLSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[Ambulance]
}

func TestAmbulanceRLSuite(t *testing.T) {
	suite.Run(t, new(AmbulanceRLSuite))
}

type DbServiceMock[DocType interface{}] struct {
	mock.Mock
}

func (m *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := m.Called(ctx, id, document)
	return args.Error(0)
}

func (m *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*DocType), args.Error(1)
}

func (m *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := m.Called(ctx, id, document)
	return args.Error(0)
}

func (m *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (suite *AmbulanceRLSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[Ambulance]{}

	// Compile time assert that the mock is of type db_service.DbService[Ambulance]
	var _ db_service.DbService[Ambulance] = suite.dbServiceMock

	consumationTime, _ := time.Parse(time.RFC3339, "2024-05-26T12:00:00Z")
	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Ambulance{
				Id: "test-ambulance",
				MealOrders: []MealOrder{
					{
						Id:              "test-entry",
						Name:            "test-patient",
						DietaryReq:      "Vegetarian",
						MedicalNeed:     "Diabetes",
						ConsumationTime: consumationTime,
					},
				},
			},
			nil,
		)
}

func (suite *AmbulanceRLSuite) Test_UpdateRL_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
		"id": "test-entry",
		"name": "Jano",
		"dietaryReq": "Vegetarian",
		"medicalNeed": "Diabetes",
		"consumationTime": "2024-05-26T12:00:00Z"
	}`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "orderId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/meal-orders/test-ambulance/records/test-entry", strings.NewReader(json))
	ctx.Request.Header.Set("Content-Type", "application/json")

	sut := implMealOrdersAPI{}

	// ACT
	sut.UpdateMealOrder(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-ambulance", mock.Anything)
	suite.Equal(200, recorder.Code)
}

func (suite *AmbulanceRLSuite) Test_UpdateRL_InvalidInput() {
	// ARRANGE
	json := `{
		"name": "Jano"
		"medicalNeed": "Diabetes",
		"consumationTime": "2024-05-26T12:00:00Z"
	}` // Missing comma and dietaryReq

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "orderId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/meal-orders/test-ambulance/records/test-entry", strings.NewReader(json))
	ctx.Request.Header.Set("Content-Type", "application/json")

	sut := implMealOrdersAPI{}

	// ACT
	sut.UpdateMealOrder(ctx)

	// ASSERT
	suite.Equal(400, recorder.Code)
}

func (suite *AmbulanceRLSuite) Test_UpdateRL_NotFound() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(nil, nil)

	json := `{
		"id": "non-existent-entry",
		"name": "Jano",
		"dietaryReq": "Vegetarian",
		"medicalNeed": "Diabetes",
		"consumationTime": "2024-05-26T12:00:00Z"
	}`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "orderId", Value: "non-existent-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/meal-orders/test-ambulance/records/non-existent-entry", strings.NewReader(json))
	ctx.Request.Header.Set("Content-Type", "application/json")

	sut := implMealOrdersAPI{}

	// ACT
	sut.UpdateMealOrder(ctx)

	// ASSERT
	suite.Equal(404, recorder.Code)
}

func (suite *AmbulanceRLSuite) Test_UpdateRL_Success() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
		"id": "test-entry",
		"name": "Jano",
		"dietaryReq": "Vegetarian",
		"medicalNeed": "Diabetes",
		"consumationTime": "2024-05-26T12:00:00Z"
	}`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "orderId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/meal-orders/test-ambulance/records/test-entry", strings.NewReader(json))
	ctx.Request.Header.Set("Content-Type", "application/json")

	sut := implMealOrdersAPI{}

	// ACT
	sut.UpdateMealOrder(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-ambulance", mock.Anything)
	suite.Equal(200, recorder.Code)
}

