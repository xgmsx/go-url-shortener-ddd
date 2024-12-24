package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
)

type CreateLinkInput struct {
	URL string `json:"url"`
}

func (i *CreateLinkInput) Validate() error {
	if i.URL == "" {
		return entity.ErrInputValidation
	}

	return nil
}

type CreateLinkOutput struct {
	URL       string    `json:"url"`
	Alias     string    `json:"alias"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (o CreateLinkOutput) Load(l entity.Link) CreateLinkOutput {
	o.URL = l.URL
	o.Alias = l.Alias
	o.ExpiredAt = l.ExpiredAt

	return o
}

func (o CreateLinkOutput) Str() string {
	b, err := json.Marshal(o)
	if err != nil {
		log.Error().Err(err).Msg("dto.Str: json.Marshal")
		return fmt.Sprintf("%v", o)
	}
	return string(b)
}
