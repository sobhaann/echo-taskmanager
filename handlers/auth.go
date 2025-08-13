package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sobhaann/echo-taskmanager/auth"
	_ "github.com/sobhaann/echo-taskmanager/docs" // Swagger generated docs
	"github.com/sobhaann/echo-taskmanager/models"
)

// Signup godoc
//
//	@Summary		Signup
//	@Description	Create a new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			signupReq	body		models.SignupReq	true	"Signup request"
//	@Success		201			{object}	models.User
//	@Failure		400			{object}	map[string]string
//	@Failure		409			{object}	map[string]string
//	@Router			/signup [post]
func (u *Handler) Signup(c echo.Context) error {
	var signupReq models.SignupReq
	if err := c.Bind(&signupReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad body"})
	}

	if existed_user, _ := u.store.GetUserByPhoneNumebr(signupReq.PhoneNumber, c.Request().Context()); existed_user != nil {
		return c.JSON(http.StatusConflict, existed_user)
	}

	hashsed_password, err := auth.HashPassword(signupReq.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	user := &models.User{
		UserName:    signupReq.Name,
		Password:    hashsed_password,
		PhoneNumber: signupReq.PhoneNumber,
	}

	err = u.store.CreateUser(user, c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, user)
}

// Login godoc
//
//	@Summary		Login
//	@Description	Authenticate user and return JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			loginReq	body		models.LoginReq	true	"Login request"
//	@Success		200			{object}	map[string]string
//	@Failure		400			{object}	map[string]string
//	@Router			/login [post]
func (u *Handler) Login(c echo.Context) error {
	var loginReq models.LoginReq

	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad body"})
	}

	user, err := u.store.GetUserByPhoneNumebr(loginReq.PhoneNumber, c.Request().Context())
	if err != nil || user == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user doesn't exist"})
	}
	if err := auth.CheckPassword(user.Password, loginReq.Password); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "password doesnt match"})
	}
	token, err := auth.CreateToken(*user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "token error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"access_token": token})

}

// GetUsers godoc
//
//	@Summary		Get all users
//	@Description	Get all users from database
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	models.User
//	@Failure		500	{object}	map[string]string
//	@Router			/users [get]
func (u *Handler) GetUsers(c echo.Context) error {
	users, err := u.store.GetUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

// Profile godoc
//
//	@Summary		Get profile
//	@Description	Get user profile from JWT token in Authorization header
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.User
//	@Failure		401	{object}	map[string]string
//	@Router			/profile [get]
func (h *Handler) Profile(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	phone_number := claims["phone_number"].(string)
	claimed_user, err := h.store.GetUserByPhoneNumebr(phone_number, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]error{"user not found": err})
	}
	return c.JSON(http.StatusOK, claimed_user)
}

//refresh token
//access token
