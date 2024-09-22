package ui

import (
	"path/filepath"
	"sort"
	"strings"

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
		return item.Index
	case 4:
		return item.NameText
	case 5:
		return item.EnglishNameText
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
			return c(a.Index < b.Index)
		case 4:
			return c(a.NameText < b.NameText)
		case 5:
			return c(a.EnglishNameText < b.EnglishNameText)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
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
	path, fileName, _ := mutils.SplitPath(model.Path())
	item := &domain.NameItem{
		Checked:         !usecase.IsJapaneseString(ks, fileName),
		Number:          len(m.Records) + 1,
		TypeText:        mi18n.T("ファイル"),
		Index:           0,
		NameText:        fileName,
		EnglishNameText: "",
	}
	m.Records = append(m.Records, item)

	for i, p := range strings.Split(path, string(filepath.Separator)) {
		if p == "" {
			continue
		}
		item := &domain.NameItem{
			Checked:         !usecase.IsJapaneseString(ks, p),
			Number:          len(m.Records) + 1,
			TypeText:        mi18n.T("ディレクトリ"),
			Index:           i,
			NameText:        p,
			EnglishNameText: "",
		}
		m.Records = append(m.Records, item)
	}

	for _, mat := range model.Materials.Data {
		item := &domain.NameItem{
			Checked:         !usecase.IsJapaneseString(ks, mat.Name()),
			Number:          len(m.Records) + 1,
			TypeText:        mi18n.T("材質"),
			Index:           mat.Index(),
			NameText:        mat.Name(),
			EnglishNameText: mat.EnglishName(),
		}
		m.Records = append(m.Records, item)
	}

	n := 0
	for _, tex := range model.Textures.Data {
		dirPath, fileName, _ := mutils.SplitPath(tex.Name())
		for _, p := range strings.Split(dirPath, string(filepath.Separator)) {
			if p == "" {
				continue
			}
			for _, p2 := range strings.Split(p, "/") {
				if p2 == "" {
					continue
				}
				item := &domain.NameItem{
					Checked:         !usecase.IsJapaneseString(ks, p2),
					Number:          len(m.Records) + 1,
					TypeText:        mi18n.T("ディレクトリ"),
					Index:           n,
					NameText:        p2,
					EnglishNameText: "",
				}
				m.Records = append(m.Records, item)
				n++
			}
		}

		{
			item := &domain.NameItem{
				Checked:         !usecase.IsJapaneseString(ks, fileName),
				Number:          len(m.Records) + 1,
				TypeText:        mi18n.T("ファイル"),
				Index:           n,
				NameText:        fileName,
				EnglishNameText: "",
			}
			m.Records = append(m.Records, item)
		}
	}

	for _, bone := range model.Bones.Data {
		item := &domain.NameItem{
			Checked:         !usecase.IsJapaneseString(ks, bone.Name()),
			Number:          len(m.Records) + 1,
			TypeText:        mi18n.T("ボーン"),
			Index:           bone.Index(),
			NameText:        bone.Name(),
			EnglishNameText: bone.EnglishName(),
		}
		m.Records = append(m.Records, item)
	}

	for _, morph := range model.Morphs.Data {
		item := &domain.NameItem{
			Checked:         !usecase.IsJapaneseString(ks, morph.Name()),
			Number:          len(m.Records) + 1,
			TypeText:        mi18n.T("モーフ"),
			Index:           morph.Index(),
			NameText:        morph.Name(),
			EnglishNameText: morph.EnglishName(),
		}
		m.Records = append(m.Records, item)
	}

	for _, disp := range model.DisplaySlots.Data {
		item := &domain.NameItem{
			Checked:         !usecase.IsJapaneseString(ks, disp.Name()),
			Number:          len(m.Records) + 1,
			TypeText:        mi18n.T("表示枠"),
			Index:           disp.Index(),
			NameText:        disp.Name(),
			EnglishNameText: disp.EnglishName(),
		}
		m.Records = append(m.Records, item)
	}

	for _, rb := range model.RigidBodies.Data {
		item := &domain.NameItem{
			Checked:         !usecase.IsJapaneseString(ks, rb.Name()),
			Number:          len(m.Records) + 1,
			TypeText:        mi18n.T("剛体"),
			Index:           rb.Index(),
			NameText:        rb.Name(),
			EnglishNameText: rb.EnglishName(),
		}
		m.Records = append(m.Records, item)
	}

	for _, joint := range model.Joints.Data {
		item := &domain.NameItem{
			Checked:         !usecase.IsJapaneseString(ks, joint.Name()),
			Number:          len(m.Records) + 1,
			TypeText:        mi18n.T("ジョイント"),
			Index:           joint.Index(),
			NameText:        joint.Name(),
			EnglishNameText: joint.EnglishName(),
		}
		m.Records = append(m.Records, item)
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
		Model:     nameModel,
	}

	return nameTableView
}

func (tv *CsvTableView) ResetModel(model *pmx.PmxModel) {
	tv.Model.ResetRows(model)
}
