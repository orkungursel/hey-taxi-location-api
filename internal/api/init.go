package api

import (
	"github.com/go-redis/redis/v8"
	"github.com/orkungursel/hey-taxi-location-api/internal/server"
	userService "github.com/orkungursel/hey-taxi-location-api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	server.Plug(func(s *server.Server, next server.Next) {
		config := s.Config()

		// Redis
		redisOptions := &redis.Options{
			Addr:       config.Redis.Addr,
			Password:   config.Redis.Password,
			PoolSize:   config.Redis.PoolSize,
			DB:         config.Redis.DB,
			MaxRetries: config.Redis.MaxRetries,
		}
		rc := redis.NewClient(redisOptions)
		defer rc.Close()

		// Set up a connection to the server.
		userServiceAddr := config.AuthGrpc.Host + ":" + config.AuthGrpc.Port
		userServiceConn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			next(err)
			return
		}
		defer userServiceConn.Close()

		// User Service GRPC Client
		usc := userService.NewUserServiceClient(userServiceConn)

		if err := Api(s, rc, usc); err != nil {
			next(err)
			return
		}

		next(nil)

		<-s.Wait() // should wait until all http handlers are closed because of the context
	})
}
