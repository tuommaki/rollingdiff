package rollingdiff

import (
	"testing"

	"github.com/tuommaki/filediff/fastcdc"
)

func Test_Signatures_Splits_Data_Into_Chunks(t *testing.T) {
	minChunkCount := 4
	data := randomBytes(t, *seed, minChunkCount*fastcdc.MaxSize)

	chunks := Signatures(data)

	if len(chunks) < minChunkCount {
		t.Fatalf("expected len(chunks) < %d, got %d", minChunkCount, len(chunks))
	}
}

func Test_Signatures_Splits_Data_Into_Chunks_With_Correct_Content(t *testing.T) {
	data := randomBytes(t, *seed, 4*fastcdc.MaxSize)

	chunks := Signatures(data)

	offset := 0
	for i, chunk := range chunks {
		for j, b := range chunk.Bytes {
			if data[offset] != b {
				t.Fatalf("expected data[%d] == chunk[%d].Bytes[%d], got %#x != %#x", offset, i, j, data[offset], b)
			}
			offset++
		}
	}

	if len(data) != offset {
		t.Fatalf("expected len(data) == offset, got %d != %d", len(data), offset)
	}
}

func Test_Signatures_Chunking_Is_Repeatable(t *testing.T) {
	oldData := randomBytes(t, *seed, 4*fastcdc.MaxSize)
	oldChunks := Signatures(oldData)

	newData := make([]byte, len(oldData))
	copy(newData, oldData)

	newChunks := Signatures(newData)

	for i := 0; i < len(oldChunks); i++ {
		if len(oldChunks[i].Bytes) != len(newChunks[i].Bytes) {
			t.Fatalf("expected len(oldChunks[%d].Bytes) == len(newChunks[%d].Bytes), got %d != %d", i, i, len(oldChunks[i].Bytes), len(newChunks[i].Bytes))
		}

		if oldChunks[i].Signature != newChunks[i].Signature {
			t.Fatalf("expected oldChunks[%d].Signature == newChunks[%d].Signature, got %#x != %#x", i, i, oldChunks[i].Signature, newChunks[i].Signature)
		}
	}
}

func Test_Signatures_Change_on_Second_Segment(t *testing.T) {
	oldData := randomBytes(t, *seed, 4*fastcdc.MaxSize)

	// Compute signatures for old data.
	origChunks := Signatures(oldData)

	// Clone data for modification.
	newData := make([]byte, len(oldData))
	copy(newData, oldData)

	// Find offset for second chunk segment.
	offset := len(origChunks[0].Bytes)

	// Modify byte from the beginning of second chunk.
	newData[offset] = (newData[offset] % 42) + 1

	updatedChunks := Signatures(newData)

	// Since first segment was not touched. It's expected to stay the say.
	if origChunks[0].Signature != updatedChunks[0].Signature {
		t.Fatalf("expected origChunks[0].Signature == updatedChunks[0].Signature, got:\n%#x != %#x", origChunks[0].Signature, updatedChunks[0].Signature)
	}

	// The second segment must not be equal.
	if origChunks[1].Signature == updatedChunks[1].Signature {
		t.Fatalf("expected origChunks[1].Signature != updatedChunks[1].Signature, got:\n%#x != %#x", origChunks[1].Signature, updatedChunks[1].Signature)
	}

	// Since there was at least four chunks and the second one was changed,
	// the last chunk must be equal in both.
	if origChunks[len(origChunks)-1].Signature != updatedChunks[len(updatedChunks)-1].Signature {
		t.Fatalf("expected origChunks[len(origChunks)-1].Signature == updatedChunks[len(updatedChunks)-1].Signature, got:\n%#x != %#x", origChunks[len(origChunks)-1].Signature, updatedChunks[len(updatedChunks)-1].Signature)
	}
}

func Test_Signatures_Change_of_Last_Segment(t *testing.T) {
	oldData := randomBytes(t, *seed, 4*fastcdc.MaxSize)

	origChunks := Signatures(oldData)

	// Clone data for modification.
	newData := make([]byte, len(oldData))
	copy(newData, oldData)

	// Find beginning of last chunk segment.
	offset := len(oldData) - len(origChunks[len(origChunks)-1].Bytes) + 1

	// Modify byte from the beginning of last chunk.
	newData[offset] = (newData[offset] % 42) + 1

	updatedChunks := Signatures(newData)

	// Assert that there's at least equal number of updated chunks.
	if len(updatedChunks) < len(origChunks) {
		t.Fatalf("expected len(updatedChunks) >= len(origChunks), got %d < %d", len(updatedChunks), len(origChunks))
	}

	// Assert that all but last chunk changed (compared from original chunks
	// point of view).
	for i := 0; i < len(origChunks)-1; i++ {
		if origChunks[i].Signature != updatedChunks[i].Signature {
			t.Fatalf("expected origChunks[%d].Signature == updatedChunks[%d].Signature, got %#x != %#x, with i == %d", i, i, origChunks[i].Signature, updatedChunks[i].Signature, i)
		}
	}

	// The last chunk must not be equal in both.
	i := len(origChunks) - 1
	if origChunks[i].Signature == updatedChunks[i].Signature {
		t.Fatalf("expected origChunks[%d].Signature != updatedChunks[%d].Signature, got: %#x != %#x, with i == %d", i, i, origChunks[i].Signature, updatedChunks[i].Signature, i)
	}
}

