package store

type UrlData struct {
	ID          uint   `json:"id" db:"id"`
	RedirectUrl string `json:"redirectUrl" db:"redirect_url"`
	LongUrl     string `json:"longUrl" db:"long_url"`
}
