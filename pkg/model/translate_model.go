package model

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type TranslateModel struct {
	Model               *pmx.PmxModel  // 処理対象モデル
	Motion              *vmd.VmdMotion // 処理対象モーション
	OutputModelPath     string         // 出力パス
	LangCsv             *core.CsvModel // 言語CSVパス
	IsAllowAscii        bool           // ASCII文字許可
	IsAllowHiragana     bool           // ひらがな文字許可
	IsAllowKatakanaHan  bool           // 半角カタカナ文字許可
	IsAllowKatakanaZen  bool           // 全角カタカナ文字許可
	IsAllowAlphanumeric bool           // 全角英数字文字許可
	IsAllowKanji        bool           // 漢字文字許可
}

func NewTranslateModel() *TranslateModel {
	return &TranslateModel{
		Motion: vmd.NewVmdMotion(""),
	}
}
