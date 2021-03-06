package main

import (
	"net/http"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/termbox/termbox/api/config"
	"github.com/termbox/termbox/api/driver"
)

type Api struct {
	log    *logrus.Logger
	config *config.Config
	echo   *echo.Echo
}

func New(config *config.Config) *Api {

	// -- Logging

	log := logrus.New()

	// -- Echo

	echo := echo.New()

	echo.Use(middleware.Gzip())

	// -- Api

	a := &Api{log, config, echo}

	echo.POST("/machines", a.createMachine)

	return a
}

func (a *Api) Run() error {

	if a.config.TLSConfig.Enable {
		if a.config.TLSConfig.Auto {
			return a.echo.StartAutoTLS(a.config.Address)
		} else {
			return a.echo.StartTLS(a.config.Address, a.config.TLSConfig.Cert, a.config.TLSConfig.Key)
		}
	} else {
		return a.echo.Start(a.config.Address)
	}
}

func (a *Api) getDriver(m *driver.Machine) (driver.Driver, error) {
	ctx := driver.DriverContext{Config: a.config, Machine: m}

	if a.config.ClusterConfig.Enable {
		return driver.NewClusterDriver(&ctx)
	} else {
		str := a.config.Read(fmt.Sprintf("%v.remote", m.Driver))
		url, err := url.ParseUrl(str)
		if err != nil {
			return err
		}

		ctx.Remote = url
		return driver.NewDriver(&ctx)
	}
}

func (a *Api) createMachine(c echo.Context) error {

	m := new(driver.Machine)
	if err := c.Bind(m); err != nil {
		return err
	}

	driver, err := a.getDriver(m)
	if err != nil {
		return err
	}

	if err := driver.Create(); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, m)
}
