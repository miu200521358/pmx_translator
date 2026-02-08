//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/pmx_translator/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/pmx_translator/pkg/domain"
)

// TranslateTableView は名称置換タブの一覧を表す。
type TranslateTableView struct {
	*walk.TableView
	translator i18n.II18n
	model      *TranslateNameTableModel
	editing    bool
}

// NewTranslateTableView は TranslateTableView を生成する。
func NewTranslateTableView(translator i18n.II18n) *TranslateTableView {
	return &TranslateTableView{
		translator: translator,
		model:      NewTranslateNameTableModel(translator),
	}
}

// Widgets はUI構成を返す。
func (tv *TranslateTableView) Widgets() declarative.Composite {
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
					{Title: tv.t(messages.LabelTableIndex), Width: 60},
					{Title: tv.t(messages.LabelTableSourceName), Width: 170},
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
				OnSelectedIndexesChanged: func() {
					tv.openEditDialog()
				},
			},
		},
	}
}

// ResetRows は一覧行を再構築する。
func (tv *TranslateTableView) ResetRows(items []domain.TranslateNameItem) {
	if tv == nil || tv.model == nil {
		return
	}
	tv.model.ResetRows(items)
}

// Rows は現在の一覧行を返す。
func (tv *TranslateTableView) Rows() []domain.TranslateNameItem {
	if tv == nil || tv.model == nil {
		return []domain.TranslateNameItem{}
	}
	return tv.model.Rows()
}

// t は翻訳済み文言を返す。
func (tv *TranslateTableView) t(key string) string {
	return i18n.TranslateOrMark(tv.translator, key)
}

