package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/entity"
	_ "github.com/xgmsx/go-url-shortener-ddd/pkg/http"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

// createLink handler
//
// @Summary Create a short link
// @Tags Links
// @Accept json
// @Produce json
// @Param input body dto.CreateLinkInput true "New link"
// @Success 201 {object} dto.CreateLinkOutput
// @Success 302 {object} dto.CreateLinkOutput
// @Failure 400 {object} http.ErrHTTP
// @Failure 404 {object} http.ErrHTTP
// @Failure 500 {object} http.ErrHTTP
// @Router /shortener/v1/link [post]
func (r *Router) createLink(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.Context(), "http/v1 CreateLink")
	defer span.End()

	var input dto.CreateLinkInput
	if err := c.BodyParser(&input); err != nil {
		log.Error().Err(err).Msg("c.BodyParser")
		return fiber.NewError(fiber.StatusBadRequest, "json error")
	}

	if err := input.Validate(); err != nil {
		log.Error().Err(err).Msg("uc.CreateLink: validate error")
		return fiber.NewError(fiber.StatusBadRequest, "validate error")
	}

	output, err := r.ucCreate.Create(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExist):
			return c.Status(fiber.StatusFound).JSON(output)
		case errors.Is(err, entity.ErrEntityValidation):
			log.Error().Err(err).Msg("uc.CreateLink: validate error")
			return fiber.NewError(fiber.StatusBadRequest, "validate error")
		default:
			log.Error().Err(err).Msg("uc.CreateLink: internal error")
			return fiber.NewError(fiber.StatusInternalServerError, "internal error")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

// fetchLink handler
//
// @Summary Fetch a short link by alias
// @Tags Links
// @Accept plain
// @Produce json
// @Param alias path string true "Link alias"
// @Success 200 {object} dto.FetchLinkOutput
// @Failure 400 {object} http.ErrHTTP
// @Failure 404 {object} http.ErrHTTP
// @Failure 500 {object} http.ErrHTTP
// @Router /shortener/v1/link/{alias} [get]
func (r *Router) fetchLink(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.Context(), "http/v1 FetchLink")
	defer span.End()

	alias := c.Params("alias")

	input := dto.FetchLinkInput{Alias: alias}
	if err := input.Validate(); err != nil {
		log.Error().Msg("uc.FetchLink: alias is required")
		return fiber.NewError(fiber.StatusBadRequest, "validate error")
	}

	output, err := r.ucFetch.Fetch(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			log.Error().Err(err).Msg("uc.FetchLink: not found")
			return fiber.NewError(fiber.StatusNotFound, "not found")
		default:
			log.Error().Err(err).Msg("uc.FetchLink: internal error")
			return fiber.NewError(fiber.StatusInternalServerError, "internal error")
		}
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

// Redirect Router.
//
// @Summary      Redirect to URL by alias
// @Tags         Links
// @Accept       plain
// @Produce      plain
// @Param        alias path string true "Link alias"
// @Success      302 "redirect to the original url"
// @Failure      400 {object} http.ErrHTTP
// @Failure      404 {object} http.ErrHTTP
// @Failure      500 {object} http.ErrHTTP
// @Router       /shortener/v1/link/{alias}/redirect [get]
func (r *Router) redirect(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.Context(), "http/v1 Redirect")
	defer span.End()

	alias := c.Params("alias")

	input := dto.FetchLinkInput{Alias: alias}
	if err := input.Validate(); err != nil {
		log.Error().Msg("uc.FetchLink: alias is required")
		return fiber.NewError(fiber.StatusBadRequest, "validate error")
	}

	output, err := r.ucFetch.Fetch(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			log.Error().Err(err).Msg("uc.FetchLink: not found")
			return fiber.NewError(fiber.StatusNotFound, "not found")
		default:
			log.Error().Err(err).Msg("uc.FetchLink: internal error")
			return fiber.NewError(fiber.StatusInternalServerError, "internal error")
		}
	}

	return c.Redirect(output.URL, fiber.StatusFound)
}
