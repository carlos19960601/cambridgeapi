package cambridgeapi

import (
	"os"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog"
)

const (
	baseUrl = "https://dictionary.cambridge.org"
)

type Client interface {
	Query(q string) (*Word, error)
	RebuildDatabase(filepath string) error
}

type client struct {
	log      zerolog.Logger
	collecor *colly.Collector
}

func New() Client {
	return &client{
		log: zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
		collecor: colly.NewCollector(
			colly.AllowedDomains("dictionary.cambridge.org"),
			colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"),
		),
	}
}
