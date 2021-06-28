package store

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

//Interface ...
type Interface struct {
	db   *pgxpool.Pool
	opts Options
}

type Options struct {
	Ip string `long:"ip" env:"IP" required:"true" description:"ip of host"`
}

//New ...
func New(dbPath string) (*Interface, error) {

	dbConnectInfo := dbPath
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, dbConnectInfo)
	return &Interface{db: db}, err
}

func (d *Interface) InsertUrlToDb(c echo.Context) error {

	ctx := context.Background()
	var urls UrlData
	err := json.NewDecoder(c.Request().Body).Decode(&urls)
	if err != nil {
		log.Println(err)
	}
	log.Println(urls.RedirectUrl)
	if isValidUrl(urls.RedirectUrl) {
		a, _ := d.db.Exec(ctx, "select * from urls where redirect_url = $1;", urls.RedirectUrl)
		if string(a) == "SELECT 0" {
			charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP1234567890"
			var output strings.Builder
			length := 60
			for i := 0; i < length; i++ {
				random := rand.Intn(len(charSet))
				randomChar := charSet[random]
				output.WriteString(string(randomChar))
			}
			urls.ShortUrl = output.String()

			d.db.Exec(ctx, "insert into urls(redirect_url, short_url) values ($1,$2);", urls.RedirectUrl, d.opts.Ip+urls.ShortUrl)
			err = c.JSON(http.StatusOK, urls)
		} else {
			a, err = d.db.Exec(ctx, "select * from urls where redirect_url = $1;", urls.RedirectUrl)
			if string(a) == "SELECT 1" {
				rows, _ := d.db.Query(ctx, "select short_url from urls where redirect_url = $1;", urls.RedirectUrl)
				for rows.Next() {
					var longUrl string
					err := rows.Scan(&longUrl)
					if err != nil {
						return err
					}
					urls.ShortUrl = longUrl
					log.Println(urls)
					err = c.JSON(http.StatusOK, urls)
					return err
				}
			}

		}
		log.Println(a)

		return err
	} else {
		return c.JSON(http.StatusBadRequest, "Bad url, input carefully please")
	}
}
func (d *Interface) insertUrl(c echo.Context) error {
	ctx := context.Background()
	var urls UrlData
	err := json.NewDecoder(c.Request().Body).Decode(&urls)
	if err != nil {
		log.Println(err)
	}
	if isValidUrl(urls.RedirectUrl) {
		d.db.Exec(ctx, "insert into urls(redirect_url, short_url) values ($1,$2);", urls.RedirectUrl, urls.ShortUrl)
		err = c.JSON(http.StatusOK, urls)
		return err
	} else {
		return c.JSON(http.StatusBadRequest, "")
	}
}

func (d *Interface) Redirect(c echo.Context, shortUrl string) error {
	ctx := context.Background()
	var urls UrlData
	urls.ShortUrl = shortUrl

	a, err := d.db.Exec(ctx, "select * from urls where short_url = $1;", urls.ShortUrl)
	log.Println(a)
	if string(a) == "SELECT 1" {
		rows, _ := d.db.Query(ctx, "select redirect_url from urls where short_url = $1;", urls.ShortUrl)
		for rows.Next() {
			var redirectUrl string
			err := rows.Scan(&redirectUrl)
			if err != nil {
				return err
			}
			urls.RedirectUrl = redirectUrl
			err = c.Redirect(http.StatusPermanentRedirect, redirectUrl)
		}

	} else {
		err = c.JSON(http.StatusBadRequest, urls)
	}

	return err
}

func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
