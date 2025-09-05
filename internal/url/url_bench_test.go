package url

import (
	"context"
	"github.com/thxhix/shortener/internal/config"
	"github.com/thxhix/shortener/internal/database/drivers"
	"strconv"
	"testing"
)

func BenchmarkShorten(b *testing.B) {
	cfg := config.Config{BaseURL: "http://localhost:8080"}
	db, _ := drivers.NewFileDatabase("./tmp_bench.json")
	uc := NewURLUseCase(db, cfg)

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = uc.Shorten(ctx, "https://example.com/bench"+strconv.Itoa(i))
	}
}
