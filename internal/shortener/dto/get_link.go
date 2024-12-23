package dto

import (
	"time"

	"url-shortener/internal/shortener/entity"
)

type GetLinkInput struct {
	Alias string `json:"alias"`
}

func (i GetLinkInput) Validate() error {
	if i.Alias == "" {
		return entity.ErrInputValidation
	}
	return nil
}

type GetLinkOutput struct {
	URL       string    `json:"url"`
	Alias     string    `json:"alias"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (o GetLinkOutput) Load(l entity.Link) GetLinkOutput {
	o.URL = l.URL
	o.Alias = l.Alias
	o.ExpiredAt = l.ExpiredAt

	return o
}
