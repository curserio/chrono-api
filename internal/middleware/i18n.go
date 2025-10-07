package middleware

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

func I18nMiddleware(defaultLangStr string) echo.MiddlewareFunc {
	tt, _, err := language.ParseAcceptLanguage(defaultLangStr)
	if err != nil || len(tt) == 0 {
		tt = []language.Tag{language.English}
	}

	defaultLang := tt[0]

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			acceptLang := c.Request().Header.Get("Accept-Language")
			if acceptLang == "" {
				c.Set("lang", defaultLang)
				return next(c)
			}

			tags, _, err := language.ParseAcceptLanguage(acceptLang)
			if err != nil || len(tags) == 0 {
				c.Set("lang", defaultLang)
				return next(c)
			}

			c.Set("lang", tags[0])
			return next(c)
		}
	}
}
