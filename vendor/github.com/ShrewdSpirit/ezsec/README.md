ezsec
=====
Bunch of security helper functions for browser javascript and go.
Basically [rage_utils](https://github.com/ShrewdSpirit/rage_utils) but uses commonjs modules.

To install:

`yarn add ezsec`

or if want to use in a go project:

`go get github.com/ShrewdSpirit/ezsec/...`

## Usage
**Note**: Some functions return an object with an `error` field that you must check whether it's null to take required action in case of errors.

I'm not gonna document everything since it's obvious what the functions do. But for JS/TS reference, refer to [definitions file](https://github.com/ShrewdSpirit/ezsec/blob/master/types/index.d.ts) and usage can be seen in [unit test file](https://github.com/ShrewdSpirit/ezsec/blob/master/test/functions.test.js).

# License
MIT
