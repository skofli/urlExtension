package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"urlShort/backend/store"
)

type Interface struct {
	store *store.Interface
}

func (i *Interface) Run() error {

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.File("/", "/home/skofli/go/src/urlShortener/frontend/index.html")
	e.POST("/set", i.setUrlToDb)
	e.GET("/:longUrl", i.redirect)

	return http.ListenAndServe(":8080", e)
}

func New(store *store.Interface) *Interface {
	return &Interface{
		store: store,
	}
}

func (i *Interface) setUrlToDb(c echo.Context) error {
	err := i.store.InsertUrlToDb(c)
	return err
}

func (i *Interface) redirect(c echo.Context) error {
	longUrl := c.Param("longUrl")
	err := i.store.Redirect(c, longUrl)

	return err
}
