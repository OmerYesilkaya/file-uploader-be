package api

import (
	"database/sql"

	"github.com/OmerYesilkaya/fileuploader/internal/config"
)

type AppContext struct {
	DB     *sql.DB
	Config *config.Config
	// Add other dependencies here, like logger, cache, etc.
}
