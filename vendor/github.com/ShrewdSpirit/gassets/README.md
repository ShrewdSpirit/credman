gassets
=====
Easy asset and resource embedder!

## Install
Get the source and install it:

```
go get -u github.com/ShrewdSpirit/gassets/cmd/gassets/...
go install github.com/ShrewdSpirit/gassets/cmd/gassets
```

## Usage
Create a `gassets.toml` file in your assets directory and write the following in it:
```toml
output-path="."

[root]
```

These are required keys for `gassets.toml`. `output-path` specifies where the generated go source will be. It's relative to the directory where `gassets.toml` is.

For each virtual directory and file you create a toml table. `root` is the root of all your vdirs and files. You can add as many vdirs and files to `root` since it's a vdir itself.

**To add virtual directory** to `root` you should add:
```toml
[root.my-virtual-dir]
```

You can add files to each vdir in two ways:
1. Use glob
2. Add a file table

#### Glob file
To use globing for adding files, you must add a `include` key to your vdir table with the value being an array of strings following go's glob format therefore globing isn't recursive.

Each matching file in the list will be added directory to current vdir without any sub virtual directories.

#### File table
A file table has 2 keys, one being optional. The `path` key is required and it's the path to the file relative to `gassets.toml`

The other key is `override-files` which is an array of strings for files to be overriden by the current file if the given path exists at runtime. The path is relative to working directory.

## Command line
Pass `-d` and your asset directory to `gassets`

## Example
```toml
output-path="."

[root]
    [root.file1]
    path="file.txt"

    [root.file2]
    path="test/file2.txt"

    [root.myvdir]
    include=["test/glob/*.jpg"]

        [root.myvdir.index_template]
        path="index.templ"
        override-files=["templates/index.templ", "index.templ"]
```

## License
MIT
