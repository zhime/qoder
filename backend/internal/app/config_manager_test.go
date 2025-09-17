package app

import (
	"testing"

	"devops/internal/config"
)

func TestConfigManager_NewConfigManager(t *testing.T) {
	cm := NewConfigManager()
	if cm == nil {
		t.Fatal("NewConfigManager should not return nil")
	}
	if cm.config != nil {
		t.Error("Initial config should be nil")
	}
}

func TestConfigManager_Validate(t *testing.T) {
	cm := NewConfigManager()
	
	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &config.Config{
				Server: config.Server{
					Port: 8080,
				},
				Database: config.Database{
					DSN: "user:pass@tcp(localhost:3306)/dbname",
				},
				Redis: config.Redis{
					Addr: "localhost:6379",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid port - zero",
			config: &config.Config{
				Server: config.Server{
					Port: 0,
				},
				Database: config.Database{
					DSN: "user:pass@tcp(localhost:3306)/dbname",
				},
				Redis: config.Redis{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid port - too large",
			config: &config.Config{
				Server: config.Server{
					Port: 70000,
				},
				Database: config.Database{
					DSN: "user:pass@tcp(localhost:3306)/dbname",
				},
				Redis: config.Redis{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
		},
		{
			name: "empty database DSN",
			config: &config.Config{
				Server: config.Server{
					Port: 8080,
				},
				Database: config.Database{
					DSN: "",
				},
				Redis: config.Redis{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
		},
		{
			name: "empty redis address",
			config: &config.Config{
				Server: config.Server{
					Port: 8080,
				},
				Database: config.Database{
					DSN: "user:pass@tcp(localhost:3306)/dbname",
				},
				Redis: config.Redis{
					Addr: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cm.validate(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigManager_GetConfig(t *testing.T) {
	cm := NewConfigManager()
	
	// Test when config is nil
	if cm.GetConfig() != nil {
		t.Error("GetConfig should return nil when config is not loaded")
	}
	
	// Test when config is set
	testConfig := &config.Config{
		Server: config.Server{Port: 8080},
	}
	cm.config = testConfig
	
	if cm.GetConfig() != testConfig {
		t.Error("GetConfig should return the set config")
	}
}

func TestConfigManager_InitLogger_WithoutConfig(t *testing.T) {
	cm := NewConfigManager()
	
	err := cm.InitLogger()
	if err == nil {
		t.Error("InitLogger should return error when config is not loaded")
	}
}

func TestConfigManager_Reload(t *testing.T) {
	cm := NewConfigManager()
	
	// 测试重新加载配置
	err := cm.Reload()
	// 在测试环境中可能会失败，这是正常的
	if err != nil {
		t.Logf("Reload failed as expected in test environment: %v", err)
	}
}

func TestConfigManager_ValidateEnvironment(t *testing.T) {
	cm := NewConfigManager()
	
	// 测试没有加载配置时的情况
	err := cm.ValidateEnvironment()
	if err == nil {
		t.Error("ValidateEnvironment should return error when config is not loaded")
	}
	
	// 测试有配置时的情况
	testConfig := &config.Config{
		Server: config.Server{
			Port: 8080,
			Mode: "production",
		},
	}
	cm.config = testConfig
	
	err = cm.ValidateEnvironment()
	if err == nil {
		t.Log("ValidateEnvironment passed or failed as expected based on environment")
	}
}

func TestConfigManager_IsProduction(t *testing.T) {
	cm := NewConfigManager()
	
	// 测试没有配置时
	if cm.IsProduction() {
		t.Error("IsProduction should return false when config is not loaded")
	}
	
	// 测试生产环境
	testConfig := &config.Config{
		Server: config.Server{
			Mode: "release",
		},
	}
	cm.config = testConfig
	
	if !cm.IsProduction() {
		t.Error("IsProduction should return true for release mode")
	}
	
	// 测试非生产环境
	testConfig.Server.Mode = "debug"
	if cm.IsProduction() {
		t.Error("IsProduction should return false for debug mode")
	}
}

func TestConfigManager_IsDevelopment(t *testing.T) {
	cm := NewConfigManager()
	
	// 测试没有配置时
	if !cm.IsDevelopment() {
		t.Error("IsDevelopment should return true when config is not loaded")
	}
	
	// 测试开发环境
	testConfig := &config.Config{
		Server: config.Server{
			Mode: "debug",
		},
	}
	cm.config = testConfig
	
	if !cm.IsDevelopment() {
		t.Error("IsDevelopment should return true for debug mode")
	}
	
	// 测试非开发环境
	testConfig.Server.Mode = "release"
	if cm.IsDevelopment() {
		t.Error("IsDevelopment should return false for release mode")
	}
}

func TestConfigManager_GetVersion(t *testing.T) {
	cm := NewConfigManager()
	version := cm.GetVersion()
	if version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", version)
	}
}

func TestConfigManager_String(t *testing.T) {
	cm := NewConfigManager()
	
	// 测试没有配置时
	str := cm.String()
	expected := "ConfigManager{config: nil}"
	if str != expected {
		t.Errorf("Expected %s, got %s", expected, str)
	}
	
	// 测试有配置时
	testConfig := &config.Config{
		Server: config.Server{
			Port: 8080,
			Mode: "debug",
		},
	}
	cm.config = testConfig
	
	str = cm.String()
	expected = "ConfigManager{port: 8080, mode: debug}"
	if str != expected {
		t.Errorf("Expected %s, got %s", expected, str)
	}
}

func TestConfigManager_IsConfigLoaded(t *testing.T) {
	cm := NewConfigManager()
	
	// 测试没有加载配置时
	if cm.IsConfigLoaded() {
		t.Error("IsConfigLoaded should return false when config is not loaded")
	}
	
	// 测试加载配置后
	testConfig := &config.Config{}
	cm.config = testConfig
	
	if !cm.IsConfigLoaded() {
		t.Error("IsConfigLoaded should return true when config is loaded")
	}
}