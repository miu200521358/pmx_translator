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

// AppendTableView はCSV追加タブの一覧を表す。
type AppendTableView struct {
	*walk.TableView
	translator i18n.II18n
	model      *AppendNameTableModel
}

// NewAppendTableView は AppendTableView を生成する。
func NewAppendTableView(translator i18n.II18n) *AppendTableView {
	return &AppendTableView{
		translator: translator,
		model:      NewAppendNameTableModel(),
	}
}

// Widgets はUI構成を返す。
func (tv *AppendTableView) Widgets() declarative.Composite {
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
					{Title: tv.t(messages.LabelTableSourceName), Width: 180},
					{Title: tv.t(messages.LabelTableJapaneseName), Width: 180},
					{Title: tv.t(messages.LabelTableEnglishName), Width: 180},
				},
				StyleCell: func(style *walk.CellStyle) {
					if tv.model.Checked(style.Row()) {
						if tv.model.IsOriginal(style.Row()) {
							style.BackgroundColor = walk.RGB(239, 255, 160)
							return
						}
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
func (tv *AppendTableView) ResetRows(items []domain.AppendNameItem) {
	if tv == nil || tv.model == nil {
		return
	}
	tv.model.ResetRows(items)
}

// Rows は現在の一覧行を返す。
func (tv *AppendTableView) Rows() []domain.AppendNameItem {
	if tv == nil || tv.model == nil {
		return []domain.AppendNameItem{}
	}
	return tv.model.Rows()
}

// t は翻訳済み文言を返す。
func (tv *AppendTableView) t(key string) string {
	return i18n.TranslateOrMark(tv.translator, key)
}

// AppendNameTableModel はCSV追加テーブルのモデルを表す。
type AppendNameTableModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	records    []*domain.AppendNameItem
}

// NewAppendNameTableModel は AppendNameTableModel を生成する。
func NewAppendNameTableModel() *AppendNameTableModel {
	return &AppendNameTableModel{
		sortColumn: 1,
		sortOrder:  walk.SortAscending,
		records:    []*domain.AppendNameItem{},
	}
}

// RowCount は行数を返す。
func (m *AppendNameTableModel) RowCount() int {
	return len(m.records)
}

// Value はセルの値を返す。
func (m *AppendNameTableModel) Value(row int, col int) interface{} {
	item := m.records[row]

	switch col {
	case 0:
		return item.Checked
	case 1:
		return item.Number
	case 2:
		return item.SourceName
	case 3:
		return item.JapaneseName
	case 4:
		return item.EnglishName
	default:
		return ""
	}
}

// Checked はチェック状態を返す。
func (m *AppendNameTableModel) Checked(row int) bool {
	return m.records[row].Checked
}

// IsOriginal は元CSV行かを返す。
func (m *AppendNameTableModel) IsOriginal(row int) bool {
	return m.records[row].IsOriginal
}

// SetChecked はチェック状態を設定する。
func (m *AppendNameTableModel) SetChecked(row int, checked bool) error {
	m.records[row].Checked = checked
	return nil
}

// ColumnSortable はソート可否を返す。
func (m *AppendNameTableModel) ColumnSortable(col int) bool {
	return col >= 0
}

// SortedColumn は現在のソート列を返す。
func (m *AppendNameTableModel) SortedColumn() int {
	return m.sortColumn
}

// SortOrder は現在のソート順を返す。
func (m *AppendNameTableModel) SortOrder() walk.SortOrder {
	return m.sortOrder
}

// Sort は指定列で行を並び替える。
func (m *AppendNameTableModel) Sort(col int, order walk.SortOrder) error {
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
			less = a.SourceName < b.SourceName
		case 3:
			less = a.JapaneseName < b.JapaneseName
		case 4:
			less = a.EnglishName < b.EnglishName
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
func (m *AppendNameTableModel) ResetRows(items []domain.AppendNameItem) {
	m.records = make([]*domain.AppendNameItem, 0, len(items))
	for i := range items {
		item := items[i]
		m.records = append(m.records, &item)
	}
	m.PublishRowsReset()
}

// Rows は行一覧のコピーを返す。
func (m *AppendNameTableModel) Rows() []domain.AppendNameItem {
	rows := make([]domain.AppendNameItem, 0, len(m.records))
	for _, item := range m.records {
		if item == nil {
			continue
		}
		rows = append(rows, *item)
	}
	return rows
}