func Test_Signatures_Swap_Second_And_Third_Segment(t *testing.T) {
	oldData := randomBytes(t, *seed, 4*fastcdc.MaxSize)

	origChunks := Signatures(oldData)

	// Clone data for modification.
	newData := make([]byte, 0, len(oldData))

	// Add first chunk
	newData = append(newData, origChunks[0].Bytes...)
	// Append third chunk
	newData = append(newData, origChunks[2].Bytes...)
	// Append second chunk
	newData = append(newData, origChunks[1].Bytes...)
	// ...and rest of them
	for _, c := range origChunks[3:] {
		newData = append(newData, c.Bytes...)
	}

	updatedChunks := Signatures(newData)

	// Since first segment was not touched. It's expected to stay the say.
	if origChunks[0].Signature != updatedChunks[0].Signature {
		t.Fatalf("expected origChunks[0].Signature == updatedChunks[0].Signature, got:\n%#x != %#x", origChunks[0].Signature, updatedChunks[0].Signature)
	}

	// The old second segment signature must match new third segment.
	if origChunks[1].Signature != updatedChunks[2].Signature {
		t.Fatalf("expected origChunks[1].Signature == updatedChunks[2].Signature, got:\n%#x != %#x", origChunks[1].Signature, updatedChunks[2].Signature)
	}

	// The old third segment signature must match new second segment.
	if origChunks[2].Signature != updatedChunks[1].Signature {
		t.Fatalf("expected origChunks[2].Signature == updatedChunks[1].Signature, got:\n%#x != %#x", origChunks[2].Signature, updatedChunks[1].Signature)
	}

	// Since there was at least four chunks and the second and third were
	// changed, the last chunk must be equal in both.
	if origChunks[len(origChunks)-1].Signature != updatedChunks[len(updatedChunks)-1].Signature {
		t.Fatalf("expected origChunks[len(origChunks)-1].Signature == updatedChunks[len(updatedChunks)-1].Signature, got:\n%#x != %#x", origChunks[len(origChunks)-1].Signature, updatedChunks[len(updatedChunks)-1].Signature)
	}
}

func Test_Signatures_Add_Data_Between_Second_And_Third_Segment(t *testing.T) {
	oldData := randomBytes(t, *seed, 5*fastcdc.MaxSize)

	origChunks := Signatures(oldData)

	// Clone data for modification.
	newData := make([]byte, len(oldData))
	copy(newData, oldData)

	for i := 0; i < len(oldData); i++ {
		if oldData[i] != newData[i] {
			t.Fatalf("expected oldData[%d] == newData[%d], got %#x != %#x", i, i, oldData[i], newData[i])
		}
	}

	// Find offset for second chunk segment.
	offsetSecond := len(origChunks[0].Bytes)
	lengthSecond := len(origChunks[1].Bytes)
	offsetThird := offsetSecond + lengthSecond

	// Insert data between second and third segment.
	newData = append(newData[:offsetThird], append(randomBytes(t, *seed, fastcdc.MaxSize), newData[offsetThird+1:]...)...)

	updatedChunks := Signatures(newData)

	// Since first segment was not touched. It's expected to stay the say.
	if origChunks[0].Signature != updatedChunks[0].Signature {
		t.Fatalf("expected origChunks[0].Signature == updatedChunks[0].Signature, got:\n%#x != %#x", origChunks[0].Signature, updatedChunks[0].Signature)
	}

	// Second segment must also match.
	if origChunks[1].Signature != updatedChunks[1].Signature {
		t.Fatalf("expected origChunks[1].Signature == updatedChunks[1].Signature, got:\n%#x != %#x", origChunks[1].Signature, updatedChunks[1].Signature)
	}

	// Third segment must differ.
	if origChunks[2].Signature == updatedChunks[2].Signature {
		t.Fatalf("expected origChunks[2].Signature != updatedChunks[2].Signature, got:\n%#x != %#x", origChunks[2].Signature, updatedChunks[2].Signature)
	}

	found := false
findMatchingTrack:
	for i := 2; i < len(origChunks); i++ {
		for j := 2; j < len(updatedChunks); j++ {
			if origChunks[i].Signature == updatedChunks[j].Signature {
				found = true
				break findMatchingTrack
			}
		}
	}

	if !found {
		t.Fatalf("after inserting data in the middle of stream, no matching chunks found afterwards")
	}
}

