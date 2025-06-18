// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package loggers_test

import (
	"bufio"
	"encoding/json"
	"testing"

	"go.uber.org/zap/zaptest"

	"backend.brokedaear.com/pkg/assert"
)

func newBufAndWriter() (*zaptest.Buffer, *bufio.Writer) {
	var b zaptest.Buffer
	bWriter := bufio.NewWriter(&b)
	return &b, bWriter
}

// parseLogOutput is a helper function that parses JSON log lines from a buffer
// into a slice of logEntry structs.
func parseLogOutput(t *testing.T, b *zaptest.Buffer) ([]logEntry, error) {
	lines := b.Lines()
	entries := make([]logEntry, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		rawEntry := logEntry{
			Level:     "",
			Message:   "",
			Timestamp: 0,
			Caller:    "",
			Fields:    make(map[string]any),
		}

		err := json.Unmarshal([]byte(line), &rawEntry)
		if err != nil {
			return nil, err
		}

		result := make(map[string]any)
		json.Unmarshal([]byte(line), &result)

		t.Logf("Entry: %#v", rawEntry)
		t.Logf("Entry.Fields: %v", result)

		for k, v := range result {
			if k != "level" && k != "msg" && k != "ts" && k != "caller" {
				rawEntry.Fields[k] = v
				t.Logf("%s: %v", k, v)
			}
		}

		entries = append(entries, rawEntry)
	}

	return entries, nil
}

func TestParseLogOutput(t *testing.T) {
	s := `{"level":"info","ts":1750057415.0888262,"caller":"loggers/zap.go:235","msg":"staging info","value":123}
{"level":"warn","ts":1750057415.088834,"caller":"loggers/zap.go:265","msg":"staging warn"}`
	b, bWriter := newBufAndWriter()
	bWriter.WriteString(s)
	bWriter.Flush()
	t.Logf("Buffer contents: %v", b.String())
	out, err := parseLogOutput(t, b)
	assert.NoError(t, err)
	v1 := out[0]
	val, ok := v1.Fields["value"]
	if ok {
		assert.Equal(t, val.(float64), 123.0)
	} else {
		t.Error("wanted val does not exist")
	}
}
