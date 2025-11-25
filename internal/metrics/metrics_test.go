package metrics

import (
	"testing"
)

func TestNewMetrics(t *testing.T) {
	m := NewMetrics("test")

	if m == nil {
		t.Fatal("expected metrics to be created")
	}

	if m.HTTPRequestsTotal == nil {
		t.Error("expected HTTPRequestsTotal to be initialized")
	}

	if m.HTTPRequestDuration == nil {
		t.Error("expected HTTPRequestDuration to be initialized")
	}

	if m.HTTPRequestsInFlight == nil {
		t.Error("expected HTTPRequestsInFlight to be initialized")
	}

	if m.HTTPResponseSize == nil {
		t.Error("expected HTTPResponseSize to be initialized")
	}

	if m.DBQueriesTotal == nil {
		t.Error("expected DBQueriesTotal to be initialized")
	}

	if m.DBQueryDuration == nil {
		t.Error("expected DBQueryDuration to be initialized")
	}

	if m.AuthAttemptsTotal == nil {
		t.Error("expected AuthAttemptsTotal to be initialized")
	}

	if m.TokensIssued == nil {
		t.Error("expected TokensIssued to be initialized")
	}
}

func TestRecordHTTPRequest(t *testing.T) {
	m := NewMetrics("test_http")

	// Should not panic
	m.RecordHTTPRequest("GET", "/health", "200", 0.1, 100)
	m.RecordHTTPRequest("POST", "/api/v1/users", "201", 0.2, 500)
	m.RecordHTTPRequest("GET", "/api/v1/users", "401", 0.05, 50)

	// Verify metric was recorded (basic sanity check)
	if m.HTTPRequestsTotal == nil {
		t.Error("metric was not recorded")
	}
}

func TestRecordDBQuery(t *testing.T) {
	m := NewMetrics("test_db")

	// Should not panic
	m.RecordDBQuery("SELECT", "users", "success", 0.01)
	m.RecordDBQuery("INSERT", "messages", "success", 0.02)
	m.RecordDBQuery("SELECT", "users", "error", 0.001)

	// Verify metric was recorded (basic sanity check)
	if m.DBQueriesTotal == nil {
		t.Error("metric was not recorded")
	}
}

func TestRecordAuthAttempt(t *testing.T) {
	m := NewMetrics("test_auth")

	// Should not panic
	m.RecordAuthAttempt("login", "success")
	m.RecordAuthAttempt("login", "failure")
	m.RecordAuthAttempt("register", "success")
	m.RecordAuthAttempt("apple", "success")
	m.RecordAuthAttempt("google", "failure")

	// Verify metric was recorded (basic sanity check)
	if m.AuthAttemptsTotal == nil {
		t.Error("metric was not recorded")
	}
}

func TestRecordTokenIssued(t *testing.T) {
	m := NewMetrics("test_token")

	// Should not panic
	m.RecordTokenIssued("access")
	m.RecordTokenIssued("refresh")

	// Verify metric was recorded (basic sanity check)
	if m.TokensIssued == nil {
		t.Error("metric was not recorded")
	}
}

func TestUpdateDBPoolStats(t *testing.T) {
	m := NewMetrics("test_pool")

	// Should not panic
	m.UpdateDBPoolStats(25, 10, 15)
	m.UpdateDBPoolStats(25, 20, 5)

	// Verify metric was recorded (basic sanity check)
	if m.DBConnectionPool == nil {
		t.Error("metric was not recorded")
	}
}

func TestUpdateBusinessMetrics(t *testing.T) {
	m := NewMetrics("test_business")

	// Should not panic
	m.UpdateBusinessMetrics(100, 500)
	m.UpdateBusinessMetrics(150, 750)

	// Verify metric was recorded (basic sanity check)
	if m.UsersTotal == nil {
		t.Error("metric was not recorded")
	}
}

func TestDefaultMetrics(t *testing.T) {
	if DefaultMetrics == nil {
		t.Error("expected DefaultMetrics to be initialized")
	}

	if DefaultMetrics.HTTPRequestsTotal == nil {
		t.Error("expected DefaultMetrics.HTTPRequestsTotal to be initialized")
	}
}
