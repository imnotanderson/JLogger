package JLogger

import "testing"

func TestNew(t *testing.T) {
	log := New("/Users/imnotanderson/", "test", 2)
	log.Info("hello world")
}