func Test_Signatures_Add_Chunk_Between_Second_And_Third_Segment(t *testing.T) {
	oldData := randomBytes(t, *seed, 5*fastcdc.MaxSize)

	origChunks := Signatures(oldData)

	// Generate temporary chunks where one whole chunk can be copied & inserted
	// in-between original data so that original chunk boundaries are
	// maintained.
	tmpData := randomBytes(t, *seed, 2*fastcdc.MaxSize)
	tmpChunks := Signatures(tmpData)

	// Build newData.
	newData := make([]byte, 0, len(oldData)+len(tmpChunks[0].Bytes))

	// Add first chunk from old data.
	newData = append(newData, origChunks[0].Bytes...)
	// Append second chunk from old data.
	newData = append(newData, origChunks[1].Bytes...)
	// Append third chunk from separate tmp data.
	newData = append(newData, tmpChunks[0].Bytes...)
	// ...and rest of the chunks from old data.
	for _, c := range origChunks[2:] {
		newData = append(newData, c.Bytes...)
	}

	updatedChunks := Signatures(newData)

	// Verify that updated chunks has exactly one chunk more than original.
	if len(updatedChunks) != len(origChunks)+1 {
		t.Fatalf("expected len(updatedChunks) == len(origChunks) + 1, got %d != %d", len(updatedChunks), len(origChunks)+1)
	}

	// Since first segment was not touched. It's expected to stay the say.
	if origChunks[0].Signature != updatedChunks[0].Signature {
		t.Fatalf("expected origChunks[0].Signature == updatedChunks[0].Signature, got:\n%#x != %#x", origChunks[0].Signature, updatedChunks[0].Signature)
	}

	// Second segment must also match.
	if origChunks[1].Signature != updatedChunks[1].Signature {
		t.Fatalf("expected origChunks[1].Signature == updatedChunks[1].Signature, got:\n%#x != %#x", origChunks[1].Signature, updatedChunks[1].Signature)
	}

	// Third segment must match the external tmp data.
	if tmpChunks[0].Signature != updatedChunks[2].Signature {
		t.Fatalf("expected tmpChunks[0].Signature == updatedChunks[2].Signature, got:\n%#x != %#x", tmpChunks[0].Signature, updatedChunks[2].Signature)
	}

	// Rest of the chunks are expected to match with old data.
	for i := 2; i < len(origChunks); i++ {
		if origChunks[i].Signature != updatedChunks[i+1].Signature {
			t.Fatalf("expected origChunks[%d].Signature == updatedChunks[%d].Signature, got:\n%#x != %#x", i, i+1, origChunks[i].Signature, updatedChunks[i+1].Signature)
		}
	}
}

func Test_Signatures_Delete_Second_Segment(t *testing.T) {
	oldData := randomBytes(t, *seed, 5*fastcdc.MaxSize)

	origChunks := Signatures(oldData)

	// Clone data for modification.
	newData := make([]byte, len(oldData))
	copy(newData, oldData)

	for i := 0; i < len(oldData); i++ {
		if oldData[i] != newData[i] {
			t.Fatalf("expected oldData[%d] == newData[%d], got %#x != %#x", i, i, oldData[i], newData[i])
		}
	}

	// Find offset for second chunk segment.
	offsetSecond := len(origChunks[0].Bytes)
	lengthSecond := len(origChunks[1].Bytes)
	offsetThird := offsetSecond + lengthSecond

	// Delete second segment.
	newData = append(newData[:offsetSecond], newData[offsetThird:]...)

	updatedChunks := Signatures(newData)

	// Since first segment was not touched. It's expected to stay the say.
	if origChunks[0].Signature != updatedChunks[0].Signature {
		t.Fatalf("expected origChunks[0].Signature == updatedChunks[0].Signature, got:\n%#x != %#x", origChunks[0].Signature, updatedChunks[0].Signature)
	}

	// Starting from third segment, rest of the chunks are expected to match,
	// but offsetting.
	for i := 2; i < len(origChunks); i++ {
		if origChunks[i].Signature != updatedChunks[i-1].Signature {
			t.Fatalf("expected origChunks[%d].Signature != updatedChunks[%d].Signature, got %#x != %#x", i, i-1, origChunks[i].Signature, updatedChunks[i-1].Signature)
		}
	}
}
