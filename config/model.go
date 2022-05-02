package config

type (
	Config struct {
		App struct {
			Name string `default:"HeyTaxi Location API"`
		}

		Server struct {
			Http struct {
				Host            string   `default:""`
				Port            string   `default:"8080"`
				BodyLimit       string   `default:"1M"`
				RequestTimeout  int      `default:"60"`
				ShutdownTimeout int      `default:"5"`
				CorsOrigins     []string `default:"*"`
			}
		}

		Redis struct {
			Addr         string `default:"localhost:6379"`
			Password     string `default:""`
			DB           int    `default:""`
			DefaultDb    string `default:""`
			MinIdleConns int    `default:""`
			PoolSize     int    `default:""`
			PoolTimeout  int    `default:""`
			MaxRetries   int    `default:"3"`
		}

		VehicleService struct {
			Host string `default:"localhost"`
			Port string `default:"50052"`
		}

		Jwt struct {
			Issuer                   string `default:"hey-taxi-identity-api"`
			AccessTokenPublicKeyFile string `default:"/etc/certs/access-token-public-key.pem"`
		}
	}
)
