package resources

import (
	"compress/flate"
	"html/template"
	"io"
	"log"
	"math"
	"math/rand"
	"strings"
)

func OperateReplyText(unsafe string) template.HTML {
	unsafe = template.HTMLEscapeString(unsafe)
	unsafe = strings.Replace(unsafe, "\n", "<br />", -1)
	return template.HTML(unsafe)
}

const (
	passScoreAggressive = 7.1
	passScoreReactive   = 0.65
	passScoreLimitMin   = 0.01
	passScoreLimitMax   = 0.99
)

var (
	blacklist = []string{
		"spam",
		"pizza",
		"buy",
		"free",
		"subscription",
		"penis",
		"nazi",
		"beemovie",
		"bee movie",
		"honey",
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

	lines := strings.Count(spam, "\n")

	if lines == 0 {
		lines = 1
	}

	score := float64(len(spam)*lines) / float64(counter.p)

	return (score * blScore) / 100, nil
}

type byteCounter struct {
	p int
}

func (b *byteCounter) Write(p []byte) (n int, err error) {
	b.p += len(p)
	return len(p), nil
}

func CaptchaProb(spamScore float64) float64 {
	return math.Max(
		passScoreLimitMin,
		math.Min(
			passScoreReactive*math.Atan(
				passScoreAggressive*spamScore,
			),
			passScoreLimitMax))
}

func CaptchaPass(spamScore float64) bool {
	chance := CaptchaProb(spamScore)
	take := rand.Float64()
	log.Printf("Chance: %f, Take %f\n", chance, take)
	return take > chance
}
