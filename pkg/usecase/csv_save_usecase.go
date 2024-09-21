package usecase

import (
	"unicode"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// 簡体中国語の範囲を定義
var simplifiedChineseRanges = []*unicode.RangeTable{
	{
		R16: []unicode.Range16{
			{Lo: 0x4E00, Hi: 0x9FFF, Stride: 1}, // CJK Unified Ideographs
			{Lo: 0x3400, Hi: 0x4DBF, Stride: 1}, // CJK Unified Ideographs Extension A
			{Lo: 0xF900, Hi: 0xFAFF, Stride: 1}, // CJK Compatibility Ideographs
		},
		R32: []unicode.Range32{
			{Lo: 0x20000, Hi: 0x2A6DF, Stride: 1}, // CJK Unified Ideographs Extension B
			{Lo: 0x2A700, Hi: 0x2B73F, Stride: 1}, // CJK Unified Ideographs Extension C
			{Lo: 0x2B740, Hi: 0x2B81F, Stride: 1}, // CJK Unified Ideographs Extension D
			{Lo: 0x2B820, Hi: 0x2CEAF, Stride: 1}, // CJK Unified Ideographs Extension E
		},
	},
}

// 文字列に簡体中国語の文字が含まれているかどうかをチェックする
func containsSimplifiedChinese(s string) bool {
	for _, r := range s {
		if unicode.IsOneOf(simplifiedChineseRanges, r) {
			return true
		}
	}
	return false
}

// Shift-JISで解釈できる文字列かどうかを判定する
func isValidShiftJis(str string) bool {
	// 簡体中国語の文字が含まれているかチェック
	if containsSimplifiedChinese(str) {
		return false
	}

	// Create a transformer for Shift-JIS encoding
	encoder := japanese.ShiftJIS.NewEncoder()

	// Transform the input string to Shift-JIS
	if _, _, err := transform.String(encoder, str); err != nil {
		return false
	}
	return true
}

func CsvSave(model *pmx.PmxModel, outputPath string) error {
	records := make([][]string, 0)

	modelName := model.Name()
	if !isValidShiftJis(modelName) {
		records = append(records, []string{modelName, modelName})
	}

	for _, mat := range model.Materials.Data {
		if !isValidShiftJis(mat.Name()) {
			records = append(records, []string{modelName, mat.Name()})
		}
	}

	for _, bone := range model.Bones.Data {
		if !isValidShiftJis(bone.Name()) {
			records = append(records, []string{modelName, bone.Name()})
		}
	}

	for _, morph := range model.Morphs.Data {
		if !isValidShiftJis(morph.Name()) {
			records = append(records, []string{modelName, morph.Name()})
		}
	}

	for _, disp := range model.DisplaySlots.Data {
		if !isValidShiftJis(disp.Name()) {
			records = append(records, []string{modelName, disp.Name()})
		}
	}

	for _, rb := range model.RigidBodies.Data {
		if !isValidShiftJis(rb.Name()) {
			records = append(records, []string{modelName, rb.Name()})
		}
	}

	for _, joint := range model.Joints.Data {
		if !isValidShiftJis(joint.Name()) {
			records = append(records, []string{modelName, joint.Name()})
		}
	}

	data := core.NewCsvModel(records)

	csvRep := repository.NewCsvRepository()
	if err := csvRep.Save(outputPath, data, false); err != nil {
		return err
	}

	mlog.IT(mi18n.T("出力成功"), mi18n.T("出力成功メッセージ", map[string]interface{}{"Path": outputPath}))

	return nil
}
