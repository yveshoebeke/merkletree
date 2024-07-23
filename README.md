<h1><img src="docs/merkletree.png" style="height:30px;width:30px;float:left;"/>&nbsp;&nbsp;Merkle Tree</h1>

## Merkle Tree functionality in Go

### General

A hash tree, also known as a Merkle tree, is a tree in which each leaf node is labeled with the cryptographic hash of a data block, and each non-leaf, or branch, node is labeled with the cryptographic hash of its child nodes' labels.

This Merkle tree function allows the user to specify certain behaviors. (See in Wiki: [General Process Flow Overview](https://github.com/yveshoebeke/merkletree/wiki/7.-General-Process-Flow-Overview))

Options available:

* Choose encoding algorithm
* Specify tree creation and Unbalanced tree process type
* ~Initial branch encoding option~

---

### Install

```shell
go get github.com/yveshoebeke/merkletree
```

---

### Import

```go
import merkletree github.com/yveshoebeke/Merkletree
```

---

### Execute

```go
root, err := merkletree.DeriveRoot(data, algorithm, processType)
```

---

### Signature

```go
type Merkletree func([][]byte, string, int) ([]byte, error)
```

---

### Test, Benchmark and Proof

#### Test

Standard evocation:

A custom flag is provided to see a more detailed test results:

```shell
go test -detail
```

#### Benchmark

The benchmark will give you results running through 10,000 element input slice for all 3 process types, using sha256.Sum256 encoding.

```shell
go test -bench=.
```

#### Proof

In the ```proof/``` directory is a facility that will show you the step-by-step encoding results (console output) of the various branches. It proves all three supporting processes.

_In order for this to be a valid proof it does not use the ```merkletree.go``` functionality, but can be cross-checked with the results of the test since it uses the same starting data._

You can invoke it as: 

```shell
go run proof/proof.go
```

Hint: For readability you might want to pipe it to a paging facility like ```less``` (in *nix/macOS): 

```shell
go run proof/proof.go | less
```

---

### Parameters

There are 3 mandatory parameters.

They are, in order:

1. data
1. algorithm
1. process type

#### Data ```[][]byte```

Expected data type: ```[][]byte```<sup>(2)</sup>

<sup>(2)</sup>If empty will raise the *empty list* error.

#### Algorithm parameter ```string```

Accepts a string that denotes the desired hashing algorithm.
The following functionalities are in the registry:

|ID<sup>(1)</sup>  |  Resulting import     | Syntax evoked |
|-------------|---------------|----------------|
|MD5          | ```crypto/md5```| ```md5.Sum(data)```|
|SHA1          | ```crypto/sha1``` | ```sha1.Sum(data)```|
|SHA3SUM256    | ```golang.org/x/crypto/sha3``` |```sha3.Sum256(data)```|
|SHA256SUM256  | ```crypto/sha256``` |```sha256.Sum256(data)```|
|SHA512SUM256  | ```crypto/sha512``` |```sha512.Sum256(data)```|
|SHA512SUM512  | ```crypto/sha512``` |```sha512.Sum512(data)```|

<sup>(1)</sup>Will raise an *unknown hash algorithm* error if no match is found.

Note:

* Other schemes can be added by editing the ```cryptofuncs.go``` source.
* Registry signature: ```var AlgorithmRegistry map[string]CryptoFunc```
* Function signature: ```type CryptoFunc func([]byte) []byte```

#### Process Type ```int```

There are 3 process types that can be specified. Each one will handle unbalanced trees in it's own manner:

|Value<sup>(3)</sup>|Process Name<sup>(4)</sup>|
|-----------|-----------|
|0| [Pass Through](https://github.com/yveshoebeke/merkletree/wiki/8.-Process-Types#pass-through)|
|1| [Duplicate and Append](https://github.com/yveshoebeke/merkletree/wiki/8.-Process-Types#duplicate-and-append)|
|2| [Binary Tree](https://github.com/yveshoebeke/merkletree/wiki/8.-Process-Types#binary-tree)|

<sup>(3)</sup>Incorrect value will raise the *invalid process type* error.

<sup>(4)</sup>See [Wiki](https://github.com/yveshoebeke/merkletree/wiki) for detailed process description.

Notes:

1. Processing the same input data and subjecting it to different Process Types will obviously result in different Merkle Root values.
1. To the best of my knowledge at time of writing this some real world process type usage:
    * *Duplicate and Append* is used in the Bitcoin cryptocurrency/blockchain.
    * *Pass Through* is used in the Monero cryptocurrency/blockchain.

---

### More Info/Details

See [Wiki](https://github.com/yveshoebeke/merkletree/wiki)

___
