package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/orkungursel/hey-taxi-location-api/config"
	"github.com/orkungursel/hey-taxi-location-api/internal/api/http/middleware"
	"github.com/orkungursel/hey-taxi-location-api/internal/app"
	"github.com/orkungursel/hey-taxi-location-api/pkg/logger"
)

type Controller struct {
	config          *config.Config
	logger          logger.ILogger
	locationService app.LocationService
	tokenService    app.TokenService
}

func NewController(config *config.Config, logger logger.ILogger,
	ls app.LocationService, ts app.TokenService) *Controller {

	return &Controller{
		config:          config,
		logger:          logger,
		tokenService:    ts,
		locationService: ls,
	}
}

// RegisterRoutes registers the routes to the echo server
func (a *Controller) RegisterRoutes(e *echo.Group) {
	e.Use(middleware.ErrorHandler())
	e.Use(middleware.Auth(a.tokenService))

	e.POST("/save/", a.saveLocation())
	e.POST("/search/", a.searchLocation())
}

// @Summary      Save Location
// @Description  Saves the driver location
// @Tags         Location Service
// @Accept       json
// @Produce      json
// @Param        payload  body      app.SaveLocationRequest  true  "Payload"
// @Success      200      {array}   app.LocationResponse
// @Failure      400      {object}  app.HTTPError
// @Failure      500      {object}  app.HTTPError
// @Router       /location/save [post]
// @Security     BearerAuth
func (a *Controller) saveLocation() echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := &app.SaveLocationRequest{}
		if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
			return err
		}

		userId, err := GetUserId(c)
		if err != nil {
			return err
		}

		if err := app.Validate(payload); err != nil {
			return err
		}

		if err := a.locationService.SaveLocation(c.Request().Context(), userId, *payload); err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

// @Summary      Search
// @Description  Searches for driver locations
// @Tags         Location Service
// @Accept       json
// @Produce      json
// @Param        payload  body      app.SearchLocationRequest  true  "Payload"
// @Success      200      {array}   app.LocationResponse
// @Failure      400      {object}  app.HTTPError
// @Failure      500      {object}  app.HTTPError
// @Router       /location/search [post]
// @Security     BearerAuth
func (a *Controller) searchLocation() echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := &app.SearchLocationRequest{}
		if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
			return err
		}

		if err := app.Validate(payload); err != nil {
			return err
		}

		res, err := a.locationService.SearchLocations(c.Request().Context(), *payload)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	}
}
