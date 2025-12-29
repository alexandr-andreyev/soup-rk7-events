package transport

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
)

// MockNotifyService is a mock implementation of NotifyEventService
type MockNotifyService struct {
	handleFunc func(*models.Rk7NotifyEvent) error
}

func (m *MockNotifyService) HandleNotification(data *models.Rk7NotifyEvent) error {
	if m.handleFunc != nil {
		return m.handleFunc(data)
	}
	return nil
}

func TestHandleEvents_Success(t *testing.T) {
	// Arrange
	var capturedEvent *models.Rk7NotifyEvent
	mockService := &MockNotifyService{
		handleFunc: func(data *models.Rk7NotifyEvent) error {
			capturedEvent = data
			return nil
		},
	}

	handler := NewHandler(mockService)

	xmlBody := `<a RestCode="123" DateTime="2024-01-15T10:30:00" Situation="1" seqnumber="42" guid="test-guid" name="Started" ShiftNum="5" ShiftDate="2024-01-15"></a>`
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(xmlBody))
	rec := httptest.NewRecorder()

	// Act
	handler.HandleEvents(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if capturedEvent == nil {
		t.Fatal("Event was not captured")
	}

	if capturedEvent.RestCode != "123" {
		t.Errorf("Expected RestCode '123', got '%s'", capturedEvent.RestCode)
	}

	if capturedEvent.Name != "Started" {
		t.Errorf("Expected Name 'Started', got '%s'", capturedEvent.Name)
	}

	if capturedEvent.SeqNumber != "42" {
		t.Errorf("Expected SeqNumber '42', got '%s'", capturedEvent.SeqNumber)
	}
}

func TestHandleEvents_WithNestedElements(t *testing.T) {
	// Arrange
	var capturedEvent *models.Rk7NotifyEvent
	mockService := &MockNotifyService{
		handleFunc: func(data *models.Rk7NotifyEvent) error {
			capturedEvent = data
			return nil
		},
	}

	handler := NewHandler(mockService)

	xmlBody := `<a RestCode="456" DateTime="2024-01-15T10:30:00" name="OrderChanged">
		<Station id="1" code="ST001" name="Station 1" guid="station-guid" NetName="NET1"></Station>
		<Server id="2" code="SRV001" name="Server 1" guid="server-guid" NetName="NET2"></Server>
		<Waiter id="3" code="W001" name="John Doe">
			<Role id="10" code="R001" name="Waiter"></Role>
		</Waiter>
		<Order visit="100" orderIdent="ORD123" guid="order-guid" orderName="Table 5" orderSum="1500.50" openTime="2024-01-15T09:00:00" locked="0"></Order>
		<Item ClassName="Dish"></Item>
	</a>`

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(xmlBody))
	rec := httptest.NewRecorder()

	// Act
	handler.HandleEvents(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if capturedEvent == nil {
		t.Fatal("Event was not captured")
	}

	// Verify nested Station
	if capturedEvent.Station == nil {
		t.Fatal("Station was not parsed")
	}
	if capturedEvent.Station.ID != "1" {
		t.Errorf("Expected Station.ID '1', got '%s'", capturedEvent.Station.ID)
	}
	if capturedEvent.Station.Name != "Station 1" {
		t.Errorf("Expected Station.Name 'Station 1', got '%s'", capturedEvent.Station.Name)
	}

	// Verify nested Server
	if capturedEvent.Server == nil {
		t.Fatal("Server was not parsed")
	}
	if capturedEvent.Server.Code != "SRV001" {
		t.Errorf("Expected Server.Code 'SRV001', got '%s'", capturedEvent.Server.Code)
	}

	// Verify nested Waiter with Role
	if capturedEvent.Waiter == nil {
		t.Fatal("Waiter was not parsed")
	}
	if capturedEvent.Waiter.Name != "John Doe" {
		t.Errorf("Expected Waiter.Name 'John Doe', got '%s'", capturedEvent.Waiter.Name)
	}
	if capturedEvent.Waiter.Role == nil {
		t.Fatal("Waiter.Role was not parsed")
	}
	if capturedEvent.Waiter.Role.Name != "Waiter" {
		t.Errorf("Expected Role.Name 'Waiter', got '%s'", capturedEvent.Waiter.Role.Name)
	}

	// Verify nested Order
	if capturedEvent.Order == nil {
		t.Fatal("Order was not parsed")
	}
	if capturedEvent.Order.OrderIdent != "ORD123" {
		t.Errorf("Expected Order.OrderIdent 'ORD123', got '%s'", capturedEvent.Order.OrderIdent)
	}
	if capturedEvent.Order.OrderSum != "1500.50" {
		t.Errorf("Expected Order.OrderSum '1500.50', got '%s'", capturedEvent.Order.OrderSum)
	}

	// Verify nested Item
	if capturedEvent.Item == nil {
		t.Fatal("Item was not parsed")
	}
	if capturedEvent.Item.ClassName != "Dish" {
		t.Errorf("Expected Item.ClassName 'Dish', got '%s'", capturedEvent.Item.ClassName)
	}
}

func TestHandleEvents_InvalidXML(t *testing.T) {
	// Arrange
	mockService := &MockNotifyService{}
	handler := NewHandler(mockService)

	invalidXML := `<invalid>not-closed`
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(invalidXML))
	rec := httptest.NewRecorder()

	// Act
	handler.HandleEvents(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid XML, got %d", http.StatusBadRequest, rec.Code)
	}

	if !bytes.Contains(rec.Body.Bytes(), []byte("Failed to parse XML")) {
		t.Errorf("Expected 'Failed to parse XML' in response body, got: %s", rec.Body.String())
	}
}

func TestHandleEvents_ServiceError(t *testing.T) {
	// Arrange
	expectedError := errors.New("service processing failed")
	serviceCalled := false
	mockService := &MockNotifyService{
		handleFunc: func(data *models.Rk7NotifyEvent) error {
			serviceCalled = true
			return expectedError
		},
	}

	handler := NewHandler(mockService)

	xmlBody := `<a RestCode="123" name="Started"></a>`
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(xmlBody))
	rec := httptest.NewRecorder()

	// Act
	handler.HandleEvents(rec, req)

	// Assert
	// Handler returns 200 OK immediately after parsing XML to acknowledge receipt
	// This prevents RK7 from retrying the event
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d (handler acknowledges immediately), got %d", http.StatusOK, rec.Code)
	}

	// Verify the service was called (error is logged, not returned to client)
	if !serviceCalled {
		t.Error("Expected service to be called despite error")
	}
}

func TestHandleEvents_EmptyBody(t *testing.T) {
	// Arrange
	mockService := &MockNotifyService{}
	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(""))
	rec := httptest.NewRecorder()

	// Act
	handler.HandleEvents(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for empty body, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestNewHandler(t *testing.T) {
	// Arrange
	mockService := &MockNotifyService{}

	// Act
	handler := NewHandler(mockService)

	// Assert
	if handler == nil {
		t.Fatal("NewHandler returned nil")
	}

	if handler.NotifyService != mockService {
		t.Error("NotifyService was not set correctly")
	}
}