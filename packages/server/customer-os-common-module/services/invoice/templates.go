package invoice

import "embed"

//go:embed pdf_template/**
var Templates embed.FS
