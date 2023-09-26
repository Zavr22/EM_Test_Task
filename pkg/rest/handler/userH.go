package handler

import (
	"github.com/Zavr22/EMTestTask/pkg/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// CreateUser is used to create user by admin
//
// @Summary Create user
// @securityDefinitions.apikey ApiKeyAuth
// @Tags users
// @Description create user
// @ID create-user
// @Accept  json
// @Produce  json
// @Param input body model.FIO true "user info"
// @Success 200 {object} model.User
// @Failure 400,404,403 {object} model.CommonResponse
// @Failure 500 {object} model.CommonResponse
// @Router /api/users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var reqBody model.FIO
	errBind := c.Bind(&reqBody)
	if errBind != nil {
		logrus.WithFields(logrus.Fields{
			"reqBody": reqBody,
		}).Errorf("Bind json, %s", errBind)
		return echo.NewHTTPError(http.StatusInternalServerError, model.CommonResponse{Message: "data not correct"})
	}
	userID, err := h.userS.CreateUser(c.Request().Context(), &reqBody)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"message": "unable to create user",
		}).Errorf("error while creating user, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, model.CommonResponse{Message: "error while creating user"})
	}
	return c.JSON(http.StatusOK, userID)
}

// GetUsers is used to get all users
//
// @Summary Get users
// @securityDefinitions.apikey ApiKeyAuth
// @Tags users
// @Description get users
// @ID get-users
// @Accept json
// @Param page query string true "current page"
// @Produce  json
// @Success 200 {array} model.User
// @Failure 400,404,403 {object} model.CommonResponse
// @Failure 500 {object} model.CommonResponse
// @Router /api/users [get]
func (h *Handler) GetUsers(c echo.Context) error {
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"page": page,
		}).Errorf("error converting page into string, %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, model.CommonResponse{Message: "incorrect page"})
	}
	users, err := h.userS.GetAllUsers(c.Request().Context(), page)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"users": users,
		}).Errorf("Get users failed, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "get users failed")
	}
	return c.JSON(http.StatusOK, users)
}

// GetUserByID is used to get user by id
//
// @Summary Get user by id
// @securityDefinitions.apikey ApiKeyAuth
// @Tags users
// @Description get user by id
// @ID get-user-by-id
// @Accept   json
// @Produce  json
// @Param userID query string true "userID"
// @Success 200 {object} model.User
// @Failure 400,404,403 {object} model.CommonResponse
// @Failure 500 {object} model.CommonResponse
// @Router /api/users/:id [get]
func (h *Handler) GetUserByID(c echo.Context) error {
	c.Param("id")
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Errorf("wrong id format, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "wrong id format")
	}
	user, err := h.userS.GetUser(c.Request().Context(), userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user": user,
		}).Errorf("cannot get user, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "cannot get user")
	}
	return c.JSON(http.StatusOK, user)
}

// UpdateUser is used to update user
//
// @Summary Update user
// @securityDefinitions.apikey ApiKeyAuth
// @Tags users
// @Description update user
// @ID update-user
// @Accept   json
// @Produce  json
// @Param input body model.EnrichedFIO true "enter new account info"
// @Success 200 {object} model.CommonResponse
// @Failure 400,404,403 {object} model.CommonResponse
// @Failure 500 {object} model.CommonResponse
// @Router /api/users/:id [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	var reqBody model.EnrichedFIO
	errBind := c.Bind(&reqBody)
	if errBind != nil {
		logrus.WithFields(logrus.Fields{
			"reqBody": reqBody,
		}).Errorf("Bind json, %s", errBind)
		return echo.NewHTTPError(http.StatusInternalServerError, model.CommonResponse{Message: "data not correct"})
	}
	c.Param("id")
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Errorf("wrong id format, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "wrong id format")
	}
	errUpdate := h.userS.UpdateProfile(c.Request().Context(), userID, reqBody)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Errorf("error while creating user, %s", errUpdate)
		return echo.NewHTTPError(http.StatusBadRequest, model.CommonResponse{Message: "error while creating user"})
	}
	return c.JSON(http.StatusOK, model.CommonResponse{Message: "user updated successfully"})
}

// DeleteUser is used to delete user
//
// @Summary Delete user
// @securityDefinitions.apikey ApiKeyAuth
// @Tags users
// @Description delete user
// @ID delete-user
// @Accept   json
// @Produce  json
// @Success 200 {object} model.CommonResponse
// @Failure 400,404,403 {object} model.CommonResponse
// @Failure 500 {object} model.CommonResponse
// @Router /api/users/:id [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	c.Param("id")
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Errorf("wrong id format, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "wrong id format")
	}
	errDelete := h.userS.DeleteProfile(c.Request().Context(), userID)
	if err != nil {
		logrus.Errorf("error while delete user, %s", errDelete)
		return echo.NewHTTPError(http.StatusBadRequest, model.CommonResponse{Message: "couldn't delete user"})
	}
	return c.JSON(http.StatusOK, model.CommonResponse{Message: "user deleted successfully"})
}
