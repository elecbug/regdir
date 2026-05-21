package parameter

import (
	"testing"
)

// TestParseParams tests the ParseParams function to ensure it correctly parses CLI input into a Parameters struct.
func TestParseParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Parameters
	}{
		{
			name:  "No parameters",
			input: "myprogram",
			expected: &Parameters{
				Program:   "myprogram",
				Params:    map[string]string{},
				MainParam: "",
			},
		},
		{
			name:  "With key-value parameters",
			input: "myprogram --key1 value1 -k2 value2",
			expected: &Parameters{
				Program: "myprogram",
				Params: map[string]string{
					"key1": "value1",
					"k2":   "value2",
				},
				MainParam: "",
			},
		},
		{
			name:  "With main parameter",
			input: "myprogram mainparam --key1 value1",
			expected: &Parameters{
				Program: "myprogram",
				Params: map[string]string{
					"key1": "value1",
				},
				MainParam: "mainparam",
			},
		},
		{
			name:  "With flag parameters",
			input: "myprogram --flag1 -f2",
			expected: &Parameters{
				Program: "myprogram",
				Params: map[string]string{
					"flag1": "true",
					"f2":    "true",
				},
				MainParam: "",
			},
		},
		{
			name:  "With mixed parameters",
			input: "myprogram --key1 value1 mainparam -f2 --flag3",
			expected: &Parameters{
				Program: "myprogram",
				Params: map[string]string{
					"key1":  "value1",
					"f2":    "true",
					"flag3": "true",
				},
				MainParam: "mainparam",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Parse(tt.input)
			if !equal(tt.expected, result) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func equal(a, b *Parameters) bool {
	if a.Program != b.Program || a.MainParam != b.MainParam {
		return false
	}

	if len(a.Params) != len(b.Params) {
		return false
	}

	for key, valueA := range a.Params {
		valueB, exists := b.Params[key]
		if !exists || valueA != valueB {
			return false
		}
	}

	return true
}
