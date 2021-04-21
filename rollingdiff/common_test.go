package rollingdiff

import (
	"flag"
	"math/rand"
	"testing"
	"time"

	"github.com/tuommaki/rollingdiff/fastcdc"
)

var seed = flag.Int64("seed", time.Now().Unix(), "seed for rng")

func randomBytes(t *testing.T, s int64, n int) []byte {
	t.Helper()
	t.Logf("randomBytes(%d): seed == %d\n", n, s)

	rng := rand.New(rand.NewSource(s))

	buf := make([]byte, n)
	for i := range buf {
		buf[i] = byte(rng.Intn(256))
	}
	return buf
}

func alignChunkIndexes(chunks []Chunk) []Chunk {
	for i := range chunks {
		chunks[i].Index = i
	}

	return chunks
}

func dropChunkAt(chunks []Chunk, index int) []Chunk {
	return append(chunks[:index], chunks[index+1:]...)
}

func randomChunk(t *testing.T, s int64, index int) Chunk {
	t.Helper()
	c := randomChunks(t, s, 1)[0]
	c.Index = index
	return c
}

func randomChunks(t *testing.T, s int64, n int) []Chunk {
	t.Helper()
	bs := randomBytes(t, s, n*fastcdc.MaxSize)
	chunks := Signatures(bs)
	return chunks[:n]
}

func replaceChunkAt(chunks []Chunk, chunk Chunk, index int) []Chunk {
	chunks[index] = chunk
	return chunks
}

func swapChunksAt(chunks []Chunk, from, to int) []Chunk {
	c := chunks[from]
	chunks[from] = chunks[to]
	chunks[to] = c
	return chunks
}
