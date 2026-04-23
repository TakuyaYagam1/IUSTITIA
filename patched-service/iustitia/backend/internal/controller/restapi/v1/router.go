package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wahrwelt-kit/go-httpkit/httperr"

	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/middleware"
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

func RegisterRoutes(r chi.Router, h *Server, jwtSecret string) {
	wrapper := &openapi.ServerInterfaceWrapper{
		Handler: h,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			var (
				requiredHeader *openapi.RequiredHeaderError
				requiredParam  *openapi.RequiredParamError
				invalidParam   *openapi.InvalidParamFormatError
			)
			switch {
			case errors.As(err, &requiredHeader):
				h.OnError(w, r, httperr.New(err, http.StatusUnauthorized, "NOT_AUTHENTICATED"),
					"openapi", "required_header")
			case errors.As(err, &requiredParam), errors.As(err, &invalidParam):
				h.OnError(w, r, httperr.New(err, http.StatusBadRequest, "BAD_REQUEST"),
					"openapi", "bad_param")
			default:
				h.OnError(w, r, httperr.New(err, http.StatusBadRequest, "BAD_REQUEST"),
					"openapi", "bad_request")
			}
		},
	}

	// Public
	r.Get("/api/health", wrapper.GetHealth)

	r.Post("/api/auth/login", wrapper.Login)

	// Authenticated (any role)
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.RequireAuth(jwtSecret))

		pr.Post("/api/auth/logout", wrapper.Logout)

		pr.Get("/api/auth/me", wrapper.GetCurrentUser)
	})

	// Complaints
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.RequireAuth(jwtSecret))

		pr.With(middleware.RequireRole(domain.RoleCitizen)).
			Post("/api/complaints", wrapper.CreateComplaint)

		pr.With(middleware.RequireRole(domain.RoleJudge, domain.RoleProsecutor)).
			Get("/api/complaints/{case_id}", wrapper.ListComplaintsByCase)

		pr.With(middleware.RequireRole(domain.RoleProsecutor)).
			Post("/api/complaints/{id}/evidence", wrapper.AttachEvidence)
	})

	// Cases
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.RequireAuth(jwtSecret))

		readRoles := middleware.RequireRole(domain.RoleJudge, domain.RoleProsecutor)

		pr.With(readRoles).Get("/api/cases", wrapper.ListCases)

		pr.With(readRoles).Get("/api/cases/{id}", wrapper.GetCaseByID)

		pr.With(middleware.RequireRole(domain.RoleCitizen)).
			Post("/api/cases", wrapper.CreateCase)

		pr.With(middleware.RequireRole(domain.RoleRegistrar)).
			Post("/api/cases/{id}/accept", wrapper.AcceptCase)

		pr.With(middleware.RequireRole(domain.RoleRegistrar)).
			Post("/api/cases/{id}/dismiss", wrapper.DismissCase)

		pr.With(middleware.RequireRole(domain.RoleProsecutor)).
			Post("/api/cases/{id}/opinion", wrapper.FileOpinion)

		pr.With(middleware.RequireRole(domain.RoleJudge, domain.RoleProsecutor)).
			Get("/api/cases/{id}/opinion", wrapper.GetCaseOpinion)

		pr.With(middleware.RequireRole(domain.RoleJudge)).
			Post("/api/cases/{id}/verdict", wrapper.SetCaseVerdict)

		pr.Post("/api/cases/search", wrapper.SearchCases)
	})

	// Hearings (judge queue) + Users lookup (for registrar)
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.RequireAuth(jwtSecret))

		pr.With(middleware.RequireRole(domain.RoleJudge)).
			Get("/api/hearings", wrapper.ListHearings)

		pr.With(middleware.RequireRole(domain.RoleRegistrar, domain.RoleJudge)).
			Get("/api/users", wrapper.ListUsersByRole)
	})

	// Documents
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.RequireAuth(jwtSecret))

		pr.With(middleware.RequireRole(domain.RoleJudge)).
			Post("/api/documents/generate", wrapper.GenerateDocument)

		pr.With(middleware.RequireRole(domain.RoleProsecutor, domain.RoleJudge)).
			Get("/api/documents/{id}", wrapper.GetDocumentByID)

		pr.With(middleware.RequireRole(domain.RoleJudge)).
			Get("/api/cases/{id}/documents", wrapper.ListCaseDocuments)
	})

	// Archive-
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.RequireAuth(jwtSecret))

		pr.Get("/api/archive", wrapper.ListArchive)

		pr.Get("/api/archive/{id}", wrapper.GetArchiveByID)

		pr.With(middleware.RequireRole(domain.RoleJudge)).
			Patch("/api/archive/{id}", wrapper.PatchArchive)
	})
}
