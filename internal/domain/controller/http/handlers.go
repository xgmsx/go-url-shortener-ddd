package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/entity"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/fetch"
	_ "github.com/xgmsx/go-url-shortener-ddd/pkg/http"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

type HandlerCreateLink struct {
	uc create.Usecase
}

func NewHandlerCreateLink(uc create.Usecase) *HandlerCreateLink {
	return &HandlerCreateLink{uc: uc}
}

// Handler CreateLink
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
func (h *HandlerCreateLink) Handler(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.Context(), "http/v1 CreateLink")
	defer span.End()

	var input dto.CreateLinkInput
	if err := c.BodyParser(&input); err != nil {
		log.Error().Err(err).Msg("c.BodyParser")
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}

	if err := input.Validate(); err != nil {
		log.Error().Err(err).Msg("uc.CreateLink: validation error")
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	output, err := h.uc.Create(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExist):
			return c.Status(fiber.StatusFound).JSON(output)
		default:
			log.Error().Err(err).Msg("uc.CreateLink: internal error")
			return fiber.NewError(fiber.StatusInternalServerError, "internal error")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

type HandlerFetchLink struct {
	uc fetch.Usecase
}

func NewHandlerFetchLink(uc fetch.Usecase) *HandlerFetchLink {
	return &HandlerFetchLink{uc: uc}
}

// Handler FetchLink
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
func (h *HandlerFetchLink) Handler(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.Context(), "http/v1 FetchLink")
	defer span.End()

	alias := c.Params("alias")

	input := dto.FetchLinkInput{Alias: alias}
	if err := input.Validate(); err != nil {
		log.Error().Msg("uc.FetchLink: alias is required")
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	output, err := h.uc.Fetch(ctx, input)
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

type HandlerRedirect struct {
	uc fetch.Usecase
}

func NewHandlerRedirect(uc fetch.Usecase) *HandlerRedirect {
	return &HandlerRedirect{uc: uc}
}

// Handler Redirect
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
func (h *HandlerRedirect) Handler(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.Context(), "http/v1 Redirect")
	defer span.End()

	alias := c.Params("alias")

	input := dto.FetchLinkInput{Alias: alias}
	if err := input.Validate(); err != nil {
		log.Error().Msg("uc.FetchLink: alias is required")
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	output, err := h.uc.Fetch(ctx, input)
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
