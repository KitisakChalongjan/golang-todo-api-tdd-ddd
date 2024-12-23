package handler

import (
	"golang-todo-api-tdd-ddd/core"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/service"
	"golang-todo-api-tdd-ddd/valueobject"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func InitializeRoleHandler(engine core.Engine) {

	roleRepo := repository.NewRoleRepository(engine.DB)
	roleService := service.NewRoleService(roleRepo)
	roleHandler := NewRoleHandler(roleService)

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{SigningKey: []byte(engine.SecretKey)})

	userGroup := engine.Echo.Group("/role")

	userGroup.POST("/", roleHandler.CreateRole, jwtMiddleware)

}

func (handler *RoleHandler) CreateRole(c echo.Context) error {

	response := core.ApiRespose{}
	roleVO := valueobject.CreateRoleVO{}

	if err := c.Bind(&roleVO); err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusBadRequest, response)
	}

	newRoleID, err := handler.roleService.CreateRole(roleVO)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Error = ""
	response.Data = map[string]string{"ID": newRoleID}

	return c.JSON(http.StatusOK, response)
}
