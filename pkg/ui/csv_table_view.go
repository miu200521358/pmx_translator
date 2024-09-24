package ui

import (
	"path/filepath"
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/pmx_translator/pkg/domain"
	"github.com/miu200521358/pmx_translator/pkg/usecase"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type CsvTableView struct {
	*declarative.TableView
	Model *CsvNameModel
}

type CsvNameModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	Records    []*domain.NameItem
}

func NewCsvNameModel(model *pmx.PmxModel) *CsvNameModel {
	m := new(CsvNameModel)
	m.ResetRows(model)
	return m
}

func (m *CsvNameModel) CheckedNames() []string {
	var names []string
	for _, item := range m.Records {
		if item.Checked {
			names = append(names, item.NameText)
		}
	}
	return names
}

func (m *CsvNameModel) RowCount() int {
	return len(m.Records)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *CsvNameModel) Value(row, col int) interface{} {
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
	}

	panic("unexpected col")
}

func (m *CsvNameModel) Checked(row int) bool {
	return m.Records[row].Checked
}

func (m *CsvNameModel) SetChecked(row int, Checked bool) error {
	m.Records[row].Checked = Checked

	return nil
}

func (m *CsvNameModel) Sort(col int, order walk.SortOrder) error {
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
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *CsvNameModel) exists(txt string) bool {
	for _, item := range m.Records {
		if item.NameText == txt {
			return true
		}
	}
	return false
}

var separators = []string{string(filepath.Separator), "_", "-", " ", "　", "/", ".", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func (m *CsvNameModel) AddRecord(ks, txt, fieldKey string) {
	for _, t := range mutils.SplitAll(txt, separators) {
		if t == "" || m.exists(t) {
			continue
		}
		item := &domain.NameItem{
			Checked:  !usecase.IsJapaneseString(ks, t),
			Number:   len(m.Records) + 1,
			TypeText: mi18n.T(fieldKey),
			NameText: t,
		}
		m.Records = append(m.Records, item)
	}
}

func (m *CsvNameModel) ResetRows(model *pmx.PmxModel) {
	m.Records = make([]*domain.NameItem, 0)

	m.PublishRowsReset()

	if model == nil {
		return
	}

	ks, err := usecase.LoadKanji()
	if err != nil {
		return
	}

	// ファイルパスの中国語もピックアップ
	m.AddRecord(ks, model.Path(), "パス")
	m.AddRecord(ks, model.Name(), "モデル")

	for _, mat := range model.Materials.Data {
		m.AddRecord(ks, mat.Name(), "材質")
		m.AddRecord(ks, mat.EnglishName(), "材質")
	}

	for _, tex := range model.Textures.Data {
		m.AddRecord(ks, tex.Name(), "テクスチャ")
		m.AddRecord(ks, tex.EnglishName(), "テクスチャ")
	}

	for _, bone := range model.Bones.Data {
		m.AddRecord(ks, bone.Name(), "ボーン")
		m.AddRecord(ks, bone.EnglishName(), "ボーン")
	}

	for _, morph := range model.Morphs.Data {
		m.AddRecord(ks, morph.Name(), "モーフ")
		m.AddRecord(ks, morph.EnglishName(), "モーフ")
	}

	for _, disp := range model.DisplaySlots.Data {
		m.AddRecord(ks, disp.Name(), "表示枠")
		m.AddRecord(ks, disp.EnglishName(), "表示枠")
	}

	for _, rb := range model.RigidBodies.Data {
		m.AddRecord(ks, rb.Name(), "剛体")
		m.AddRecord(ks, rb.EnglishName(), "剛体")
	}

	for _, joint := range model.Joints.Data {
		m.AddRecord(ks, joint.Name(), "ジョイント")
		m.AddRecord(ks, joint.EnglishName(), "ジョイント")
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
			{Title: mi18n.T("日本語名称(分割)"), Width: 200},
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
		Model:     nameModel,
	}

	return nameTableView
}

func (tv *CsvTableView) ResetModel(model *pmx.PmxModel) {
	tv.Model.ResetRows(model)
}
