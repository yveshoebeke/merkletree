<h1><img src="docs/merkletree.png" style="height:25px;width:30px;float:left;margin-right:15px;border:1px solid black;"/>Merkle Tree</h1>

## Merkle Tree functionality in Go

### General

A hash tree, also known as a Merkle tree, is a tree in which each leaf node is labeled with the cryptographic hash of a data block, and each non-leaf, or branch, node is labeled with the cryptographic hash of its child nodes' labels.

This Merkel tree function allows the user to specify certain behaviors.

Options available:

* Choose encoding algorithm
* Specify tree creation and Unbalanced tree process type
* Initial branch encoding option

---

### Install

```shell
go get github.com/yveshoebeke/merkletree
```

---

### Import

```go
import github.com/yveshoebeke/Merkletree
```

---

### Execute

```go
root, err := merkletree.DeriveRoot(algorithm, data, processType, initialEncoding)
```

---

### Signature

```go
type Merkletree func(string, [][]byte, int, ...bool) ([]byte, error)
```

---

### Test and Benchmark

Standard evocation:

A custom flag is provided to see a more detailed test result:

```shell
go test -detail
```

The benchmark will give you results running through 10,000 element input slice for all 3 process types, using sha256.Sum256 encoding.

```shell
go test -bench=.
```

---

### Parameters

There are 3 mandatory parameters and 1 optional.

They are, in order:

1. algorithm
1. data
1. process type
1. initiate with encoding ~(optional)~

#### Algorithm parameter

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

#### Data

Expected data type: ```[][]byte```<sup>(2)</sup>

<sup>(2)</sup>If empty will raise the *empty list* error.

#### Process Type

There are 3 process types that can be specified. Each one will handle unbalanced trees in it's own manner:

|Value<sup>(3)</sup>|Process Name|
|-----------|-----------|
|0| [Pass Through](#pass-through)|
|1| [Duplicate and Append](#duplicate-and-append)|
|2| [Binary Tree](#binary-tree)|

<sup>(3)</sup>Incorrect value will raise the *invalid process type* error.

Notes:

1. Processing the same input data and subjecting it to different Process Types will obviously result in different Merkle Root values.
1. Some real world usages:
    * *Duplicate and Append* is used in the Bitcoin cryptocurrency/blockchain.
    * *Pass Through* is used in the Monero cryptocurrency/blockchain.

#### Initiate with encodeing

Boolean. ~If not provided, will default to ```false```.~ If set to ```true``` will direct the function to initialize all data in the data parameter to be encoded with the specified hashing algorithm. To be used in case your input contains unencoded data.

---

### More Info/Details

See [Wiki](https://github.com/yveshoebeke/merkletree/wiki)

___
