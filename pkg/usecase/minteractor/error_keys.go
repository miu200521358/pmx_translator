// 指示: miu200521358
package minteractor

const (
	// errorCsvReaderMissing はCSV読み込みリポジトリ未設定時のメッセージキー。
	errorCsvReaderMissing = "error_csv_reader_missing"
	// errorCsvModelConvertFailed はCSVモデル変換失敗時のメッセージキー。
	errorCsvModelConvertFailed = "error_csv_model_convert_failed"
	// errorCsvModelNil はCSVモデルnil時のメッセージキー。
	errorCsvModelNil = "error_csv_model_nil"
	// errorCsvHeaderNotFound はCSVヘッダ未検出時のメッセージキー。
	errorCsvHeaderNotFound = "error_csv_header_not_found"
	// errorTranslationCsvColumnsInvalid は翻訳CSV列数不正時のメッセージキー。
	errorTranslationCsvColumnsInvalid = "error_translation_csv_columns_invalid"
	// errorTranslateTargetModelRequired は名称置換対象モデル未読込時のメッセージキー。
	errorTranslateTargetModelRequired = "error_translate_target_model_required"
	// errorOutputModelPathRequired は出力モデルパス未指定時のメッセージキー。
	errorOutputModelPathRequired = "error_output_model_path_required"
	// errorModelWriterMissing はモデル保存リポジトリ未設定時のメッセージキー。
	errorModelWriterMissing = "error_model_writer_missing"
	// errorOutputDirCreateFailed は出力ディレクトリ作成失敗時のメッセージキー。
	errorOutputDirCreateFailed = "error_output_dir_create_failed"
	// errorCsvTargetModelRequired はCSV出力対象モデル未読込時のメッセージキー。
	errorCsvTargetModelRequired = "error_csv_target_model_required"
	// errorOutputCsvPathRequired は出力CSVパス未指定時のメッセージキー。
	errorOutputCsvPathRequired = "error_output_csv_path_required"
	// errorCsvWriterMissing はCSV保存リポジトリ未設定時のメッセージキー。
	errorCsvWriterMissing = "error_csv_writer_missing"
	// errorTextureOutputDirCreateFailed はテクスチャ出力先ディレクトリ作成失敗時のメッセージキー。
	errorTextureOutputDirCreateFailed = "error_texture_output_dir_create_failed"
	// errorTextureSaveFailed はテクスチャ保存失敗時のメッセージキー。
	errorTextureSaveFailed = "error_texture_save_failed"
	// errorPrerequisiteMissing は前提不足の既定メッセージキー。
	errorPrerequisiteMissing = "error_prerequisite_missing"
)
