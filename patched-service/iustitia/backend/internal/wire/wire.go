//go:build wireinject

package wire

import (
	"database/sql"

	"github.com/google/wire"
	logkit "github.com/wahrwelt-kit/go-logkit"

	"github.com/TakuyaYagam1/iustitia/config"
)

//go:generate wire

func InitializeApp(
	cfg *config.Config,
	logger logkit.Logger,
	db *sql.DB,
) (*App, error) {
	wire.Build(RepoSet, UseCaseSet, HTTPSet)
	return nil, nil
}
