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

These are required keys for `gassets.toml`. `output-path` specifies where the generated file is placed. It's relative to the given directory (`-d` commandline option).

For each virtual directory and file, a toml table is required. `root` is the root of all your virtual directories and files. You can add as many vdirs and files to `root` and any other vdir.

**To add virtual directory** to `root` you should add:
```toml
[root.my-virtual-dir]
```

You can add files to each vdir in two ways:
1. Use glob
2. Add a file table

#### Glob file
To add files by globs, `include` key must be supplied with an array of globs. See the example

Only matching files will be included in current vdir.

#### File table
A file table has 2 keys, one being optional. The `path` key is required and it's the path to the file relative to given directory.

The other key is `override-files` which is an array of strings for files to be loaded instead if the given path exists at runtime. The path is relative to working directory.

## Command line
- `-d`: Assets directory (where `gassets.toml` is)

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
