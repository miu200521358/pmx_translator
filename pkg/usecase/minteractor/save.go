// 指示: miu200521358
package minteractor

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_csv"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/file/mfile"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/pmx_translator/pkg/domain"
)

const (
	prerequisiteMissingErrorID = "13001"
)

// SaveTranslatedModel は名称置換結果をPMXとして保存する。
func (uc *PmxTranslatorUsecase) SaveTranslatedModel(outputPath string, sourceModel *model.PmxModel, nameItems []domain.TranslateNameItem, options SaveOptions) error {
	if sourceModel == nil {
		return newPrerequisiteMissingError(errorTranslateTargetModelRequired)
	}
	if strings.TrimSpace(outputPath) == "" {
		return newPrerequisiteMissingError(errorOutputModelPathRequired)
	}
	if uc.modelWriter == nil {
		return newPrerequisiteMissingError(errorModelWriterMissing)
	}

	copiedModel, err := sourceModel.Copy()
	if err != nil {
		return err
	}
	modelData := &copiedModel

	number := 1
	number++ // パス項目は保存時に直接反映しない。
	jpName, enName := translatedNamesByNumber(number, modelData.Name(), modelData.EnglishName, nameItems)
	modelData.SetName(jpName)
	modelData.EnglishName = enName
	number++

	if modelData.Materials != nil {
		for _, materialData := range modelData.Materials.Values() {
			if materialData == nil {
				continue
			}
			jpName, enName = translatedNamesByNumber(number, materialData.Name(), materialData.EnglishName, nameItems)
			materialData.SetName(jpName)
			materialData.EnglishName = enName
			number++
		}
	}

	sourceDir, _, _ := mfile.SplitPath(sourceModel.Path())
	outputDir, _, _ := mfile.SplitPath(outputPath)
	if strings.TrimSpace(outputDir) != "" {
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			return io_common.NewIoSaveFailed(errorOutputDirCreateFailed, err)
		}
	}
	if modelData.Textures != nil {
		for _, textureData := range modelData.Textures.Values() {
			if textureData == nil {
				continue
			}
			originalName := textureData.Name()
			jpName, _ := translatedNamesByNumber(number, originalName, "", nameItems)
			textureData.SetName(jpName)
			number++

			if strings.TrimSpace(originalName) == "" || strings.TrimSpace(sourceDir) == "" || strings.TrimSpace(outputDir) == "" {
				continue
			}
			sourceTexturePath := filepath.Join(sourceDir, originalName)
			outputTexturePath := filepath.Join(outputDir, jpName)
			if filepath.Clean(sourceTexturePath) == filepath.Clean(outputTexturePath) {
				continue
			}
			if err := copyTextureFile(sourceTexturePath, outputTexturePath); err != nil {
				return err
			}
		}
	}

	if modelData.Bones != nil {
		for _, boneData := range modelData.Bones.Values() {
			if boneData == nil {
				continue
			}
			jpName, enName = translatedNamesByNumber(number, boneData.Name(), boneData.EnglishName, nameItems)
			boneData.SetName(jpName)
			boneData.EnglishName = enName
			number++
		}
	}
	if modelData.Morphs != nil {
		for _, morphData := range modelData.Morphs.Values() {
			if morphData == nil {
				continue
			}
			jpName, enName = translatedNamesByNumber(number, morphData.Name(), morphData.EnglishName, nameItems)
			morphData.SetName(jpName)
			morphData.EnglishName = enName
			number++
		}
	}
	if modelData.DisplaySlots != nil {
		for _, displaySlotData := range modelData.DisplaySlots.Values() {
			if displaySlotData == nil {
				continue
			}
			jpName, enName = translatedNamesByNumber(number, displaySlotData.Name(), displaySlotData.EnglishName, nameItems)
			displaySlotData.SetName(jpName)
			displaySlotData.EnglishName = enName
			number++
		}
	}
	if modelData.RigidBodies != nil {
		for _, rigidBodyData := range modelData.RigidBodies.Values() {
			if rigidBodyData == nil {
				continue
			}
			jpName, enName = translatedNamesByNumber(number, rigidBodyData.Name(), rigidBodyData.EnglishName, nameItems)
			rigidBodyData.SetName(jpName)
			rigidBodyData.EnglishName = enName
			number++
		}
	}
	if modelData.Joints != nil {
		for _, jointData := range modelData.Joints.Values() {
			if jointData == nil {
				continue
			}
			jpName, enName = translatedNamesByNumber(number, jointData.Name(), jointData.EnglishName, nameItems)
			jointData.SetName(jpName)
			jointData.EnglishName = enName
			number++
		}
	}

	return uc.modelWriter.Save(outputPath, modelData, options)
}

