// 指示: miu200521358
package minteractor

import (
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_csv"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/file/mfile"
	"github.com/miu200521358/mlib_go/pkg/usecase"
	"github.com/miu200521358/mlib_go/pkg/usecase/port/io"
	"github.com/miu200521358/pmx_translator/pkg/domain"
)

const (
	translationCsvColumns = 4
)

// LoadModel はモデルを読み込む。
func (uc *PmxTranslatorUsecase) LoadModel(reader io.IFileReader, path string) (*model.PmxModel, error) {
	repo := reader
	if repo == nil {
		repo = uc.modelReader
	}
	return usecase.LoadModel(repo, path)
}

// LoadTranslationCsv は翻訳CSVを読み込み、4列レコードとして返す。
func (uc *PmxTranslatorUsecase) LoadTranslationCsv(reader io.IFileReader, path string) ([]domain.TranslationCsvRecord, error) {
	repo := reader
	if repo == nil {
		repo = uc.csvReader
	}
	if repo == nil {
		return nil, newPrerequisiteMissingError("CSV読み込みリポジトリがありません")
	}

	data, err := repo.Load(path)
	if err != nil {
		return nil, err
	}
	csvModel, ok := data.(*io_csv.CsvModel)
	if !ok || csvModel == nil {
		return nil, io_common.NewIoParseFailed("CSVモデル変換に失敗しました", nil)
	}
	if err := validateTranslationCsvModel(csvModel); err != nil {
		return nil, err
	}

	records := make([]domain.TranslationCsvRecord, 0)
	if err := io_csv.UnmarshalWithOptions(csvModel, &records, io_csv.CsvUnmarshalOptions{
		ColumnMapping: io_csv.CsvColumnMappingOrder,
	}); err != nil {
		return nil, err
	}
	return records, nil
}

// BuildTranslateNameItems はモデルと辞書から名称置換一覧を生成する。
func (uc *PmxTranslatorUsecase) BuildTranslateNameItems(modelData *model.PmxModel, records []domain.TranslationCsvRecord) []domain.TranslateNameItem {
	if modelData == nil {
		return []domain.TranslateNameItem{}
	}

	_, fileName, _ := mfile.SplitPath(modelData.Path())
	items := make([]domain.TranslateNameItem, 0)

	appendItem := func(index int, nameText string, englishName string, typeKey string) {
		jpTranslated, enTranslated := translateByRecords(nameText, englishName, records, fileName)
		items = append(items, domain.TranslateNameItem{
			Number:           len(items) + 1,
			Checked:          nameText != jpTranslated || englishName != enTranslated,
			TypeKey:          typeKey,
			Index:            index,
			NameText:         nameText,
			JapaneseNameText: jpTranslated,
			EnglishNameText:  enTranslated,
		})
	}

	appendItem(0, modelData.Path(), "", domain.NameTypePath)
	appendItem(0, modelData.Name(), modelData.EnglishName, domain.NameTypeModel)

	if modelData.Materials != nil {
		for _, materialData := range modelData.Materials.Values() {
			if materialData == nil {
				continue
			}
			appendItem(materialData.Index(), materialData.Name(), materialData.EnglishName, domain.NameTypeMaterial)
		}
	}
	if modelData.Textures != nil {
		for _, textureData := range modelData.Textures.Values() {
			if textureData == nil {
				continue
			}
			appendItem(textureData.Index(), textureData.Name(), "", domain.NameTypeTexture)
		}
	}
	if modelData.Bones != nil {
		for _, boneData := range modelData.Bones.Values() {
			if boneData == nil {
				continue
			}
			appendItem(boneData.Index(), boneData.Name(), boneData.EnglishName, domain.NameTypeBone)
		}
	}
	if modelData.Morphs != nil {
		for _, morphData := range modelData.Morphs.Values() {
			if morphData == nil {
				continue
			}
			appendItem(morphData.Index(), morphData.Name(), morphData.EnglishName, domain.NameTypeMorph)
		}
	}
	if modelData.DisplaySlots != nil {
		for _, displaySlotData := range modelData.DisplaySlots.Values() {
			if displaySlotData == nil {
				continue
			}
			appendItem(displaySlotData.Index(), displaySlotData.Name(), displaySlotData.EnglishName, domain.NameTypeDisplaySlot)
		}
	}
	if modelData.RigidBodies != nil {
		for _, rigidBodyData := range modelData.RigidBodies.Values() {
			if rigidBodyData == nil {
				continue
			}
			appendItem(rigidBodyData.Index(), rigidBodyData.Name(), rigidBodyData.EnglishName, domain.NameTypeRigidBody)
		}
	}
	if modelData.Joints != nil {
		for _, jointData := range modelData.Joints.Values() {
			if jointData == nil {
				continue
			}
			appendItem(jointData.Index(), jointData.Name(), jointData.EnglishName, domain.NameTypeJoint)
		}
	}

	return items
}

// validateTranslationCsvModel は翻訳CSVの列仕様を検証する。
func validateTranslationCsvModel(csvModel *io_csv.CsvModel) error {
	if csvModel == nil {
		return io_common.NewIoParseFailed("CSVモデルがnilです", nil)
	}
	records := csvModel.Records()
	if len(records) == 0 {
		return io_common.NewIoParseFailed("CSVヘッダが見つかりません", nil)
	}
	if len(records[0]) != translationCsvColumns {
		return io_common.NewIoParseFailed("翻訳CSVの列数が不正です(行:1 期待:%d 実際:%d)", nil, translationCsvColumns, len(records[0]))
	}
	for i := 1; i < len(records); i++ {
		if len(records[i]) != translationCsvColumns {
			return io_common.NewIoParseFailed(
				"翻訳CSVの列数が不正です(行:%d 期待:%d 実際:%d)",
				nil,
				i+1,
				translationCsvColumns,
				len(records[i]),
			)
		}
	}
	return nil
}
