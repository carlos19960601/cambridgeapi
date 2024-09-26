package cambridgeapi

import "gorm.io/datatypes"

type Word struct {
	Word    string                          `gorm:"primarykey;column:word;size:128" json:"word"`
	URL     string                          `gorm:"column:url;size:256" json:"url"`
	Entries datatypes.JSONSlice[*WordEntry] `gorm:"column:entries" json:"entries"`
}

func (u *Word) TableName() string {
	return "words"
}

type WordEntry struct {
	Entry          string           `json:"entry"`
	POS            string           `json:"pos"`
	Pronunciations []*Pronunciation `json:"pronunciations"`
	Dsenses        []*Dsense        `json:"dsenses"`
}

type Pronunciation struct {
	Lang string `json:"lang"`
	URL  string `json:"url"`
	Pron string `json:"pron"`
}

type Dsense struct {
	Guide     string      `json:"guide"`
	DefBlocks []*DefBlock `json:"def_blocks"`
}

type Example struct {
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

type DefBlock struct {
	CEFRLevel   string     `json:"cefr_level"`
	Gram        string     `json:"gram"`
	Usage       string     `json:"usage"`
	Text        string     `json:"text"`
	Translation string     `json:"translation"`
	Examples    []*Example `json:"examples"`
}

// func (w *Word) FromWordEntries(entries []*WordEntryModel) {
// 	if len(entries) == 0 {
// 		return
// 	}

// 	w.Entry = entries[0].Entry

// 	pronunciations := make([]*Pronunciation, 0)
// 	for _, p := range entries[0].Pronunciations {
// 		pronunciation := &Pronunciation{
// 			Lang: p.Lang,
// 			URL:  p.URL,
// 			Pron: p.Pron,
// 		}
// 		pronunciations = append(pronunciations, pronunciation)
// 	}
// 	w.Pronunciations = pronunciations

// 	definitions := make([]*Definition, 0)
// 	for _, entry := range entries {
// 		definition := &Definition{}
// 		definition.POS = entry.POS

// 		examples := make([]*Example, 0)
// 		for _, dsense := range entry.Dsenses {
// 			for _, defBlock := range dsense.DefBlocks {
// 				definition.Text = defBlock.Text
// 				definition.Translation = defBlock.Translation
// 				for _, e := range defBlock.Examples {
// 					example := &Example{
// 						Text:        e.Text,
// 						Translation: e.Translation,
// 					}
// 					examples = append(examples, example)
// 				}
// 			}
// 		}
// 		definitions = append(definitions, definition)
// 	}
// 	w.Definitions = definitions
// }
