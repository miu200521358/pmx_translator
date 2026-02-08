// 指示: miu200521358
package minteractor

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/file/mstring"
	"github.com/miu200521358/pmx_translator/pkg/domain"
)

var (
	csvCandidateSeparators = []string{
		string(filepath.Separator), "_", "-", " ", "　", "/", ".",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
	ignoredSegmentNames = map[string]struct{}{
		"png": {}, "bmp": {}, "jpg": {}, "gif": {}, "tga": {}, "jpeg": {},
	}
)

//go:embed chara/*.txt
var charaFiles embed.FS

// BuildCsvCandidates はCSV出力対象候補を生成する。
func (uc *PmxTranslatorUsecase) BuildCsvCandidates(modelData *model.PmxModel) ([]domain.CsvCandidateItem, error) {
	if modelData == nil {
		return []domain.CsvCandidateItem{}, nil
	}

	kanjiSet, err := LoadKanji()
	if err != nil {
		return nil, err
	}

	items := make([]domain.CsvCandidateItem, 0)
	seen := map[string]struct{}{}

	addRecord := func(nameText string, englishName string, typeKey string, segmented bool) {
		if nameText == "" {
			return
		}
		if _, exists := seen[nameText]; exists {
			return
		}
		seen[nameText] = struct{}{}
		items = append(items, domain.CsvCandidateItem{
			Number:          len(items) + 1,
			Checked:         !IsJapaneseString(kanjiSet, nameText),
			TypeKey:         typeKey,
			Segmented:       segmented,
			NameText:        nameText,
			EnglishNameText: englishName,
		})
	}

	addCandidates := func(nameText string, englishName string, typeKey string, includeFullName bool) {
		if includeFullName {
			addRecord(nameText, englishName, typeKey, false)
		}
		for _, segment := range mstring.SplitAll(nameText, csvCandidateSeparators) {
			trimmed := strings.TrimSpace(segment)
			if trimmed == "" {
				continue
			}
			if _, exists := seen[trimmed]; exists {
				continue
			}
			if len([]rune(trimmed)) <= 1 && IsJapaneseString(kanjiSet, trimmed) {
				continue
			}
			if _, ignored := ignoredSegmentNames[strings.ToLower(trimmed)]; ignored {
				continue
			}
			addRecord(trimmed, "", typeKey, true)
		}
	}

	addCandidates(modelData.Path(), "", domain.NameTypePath, false)
	addCandidates(modelData.Name(), modelData.EnglishName, domain.NameTypeModel, true)

	if modelData.Materials != nil {
		for _, materialData := range modelData.Materials.Values() {
			if materialData == nil {
				continue
			}
			addCandidates(materialData.Name(), materialData.EnglishName, domain.NameTypeMaterial, true)
		}
	}
	if modelData.Textures != nil {
		for _, textureData := range modelData.Textures.Values() {
			if textureData == nil {
				continue
			}
			addCandidates(textureData.Name(), textureData.EnglishName, domain.NameTypeTexture, false)
		}
	}
	if modelData.Bones != nil {
		for _, boneData := range modelData.Bones.Values() {
			if boneData == nil {
				continue
			}
			addCandidates(boneData.Name(), boneData.EnglishName, domain.NameTypeBone, true)
		}
	}
	if modelData.Morphs != nil {
		for _, morphData := range modelData.Morphs.Values() {
			if morphData == nil {
				continue
			}
			addCandidates(morphData.Name(), morphData.EnglishName, domain.NameTypeMorph, true)
		}
	}
	if modelData.DisplaySlots != nil {
		for _, displaySlotData := range modelData.DisplaySlots.Values() {
			if displaySlotData == nil {
				continue
			}
			addCandidates(displaySlotData.Name(), displaySlotData.EnglishName, domain.NameTypeDisplaySlot, true)
		}
	}
	if modelData.RigidBodies != nil {
		for _, rigidBodyData := range modelData.RigidBodies.Values() {
			if rigidBodyData == nil {
				continue
			}
			addCandidates(rigidBodyData.Name(), rigidBodyData.EnglishName, domain.NameTypeRigidBody, true)
		}
	}
	if modelData.Joints != nil {
		for _, jointData := range modelData.Joints.Values() {
			if jointData == nil {
				continue
			}
			addCandidates(jointData.Name(), jointData.EnglishName, domain.NameTypeJoint, true)
		}
	}

	return items, nil
}

// LoadKanji は Shift-JIS 判定用文字集合を読み込む。
func LoadKanji() (string, error) {
	buf, err := fs.ReadFile(charaFiles, "chara/shiftjis.txt")
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// IsJapaneseString は文字列が Shift-JIS 互換文字のみで構成されるかを判定する。
func IsJapaneseString(kanjiSet string, text string) bool {
	for _, r := range text {
		if isAllowedCharacter(r) {
			continue
		}
		if strings.Contains(kanjiSet, string(r)) {
			continue
		}
		return false
	}
	return true
}

// isAllowedCharacter は Shift-JIS 互換として扱うUnicode範囲を判定する。
func isAllowedCharacter(r rune) bool {
	switch {
	case r >= 0x0000 && r <= 0x007F:
		return true
	case r >= 0x00A2 && r <= 0x00F7:
		return true
	case r >= 0x3000 && r <= 0x309F:
		return true
	case r >= 0xFF61 && r <= 0xFF9F:
		return true
	case r >= 0x30A0 && r <= 0x30FF:
		return true
	case r >= 0xFF01 && r <= 0xFF5D:
		return true
	default:
		return false
	}
}
