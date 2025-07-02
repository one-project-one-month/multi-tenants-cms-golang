package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/service"
	"github.com/sirupsen/logrus"
)

type Handler interface{}

type AuthHandle interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}

type AuthHandler struct {
	log     *logrus.Logger
	service service.AuthService
}

var _ Handler = (*AuthHandler)(nil)
var _ AuthHandle = (*AuthHandler)(nil)

func NewAuthHandler(
	log *logrus.Logger,
	service service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		log:     log,
		service: service,
	}
}
func (a AuthHandler) Login(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (a AuthHandler) Register(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (a AuthHandler) Logout(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (a AuthHandler) Refresh(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
