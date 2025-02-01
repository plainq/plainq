package mutations

import (
	"embed"
	"io/fs"
)

var (
	//go:embed storage/*.sql
	storage embed.FS

	//go:embed telemetry/*.sql
	telemetry embed.FS
)

// StorageMutations returns all embedded storage files as embed.FS.
func StorageMutations() fs.FS {
	d, err := fs.Sub(storage, "storage")
	if err != nil {
		panic(err)
	}

	return d
}

// TelemetryMutation returns all embedded storage files as embed.FS.
func TelemetryMutation() fs.FS {
	d, err := fs.Sub(telemetry, "telemetry")
	if err != nil {
		panic(err)
	}

	return d
}
