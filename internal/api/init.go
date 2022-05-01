package api

import (
	"github.com/go-redis/redis/v8"
	"github.com/orkungursel/hey-taxi-location-api/config"
	"github.com/orkungursel/hey-taxi-location-api/internal/server"
	"github.com/orkungursel/hey-taxi-location-api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	server.Plug(func(s *server.Server, next server.Next) {
		config := s.Config()

		// Redis
		rc := NewRedisClientWithConfig(config)
		defer rc.Close()

		// Set up a connection to vehicle grpc service.
		vehicleServiceAddr := config.VehicleService.Host + ":" + config.VehicleService.Port
		vehicleServiceConn, err := grpc.Dial(vehicleServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			next(err)
			return
		}
		defer vehicleServiceConn.Close()

		// User Service GRPC Client
		vs := proto.NewVehicleServiceClient(vehicleServiceConn)

		if err := Api(s, rc, vs); err != nil {
			next(err)
			return
		}

		next(nil)

		<-s.Wait() // should wait until all http handlers are closed because of the context
	})
}

func NewRedisClientWithConfig(config *config.Config) *redis.Client {
	redisOptions := &redis.Options{
		Addr:       config.Redis.Addr,
		Password:   config.Redis.Password,
		PoolSize:   config.Redis.PoolSize,
		DB:         config.Redis.DB,
		MaxRetries: config.Redis.MaxRetries,
	}

	return redis.NewClient(redisOptions)
}
