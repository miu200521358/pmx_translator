// 指示: miu200521358
package domain

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

const (
	// NameTypePath はパス項目を表す。
	NameTypePath = "path"
	// NameTypeModel はモデル項目を表す。
	NameTypeModel = "model"
	// NameTypeMaterial は材質項目を表す。
	NameTypeMaterial = "material"
	// NameTypeTexture はテクスチャ項目を表す。
	NameTypeTexture = "texture"
	// NameTypeBone はボーン項目を表す。
	NameTypeBone = "bone"
	// NameTypeMorph はモーフ項目を表す。
	NameTypeMorph = "morph"
	// NameTypeDisplaySlot は表示枠項目を表す。
	NameTypeDisplaySlot = "display_slot"
	// NameTypeRigidBody は剛体項目を表す。
	NameTypeRigidBody = "rigid_body"
	// NameTypeJoint はジョイント項目を表す。
	NameTypeJoint = "joint"
)

// TranslationCsvRecord は pmx_translator 専用CSVの1行を表す。
type TranslationCsvRecord struct {
	FileName     string `csv:"ファイル名"`
	SourceName   string `csv:"元名称"`
	JapaneseName string `csv:"日本語名称"`
	EnglishName  string `csv:"英語名称"`
}

// TranslateModel はツール全体の状態を保持する。
type TranslateModel struct {
	Model             *model.PmxModel
	Motion            *motion.VmdMotion
	OutputModelPath   string
	LangCsv           []TranslationCsvRecord
	AppendOriginalCsv []TranslationCsvRecord
	AppendCsv         []TranslationCsvRecord
}

// NewTranslateModel は TranslateModel を生成する。
func NewTranslateModel() *TranslateModel {
	return &TranslateModel{
		Motion: motion.NewVmdMotion(""),
	}
}

// TranslateNameItem は名称置換タブの1行を表す。
type TranslateNameItem struct {
	Number           int
	Checked          bool
	TypeKey          string
	Index            int
	NameText         string
	JapaneseNameText string
	EnglishNameText  string
}

// CsvCandidateItem はCSV出力タブの1行を表す。
type CsvCandidateItem struct {
	Number          int
	Checked         bool
	TypeKey         string
	Segmented       bool
	NameText        string
	EnglishNameText string
}

// AppendNameItem はCSV追加タブの1行を表す。
type AppendNameItem struct {
	Number       int
	Checked      bool
	SourceName   string
	JapaneseName string
	EnglishName  string
	IsOriginal   bool
}
