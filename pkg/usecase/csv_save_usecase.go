package usecase

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

//go:embed chara/*.txt
var charaFiles embed.FS

func isJapaneseString(ks string, s string) bool {
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

func loadKanji() (string, error) {
	buf, err := fs.ReadFile(charaFiles, "chara/shiftjis.txt")
	if err != nil {
		mlog.E("LoadShiftjis error: %v", err)
		return "", err
	}

	return string(buf), nil
}

func existText(records [][]string, txt string) bool {
	for _, row := range records {
		if row[1] == txt {
			return true
		}
	}
	return false
}

func CsvSave(model *pmx.PmxModel, outputPath string) error {
	records := make([][]string, 0)
	ks, err := loadKanji()
	if err != nil {
		return err
	}

	// ファイルパスの中国語もピックアップ
	path, fileName, _ := mutils.SplitPath(model.Path())
	if !isJapaneseString(ks, fileName) && !existText(records, fileName) {
		records = append(records, []string{fileName, fileName})
	}

	for _, p := range strings.Split(path, string(filepath.Separator)) {
		if !isJapaneseString(ks, p) && !existText(records, p) {
			records = append(records, []string{fileName, p})
		}
	}

	modelName := model.Name()
	if !isJapaneseString(ks, modelName) && !existText(records, modelName) {
		records = append(records, []string{fileName, modelName})
	}

	for _, mat := range model.Materials.Data {
		if !isJapaneseString(ks, mat.Name()) && !existText(records, mat.Name()) {
			records = append(records, []string{fileName, mat.Name()})
		}
	}

	for _, bone := range model.Bones.Data {
		if !isJapaneseString(ks, bone.Name()) && !existText(records, bone.Name()) {
			records = append(records, []string{fileName, bone.Name()})
		}
	}

	for _, morph := range model.Morphs.Data {
		if !isJapaneseString(ks, morph.Name()) && !existText(records, morph.Name()) {
			records = append(records, []string{fileName, morph.Name()})
		}
	}

	for _, disp := range model.DisplaySlots.Data {
		if !isJapaneseString(ks, disp.Name()) && !existText(records, disp.Name()) {
			records = append(records, []string{fileName, disp.Name()})
		}
	}

	for _, rb := range model.RigidBodies.Data {
		if !isJapaneseString(ks, rb.Name()) && !existText(records, rb.Name()) {
			records = append(records, []string{fileName, rb.Name()})
		}
	}

	for _, joint := range model.Joints.Data {
		if !isJapaneseString(ks, joint.Name()) && !existText(records, joint.Name()) {
			records = append(records, []string{fileName, joint.Name()})
		}
	}

	data := core.NewCsvModel(records)

	csvRep := repository.NewCsvRepository()
	if err := csvRep.Save(outputPath, data, false); err != nil {
		return err
	}

	mlog.IT(mi18n.T("出力成功"), mi18n.T("Csv出力成功メッセージ", map[string]interface{}{"Path": outputPath}))

	return nil
}
