//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"sort"

	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/pmx_translator/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/pmx_translator/pkg/domain"
)

// CsvTableView はCSV出力タブの一覧を表す。
type CsvTableView struct {
	*walk.TableView
	translator i18n.II18n
	model      *CsvNameTableModel
}

// NewCsvTableView は CsvTableView を生成する。
func NewCsvTableView(translator i18n.II18n) *CsvTableView {
	return &CsvTableView{
		translator: translator,
		model:      NewCsvNameTableModel(translator),
	}
}

// Widgets はUI構成を返す。
func (tv *CsvTableView) Widgets() declarative.Composite {
	return declarative.Composite{
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.TableView{
				AssignTo:         &tv.TableView,
				AlternatingRowBG: true,
				CheckBoxes:       true,
				ColumnsOrderable: true,
				MultiSelection:   true,
				Model:            tv.model,
				MinSize:          declarative.Size{Width: 400, Height: 250},
				Columns: []declarative.TableViewColumn{
					{Title: "#", Width: 50},
					{Title: "No.", Width: 50},
					{Title: tv.t(messages.LabelTableType), Width: 80},
					{Title: tv.t(messages.LabelTableSegmented), Width: 60},
					{Title: tv.t(messages.LabelTableJapaneseName), Width: 170},
					{Title: tv.t(messages.LabelTableEnglishName), Width: 170},
				},
				StyleCell: func(style *walk.CellStyle) {
					if tv.model.Checked(style.Row()) {
						style.BackgroundColor = walk.RGB(159, 255, 243)
						return
					}
					style.BackgroundColor = walk.RGB(255, 255, 255)
				},
			},
		},
	}
}

// ResetRows は一覧行を再構築する。
func (tv *CsvTableView) ResetRows(items []domain.CsvCandidateItem) {
	if tv == nil || tv.model == nil {
		return
	}
	tv.model.ResetRows(items)
}

// CheckedNames はチェック済み名称を返す。
func (tv *CsvTableView) CheckedNames() []string {
	if tv == nil || tv.model == nil {
		return []string{}
	}
	return tv.model.CheckedNames()
}

// t は翻訳済み文言を返す。
func (tv *CsvTableView) t(key string) string {
	return i18n.TranslateOrMark(tv.translator, key)
}

// CsvNameTableModel はCSV出力テーブルのモデルを表す。
type CsvNameTableModel struct {
	walk.TableModelBase
	walk.SorterBase
	translator i18n.II18n
	sortColumn int
	sortOrder  walk.SortOrder
	records    []*domain.CsvCandidateItem
}

// NewCsvNameTableModel は CsvNameTableModel を生成する。
func NewCsvNameTableModel(translator i18n.II18n) *CsvNameTableModel {
	return &CsvNameTableModel{
		translator: translator,
		sortColumn: 1,
		sortOrder:  walk.SortAscending,
		records:    []*domain.CsvCandidateItem{},
	}
}

// RowCount は行数を返す。
func (m *CsvNameTableModel) RowCount() int {
	return len(m.records)
}

// Value はセルの値を返す。
func (m *CsvNameTableModel) Value(row int, col int) interface{} {
	item := m.records[row]

	switch col {
	case 0:
		return item.Checked
	case 1:
		return item.Number
	case 2:
		return typeLabelByKey(m.translator, item.TypeKey)
	case 3:
		return item.Segmented
	case 4:
		return item.NameText
	case 5:
		return item.EnglishNameText
	default:
		return ""
	}
}

// Checked はチェック状態を返す。
func (m *CsvNameTableModel) Checked(row int) bool {
	return m.records[row].Checked
}

// SetChecked はチェック状態を設定する。
func (m *CsvNameTableModel) SetChecked(row int, checked bool) error {
	m.records[row].Checked = checked
	return nil
}

// ColumnSortable はソート可否を返す。
func (m *CsvNameTableModel) ColumnSortable(col int) bool {
	return col >= 0
}

// SortedColumn は現在のソート列を返す。
func (m *CsvNameTableModel) SortedColumn() int {
	return m.sortColumn
}

// SortOrder は現在のソート順を返す。
func (m *CsvNameTableModel) SortOrder() walk.SortOrder {
	return m.sortOrder
}

// Sort は指定列で行を並び替える。
func (m *CsvNameTableModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn = col
	m.sortOrder = order

	sort.SliceStable(m.records, func(i int, j int) bool {
		a := m.records[i]
		b := m.records[j]

		less := false
		switch m.sortColumn {
		case 0:
			less = boolToInt(a.Checked) < boolToInt(b.Checked)
		case 1:
			less = a.Number < b.Number
		case 2:
			less = typeLabelByKey(m.translator, a.TypeKey) < typeLabelByKey(m.translator, b.TypeKey)
		case 3:
			less = boolToInt(a.Segmented) < boolToInt(b.Segmented)
		case 4:
			less = a.NameText < b.NameText
		case 5:
			less = a.EnglishNameText < b.EnglishNameText
		default:
			less = a.Number < b.Number
		}

		if m.sortOrder == walk.SortAscending {
			return less
		}
		return !less
	})

	return m.SorterBase.Sort(col, order)
}

// ResetRows は行一覧を置き換える。
func (m *CsvNameTableModel) ResetRows(items []domain.CsvCandidateItem) {
	m.records = make([]*domain.CsvCandidateItem, 0, len(items))
	for i := range items {
		item := items[i]
		m.records = append(m.records, &item)
	}
	m.PublishRowsReset()
}

// CheckedNames はチェック済み名称を返す。
func (m *CsvNameTableModel) CheckedNames() []string {
	names := make([]string, 0, len(m.records))
	for _, item := range m.records {
		if item == nil || !item.Checked {
			continue
		}
		names = append(names, item.NameText)
	}
	return names
}
