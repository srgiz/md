package route

import (
	"log/slog"
	"md/internal/lib/kernel"
	"md/internal/userctx/domain/login"
	"net/http"
	"os"
	"time"
)

func NewPostLogin(app *kernel.App) (string, http.Handler) {
	return "POST /login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := app.Handle(r.Context(), &login.Command{
			Id:       r.FormValue("id"),
			Password: r.FormValue("password"),
		})

		if err != nil {
			slog.DebugContext(r.Context(), err.Error(), "form", r.Form)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		cookieName := os.Getenv("APP_JWT_COOKIE_NAME")

		if cookieName == "" {
			cookieName = "i"
		}

		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    *token.(*string),
			Expires:  time.Now().Add(time.Hour * 24),
			Path:     "/",
			HttpOnly: true,
			//Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})

		w.WriteHeader(http.StatusOK)
	})
}
