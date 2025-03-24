package domain

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mcsv"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type TranslateModel struct {
	Model             *pmx.PmxModel  // 処理対象モデル
	Motion            *vmd.VmdMotion // 処理対象モーション
	OutputModelPath   string         // 出力パス
	LangCsv           *mcsv.CsvModel // 言語Csvデータ
	AppendOriginalCsv *mcsv.CsvModel // 追加元Csvデータ
	AppendCsv         *mcsv.CsvModel // 追加Csvデータ
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
