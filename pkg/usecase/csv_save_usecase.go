package usecase

import (
	"embed"
	"io/fs"
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

func LoadKanji() (string, error) {
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

func CsvSave(model *pmx.PmxModel, checkedNames []string, outputPath string) error {
	records := make([][]string, 0)
	records = append(records, []string{
		mi18n.T("ファイル名"), mi18n.T("元名称"), mi18n.T("日本語名称"), mi18n.T("英語名称")})

	_, fileName, _ := mutils.SplitPath(model.Path())

	for _, name := range checkedNames {
		if !existText(records, name) {
			records = append(records, []string{fileName, name, "", ""})
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
