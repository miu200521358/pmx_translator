package domain

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type TranslateModel struct {
	Model             *pmx.PmxModel  // 処理対象モデル
	Motion            *vmd.VmdMotion // 処理対象モーション
	OutputModelPath   string         // 出力パス
	LangCsv           *core.CsvModel // 言語CSVデータ
	AppendOriginalCsv *core.CsvModel // 追加元CSVデータ
	AppendCsv         *core.CsvModel // 追加CSVデータ
}

func NewTranslateModel() *TranslateModel {
	return &TranslateModel{
		Motion: vmd.NewVmdMotion(""),
	}
}

type NameItem struct {
	Number           int
	Checked          bool
	TypeText         string
	Index            int
	NameText         string
	JapaneseNameText string
	EnglishNameText  string
	Segmented        bool
	IsOriginal       bool
}