// openEditDialog は選択行の名称編集ダイアログを表示する。
func (tv *TranslateTableView) openEditDialog() {
	if tv == nil || tv.TableView == nil || tv.model == nil || tv.editing {
		return
	}
	row := tv.CurrentIndex()
	if row < 0 || row >= tv.model.RowCount() {
		return
	}
	item := tv.model.Record(row)
	if item == nil {
		return
	}

	tv.editing = true
	defer func() {
		tv.editing = false
	}()

	var dlg *walk.Dialog
	var okButton *walk.PushButton
	var cancelButton *walk.PushButton
	var jpEdit *walk.TextEdit
	var enEdit *walk.TextEdit

	dialog := declarative.Dialog{
		AssignTo:      &dlg,
		CancelButton:  &cancelButton,
		DefaultButton: &okButton,
		Title:         tv.t(messages.LabelNameEditDialog),
		Layout:        declarative.VBox{},
		MinSize:       declarative.Size{Width: 420, Height: 220},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.Grid{Columns: 2},
				Children: []declarative.Widget{
					declarative.Label{Text: tv.t(messages.LabelTableType)},
					declarative.Label{Text: typeLabelByKey(tv.translator, item.TypeKey)},
					declarative.Label{Text: tv.t(messages.LabelTableSourceName)},
					declarative.Label{Text: item.NameText},
					declarative.Label{Text: tv.t(messages.LabelTableJapaneseName)},
					declarative.TextEdit{AssignTo: &jpEdit, Text: item.JapaneseNameText},
					declarative.Label{Text: tv.t(messages.LabelTableEnglishName)},
					declarative.TextEdit{AssignTo: &enEdit, Text: item.EnglishNameText},
				},
			},
			declarative.Composite{
				Layout: declarative.HBox{Alignment: declarative.AlignHFarVCenter},
				Children: []declarative.Widget{
					declarative.PushButton{
						AssignTo: &okButton,
						Text:     tv.t(messages.LabelOK),
						OnClicked: func() {
							jp := strings.TrimSpace(jpEdit.Text())
							if jp == "" {
								walk.MsgBox(dlg, tv.t(messages.MessageTextRequired), tv.t(messages.MessageTextRequired), walk.MsgBoxIconWarning)
								return
							}
							item.JapaneseNameText = jp
							item.EnglishNameText = strings.TrimSpace(enEdit.Text())
							dlg.Accept()
						},
					},
					declarative.PushButton{
						AssignTo: &cancelButton,
						Text:     tv.t(messages.LabelCancel),
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}

	owner := tv.TableView.Form()
	if owner == nil {
		owner = walk.App().ActiveForm()
	}
	if owner == nil {
		return
	}

	cmd, err := dialog.Run(owner)
	if err != nil {
		return
	}
	if cmd != walk.DlgCmdOK {
		return
	}

	item.Checked = true
	tv.model.PublishRowChanged(row)
}

// TranslateNameTableModel は名称置換テーブルのモデルを表す。
type TranslateNameTableModel struct {
	walk.TableModelBase
	walk.SorterBase
	translator i18n.II18n
	sortColumn int
	sortOrder  walk.SortOrder
	records    []*domain.TranslateNameItem
}

// NewTranslateNameTableModel は TranslateNameTableModel を生成する。
func NewTranslateNameTableModel(translator i18n.II18n) *TranslateNameTableModel {
	return &TranslateNameTableModel{
		translator: translator,
		sortColumn: 1,
		sortOrder:  walk.SortAscending,
		records:    []*domain.TranslateNameItem{},
	}
}

// RowCount は行数を返す。
func (m *TranslateNameTableModel) RowCount() int {
	return len(m.records)
}

// Value はセルの値を返す。
func (m *TranslateNameTableModel) Value(row int, col int) interface{} {
	item := m.records[row]

	switch col {
	case 0:
		return item.Checked
	case 1:
		return item.Number
	case 2:
		return typeLabelByKey(m.translator, item.TypeKey)
	case 3:
		return item.Index
	case 4:
		return item.NameText
	case 5:
		return item.JapaneseNameText
	case 6:
		return item.EnglishNameText
	default:
		return ""
	}
}

// Checked はチェック状態を返す。
func (m *TranslateNameTableModel) Checked(row int) bool {
	return m.records[row].Checked
}

// SetChecked はチェック状態を設定する。
func (m *TranslateNameTableModel) SetChecked(row int, checked bool) error {
	m.records[row].Checked = checked
	return nil
}

// ColumnSortable はソート可否を返す。
func (m *TranslateNameTableModel) ColumnSortable(col int) bool {
	return col >= 0
}

// SortedColumn は現在のソート列を返す。
func (m *TranslateNameTableModel) SortedColumn() int {
	return m.sortColumn
}

// SortOrder は現在のソート順を返す。
func (m *TranslateNameTableModel) SortOrder() walk.SortOrder {
	return m.sortOrder
}

// Sort は指定列で行を並び替える。
func (m *TranslateNameTableModel) Sort(col int, order walk.SortOrder) error {
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
			less = a.Index < b.Index
		case 4:
			less = a.NameText < b.NameText
		case 5:
			less = a.JapaneseNameText < b.JapaneseNameText
		case 6:
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
func (m *TranslateNameTableModel) ResetRows(items []domain.TranslateNameItem) {
	m.records = make([]*domain.TranslateNameItem, 0, len(items))
	for i := range items {
		item := items[i]
		m.records = append(m.records, &item)
	}
	m.PublishRowsReset()
}

// Rows は行一覧のコピーを返す。
func (m *TranslateNameTableModel) Rows() []domain.TranslateNameItem {
	rows := make([]domain.TranslateNameItem, 0, len(m.records))
	for _, item := range m.records {
		if item == nil {
			continue
		}
		rows = append(rows, *item)
	}
	return rows
}

// Record は指定行の参照を返す。
func (m *TranslateNameTableModel) Record(row int) *domain.TranslateNameItem {
	if row < 0 || row >= len(m.records) {
		return nil
	}
	return m.records[row]
}

// boolToInt は bool を 0/1 に変換する。
func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
