package config

import (
	"testing"
)

func shouldPanic(t *testing.T, f func()) {
	defer func() { _ = recover() }()
	f()
	t.Errorf("should have panicked")
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		want      func(t *testing.T, c *Config)
		wantPanic bool
	}{
		{
			name: "should set default values",
			want: func(t *testing.T, c *Config) {
				if c.App.Name != "HeyTaxi Location API" {
					t.Errorf("want App.Name = %q, got %q", "HeyTaxi Location API", c.App.Name)
				}
				if c.Server.Http.Port != "8080" {
					t.Errorf("want Server.Http.Port = %q, got %q", "8080", c.Server.Http.Port)
				}
				if c.Server.Http.BodyLimit != "1M" {
					t.Errorf("want Server.Http.BodyLimit = %q, got %q", "1M", c.Server.Http.BodyLimit)
				}
				if c.Server.Http.RequestTimeout != 60 {
					t.Errorf("want Server.Http.RequestTimeout = %d, got %d", 60, c.Server.Http.RequestTimeout)
				}
				if c.Server.Http.ShutdownTimeout != 5 {
					t.Errorf("want Server.Http.ShutdownTimeout = %d, got %d", 5, c.Server.Http.ShutdownTimeout)
				}
			},
		},
		{
			name: "should set env values",
			env: map[string]string{
				"APP_NAME":               "Custom App Name",
				"SERVER_HTTP_HOST":       "localhost",
				"SERVER_HTTP_PORT":       "8081",
				"SERVER_HTTP_BODY_LIMIT": "2M",
			},
			want: func(t *testing.T, c *Config) {
				if c.App.Name != "Custom App Name" {
					t.Errorf("want App.Name = %q, got %q", "Custom App Name", c.App.Name)
				}
				if c.Server.Http.Host != "localhost" {
					t.Errorf("want Server.Http.Host = %q, got %q", "localhost", c.Server.Http.Host)
				}
				if c.Server.Http.Port != "8081" {
					t.Errorf("want Server.Http.Port = %q, got %q", "8081", c.Server.Http.Port)
				}
				if c.Server.Http.BodyLimit != "2M" {
					t.Errorf("want Server.Http.BodyLimit = %q, got %q", "2M", c.Server.Http.BodyLimit)
				}
			},
		},
		{
			name: "should panic if env value is invalid",
			env: map[string]string{
				"SERVER_HTTP_REQUEST_TIMEOUT": "aa",
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			if tt.wantPanic {
				shouldPanic(t, func() {
					New()
				})
				return
			}

			got := New()
			tt.want(t, got)
		})
	}
}

func TestConfig_GetProfile(t *testing.T) {
	tests := []struct {
		name       string
		profileEnv string
		want       string
	}{
		{
			name: "should return local profile by default",
			want: "local",
		},
		{
			name:       "should return profile from env",
			profileEnv: "production",
			want:       "production",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.profileEnv != "" {
				t.Setenv("ACTIVE_PROFILE", tt.profileEnv)
			}

			c := New()
			if got := c.GetProfile(); got != tt.want {
				t.Errorf("Config.GetProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsLocal(t *testing.T) {
	tests := []struct {
		name       string
		profileEnv string
		want       bool
	}{
		{
			name: "should return true by default",
			want: true,
		},
		{
			name:       "should return false if profile is not local",
			profileEnv: "production",
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.profileEnv != "" {
				t.Setenv("ACTIVE_PROFILE", tt.profileEnv)
			}

			c := New()
			if got := c.IsLocal(); got != tt.want {
				t.Errorf("Config.IsLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name       string
		profileEnv string
		want       bool
	}{
		{
			name: "should return false by default",
			want: false,
		},
		{
			name:       "should return true if profile is production",
			profileEnv: "production",
			want:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.profileEnv != "" {
				t.Setenv("ACTIVE_PROFILE", tt.profileEnv)
			}

			c := New()
			if got := c.IsProduction(); got != tt.want {
				t.Errorf("Config.IsProduction() = %v, want %v", got, tt.want)
			}
		})
	}
}
