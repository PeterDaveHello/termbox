package api

import "github.com/labstack/echo"
import "net/http"

type Status struct {
	ServerOk   bool
	HyperOk    bool
	Containers int
}

func (a *Api) GetStatus(c echo.Context) error {
	status := Status{ServerOk: true}

	if env, err := a.Hyper.Info(); err == nil {
		status.HyperOk = true
		status.Containers = env.GetInt("Containers")
	}

	return c.JSON(http.StatusOK, status)
}
