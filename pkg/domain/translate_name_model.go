package domain

import (
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/domain/mcsv"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
	"github.com/miu200521358/walk/pkg/walk"
)

type TranslateNameModel struct {
	walk.TableModelBase
	walk.SorterBase
	tv         *walk.TableView
	sortColumn int
	sortOrder  walk.SortOrder
	Records    []*NameItem
}

func (m *TranslateNameModel) RowCount() int {
	return len(m.Records)
}

func (m *TranslateNameModel) SetParent(parent *walk.TableView) {
	m.tv = parent
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

func (m *TranslateNameModel) AddRecord(
	fileName string, index int, jpTxt, enTxt, fieldKey string, charaCsv *mcsv.CsvModel,
) {
	jpTransTxt, enTransTxt := Translate(jpTxt, enTxt, charaCsv, fileName)

	item := &NameItem{
		Checked:          jpTxt != jpTransTxt || enTxt != enTransTxt,
		Number:           len(m.Records) + 1,
		TypeText:         mi18n.T(fieldKey),
		Index:            index,
		NameText:         jpTxt,
		JapaneseNameText: jpTransTxt,
		EnglishNameText:  enTransTxt,
	}
	m.Records = append(m.Records, item)

}

func (m *TranslateNameModel) ResetRows(model *pmx.PmxModel, charaCsv *mcsv.CsvModel) {
	m.Records = make([]*NameItem, 0)

	m.PublishRowsReset()

	if model == nil || charaCsv == nil {
		return
	}

	_, fileName, _ := mfile.SplitPath(model.Path())

	m.AddRecord(fileName, 0, model.Path(), "", "パス", charaCsv)
	m.AddRecord(fileName, 0, model.Name(), model.EnglishName(), "モデル", charaCsv)

	model.Materials.ForEach(func(i int, mat *pmx.Material) {
		m.AddRecord(fileName, mat.Index(), mat.Name(), mat.EnglishName(), "材質", charaCsv)
	})

	model.Textures.ForEach(func(i int, tex *pmx.Texture) {
		m.AddRecord(fileName, tex.Index(), tex.Name(), "", "テクスチャ", charaCsv)
	})

	model.Bones.ForEach(func(i int, bone *pmx.Bone) {
		m.AddRecord(fileName, bone.Index(), bone.Name(), bone.EnglishName(), "ボーン", charaCsv)
	})

	model.Morphs.ForEach(func(i int, morph *pmx.Morph) {
		m.AddRecord(fileName, morph.Index(), morph.Name(), morph.EnglishName(), "モーフ", charaCsv)
	})

	model.DisplaySlots.ForEach(func(i int, disp *pmx.DisplaySlot) {
		m.AddRecord(fileName, disp.Index(), disp.Name(), disp.EnglishName(), "表示枠", charaCsv)
	})

	model.RigidBodies.ForEach(func(i int, rb *pmx.RigidBody) {
		m.AddRecord(fileName, rb.Index(), rb.Name(), rb.EnglishName(), "剛体", charaCsv)
	})

	model.Joints.ForEach(func(i int, joint *pmx.Joint) {
		m.AddRecord(fileName, joint.Index(), joint.Name(), joint.EnglishName(), "ジョイント", charaCsv)
	})

	m.PublishRowsReset()
}

func Translate(text, enText string, langDict *mcsv.CsvModel, modelName string) (string, string) {
	newJpText := text
	newEnText := enText

	// モデル名一致＆完全一致の翻訳を行う
	for n, row := range langDict.Records() {
		if n > 0 && row[0] == modelName && row[1] == newJpText {
			newJpText = strings.ReplaceAll(newJpText, row[1], row[2])
			if enText != "" && row[3] != "" {
				newEnText = row[3]
			}
		}
	}

	// モデル名不問＆完全一致の翻訳を行う
	for n, row := range langDict.Records() {
		if n > 0 && row[0] == "" && row[1] == newJpText {
			newJpText = strings.ReplaceAll(newJpText, row[1], row[2])
			if enText != "" && row[3] != "" {
				newEnText = row[3]
			}
		}
	}

	// モデル名一致＆部分一致翻訳を行う
	for n, row := range langDict.Records() {
		if n > 0 && row[0] == modelName && row[1] != "" {
			newJpText = strings.ReplaceAll(newJpText, row[1], row[2])
			if enText != "" && row[3] != "" {
				newEnText = strings.ReplaceAll(newEnText, row[1], row[3])
			}
		}
	}

	// モデル名不問＆部分一致翻訳
	for n, row := range langDict.Records() {
		if n > 0 && row[0] == "" && row[1] != "" {
			newJpText = strings.ReplaceAll(newJpText, row[1], row[2])
			if enText != "" && row[3] != "" {
				newEnText = strings.ReplaceAll(newEnText, row[1], row[3])
			}
		}
	}

	return newJpText, newEnText
}
