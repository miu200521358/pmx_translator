package domain

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mcsv"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type TranslateState struct {
	Model     *pmx.PmxModel       // 処理対象モデル
	CsvData   *mcsv.CsvModel      // 言語CSVデータ
	NameModel *TranslateNameModel // 名称モデル
}

func NewTranslateState() *TranslateState {
	return &TranslateState{
		NameModel: new(TranslateNameModel),
	}
}

func (t *TranslateState) LoadData() {
	if t.Model == nil || t.CsvData == nil {
		return
	}

	t.NameModel.ResetRows(t.Model, t.CsvData)
}
