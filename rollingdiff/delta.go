package rollingdiff

import "crypto/sha256"

type Operation int

const (
	Nop    Operation = iota
	Delete Operation = iota
	Add    Operation = iota
	Move   Operation = iota
)

type Change struct {
	Op    Operation
	From  int
	To    int
	Bytes []byte
}

// Delta computes difference between two lists of chunks. It returns list of
// changes that need to be performed to src in order to result with dst.
func Delta(src, dst []Chunk) []Change {
	changes := make([]Change, 0)

	mSrc := make(map[[sha256.Size]byte]Chunk, len(src))
	for _, c := range src {
		mSrc[c.Signature] = c
	}

	mDst := make(map[[sha256.Size]byte]Chunk, len(dst))
	for _, c := range dst {
		mDst[c.Signature] = c
	}

	// Add changes for deleted chunks.
	for i := 0; i < len(src); i++ {
		_, exists := mDst[src[i].Signature]
		if !exists {
			c := Change{
				Op:   Delete,
				From: src[i].Index,
			}
			changes = append(changes, c)
			delete(mSrc, src[i].Signature)
			src = append(src[:i], src[i+1:]...)
			i--
		}
	}

	// Add changes for added chunks.
	for i := 0; i < len(dst); i++ {
		_, exists := mSrc[dst[i].Signature]
		if !exists {
			c := Change{
				Op:    Add,
				To:    dst[i].Index,
				Bytes: dst[i].Bytes,
			}
			changes = append(changes, c)
			delete(mDst, dst[i].Signature)
			dst = append(dst[:i], dst[i+1:]...)
			i--
		}
	}

	// Finally check Moved chunks.
	for i, c := range src {
		if dst[i].Signature == c.Signature {
			continue
		}

		v, exists := mDst[c.Signature]
		if exists {
			// Create Move change on first occurrance.
			chg := Change{
				Op:   Move,
				From: c.Index,
				To:   v.Index,
			}
			changes = append(changes, chg)
			delete(mDst, c.Signature)
		}
	}

	return changes
}
