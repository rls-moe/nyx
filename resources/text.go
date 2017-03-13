package resources

import (
	"compress/flate"
	"fmt"
	"html/template"
	"io"
	"math"
	"math/rand"
	"strings"
)

func OperateReplyText(unsafe string) template.HTML {
	unsafe = template.HTMLEscapeString(unsafe)
	unsafe = strings.Replace(unsafe, "\n", "<br />", -1)
	return template.HTML(unsafe)
}

var (
	blacklist = []string{
		"spam",
		"pizza",
		"buy",
		"free",
		"subscription",
		"penis",
		"nazi",
	}
)

func SpamScore(spam string) (float64, error) {
	spam = strings.ToLower(spam)

	counter := &byteCounter{1}
	compressor, err := flate.NewWriter(counter, flate.BestSpeed)
	if err != nil {
		return 0.0, err
	}
	_, err = io.WriteString(compressor, spam)
	if err != nil {
		return 0.0, err
	}
	compressor.Flush()
	compressor.Close()
	blScore := 1.0
	for _, v := range blacklist {
		blScore += float64(strings.Count(spam, v))
	}

	score := float64(len(spam)) / float64(counter.p)

	return (score * blScore) / 100, nil
}

type byteCounter struct {
	p int
}

func (b *byteCounter) Write(p []byte) (n int, err error) {
	b.p += len(p)
	return len(p), nil
}

func CaptchaPass(spamScore float64) bool {
	chance := math.Max(0, math.Min(0.65*math.Atan(7.1*spamScore), 0.99))
	take := rand.Float64()
	fmt.Printf("Chance: %f, Take %f", chance, take)
	return take > chance
}
