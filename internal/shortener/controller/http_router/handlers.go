package http_router

import (
	"errors"

	"url-shortener/internal/shortener/dto"
	"url-shortener/internal/shortener/entity"
	_ "url-shortener/pkg/http"
	"url-shortener/pkg/observability/otel/tracer"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
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
func (r Router) createLink(c *fiber.Ctx) error {
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

	output, err := r.uc.CreateLink(ctx, input)
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

// getLink handler
//
// @Summary Get a short link by alias
// @Tags Links
// @Accept plain
// @Produce json
// @Param alias path string true "Link alias"
// @Success 200 {object} dto.GetLinkOutput
// @Failure 400 {object} http.ErrHTTP
// @Failure 404 {object} http.ErrHTTP
// @Failure 500 {object} http.ErrHTTP
// @Router /shortener/v1/link/{alias} [get]
func (r Router) getLink(c *fiber.Ctx) (err error) {
	ctx, span := tracer.Start(c.Context(), "http/v1 GetLink")
	defer span.End()

	alias := c.Params("alias")

	input := dto.GetLinkInput{Alias: alias}
	if err := input.Validate(); err != nil {
		log.Error().Msg("uc.GetLink: alias is required")
		return fiber.NewError(fiber.StatusBadRequest, "validate error")
	}

	output, err := r.uc.GetLink(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			log.Error().Err(err).Msg("uc.GetLink: not found")
			return fiber.NewError(fiber.StatusNotFound, "not found")
		default:
			log.Error().Err(err).Msg("uc.GetLink: internal error")
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
func (r Router) redirect(c *fiber.Ctx) (err error) {
	ctx, span := tracer.Start(c.Context(), "http/v1 Redirect")
	defer span.End()

	alias := c.Params("alias")

	input := dto.GetLinkInput{Alias: alias}
	if err := input.Validate(); err != nil {
		log.Error().Msg("uc.GetLink: alias is required")
		return fiber.NewError(fiber.StatusBadRequest, "validate error")
	}

	output, err := r.uc.GetLink(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			log.Error().Err(err).Msg("uc.GetLink: not found")
			return fiber.NewError(fiber.StatusNotFound, "not found")
		default:
			log.Error().Err(err).Msg("uc.GetLink: internal error")
			return fiber.NewError(fiber.StatusInternalServerError, "internal error")
		}
	}

	return c.Redirect(output.URL, fiber.StatusFound)
}
