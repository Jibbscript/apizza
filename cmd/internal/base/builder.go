package base

import (
	"io"

	"github.com/harrybrwn/apizza/dawg"
	"github.com/harrybrwn/apizza/pkg/cache"
)

// Builder defines an interface for an object that builds commands.
type Builder interface {
	Build(use, short string, r Runner) *Command
	Output() io.Writer
	DBBuilder
	ConfigBuilder
	AddressBuilder
}

// DBBuilder is a cli builder that can give away a database.
type DBBuilder interface {
	DB() *cache.DataBase
}

// ConfigBuilder is a cli builder that can give away a config struct.
type ConfigBuilder interface {
	Config() *Config
}

// AddressBuilder is a builder interface that should be able to get an
// address.
type AddressBuilder interface {
	Address() dawg.Address
}

// AddrDBBuilder is an anddress-builder and a db-builder.
type AddrDBBuilder interface {
	AddressBuilder
	DBBuilder
}
