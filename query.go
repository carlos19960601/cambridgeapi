package cambridgeapi

import (
	"fmt"
	"net/url"

	"github.com/gocolly/colly/v2"
)

const (
	detailPath = "/dictionary/english-chinese-simplified/%s"
)

func (c *Client) Query(q string) (*Word, error) {
	urlResult, err := url.JoinPath(baseUrl, fmt.Sprintf(detailPath, q))
	if err != nil {
		return nil, err
	}

	res := &Word{}
	res.Word = q

	wordEntries := make([]*WordEntry, 0)
	queryCollector := c.collecor.Clone()
	queryCollector.OnHTML("div.entry-body", func(e *colly.HTMLElement) {
		res.URL = e.Request.URL.String()
		e.ForEach("div.entry-body__el", func(i int, h *colly.HTMLElement) {
			wordEntry := &WordEntry{}
			c.handleWordEntry(h, wordEntry)
			wordEntries = append(wordEntries, wordEntry)
		})
	})

	err = queryCollector.Visit(urlResult)
	if err != nil {
		return nil, err
	}
	if len(wordEntries) == 0 {
		return nil, fmt.Errorf("word %s not found", q)
	}

	res.Entries = wordEntries

	return res, nil
}

func (c *Client) handleWordEntry(e *colly.HTMLElement, wordEntry *WordEntry) {
	wordEntry.Entry = e.ChildText("span.hw.dhw")
	wordEntry.POS = e.ChildText("span.pos.dpos")

	ukPron := &Pronunciation{}
	ukPron.Lang = e.ChildText("span.uk.dpron-i span.region.dreg")
	ukPron.Pron = e.ChildText("span.uk.dpron-i span.pron.dpron")
	ukPron.URL = e.ChildAttr("span.uk.dpron-i source", "src")

	usPron := &Pronunciation{}
	usPron.Lang = e.ChildText("span.us.dpron-i span.region.dreg")
	usPron.Pron = e.ChildText("span.us.dpron-i span.pron.dpron")
	usPron.URL = e.ChildAttr("span.us.dpron-i source", "src")

	wordEntry.Pronunciations = []*Pronunciation{ukPron, usPron}

	dsenses := make([]*Dsense, 0)
	e.ForEach("div.pr.dsense", func(i int, h *colly.HTMLElement) {
		dsense := &Dsense{}
		dsense.Guide = h.ChildText("span.guideword.dsense_gw span")

		defBlocks := make([]*DefBlock, 0)
		h.ForEach("div.def-block.ddef_block", func(i int, h *colly.HTMLElement) {
			defBlock := &DefBlock{}
			defBlock.CEFRLevel = h.ChildText("span.epp-xref.dxref")
			defBlock.Gram = h.ChildText("span.gram.dgram")
			defBlock.Usage = h.ChildText("span.usage.dusage")
			defBlock.Text = h.ChildText("div.def.ddef_d")
			defBlock.Translation = h.ChildText("span.trans.dtrans")

			examples := make([]*Example, 0)
			h.ForEach("div.examp.dexamp", func(i int, h *colly.HTMLElement) {
				example := &Example{}
				example.Text = h.ChildText("span.eg.deg")
				example.Translation = h.ChildText("span.trans.dtrans")
				examples = append(examples, example)
			})
			defBlock.Examples = examples

			defBlocks = append(defBlocks, defBlock)
		})
		dsense.DefBlocks = defBlocks
		dsenses = append(dsenses, dsense)
	})

	wordEntry.Dsenses = dsenses
}
