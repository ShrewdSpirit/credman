package gassets

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

var ErrInvalidConfig = errors.New("Invalid config")
var ErrInvalidNode = errors.New("Invalid node")

func Build(dir string) error {
	gassetFilename := filepath.Join(dir, "gassets.toml")

	if _, err := os.Stat(gassetFilename); os.IsNotExist(err) {
		return err
	}

	gassetFile, err := os.Open(gassetFilename)
	if err != nil {
		return err
	}

	config, err := toml.LoadReader(gassetFile)
	if err != nil {
		return err
	}

	outputPath, ok := config.Get("output-path").(string)
	if !ok {
		return ErrInvalidConfig
	}

	ac := &assetConfig{
		AssetsDir:  dir,
		OutputPath: outputPath,
	}

	root, ok := config.Get("root").(*toml.Tree)
	if !ok {
		return ErrInvalidConfig
	}

	if err := readTree(ac, nil, root, "root"); err != nil {
		return err
	}

	return buildTemplate(ac)
}

func readTree(ac *assetConfig, parent *assetEntity, tree *toml.Tree, name string) error {
	// check if it has path
	filePath, ok := tree.Get("path").(string)
	if !ok {
		// it's a vdir
		entity := &assetEntity{
			IsDir:    true,
			Name:     name,
			Parent:   parent,
			Entities: make([]*assetEntity, 0),
		}

		include, ok := tree.Get("include").([]interface{})
		if ok {
			if err := addIncludes(ac, entity, include); err != nil {
				return err
			}
		}

		keys := tree.Keys()
		for _, key := range keys {
			if key == "include" {
				continue
			}

			newTree, ok := tree.Get(key).(*toml.Tree)
			if !ok {
				return ErrInvalidNode
			}

			if err := readTree(ac, entity, newTree, key); err != nil {
				return err
			}
		}

		if name == "root" {
			ac.Root = entity
		} else {
			parent.Entities = append(parent.Entities, entity)
		}
	} else {
		// it's a file
		filePath = filepath.Join(ac.AssetsDir, filePath)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return err
		}

		entity := &assetEntity{
			Name:     name,
			FilePath: filePath,
			Parent:   parent,
		}

		if err := readFile(ac, entity); err != nil {
			return err
		}

		overrideAliases, ok := tree.Get("override-files").([]interface{})
		if ok {
			entity.OverrideFiles = make([]string, 0)
			for _, a := range overrideAliases {
				alias, ok := a.(string)
				if ok {
					entity.OverrideFiles = append(entity.OverrideFiles, alias)
				}
			}
		}

		parent.Entities = append(parent.Entities, entity)
	}

	return nil
}

func addIncludes(ac *assetConfig, vdir *assetEntity, includes []interface{}) error {
	for _, inc := range includes {
		incPath, ok := inc.(string)
		if !ok {
			return ErrInvalidConfig
		}

		files, err := filepath.Glob(filepath.Join(ac.AssetsDir, incPath))
		if err != nil {
			return err
		}

		if files == nil {
			continue
		}

		for _, filename := range files {
			s, err := os.Stat(filename)
			if os.IsNotExist(err) {
				return err
			}

			if s.IsDir() {
				continue
			}

			entity := &assetEntity{
				Name:     filepath.Base(filename),
				FilePath: filename,
				Parent:   vdir,
			}

			if err := readFile(ac, entity); err != nil {
				return err
			}

			vdir.Entities = append(vdir.Entities, entity)
		}
	}
	return nil
}

func readFile(ac *assetConfig, entity *assetEntity) error {
	data, err := ioutil.ReadFile(entity.FilePath)
	if err != nil {
		return err
	}

	entity.Data = data

	return nil
}
