package rollingdiff

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Delta(t *testing.T) {
	testCases := []struct {
		name      string
		oldChunks []Chunk
		newChunks []Chunk
		expected  []Change
	}{
		{
			name:      "append chunk to end of new chunks",
			oldChunks: randomChunks(t, *seed, 4),
			newChunks: append(randomChunks(t, *seed, 4), randomChunk(t, *seed+1, 4)),
			expected: []Change{
				{
					Op:    Add,
					To:    4,
					Bytes: randomChunk(t, *seed+1, 4).Bytes,
				},
			},
		},
		{
			name:      "prepend chunk to beginning of new chunks",
			oldChunks: randomChunks(t, *seed, 4),
			newChunks: alignChunkIndexes(append([]Chunk{randomChunk(t, *seed+1, 0)}, randomChunks(t, *seed, 4)...)),
			expected: []Change{
				{
					Op:    Add,
					To:    0,
					Bytes: randomChunk(t, *seed+1, 0).Bytes,
				},
			},
		},
		{
			name:      "prepend chunk to beginning of new chunks and delete one in the middle",
			oldChunks: randomChunks(t, *seed, 4),
			newChunks: alignChunkIndexes(dropChunkAt(append([]Chunk{randomChunk(t, *seed+1, 0)}, randomChunks(t, *seed, 4)...), 3)),
			expected: []Change{
				{
					Op:   Delete,
					From: 2,
				},
				{
					Op:    Add,
					To:    0,
					Bytes: randomChunk(t, *seed+1, 0).Bytes,
				},
			},
		},
		{
			name:      "swap chunk in the middle",
			oldChunks: randomChunks(t, *seed, 4),
			newChunks: alignChunkIndexes(swapChunksAt(randomChunks(t, *seed, 4), 1, 2)),
			expected: []Change{
				{
					Op:   Move,
					From: 1,
					To:   2,
				},
				{
					Op:   Move,
					From: 2,
					To:   1,
				},
			},
		},
		{
			name:      "swap chunk in the middle and append a chunk",
			oldChunks: randomChunks(t, *seed, 4),
			newChunks: alignChunkIndexes(append(swapChunksAt(randomChunks(t, *seed, 4), 1, 2), randomChunk(t, *seed+2, 4))),
			expected: []Change{
				{
					Op:    Add,
					To:    4,
					Bytes: randomChunk(t, *seed+2, 4).Bytes,
				},
				{
					Op:   Move,
					From: 1,
					To:   2,
				},
				{
					Op:   Move,
					From: 2,
					To:   1,
				},
			},
		},
		{
			name:      "prepend a chunk and swap chunk in the middle",
			oldChunks: randomChunks(t, *seed, 4),
			newChunks: alignChunkIndexes(append([]Chunk{randomChunk(t, *seed+2, 0)}, swapChunksAt(randomChunks(t, *seed, 4), 1, 2)...)),
			expected: []Change{
				{
					Op:    Add,
					To:    0,
					Bytes: randomChunk(t, *seed+2, 4).Bytes,
				},
				{
					Op:   Move,
					From: 1,
					To:   3,
				},
				{
					Op:   Move,
					From: 2,
					To:   2,
				},
			},
		},
		{
			name:      "delete a chunk in the middle and swap chunks around it",
			oldChunks: randomChunks(t, *seed, 4),
			newChunks: alignChunkIndexes(swapChunksAt(dropChunkAt(randomChunks(t, *seed, 4), 2), 1, 2)),
			expected: []Change{
				{
					Op:   Delete,
					From: 2,
				},
				{
					Op:   Move,
					From: 1,
					To:   2,
				},
				{
					Op:   Move,
					From: 3,
					To:   1,
				},
			},
		},
		{
			name:      "replace a chunk with a new one and swap chunks around it",
			oldChunks: randomChunks(t, *seed, 4),
			newChunks: alignChunkIndexes(
				swapChunksAt(
					replaceChunkAt(
						randomChunks(t, *seed, 4),
						randomChunk(t, *seed+2, 2), 2),
					1, 3),
			),
			expected: []Change{
				{
					Op:   Delete,
					From: 2,
				},
				{
					Op:    Add,
					To:    2,
					Bytes: randomChunk(t, *seed+2, 2).Bytes,
				},
				{
					Op:   Move,
					From: 1,
					To:   3,
				},
				{
					Op:   Move,
					From: 3,
					To:   1,
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Logf(tc.name)

			changes := Delta(tc.oldChunks, tc.newChunks)
			if !cmp.Equal(changes, tc.expected) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.expected, changes))
			}
		})
	}
}