// SaveCsvDictionary はCSV出力タブでチェックされた名称を辞書CSVとして保存する。
func (uc *PmxTranslatorUsecase) SaveCsvDictionary(modelData *model.PmxModel, checkedNames []string, outputPath string, options SaveOptions) error {
	if modelData == nil {
		return newPrerequisiteMissingError(errorCsvTargetModelRequired)
	}
	if strings.TrimSpace(outputPath) == "" {
		return newPrerequisiteMissingError(errorOutputCsvPathRequired)
	}
	if uc.csvWriter == nil {
		return newPrerequisiteMissingError(errorCsvWriterMissing)
	}

	_, fileName, _ := mfile.SplitPath(modelData.Path())
	seen := map[string]struct{}{}
	records := make([]domain.TranslationCsvRecord, 0, len(checkedNames))
	for _, name := range checkedNames {
		if strings.TrimSpace(name) == "" {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		records = append(records, domain.TranslationCsvRecord{
			FileName:     fileName,
			SourceName:   name,
			JapaneseName: "",
			EnglishName:  "",
		})
	}

	csvModel, err := io_csv.Marshal(records)
	if err != nil {
		return err
	}
	return uc.csvWriter.Save(outputPath, csvModel, options)
}

// BuildAppendNameItems は追加元CSVと追加CSVの差分行一覧を生成する。
func (uc *PmxTranslatorUsecase) BuildAppendNameItems(originalRows []domain.TranslationCsvRecord, appendRows []domain.TranslationCsvRecord) []domain.AppendNameItem {
	items := make([]domain.AppendNameItem, 0, len(originalRows)+len(appendRows))
	seen := map[string]struct{}{}

	appendItem := func(row domain.TranslationCsvRecord, isOriginal bool) {
		_, exists := seen[row.SourceName]
		seen[row.SourceName] = struct{}{}
		items = append(items, domain.AppendNameItem{
			Number:       len(items) + 1,
			Checked:      !exists,
			SourceName:   row.SourceName,
			JapaneseName: row.JapaneseName,
			EnglishName:  row.EnglishName,
			IsOriginal:   isOriginal,
		})
	}

	for _, row := range originalRows {
		appendItem(row, true)
	}
	for _, row := range appendRows {
		appendItem(row, false)
	}

	return items
}

// SaveAppendCsv はCSV追加タブのチェック行を結合して保存する。
func (uc *PmxTranslatorUsecase) SaveAppendCsv(originalRows []domain.TranslationCsvRecord, appendRows []domain.TranslationCsvRecord, checkedItems []domain.AppendNameItem, outputPath string, options SaveOptions) error {
	if strings.TrimSpace(outputPath) == "" {
		return newPrerequisiteMissingError(errorOutputCsvPathRequired)
	}
	if uc.csvWriter == nil {
		return newPrerequisiteMissingError(errorCsvWriterMissing)
	}

	checkedByNumber := map[int]domain.AppendNameItem{}
	for _, item := range checkedItems {
		if !item.Checked {
			continue
		}
		checkedByNumber[item.Number] = item
	}

	records := make([]domain.TranslationCsvRecord, 0)
	number := 1
	for _, row := range originalRows {
		if item, exists := checkedByNumber[number]; exists && item.SourceName == row.SourceName {
			records = append(records, row)
		}
		number++
	}
	for _, row := range appendRows {
		if item, exists := checkedByNumber[number]; exists && item.SourceName == row.SourceName {
			records = append(records, row)
		}
		number++
	}

	csvModel, err := io_csv.Marshal(records)
	if err != nil {
		return err
	}
	return uc.csvWriter.Save(outputPath, csvModel, options)
}

// BuildTranslationOutputPath は置換一覧からPMX出力パス既定値を生成する。
func BuildTranslationOutputPath(items []domain.TranslateNameItem, fallbackPath string) string {
	if len(items) > 0 && strings.TrimSpace(items[0].JapaneseNameText) != "" {
		return mfile.CreateOutputPath(items[0].JapaneseNameText, "")
	}
	if strings.TrimSpace(fallbackPath) == "" {
		return ""
	}
	return mfile.CreateOutputPath(fallbackPath, "")
}

// BuildCsvOutputPath はCSV出力ファイル既定値を生成する。
func BuildCsvOutputPath(modelPath string) string {
	if strings.TrimSpace(modelPath) == "" {
		return ""
	}
	path := mfile.CreateOutputPath(modelPath, "")
	ext := filepath.Ext(path)
	if strings.EqualFold(ext, ".pmx") || strings.EqualFold(ext, ".pmd") || strings.EqualFold(ext, ".x") {
		return strings.TrimSuffix(path, ext) + ".csv"
	}
	return path
}

// translatedNamesByNumber は番号一致の置換値を返す。
func translatedNamesByNumber(number int, jpText string, enText string, nameItems []domain.TranslateNameItem) (string, string) {
	for _, item := range nameItems {
		if item.Number != number {
			continue
		}
		if item.Checked {
			return item.JapaneseNameText, item.EnglishNameText
		}
		break
	}
	return jpText, enText
}

// translateByRecords は4パス翻訳で名称を置換する。
func translateByRecords(text string, englishText string, records []domain.TranslationCsvRecord, modelName string) (string, string) {
	jpText := text
	enText := englishText

	for _, row := range records {
		if row.FileName == modelName && row.SourceName == jpText {
			jpText = strings.ReplaceAll(jpText, row.SourceName, row.JapaneseName)
			if englishText != "" && row.EnglishName != "" {
				enText = row.EnglishName
			}
		}
	}
	for _, row := range records {
		if row.FileName == "" && row.SourceName == jpText {
			jpText = strings.ReplaceAll(jpText, row.SourceName, row.JapaneseName)
			if englishText != "" && row.EnglishName != "" {
				enText = row.EnglishName
			}
		}
	}
	for _, row := range records {
		if row.FileName == modelName && row.SourceName != "" {
			jpText = strings.ReplaceAll(jpText, row.SourceName, row.JapaneseName)
			if englishText != "" && row.EnglishName != "" {
				enText = strings.ReplaceAll(enText, row.SourceName, row.EnglishName)
			}
		}
	}
	for _, row := range records {
		if row.FileName == "" && row.SourceName != "" {
			jpText = strings.ReplaceAll(jpText, row.SourceName, row.JapaneseName)
			if englishText != "" && row.EnglishName != "" {
				enText = strings.ReplaceAll(enText, row.SourceName, row.EnglishName)
			}
		}
	}

	return jpText, enText
}

// copyTextureFile はテクスチャを出力先へコピーする。
func copyTextureFile(sourcePath string, outputPath string) error {
	textureBytes, err := os.ReadFile(sourcePath)
	if err != nil {
		return io_common.NewIoFileNotFound(sourcePath, err)
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return io_common.NewIoSaveFailed(errorTextureOutputDirCreateFailed, err)
	}
	if err := os.WriteFile(outputPath, textureBytes, 0o644); err != nil {
		return io_common.NewIoSaveFailed(errorTextureSaveFailed, err)
	}
	return nil
}

// newPrerequisiteMissingError は前提不足エラーを生成する。
func newPrerequisiteMissingError(message string) error {
	if strings.TrimSpace(message) == "" {
		message = errorPrerequisiteMissing
	}
	return merr.NewCommonError(prerequisiteMissingErrorID, merr.ErrorKindValidate, message, nil)
}
