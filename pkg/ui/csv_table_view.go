package ui

import (
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type CsvTableView struct {
	*declarative.TableView
	model *CsvNameModel
}

type CsvNameModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*NameItem
}

type NameItem struct {
	number          int
	checked         bool
	typeText        string
	index           int
	nameText        string
	englishNameText string
}

func NewCsvNameModel(model *pmx.PmxModel) *CsvNameModel {
	m := new(CsvNameModel)
	m.ResetRows(model)
	return m
}

func (m *CsvNameModel) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *CsvNameModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.checked
	case 1:
		return item.number
	case 2:
		return item.typeText
	case 3:
		return item.index
	case 4:
		return item.nameText
	case 5:
		return item.englishNameText
	}

	panic("unexpected col")
}

func (m *CsvNameModel) Checked(row int) bool {
	return m.items[row].checked
}

func (m *CsvNameModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

func (m *CsvNameModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			av := 0
			if a.checked {
				av = 1
			}
			bv := 0
			if b.checked {
				bv = 1
			}
			return c(av < bv)
		case 1:
			return c(a.number < b.number)
		case 2:
			return c(a.typeText < b.typeText)
		case 3:
			return c(a.index < b.index)
		case 4:
			return c(a.nameText < b.nameText)
		case 5:
			return c(a.englishNameText < b.englishNameText)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *CsvNameModel) ResetRows(model *pmx.PmxModel) {
	m.items = make([]*NameItem, 0)

	if model == nil {
		m.PublishRowsReset()
		return
	}

	for _, mat := range model.Materials.Data {
		item := &NameItem{
			checked:         false,
			number:          len(m.items) + 1,
			typeText:        mi18n.T("材質"),
			index:           mat.Index(),
			nameText:        mat.Name(),
			englishNameText: mat.EnglishName(),
		}
		m.items = append(m.items, item)
	}

	for _, tex := range model.Textures.Data {
		item := &NameItem{
			checked:         false,
			number:          len(m.items) + 1,
			typeText:        mi18n.T("テクスチャ"),
			index:           tex.Index(),
			nameText:        tex.Name(),
			englishNameText: tex.EnglishName(),
		}
		m.items = append(m.items, item)
	}

	for _, bone := range model.Bones.Data {
		item := &NameItem{
			checked:         false,
			number:          len(m.items) + 1,
			typeText:        mi18n.T("ボーン"),
			index:           bone.Index(),
			nameText:        bone.Name(),
			englishNameText: bone.EnglishName(),
		}
		m.items = append(m.items, item)
	}

	for _, morph := range model.Morphs.Data {
		item := &NameItem{
			checked:         false,
			number:          len(m.items) + 1,
			typeText:        mi18n.T("モーフ"),
			index:           morph.Index(),
			nameText:        morph.Name(),
			englishNameText: morph.EnglishName(),
		}
		m.items = append(m.items, item)
	}

	for _, disp := range model.DisplaySlots.Data {
		item := &NameItem{
			checked:         false,
			number:          len(m.items) + 1,
			typeText:        mi18n.T("表示枠"),
			index:           disp.Index(),
			nameText:        disp.Name(),
			englishNameText: disp.EnglishName(),
		}
		m.items = append(m.items, item)
	}

	for _, rb := range model.RigidBodies.Data {
		item := &NameItem{
			checked:         false,
			number:          len(m.items) + 1,
			typeText:        mi18n.T("剛体"),
			index:           rb.Index(),
			nameText:        rb.Name(),
			englishNameText: rb.EnglishName(),
		}
		m.items = append(m.items, item)
	}

	for _, joint := range model.Joints.Data {
		item := &NameItem{
			checked:         false,
			number:          len(m.items) + 1,
			typeText:        mi18n.T("ジョイント"),
			index:           joint.Index(),
			nameText:        joint.Name(),
			englishNameText: joint.EnglishName(),
		}
		m.items = append(m.items, item)
	}

	m.PublishRowsReset()
}

func NewCsvTableView(parent walk.Container, model *pmx.PmxModel) *CsvTableView {
	nameModel := NewCsvNameModel(model)

	var tv *walk.TableView
	builder := declarative.NewBuilder(parent)

	dTableView := &declarative.TableView{
		AssignTo:         &tv,
		AlternatingRowBG: true,
		CheckBoxes:       true,
		ColumnsOrderable: true,
		MultiSelection:   true,
		Model:            nameModel,
		Columns: []declarative.TableViewColumn{
			{Title: "#", Width: 30},
			{Title: "No.", Width: 50},
			{Title: mi18n.T("種類"), Width: 80},
			{Title: mi18n.T("インデックス"), Width: 40},
			{Title: mi18n.T("日本語名称"), Width: 200},
			{Title: mi18n.T("英語名称"), Width: 200},
		},
		StyleCell: func(style *walk.CellStyle) {
			if nameModel.Checked(style.Row()) {
				style.BackgroundColor = walk.RGB(159, 255, 243)
			} else {
				style.BackgroundColor = walk.RGB(255, 255, 255)
			}
		},
	}

	if err := dTableView.Create(builder); err != nil {
		panic(err)
	}

	nameTableView := &CsvTableView{
		TableView: dTableView,
		model:     nameModel,
	}

	return nameTableView
}

func (tv *CsvTableView) ResetModel(model *pmx.PmxModel) {
	tv.model.ResetRows(model)
}
