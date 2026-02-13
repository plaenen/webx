package webx

import "embed"

//go:embed static
var Static embed.FS

//go:embed byol
var Byol embed.FS
