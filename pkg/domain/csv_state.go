package domain

import (
	"embed"
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mstring"
	"github.com/miu200521358/walk/pkg/walk"
)

type CsvState struct {
	Model            *pmx.PmxModel // 処理対象モデル
	NameModel        *CsvNameModel // 名称モデル
	TextChangeDialog *walk.Dialog  // テキスト変更ダイアログ
	OutputPath       string        // 出力パス
}

func NewCsvState() *CsvState {
	return &CsvState{
		NameModel: new(CsvNameModel),
	}
}

// 拡張子変更(大文字小文字無視)
var pmxToCsvRegex = regexp.MustCompile(`(?i)\.pmx$`)

func (c *CsvState) LoadData() {
	if c.Model == nil {
		return
	}

	c.NameModel.ResetRows(c.Model)
	outputPath := mfile.CreateOutputPath(c.Model.Path(), "")
	outputPath = pmxToCsvRegex.ReplaceAllString(outputPath, ".csv")
	c.OutputPath = outputPath
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

// --------------------------------------------------

type CsvNameModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	Records    []*NameItem
}

func (m *CsvNameModel) RowCount() int {
	return len(m.Records)
}

func (m *CsvNameModel) Value(row, col int) any {
	item := m.Records[row]

	switch col {
	case 0:
		return item.Checked
	case 1:
		return item.Number
	case 2:
		return item.TypeText
	case 3:
		return item.Segmented
	case 4:
		return item.JapaneseNameText
	case 5:
		return item.EnglishNameText
	}

	panic("unexpected col")
}

func (m *CsvNameModel) Checked(row int) bool {
	return m.Records[row].Checked
}

func (m *CsvNameModel) SetChecked(row int, checked bool) error {
	m.Records[row].Checked = checked

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
			av := 0
			if a.Segmented {
				av = 1
			}
			bv := 0
			if b.Segmented {
				bv = 1
			}
			return c(av < bv)
		case 4:
			return c(a.JapaneseNameText < b.JapaneseNameText)
		case 5:
			return c(a.EnglishNameText < b.EnglishNameText)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

// --------------------------------------------------

func (m *CsvNameModel) exists(txt string) bool {
	for _, item := range m.Records {
		if item.NameText == txt {
			return true
		}
	}
	return false
}

func IsJapaneseString(ks string, s string) bool {
	for _, r := range s {
		if isAllowedCharacter(r) {
			continue
		} else {
			if strings.Contains(ks, string(r)) {
				continue
			}
			return false
		}
	}
	return true
}

func isAllowedCharacter(r rune) bool {
	switch {
	case r >= 0x0000 && r <= 0x007F: // ASCII
		return true
	case r >= 0x00A2 && r <= 0x00F7: // ASCII
		return true
	case r >= 0x3000 && r <= 0x309F: // Hiragana
		return true
	case r >= 0xFF61 && r <= 0xFF9F: // Half-width Katakana
		return true
	case r >= 0x30A0 && r <= 0x30FF: // Full-width Katakana
		return true
	case r >= 0xFF01 && r <= 0xFF5D: // Full-width Alphanumeric
		return true
	default:
		return false
	}
}

//go:embed chara/*.txt
var charaFiles embed.FS

func LoadKanji() (string, error) {
	buf, err := fs.ReadFile(charaFiles, "chara/shiftjis.txt")
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func ExistText(records [][]string, txt string) bool {
	for _, row := range records {
		if row[1] == txt {
			return true
		}
	}
	return false
}

var separators = []string{string(filepath.Separator), "_", "-", " ", "　", "/", ".", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func (m *CsvNameModel) AddRecord(ks, jpTxt, enTxt, fieldKey string) {
	if !m.exists(jpTxt) && fieldKey != "パス" && fieldKey != "テクスチャ" {
		item := &NameItem{
			Checked:          !IsJapaneseString(ks, jpTxt),
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T(fieldKey),
			NameText:         jpTxt,
			JapaneseNameText: jpTxt,
			EnglishNameText:  enTxt,
			Segmented:        false,
		}
		m.Records = append(m.Records, item)
	}

	for _, t := range mstring.SplitAll(jpTxt, separators) {
		if t == "" || m.exists(t) || (len(t) <= 1 && IsJapaneseString(ks, t)) ||
			slices.Contains([]string{"png", "bmp", "jpg", "gif", "tga", "jpeg"}, strings.ToLower(t)) {
			continue
		}
		item := &NameItem{
			Checked:          !IsJapaneseString(ks, t),
			Number:           len(m.Records) + 1,
			TypeText:         mi18n.T(fieldKey),
			NameText:         t,
			JapaneseNameText: t,
			EnglishNameText:  "",
			Segmented:        true,
		}
		m.Records = append(m.Records, item)
	}
}

func (m *CsvNameModel) ResetRows(model *pmx.PmxModel) {
	m.Records = make([]*NameItem, 0)

	m.PublishRowsReset()

	if model == nil {
		return
	}

	ks, err := LoadKanji()
	if err != nil {
		return
	}

	// ファイルパスの中国語もピックアップ
	m.AddRecord(ks, model.Path(), "", "パス")
	m.AddRecord(ks, model.Name(), model.EnglishName(), "モデル")

	model.Materials.ForEach(func(i int, mat *pmx.Material) bool {
		m.AddRecord(ks, mat.Name(), mat.EnglishName(), "材質")
		return true
	})

	model.Textures.ForEach(func(i int, tex *pmx.Texture) bool {
		m.AddRecord(ks, tex.Name(), tex.EnglishName(), "テクスチャ")
		return true
	})

	model.Bones.ForEach(func(i int, bone *pmx.Bone) bool {
		m.AddRecord(ks, bone.Name(), bone.EnglishName(), "ボーン")
		return true
	})

	model.Morphs.ForEach(func(i int, morph *pmx.Morph) bool {
		m.AddRecord(ks, morph.Name(), morph.EnglishName(), "モーフ")
		return true
	})

	model.DisplaySlots.ForEach(func(i int, disp *pmx.DisplaySlot) bool {
		m.AddRecord(ks, disp.Name(), disp.EnglishName(), "表示枠")
		return true
	})

	model.RigidBodies.ForEach(func(i int, rb *pmx.RigidBody) bool {
		m.AddRecord(ks, rb.Name(), rb.EnglishName(), "剛体")
		return true
	})

	model.Joints.ForEach(func(i int, joint *pmx.Joint) bool {
		m.AddRecord(ks, joint.Name(), joint.EnglishName(), "ジョイント")
		return true
	})

	m.PublishRowsReset()
}
