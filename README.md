# rollingdiff

Rolling hash implementation

## Introduction

This repository contains a rolling hash implementation to split input data into
chunks with accompanied SHA-256 signatures and a function to compute delta of
two lists of chunks.

## Chunking algorithm

[FastCDC algorithm](https://www.usenix.org/conference/atc16/technical-sessions/presentation/xia)
is used to split input data into chunks and it follows the default chunk size
limits of _min 2KB - max 64KB_.


## Performance characteristics

Present implementation provides an API that expects input data to be fully
in-memory. Similarly the data structures contain references to original input
data slice.

This is mainly to focus on correctness, ease of testability and lack of
concrete use case at this point.

In case this was to be used for e.g. making incremental backups of large files,
the API should be refactored into such form that it either works with stream of
data (`io.Reader`) or directly utilizes applicable system APIs such as
`mmap(2)`.

To further optimize the produced delta, consecutive repeating operations could
be coalesced into one.
