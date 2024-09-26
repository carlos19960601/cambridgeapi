package cambridgeapi

import (
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (c *Client) RebuildDatabase(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil {
		_ = os.Remove(filepath)
	}

	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.Migrator().AutoMigrate(&Word{})
	if err != nil {
		return err
	}

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			c.log.Error().Err(err).Msg("get db")
			return
		}

		err = sqlDB.Close()
		if err != nil {
			c.log.Error().Err(err).Msg("close db")
		}
	}()

	rangeCollector := colly.NewCollector(
		colly.AllowedDomains("dictionary.cambridge.org"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"),
	)

	listCollector := rangeCollector.Clone()
	detailCollector := rangeCollector.Clone()

	detailCollector.Async = true
	detailCollector.Limit(&colly.LimitRule{Parallelism: 10})

	rangeCollector.OnHTML("a.hlh32.hdb.dil.tcbd", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.log.Info().Str("text", strings.TrimSpace(e.Text)).Str("link", link).Msg("range")

		listCollector.Visit(link)
	})

	listCollector.OnHTML("div.hlh32.han a.tc-bd", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// c.log.Info().Str("text", strings.TrimSpace(e.Text)).Str("link", link).Msg("list")

		detailCollector.Visit(e.Request.AbsoluteURL(link))
	})

	detailCollector.OnHTML("div.entry-body", func(e *colly.HTMLElement) {
		res := &Word{}
		res.URL = e.Request.URL.String()

		wordEntries := make([]*WordEntry, 0)
		e.ForEach("div.entry-body__el", func(i int, h *colly.HTMLElement) {
			wordEntry := &WordEntry{}
			c.handleWordEntry(h, wordEntry)
			wordEntries = append(wordEntries, wordEntry)
		})

		if len(wordEntries) == 0 {
			return
		}

		res.Word = wordEntries[0].Entry
		res.Entries = wordEntries

		if err := db.Save(res).Error; err != nil {
			c.log.Error().Err(err).Msg("save database")
		}
	})

	for i := 'a'; i <= 'z'; i++ {
		err = rangeCollector.Visit("https://dictionary.cambridge.org/browse/english-chinese-simplified/" + string(i))
		if err != nil {
			c.log.Error().Err(err).Msgf("collector visit")
		}
	}

	rangeCollector.Wait()

	return nil
}
