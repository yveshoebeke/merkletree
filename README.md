# Merkel Tree

## Merkel Tree functionality in Go

### General

*"Hash trees can be used to verify any kind of data stored, handled and transferred in and between computers. They can help ensure that data blocks received from others and even to check that the other peers do not lie or altered and send fake blocks."*

This functionality will deliver the root hash of a data collection.

Options available:

* Choose encoding algorithm
* Specify tree creation and Unbalanced tree process
* Initial branch encoding option

---

### Usage

```go
import github.com/yveshoebeke/Merkletree
```

```go
myroot := merkletree.DeriveRoot(algorithm, data, processType, initialEncoding)
```

### General Process Flow Overview

![Process Overview](docs/ProcessOverview.png)

---

### Parameters

There are 3 mandatory parameters and 1 optional.

They are, in order:

1. algorithm
1. data
1. process type
1. initiate with encoding (optional)

#### Algorithm parameter

Accepts a string that denotes the desired hashing algorithm.
The following functionalities are in the registry:

|<mark>Ident<sup>(1)</sup></mark>  |  import     | syntax|
|-------------|---------------|----------------|
|<mark>MD5</mark>          | ```crypto/md5```| ```md5.Sum(data)```|
|<mark>SHA1</mark>          | ```crypto/sha1``` | ```sha1.Sum(data)```|
|<mark>SHA3SUM256</mark>    | ```golang.org/x/crypto/sha3``` |```sha3.Sum256(data)```|
|<mark>SHA256SUM256</mark>  | ```crypto/sha256``` |```sha256.Sum256(data)```|
|<mark>SHA512SUM256</mark>  | ```crypto/sha512``` |```sha512.Sum256(data)```|
|<mark>SHA512SUM512</mark>  | ```crypto/sha512``` |```sha512.Sum512(data)```|

<sup>(1)</sup>Will raise an *unknown hash algorithm* error if no match is found.

#### Data

Expected data type: ```[][]byte```

If empty will raise the *empty list* error.

#### Process Type

There are 3 process types that can be specified. Each one will handle unbalanced trees in it's own manner:

|<mark>Value</mark>|Process Name|
|-----------|-----------|
|<mark>0</mark>| Duplicate and Append|
|<mark>1</mark>| Pass Through|
|<mark>2</mark>| Binary Tree|

#### Initiate with encodeing

Boolean. If not provided, will default to ```false```. If set to ```true``` will direct the function to initialize all data in the data parameter to be encoded with the specified hashing algorithm. To be used in case your input contains unencoded data.

---

### Process Types

Processing of a data collection into a tree can result in a so called Unbalanced situation. This will manifest itself when you encounter an odd number of elements in your data colletion (slice).

The following will describe how this can be handled. If you are working with dta that is part of an established scheme, like a blockchain for example, you should select the one that is used in that particular scheme.

Note: You _must_ declare one. There is no default.

#### Duplicate and Append

When we encounter an odd number of elements, this process will simply create a duplicate of the last element and append it to the end of the slice, thus having in essence 2 equal data elements at the end, making the number of elements even.

![Process Overview](docs/DupeAppend.png)


#### Pass Through

Here we will ignore the last element and just pass it through to the next branch.

![Process Overview](docs/PassThrough.png)

#### Binary Tree

An initial index will be calculated to start the process, after which we will have a binary tree. A Binary tree has branches where the number of elements are of the 2<sup>x</sup> order.

To accomplish this we will first find the lowest exponent needed to raise 2 by to obtain a value that is greater then the number of elements in the slice.

This exponential result is subtracted by the number of elements to obtain the index where the first iteration will start.

We could write a loop that will increment the exponent by one untill the result is greater then the aforementioned number, but by applying the binary logarithm on the length of the slice we can basically avoid this.

So I have employed the following:

```text
L = number of elements in the slice.

X = log2(L); X is converted to an integer after rounding it up.

I = 2^X - L; this is our starting index.
```

All this can be accomplished as so:

```go
startIdx := int(math.Pow(2, math.Ceil(math.Log2(float64(len(ms.Leaves)))))) - len(ms.Leaves)
```