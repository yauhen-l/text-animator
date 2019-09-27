package gif

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	start := time.Now()

	f, _ := os.OpenFile("rgb1.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()

	//input := "(_).(_)\n(_)o(_)\n(_)0(_)"
	//input := "х\nу\nй"
	//input := "8=э\n8==э\n8===э"
	input := "ba \n   DUM!\n Tss"
	//input := "хорошо\nжорошо\nокее"
	//input := "э\nе"
	//input := ":|\n:|\n:P"
	//input := "a\na\nb"
	err := DrawGif(DefaultFace(), strings.Split(input, "\n"), []int{10, 50, 100}, f)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	log.Printf("Duration: %v", time.Since(start))

	PrintMemUsage()
}
