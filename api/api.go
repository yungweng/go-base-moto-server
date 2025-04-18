// Package api configures an http server for administration and application resources.
package api

import (
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dhax/go-base/api/admin"
	"github.com/dhax/go-base/api/app"
	"github.com/dhax/go-base/api/group"
	"github.com/dhax/go-base/api/rfid"
	"github.com/dhax/go-base/api/room"
	"github.com/dhax/go-base/api/student"
	"github.com/dhax/go-base/api/user"
	"github.com/dhax/go-base/auth/jwt"
	"github.com/dhax/go-base/auth/pwdless"
	"github.com/dhax/go-base/database"
	"github.com/dhax/go-base/email"
	"github.com/dhax/go-base/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// New configures application resources and routes.
func New(enableCORS bool) (*chi.Mux, error) {
	logger := logging.NewLogger()

	db, err := database.DBConn()
	if err != nil {
		logger.WithField("module", "database").Error(err)
		return nil, err
	}

	mailer, err := email.NewMailer()
	if err != nil {
		logger.WithField("module", "email").Error(err)
		return nil, err
	}

	authStore := database.NewAuthStore(db)
	authResource, err := pwdless.NewResource(authStore, mailer)
	if err != nil {
		logger.WithField("module", "auth").Error(err)
		return nil, err
	}

	adminAPI, err := admin.NewAPI(db)
	if err != nil {
		logger.WithField("module", "admin").Error(err)
		return nil, err
	}

	appAPI, err := app.NewAPI(db)
	if err != nil {
		logger.WithField("module", "app").Error(err)
		return nil, err
	}

	rfidAPI, err := rfid.NewAPI(db)
	if err != nil {
		logger.WithField("module", "rfid").Error(err)
		return nil, err
	}

	roomAPI, err := room.NewAPI(db)
	if err != nil {
		logger.WithField("module", "room").Error(err)
		return nil, err
	}

	// Initialize stores
	userStore := database.NewUserStore(db)
	studentStore := database.NewStudentStore(db)
	timespanStore := database.NewTimespanStore(db)

	// Create API resources
	userAPI := user.NewResource(userStore, authStore)
	studentAPI := student.NewResource(studentStore, authStore)

	// Connect RFID API with User, Student, and Timespan stores for tag tracking
	rfidAPI.SetUserStore(userStore)
	rfidAPI.SetStudentStore(studentStore)
	rfidAPI.SetTimespanStore(timespanStore)

	groupStore := database.NewGroupStore(db)
	groupAPI := group.NewResource(groupStore, authStore)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Use(logging.NewStructuredLogger(logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// use CORS middleware if client is not served by this api, e.g. from other domain or CDN
	if enableCORS {
		r.Use(corsConfig().Handler)
	}

	r.Mount("/auth", authResource.Router())

	// RFID endpoint doesn't require auth
	r.Mount("/rfid", rfidAPI.Router())

	r.Group(func(r chi.Router) {
		r.Use(authResource.TokenAuth.Verifier())
		r.Use(jwt.Authenticator)
		r.Mount("/admin", adminAPI.Router())
		r.Mount("/api", appAPI.Router())
		r.Mount("/rooms", roomAPI.Router())
		r.Mount("/users", userAPI.Router())
		r.Mount("/students", studentAPI.Router())
		r.Mount("/groups", groupAPI.Router())
	})

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Get("/*", SPAHandler("public"))

	return r, nil
}

func corsConfig() *cors.Cors {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	})
}

// SPAHandler serves the public Single Page Application.
func SPAHandler(publicDir string) http.HandlerFunc {
	handler := http.FileServer(http.Dir(publicDir))
	return func(w http.ResponseWriter, r *http.Request) {
		indexPage := path.Join(publicDir, "index.html")
		serviceWorker := path.Join(publicDir, "service-worker.js")

		requestedAsset := path.Join(publicDir, r.URL.Path)
		if strings.Contains(requestedAsset, "service-worker.js") {
			requestedAsset = serviceWorker
		}
		if _, err := os.Stat(requestedAsset); err != nil {
			http.ServeFile(w, r, indexPage)
			return
		}
		handler.ServeHTTP(w, r)
	}
}
