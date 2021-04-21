package rollingdiff

import (
	"crypto/sha256"

	"github.com/tuommaki/rollingdiff/fastcdc"
)

type Chunk struct {
	Bytes     []byte
	Index     int
	Signature [sha256.Size]byte
}

func Signatures(buf []byte) []Chunk {
	var chunks []Chunk

	for counter, offset := 0, 0; offset < len(buf); counter++ {
		// FastCDC computes the next chunk boundary.
		idx := fastcdc.Compute(buf[offset:len(buf)])

		if (offset + idx) < len(buf) {
			// Returned index from fastcdc points to last byte of chunk.
			// Increase it by one to account for Go slice operation & offset
			// pointing to beginning of next slice.
			idx++
		}

		c := Chunk{
			Bytes:     buf[offset : offset+idx],
			Index:     counter,
			Signature: sha256.Sum256(buf[offset : offset+idx]),
		}

		chunks = append(chunks, c)
		offset += idx
	}

	return chunks
}
