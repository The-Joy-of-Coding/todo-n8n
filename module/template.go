package module

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed template/*.sh
var scripts embed.FS

func getTemplate() (string, error) {
	var builder strings.Builder
	err := fs.WalkDir(
		scripts, "template",
		func(path string, d fs.DirEntry, err error) error {
			return dirFunc(path, d, &builder, err)
		})
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func dirFunc(path string, d fs.DirEntry, builder *strings.Builder, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() && strings.HasSuffix(path, ".sh") {
		content, err := scripts.ReadFile(path)
		if err != nil {
			return err
		}
		builder.Write(content)
		builder.WriteString("\n")
	}
	return nil
}
