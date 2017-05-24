package main

import (
	"testing"
	"strings"
	"bytes"
	log "github.com/sirupsen/logrus"
)

func TestAppendToFstab(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	if err := AppendToFstab("test_label", "ext999", "/dummy/path", true); err != nil {
		t.Errorf("Error: %v", err)
	}

	bufString := buf.String()
	if ! strings.Contains(bufString, "FSTAB: would have appended: LABEL=test_label /dummy/path ext999 defaults 0 1") {
	    t.Errorf("printed wrong thing to stderr. Actual: %s", bufString)
	}
}