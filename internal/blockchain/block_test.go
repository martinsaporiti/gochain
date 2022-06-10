package blockchain

import (
	"testing"
	"time"
)

func TestHash(t *testing.T) {

	dateString := "2021-11-22"
	date, _ := time.Parse("2006-01-02", dateString)

	timestamp := date.UnixNano()
	block := NewBlock(1, 0, [32]byte{}, []*Transaction{})
	block.timestamp = timestamp
	tests := map[string]struct {
		input *Block
		want  [32]byte
	}{
		"hash sould return a correct hash array": {
			input: block,
			want:  [32]byte{217, 253, 26, 86, 184, 176, 55, 115, 121, 40, 209, 126, 49, 176, 168, 27, 194, 43, 24, 246, 115, 224, 202, 180, 3, 72, 95, 170, 239, 121, 118, 90},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input.Hash()
			if got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
