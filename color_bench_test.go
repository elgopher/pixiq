package pixiq_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/jacekolszak/pixiq"
)

func BenchmarkRGBAi(b *testing.B) {
	var c pixiq.Color
	b.StopTimer()
	red := rand.Int()
	green := rand.Int()
	blue := rand.Int()
	alpha := rand.Int()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c = pixiq.RGBAi(red, green, blue, alpha)
	}
	b.StopTimer()
	fmt.Println(c.R() + c.G() + c.B() + c.A())
}
