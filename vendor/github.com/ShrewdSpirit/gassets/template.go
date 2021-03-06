package gassets

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

func buildTemplate(ac *assetConfig) error {
	t, err := template.New("gasset").Parse(assetBodyTemplate)
	if err != nil {
		return err
	}

	packageDir, err := filepath.Abs(filepath.Join(ac.AssetsDir, ac.OutputPath))
	if err != nil {
		return err
	}

	packageName := filepath.Base(packageDir)

	buffer := &bytes.Buffer{}
	if err := t.Execute(buffer, map[string]interface{}{
		"PackageName": packageName,
		"Entities":    buildEntitiesTemplate(ac),
	}); err != nil {
		return err
	}

	src, err := format.Source(buffer.Bytes())
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(packageDir, "gassets.go"), src, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func buildEntitiesTemplate(ac *assetConfig) string {
	builder := &strings.Builder{}

	for _, entity := range ac.Root.Entities {
		buildEntity(builder, entity)
	}

	return builder.String()
}

func buildEntity(builder *strings.Builder, entity *assetEntity) {
	if entity.IsDir {
		builder.WriteString(fmt.Sprintf("%s := %s.AddChild(gassets.NewEntry(true, \"%s\", nil, nil))\n", entity.Name, entity.Parent.Name, entity.Name))
		for _, e := range entity.Entities {
			buildEntity(builder, e)
		}
	} else {
		aliasesBuffer := strings.Builder{}
		if entity.OverrideFiles != nil {
			aliasesBuffer.WriteString("[]string{")
			for i, alias := range entity.OverrideFiles {
				aliasesBuffer.WriteString(fmt.Sprintf(`"%s"`, alias))
				if i != len(entity.OverrideFiles)-1 {
					aliasesBuffer.WriteString(", ")
				}
			}
			aliasesBuffer.WriteString("}")
		} else {
			aliasesBuffer.WriteString("nil")
		}

		dataBuffer := strings.Builder{}
		dataBuffer.WriteString("[]byte{")
		for i, b := range entity.Data {
			dataBuffer.WriteString(strconv.FormatUint(uint64(b), 10))
			if i != len(entity.Data)-1 {
				dataBuffer.WriteString(", ")
			}
		}
		dataBuffer.WriteString("}")

		builder.WriteString(fmt.Sprintf("%s.AddChild(gassets.NewEntry(false, \"%s\", %s, %s))\n",
			entity.Parent.Name, entity.Name, dataBuffer.String(), aliasesBuffer.String()))
	}
}

var assetBodyTemplate = `
// This file has been generated by gassets
// Any modification to this file will be overwritten by gassets tool

package {{.PackageName}}

import (
	"github.com/ShrewdSpirit/gassets"
)

var Root *gassets.Entry

func init() {
	root := gassets.NewEntry(true, "root", nil, nil)
	{{.Entities}}
	Root = root
}
`
