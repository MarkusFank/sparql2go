package assets

import "embed"

/*
	Embedding expects the files to be found in ./dist subfolder
	The frontend files have to be created with "npm run build" in the frontend project
	after that, the files from frontend's dist folder have to copied into ./dist folder (recursively)

	Or just execute copy-assets.sh file, it does everything you need
*/

//go:embed all:dist/*
var Dist embed.FS
