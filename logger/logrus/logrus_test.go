package logrus

import "testing"

func TestLogrus(t *testing.T) {
	f := WithLevel("")
	_, err := New(f)
	if err == nil {
		t.Errorf("Not testing empty log level")
	}
}
