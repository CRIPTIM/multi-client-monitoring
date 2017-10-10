# Cryptographic Monitoring System #

Proof-of-concept implementation of the multi-client predicate-only encryption scheme from “Multi-client
Predicate-only Encryption for Conjunctive Equality Tests” in [Go](https://golang.org).

**Please note that the implementation is not designed to be secure against side-channel attacks.**

Short version on how to build & run:

- Install the dependencies for PBC: `GMP`, `flex`, and `bison`

```
sudo apt-get install libgmp-dev flex bison
```

- Install [PBC](https://crypto.stanford.edu/pbc/download.html):

```
./configure
make
sudo make install
```

- Install Go
- Create a directory for your Go code:

```
mkdir ~/go
```

- Set some environment variables:

```
export GOPATH=~/go
```

- Get the pbc binding for Go:

```
go get github.com/Nik-U/pbc
```

- Clone this repository into your `$GOPATH` directory:

```
cd $GOPATH
git clone <url of repository> src/crypmonsys
```

- Build a simple test program

```
cd src/crypmonsys/cmd/eval
go build .
```

- Run the resulting executable. (You might need to run `sudo ldconfig` first to find the freshly installed PBC
  library.)

```
./eval
```

## Benchmarking ##

You may want to check how fast the implementation runs on your machine.
To run the performance evaluation experiments, use the shell script `evaluations.sh` with one of the curves in
the `param` folder.

**Be sure to use only Type 3 pairings otherwise the construction is insecure!**

* `d159.param`: [MNT curve](https://crypto.stanford.edu/pbc/manual/ch08s06.html)
* `d201.param`: MNT curve
* `d224.param`: MTN curve
* `f.param`: [BN curve](https://crypto.stanford.edu/pbc/manual/ch08s08.html)

Note that no preprocessing is used in the benchmark.
Preprocessing would speedup the testing of multiple rules against a fixed set of ciphertexts.

## Testing ##

Some simple test & benchmark code is included. To run these tests, first go to the `crypmonsys` package
directory:
```
cd $GOPATH/src/crypmonsys
```

To run the tests:
```
go test .
```

To run the tests with additional output (more verbose):
```
go test . -v
```

To run benchmarks:
```
go test -bench .
```
