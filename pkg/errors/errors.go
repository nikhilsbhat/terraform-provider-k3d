package errors

import stdErrors "errors"

var (
	ErrClusterAlreadyExists    = stdErrors.New("cluster already exists")
	ErrConfigFileReference     = stdErrors.New("for more info refer 'https://k3d.io/usage/configfile/'")
	ErrCreateNodesFailed       = stdErrors.New("creating nodes failed")
	ErrDeleteNodesFailed       = stdErrors.New("deleting nodes failed")
	ErrGenerateRandomBytes     = stdErrors.New("error generating random bytes")
	ErrImportImagesFailed      = stdErrors.New("importing images to clusters errored")
	ErrInsufficientRandomBytes = stdErrors.New("generated insufficient random bytes")
	ErrInvalidMemoryLimit      = stdErrors.New("provided memory limit value is invalid")
	ErrNodeNotFound            = stdErrors.New("nodes not found to start/stop them")
	ErrUnsupportedKind         = stdErrors.New("unsupported kind, only supported value is Simple")
)
