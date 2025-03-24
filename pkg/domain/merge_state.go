package domain

import (
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/mcsv"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
	"github.com/miu200521358/walk/pkg/walk"
)

type MergeState struct {
	OriginalCsvModel *mcsv.CsvModel  // 追加元Csvモデル
	MergedCsvModel   *mcsv.CsvModel  // 追加先Csvモデル
	NameModel        *MergeNameModel // 名称モデル
	TextChangeDialog *walk.Dialog    // テキスト変更ダイアログ
	OutputPath       string          // 出力パス
}

func NewMergeState() *MergeState {
	return &MergeState{
		NameModel: new(MergeNameModel),
	}
}

func (c *MergeState) LoadData() {
	if c.OriginalCsvModel == nil || c.MergedCsvModel == nil {
		return
	}

	c.NameModel.ResetRows(c.OriginalCsvModel, c.MergedCsvModel)
	c.OutputPath = mfile.CreateOutputPath(c.OriginalCsvModel.Path(), "")
}

func (m *MergeNameModel) CheckedNames() []string {
	var names []string
	for _, item := range m.Records {
		if item.Checked {
			names = append(names, item.NameText)
		}
	}
	return names
}

// --------------------------------------------------

type MergeNameModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	Records    []*NameItem
}

func (m *MergeNameModel) RowCount() int {
	return len(m.Records)
}

func (m *MergeNameModel) Value(row, col int) any {
	item := m.Records[row]

	switch col {
	case 0:
		return item.Checked
	case 1:
		return item.Number
	case 2:
		return item.NameText
	case 3:
		return item.JapaneseNameText
	case 4:
		return item.EnglishNameText
	}

	panic("unexpected col")

}

func (m *MergeNameModel) Checked(row int) bool {
	return m.Records[row].Checked
}

func (m *MergeNameModel) SetChecked(row int, checked bool) error {
	m.Records[row].Checked = checked

	return nil
}

func (m *MergeNameModel) IsOriginal(row int) bool {
	return m.Records[row].IsOriginal
}

func (m *MergeNameModel) CheckedItems() []*NameItem {
	var items []*NameItem
	for _, item := range m.Records {
		if item.Checked {
			items = append(items, item)
		}
	}
	return items
}

func (m *MergeNameModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.Records, func(i, j int) bool {
		a, b := m.Records[i], m.Records[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			av := 0
			if a.Checked {
				av = 1
			}
			bv := 0
			if b.Checked {
				bv = 1
			}
			return c(av < bv)
		case 1:
			return c(a.Number < b.Number)
		case 2:
			return c(a.NameText < b.NameText)
		case 3:
			return c(a.JapaneseNameText < b.JapaneseNameText)
		case 4:
			return c(a.EnglishNameText < b.EnglishNameText)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)

}

// --------------------------------------------------

func (m *MergeNameModel) exists(txt string) bool {
	for _, item := range m.Records {
		if item.TypeText == txt {
			return true
		}
	}
	return false
}

func (m *MergeNameModel) AddRecord(record []string, isOriginal bool) {
	item := &NameItem{
		Checked:          !m.exists(record[1]),
		Number:           len(m.Records) + 1,
		NameText:         record[1],
		JapaneseNameText: record[2],
		EnglishNameText:  record[3],
		IsOriginal:       isOriginal,
	}
	m.Records = append(m.Records, item)
}

func (m *MergeNameModel) ResetRows(originalCsvModel, mergedCsvModel *mcsv.CsvModel) {
	m.Records = make([]*NameItem, 0)

	m.PublishRowsReset()

	if originalCsvModel == nil || mergedCsvModel == nil {
		return
	}

	for n, record := range originalCsvModel.Records() {
		if n == 0 {
			continue
		}
		m.AddRecord(record, true)
	}

	for n, record := range mergedCsvModel.Records() {
		if n == 0 {
			continue
		}
		m.AddRecord(record, false)
	}

	m.PublishRowsReset()

}
