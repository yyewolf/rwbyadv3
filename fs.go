package rwbyadv3

import (
	"embed"
)

//go:embed cards/img/*
var cardFS embed.FS

func GetCardFS() embed.FS {
	return cardFS
}

//go:embed static/dist/*
var staticFS embed.FS

func GetStaticFS() embed.FS {
	return staticFS
}

//go:embed dungeons/dist/*
var dungeonFS embed.FS

func GetDungeonFS() embed.FS {
	return dungeonFS
}
