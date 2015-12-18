package syncer

import "database/sql"

// Lookup holds pointers to the dependencies of this module
type Lookup struct {
	*sql.DB
}
