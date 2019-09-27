package gif

import (
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()

	err := DrawGif(DefaultFace(), []string{
		"ba ",
		"   DUM!",
		" Tss",
	}, []int{10, 50, 100}, f)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
