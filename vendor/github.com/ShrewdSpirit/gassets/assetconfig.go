package gassets

type assetEntity struct {
	IsDir bool
	Name  string

	FilePath      string
	OverrideFiles []string
	Data          []byte

	Parent   *assetEntity
	Entities []*assetEntity
}

type assetConfig struct {
	AssetsDir  string
	OutputPath string
	Root       *assetEntity
}
