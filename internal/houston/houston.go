package houston

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed ui/dist/*
var bundle embed.FS

func Bundle() fs.FS {
	build, subErr := fs.Sub(bundle, "ui/dist")
	if subErr != nil {
		panic(fmt.Errorf("unable to find build folder in the bundle: %w", subErr))
	}

	return build
}
