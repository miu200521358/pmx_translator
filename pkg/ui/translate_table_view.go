package ui

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/pmx_translator/pkg/domain"
	"github.com/miu200521358/pmx_translator/pkg/usecase"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type TranslateTableView struct {
	*declarative.TableView
	NameModel *TranslateNameModel
}

type TranslateNameModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	Records    []*domain.NameItem
}

func NewTranslateNameModel(model *pmx.PmxModel, charaCsv *core.CsvModel) *TranslateNameModel {
	m := new(TranslateNameModel)
	m.ResetRows(model, charaCsv)
	return m
}

func (m *TranslateNameModel) RowCount() int {
	return len(m.Records)
}

func (m *TranslateNameModel) Value(row, col int) interface{} {
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
		return item.JapaneseNameText
	case 6:
		return item.EnglishNameText
	}

	panic("unexpected col")
}

func (m *TranslateNameModel) Checked(row int) bool {
	return m.Records[row].Checked
}

func (m *TranslateNameModel) SetChecked(row int, checked bool) error {
	m.Records[row].Checked = checked

	return nil
}

func (m *TranslateNameModel) Sort(col int, order walk.SortOrder) error {
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
			return c(a.JapaneseNameText < b.JapaneseNameText)
		case 6:
			return c(a.EnglishNameText < b.EnglishNameText)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *TranslateNameModel) ResetRows(model *pmx.PmxModel, charaCsv *core.CsvModel) {
	m.Records = make([]*domain.NameItem, 0)

	m.PublishRowsReset()

	if model == nil || charaCsv == nil {
		return
	}

	// ファイルパスの中国語もピックアップ
	path, fileName, _ := mutils.SplitPath(model.Path())
	{
		jpName := usecase.Translate(fileName, charaCsv, 2, fileName)
		item := &domain.NameItem{
			Checked:          fileName != jpName,
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T("ファイル"),
			Index:            0,
			NameText:         fileName,
			JapaneseNameText: jpName,
			EnglishNameText:  "",
		}
		m.Records = append(m.Records, item)
	}

	for i, p := range strings.Split(path, string(filepath.Separator)) {
		if p == "" {
			continue
		}

		jpName := usecase.Translate(p, charaCsv, 2, fileName)
		enName := usecase.Translate(p, charaCsv, 3, fileName)
		item := &domain.NameItem{
			Checked:          p != jpName,
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T("ディレクトリ"),
			Index:            i,
			NameText:         p,
			JapaneseNameText: jpName,
			EnglishNameText:  enName,
		}
		m.Records = append(m.Records, item)
	}

	{
		jpName := usecase.Translate(model.Name(), charaCsv, 2, fileName)
		enName := usecase.Translate(model.Name(), charaCsv, 3, fileName)
		item := &domain.NameItem{
			Checked:          fileName != jpName,
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T("モデル"),
			Index:            0,
			NameText:         fileName,
			JapaneseNameText: jpName,
			EnglishNameText:  enName,
		}
		m.Records = append(m.Records, item)
	}

	for _, mat := range model.Materials.Data {
		jpName := usecase.Translate(mat.Name(), charaCsv, 2, fileName)
		enName := usecase.Translate(mat.Name(), charaCsv, 3, fileName)
		item := &domain.NameItem{
			Checked:          mat.Name() != jpName,
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T("材質"),
			Index:            mat.Index(),
			NameText:         mat.Name(),
			JapaneseNameText: jpName,
			EnglishNameText:  enName,
		}
		m.Records = append(m.Records, item)
	}

	for _, tex := range model.Textures.Data {
		jpName := usecase.Translate(tex.Name(), charaCsv, 2, fileName)
		enName := usecase.Translate(tex.Name(), charaCsv, 3, fileName)
		item := &domain.NameItem{
			Checked:          tex.Name() != jpName,
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T("テクスチャ"),
			Index:            tex.Index(),
			NameText:         tex.Name(),
			JapaneseNameText: jpName,
			EnglishNameText:  enName,
		}
		m.Records = append(m.Records, item)
	}

	for _, bone := range model.Bones.Data {
		jpName := usecase.Translate(bone.Name(), charaCsv, 2, fileName)
		enName := usecase.Translate(bone.Name(), charaCsv, 3, fileName)
		item := &domain.NameItem{
			Checked:          bone.Name() != jpName,
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T("ボーン"),
			Index:            bone.Index(),
			NameText:         bone.Name(),
			JapaneseNameText: jpName,
			EnglishNameText:  enName,
		}
		m.Records = append(m.Records, item)
	}

	for _, morph := range model.Morphs.Data {
		jpName := usecase.Translate(morph.Name(), charaCsv, 2, fileName)
		enName := usecase.Translate(morph.Name(), charaCsv, 3, fileName)
		item := &domain.NameItem{
			Checked:          morph.Name() != jpName,
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T("モーフ"),
			Index:            morph.Index(),
			NameText:         morph.Name(),
			JapaneseNameText: jpName,
			EnglishNameText:  enName,
		}
		m.Records = append(m.Records, item)
	}

	for _, disp := range model.DisplaySlots.Data {
		jpName := usecase.Translate(disp.Name(), charaCsv, 2, fileName)
		enName := usecase.Translate(disp.Name(), charaCsv, 3, fileName)
		item := &domain.NameItem{
			Checked:          disp.Name() != jpName,
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T("表示枠"),
			Index:            disp.Index(),
			NameText:         disp.Name(),
			JapaneseNameText: jpName,
			EnglishNameText:  enName,
		}
		m.Records = append(m.Records, item)
	}

	for _, rb := range model.RigidBodies.Data {
		jpName := usecase.Translate(rb.Name(), charaCsv, 2, fileName)
		enName := usecase.Translate(rb.Name(), charaCsv, 3, fileName)
		item := &domain.NameItem{
			Checked:          rb.Name() != jpName,
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T("剛体"),
			Index:            rb.Index(),
			NameText:         rb.Name(),
			JapaneseNameText: jpName,
			EnglishNameText:  enName,
		}
		m.Records = append(m.Records, item)
	}

	for _, joint := range model.Joints.Data {
		jpName := usecase.Translate(joint.Name(), charaCsv, 2, fileName)
		enName := usecase.Translate(joint.Name(), charaCsv, 3, fileName)
		item := &domain.NameItem{
			Checked:          joint.Name() != jpName,
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T("ジョイント"),
			Index:            joint.Index(),
			NameText:         joint.Name(),
			JapaneseNameText: jpName,
			EnglishNameText:  enName,
		}
		m.Records = append(m.Records, item)
	}

	m.PublishRowsReset()
}

func NewTranslateTableView(parent walk.Container, model *pmx.PmxModel, charaCsv *core.CsvModel) *TranslateTableView {
	nameModel := NewTranslateNameModel(model, charaCsv)

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
			{Title: mi18n.T("元名称"), Width: 200},
			{Title: mi18n.T("翻訳名称"), Width: 200},
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

	nameTableView := &TranslateTableView{
		TableView: dTableView,
		NameModel: nameModel,
	}

	return nameTableView
}

func (tv *TranslateTableView) ResetModel(model *pmx.PmxModel, charaCsv *core.CsvModel) {
	tv.NameModel.ResetRows(model, charaCsv)
}
