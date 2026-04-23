package wire

import (
	"github.com/google/wire"

	"github.com/TakuyaYagam1/iustitia/internal/repo"
)

var RepoSet = wire.NewSet(
	ProvideSQLStore,
	wire.Bind(new(repo.Store), new(*repo.SQLStore)),
)

var UseCaseSet = wire.NewSet(
	ProvideUserUseCase,
	ProvideComplaintUseCase,
	ProvideCaseUseCase,
	ProvideDocumentUseCase,
	ProvideArchiveUseCase,
)

var HTTPSet = wire.NewSet(
	ProvideServerDeps,
	ProvideRESTServer,
	ProvideRouter,
	ProvideHTTPServer,
	ProvideApp,
)
