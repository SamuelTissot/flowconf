package test

import "embed"

//go:embed data/*
var FileSystem embed.FS
