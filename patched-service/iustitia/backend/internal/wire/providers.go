package wire

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	logkit "github.com/wahrwelt-kit/go-logkit"

	"github.com/TakuyaYagam1/iustitia/config"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/middleware"
	v1 "github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/helper"
	"github.com/TakuyaYagam1/iustitia/internal/repo"
	"github.com/TakuyaYagam1/iustitia/internal/usecase"
)

func ProvideSQLStore(db *sql.DB) *repo.SQLStore {
	return repo.NewSQLStore(db)
}

func ProvideUserUseCase(store repo.Store, cfg *config.Config, logger logkit.Logger) *usecase.User {
	return usecase.NewUser(store, cfg.JWTSecret, cfg.JWTTTL, logger)
}

func ProvideComplaintUseCase(store repo.Store, logger logkit.Logger) *usecase.Complaint {
	return usecase.NewComplaint(store, logger)
}

func ProvideCaseUseCase(store repo.Store, logger logkit.Logger) *usecase.Case {
	return usecase.NewCase(store, logger)
}

func ProvideDocumentUseCase(store repo.Store, logger logkit.Logger) *usecase.Document {
	return usecase.NewDocument(store, logger)
}

func ProvideArchiveUseCase(store repo.Store, logger logkit.Logger) *usecase.Archive {
	return usecase.NewArchive(store, logger)
}

func ProvideServerDeps(
	user *usecase.User,
	complaint *usecase.Complaint,
	caseUC *usecase.Case,
	document *usecase.Document,
	archive *usecase.Archive,
	logger logkit.Logger,
) *helper.ServerDeps {
	return &helper.ServerDeps{
		User:      helper.UserDeps{UserUC: user},
		Complaint: helper.ComplaintDeps{ComplaintUC: complaint},
		Case:      helper.CaseDeps{CaseUC: caseUC},
		Document:  helper.DocumentDeps{DocumentUC: document},
		Archive:   helper.ArchiveDeps{ArchiveUC: archive},
		Infra:     helper.InfraDeps{Logger: logger},
	}
}

func ProvideRESTServer(deps *helper.ServerDeps) *v1.Server {
	return v1.NewServer(deps)
}

func ProvideRouter(cfg *config.Config, logger logkit.Logger, restServer *v1.Server) chi.Router {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.CORS(cfg.CORSOrigins))

	v1.RegisterRoutes(r, restServer, cfg.JWTSecret)
	return r
}

func ProvideHTTPServer(cfg *config.Config, router chi.Router) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       60 * time.Second,
	}
}

func ProvideApp(server *http.Server) *App {
	return &App{Server: server}
}
