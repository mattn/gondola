package app

import (
	"gnd.la/orm"
)

// Orm is just a very thin wrapper around
// orm.Orm, which disables the Close method
// when running in production mode, since
// the App is always reusing the same ORM
// instance.
type Orm struct {
	*orm.Orm
}

// Close is a no-op. It prevents the App shared
// orm.Orm from being closed.
func (o *Orm) Close() error {
	return nil
}
