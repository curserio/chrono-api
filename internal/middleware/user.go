package middleware

import (
	"github.com/curserio/chrono-api/pkg/timeutil"
	"github.com/labstack/echo/v4"
)

func WithUserContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// userID := extractUserID(c) // из токена / хидера / сессии
		//user, err := userRepo.GetByID(c.Request().Context(), userID)
		//if err != nil {
		//	return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
		//}

		timeutil.SetTZ(c, "Asia/Irkutsk") // TODO for tests

		return next(c)
	}
}
