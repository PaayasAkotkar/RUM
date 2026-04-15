package searchmanager

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/milvus-io/milvus/client/v2/column"
)

type IPushReport struct {
	// same shape as your s3 CommitEntry
	Hash      string   // short random id, same idea as s3manager.GenerateHash()
	Bucket    string   // milvus collection == s3 bucket
	Branch    string   // branch filter value
	Objects   []string // every id that was inserted, format: "<bucket>/<branch>/<id>"
	Total     int
	Dim       int
	CreatedAt time.Time
	Err       error
}

// IsSucced keeps the same API your code already uses
func (r *IPushReport) IsSucced() bool { return r.Err == nil }

// Report prints a human-readable summary identical in spirit to insertAndReport
func (r *IPushReport) Report() string {
	if r.Err != nil {
		return fmt.Sprintf("milvus push failed: %v", r.Err)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("--- Milvus Push Report [%s] ---\n", r.Bucket))
	sb.WriteString(fmt.Sprintf("Hash      : %s\n", r.Hash))
	sb.WriteString(fmt.Sprintf("Branch    : %s\n", r.Branch))
	sb.WriteString(fmt.Sprintf("Total     : %d\n", r.Total))
	sb.WriteString(fmt.Sprintf("Dim       : %d\n", r.Dim))
	sb.WriteString(fmt.Sprintf("CreatedAt : %s\n\n", r.CreatedAt.Format(time.RFC3339)))
	for i, obj := range r.Objects {
		sb.WriteString(fmt.Sprintf("  [%d] %s\n", i, obj))
	}
	sb.WriteString("\nsucceed 🤗\n")
	return sb.String()
}

func generateMilvusHash() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("mlv-%x", b)
}

func dotPrd(a, b []float32) float32 {
	var dot float32
	for i := range a {
		dot += a[i] * b[i]
	}
	return dot
}

// fieldIndex returns the index of the named field, or -1 if not found.
func fieldIndex(fields []column.Column, name string) int {
	for i, f := range fields {
		if f.Name() == name {
			return i
		}
	}
	return -1
}

// branchFilter returns query for given branches
func branchFilter(branches []string) string {
	if len(branches) == 0 {
		return ""
	}
	if len(branches) == 1 {
		return fmt.Sprintf(`%s == "%s"`, FieldBranch, branches[0])
	}
	quoted := make([]string, len(branches))
	for i, b := range branches {
		quoted[i] = fmt.Sprintf(`"%s"`, b)
	}
	return fmt.Sprintf(`%s in [%s]`, FieldBranch, strings.Join(quoted, ", "))
}
