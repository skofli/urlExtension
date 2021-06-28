package api

type UrlData struct {
	ID          uint   `json:"id" db:"id"`
	RedirectUrl string `json:"redirectUrl" db:"redirect_url"`
	ShortUrl    string `json:"shortUrl" db:"short_url"`
}
