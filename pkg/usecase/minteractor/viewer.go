// 指示: miu200521358
package minteractor

import "github.com/miu200521358/mlib_go/pkg/usecase/port/io"

// PmxTranslatorUsecaseDeps は pmx_translator 用ユースケースの依存を表す。
type PmxTranslatorUsecaseDeps struct {
	ModelReader io.IFileReader
	ModelWriter io.IFileWriter
	CsvReader   io.IFileReader
	CsvWriter   io.IFileWriter
}

// PmxTranslatorUsecase は pmx_translator の入出力処理をまとめたユースケースを表す。
type PmxTranslatorUsecase struct {
	modelReader io.IFileReader
	modelWriter io.IFileWriter
	csvReader   io.IFileReader
	csvWriter   io.IFileWriter
}

// NewPmxTranslatorUsecase は pmx_translator 用ユースケースを生成する。
func NewPmxTranslatorUsecase(deps PmxTranslatorUsecaseDeps) *PmxTranslatorUsecase {
	return &PmxTranslatorUsecase{
		modelReader: deps.ModelReader,
		modelWriter: deps.ModelWriter,
		csvReader:   deps.CsvReader,
		csvWriter:   deps.CsvWriter,
	}
}
