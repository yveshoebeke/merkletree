# Merkel Tree

## Merkel Tree function written in Go

### General
Hash trees can be used to verify any kind of data stored, handled and transferred in and between computers. They can help ensure that data blocks received from others and even to check that the other peers do not lie , altered and send fake blocks.

This functionality will deliver the root hash of a hash slice.

* Sequence 1 (optionl)
* Sequence 2
* Sequence 3

## Watch it

<p align="center">
<img src="https://upload.wikimedia.org/wikipedia/commons/9/95/Hash_Tree.svg" alt="watch this" height="300" width="500" border="1" />
<p>

### Process schematic of Bitcoin method:
Note: Adjusting for odd number of leaves/branches by appending duplicate of last one.
<p>
 stop ------------> [12345555]          -> Merkle tree root(value out).
                    /        \
                [1234]       [5555]
                /    \       /   \
             [12]   [34]   [55] [55]    -> Branches.
             [12]   [34]   [55]  ^
             /  \   /  \   /  \
            [1][2] [3][4] [5][5]
 start ---> [1][2] [3][4] [5] ^         -> Leaves (array data input).

 ^ = Duplicate hash to make leaf count even.
 A single leaf input will be processed as an odd number of leaves, ie: appended to itself.
</p>

### Process schematic of Keccak method:
Note: Adjusting for odd number of leaves/branches by appending duplicate of last one.
<p>
 stop ------------> [12345555]          -> Merkle tree root(value out).
                    /        \
                [1234]       [5555]
                /    \       /   \
             [12]   [34]   [55] [55]    -> Branches.
             [12]   [34]   [55]  ^
             /  \   /  \   /  \
            [1][2] [3][4] [5][5]
 start ---> [1][2] [3][4] [5] ^         -> Leaves (array data input).

 ^ = Duplicate hash to make leaf count even.
 A single leaf input will be processed as an odd number of leaves, ie: appended to itself.
</p>

## Unit test

## How to use this

1. import github.com/bytesupply/Merkletree
2. provide the hash function to be used and the input slice of byte data
3. more to come

## How to tweak and contribute to this


