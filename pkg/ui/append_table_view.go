package ui

import (
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/pmx_translator/pkg/domain"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type AppendTableView struct {
	*declarative.TableView
	Model *AppendNameModel
}

type AppendNameModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	Records    []*domain.NameItem
}

func NewAppendNameModel(originalCsv, appendCsv *core.CsvModel) *AppendNameModel {
	m := new(AppendNameModel)
	m.ResetRows(originalCsv, appendCsv)
	return m
}

func (m *AppendNameModel) CheckedItems() []*domain.NameItem {
	var items []*domain.NameItem
	for _, item := range m.Records {
		if item.Checked {
			items = append(items, item)
		}
	}
	return items
}

func (m *AppendNameModel) RowCount() int {
	return len(m.Records)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *AppendNameModel) Value(row, col int) interface{} {
	item := m.Records[row]

	switch col {
	case 0:
		return item.Checked
	case 1:
		return item.Number
	case 2:
		return item.TypeText
	case 3:
		return item.NameText
	case 4:
		return item.EnglishNameText
	}

	panic("unexpected col")
}

func (m *AppendNameModel) Checked(row int) bool {
	return m.Records[row].Checked
}

func (m *AppendNameModel) IsOriginal(row int) bool {
	return m.Records[row].IsOriginal
}

func (m *AppendNameModel) SetChecked(row int, Checked bool) error {
	m.Records[row].Checked = Checked

	return nil
}

func (m *AppendNameModel) Sort(col int, order walk.SortOrder) error {
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
			return c(a.TypeText < b.TypeText)
		case 3:
			return c(a.NameText < b.NameText)
		case 4:
			return c(a.EnglishNameText < b.EnglishNameText)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *AppendNameModel) exists(txt string) bool {
	for _, item := range m.Records {
		if item.TypeText == txt {
			return true
		}
	}
	return false
}

func (m *AppendNameModel) AddRecord(record []string, isOriginal bool) {
	item := &domain.NameItem{
		Checked:         !m.exists(record[1]),
		Number:          len(m.Records) + 1,
		TypeText:        record[1],
		NameText:        record[2],
		EnglishNameText: record[3],
		IsOriginal:      isOriginal,
	}
	m.Records = append(m.Records, item)
}

func (m *AppendNameModel) ResetRows(originalCsv, appendCsv *core.CsvModel) {
	m.Records = make([]*domain.NameItem, 0)

	m.PublishRowsReset()

	if originalCsv == nil || appendCsv == nil {
		return
	}

	for n, record := range originalCsv.Records() {
		if n == 0 {
			continue
		}
		m.AddRecord(record, true)
	}

	for n, record := range appendCsv.Records() {
		if n == 0 {
			continue
		}
		m.AddRecord(record, false)
	}

	m.PublishRowsReset()
}

func NewAppendTableView(parent walk.Container, originalCsv, appendCsv *core.CsvModel) *AppendTableView {
	nameModel := NewAppendNameModel(originalCsv, appendCsv)

	var tv *walk.TableView
	builder := declarative.NewBuilder(parent)

	dTableView := &declarative.TableView{
		AssignTo:         &tv,
		AlternatingRowBG: true,
		CheckBoxes:       true,
		ColumnsOrderable: true,
		MultiSelection:   true,
		Model:            nameModel,
		MinSize:          declarative.Size{Width: 512, Height: 250},
		Columns: []declarative.TableViewColumn{
			{Title: "#", Width: 50},
			{Title: "No.", Width: 50},
			{Title: mi18n.T("元名称"), Width: 200},
			{Title: mi18n.T("日本語名称"), Width: 200},
			{Title: mi18n.T("英語名称"), Width: 200},
		},
		StyleCell: func(style *walk.CellStyle) {
			if nameModel.Checked(style.Row()) {
				if nameModel.IsOriginal(style.Row()) {
					// 黄色(元データ)
					style.BackgroundColor = walk.RGB(239, 255, 160)
				} else {
					// 水色(追加データで追加対象)
					style.BackgroundColor = walk.RGB(159, 255, 243)
				}
			} else {
				style.BackgroundColor = walk.RGB(255, 255, 255)
			}
		},
	}

	if err := dTableView.Create(builder); err != nil {
		panic(err)
	}

	nameTableView := &AppendTableView{
		TableView: dTableView,
		Model:     nameModel,
	}

	return nameTableView
}

func (tv *AppendTableView) ResetModel(originalCsv, appendCsv *core.CsvModel) {
	tv.Model.ResetRows(originalCsv, appendCsv)
}
