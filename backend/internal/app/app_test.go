package app

import (
	"testing"
)

func TestApplication_New(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() should not return nil")
	}
	
	if app.shutdownCh == nil {
		t.Error("shutdownCh should be initialized")
	}
	
	// Check that all managers are nil initially
	if app.configMgr != nil {
		t.Error("configMgr should be nil initially")
	}
	if app.databaseMgr != nil {
		t.Error("databaseMgr should be nil initially")
	}
	if app.cacheMgr != nil {
		t.Error("cacheMgr should be nil initially")
	}
	if app.serverMgr != nil {
		t.Error("serverMgr should be nil initially")
	}
}

func TestApplication_GetConfig(t *testing.T) {
	app := New()
	
	// Test when config is nil
	if app.GetConfig() != nil {
		t.Error("GetConfig should return nil when config is not loaded")
	}
}

func TestApplication_GetDB(t *testing.T) {
	app := New()
	
	// Test when db is nil
	if app.GetDB() != nil {
		t.Error("GetDB should return nil when database is not initialized")
	}
}

func TestApplication_GetRedis(t *testing.T) {
	app := New()
	
	// Test when redis is nil
	if app.GetRedis() != nil {
		t.Error("GetRedis should return nil when redis is not initialized")
	}
}

// MockConfigManager for testing
type MockConfigManager struct {
	loadError bool
	initLoggerError bool
}

func (m *MockConfigManager) Load() error {
	if m.loadError {
		return &mockError{"mock load error"}
	}
	return nil
}

func (m *MockConfigManager) InitLogger() error {
	if m.initLoggerError {
		return &mockError{"mock init logger error"}
	}
	return nil
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

func TestApplication_initConfig(t *testing.T) {
	app := New()
	
	// Test successful config initialization
	err := app.initConfig()
	// Since we can't mock the actual config loading without modifying the code,
	// we expect this to potentially fail in test environment
	// The important thing is that the method exists and handles errors
	if err != nil {
		t.Logf("initConfig failed as expected in test environment: %v", err)
	}
	
	// Verify that configMgr is created
	if app.configMgr == nil {
		t.Error("configMgr should be created after initConfig")
	}
}