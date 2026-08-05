package assets

import "embed"

//go:embed build/*
var Assets embed.FS

//go:embed src/templates/*
var TemplateAssets embed.FS
