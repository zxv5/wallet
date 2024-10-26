package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOffsetLimit(t *testing.T) {
	tests := []struct {
		page           int64
		size           int64
		expectedOffset int64
		expectedLimit  int64
	}{
		{0, 0, 0, 10},   // Default values
		{1, 0, 0, 10},   // Default size
		{0, 5, 0, 5},    // Default page
		{2, 5, 5, 5},    // Normal case
		{3, 10, 20, 10}, // Larger size
		{-1, 10, 0, 10}, // Negative page
		{1, -5, 0, 10},  // Negative size
		{5, 0, 40, 10},  // Page > 1 with default size
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("page=%d,size=%d", tt.page, tt.size), func(t *testing.T) {
			offset, limit := GetOffsetLimit(tt.page, tt.size)
			assert.Equal(t, tt.expectedOffset, offset)
			assert.Equal(t, tt.expectedLimit, limit)
		})
	}
}
