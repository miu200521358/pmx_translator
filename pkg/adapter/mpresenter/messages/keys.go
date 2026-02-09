// 指示: miu200521358
// Package messages はUI表示に使うメッセージキーを提供する。
package messages

// メッセージキー一覧。
const (
	HelpOverviewTitle   = "概要"
	HelpOverviewMessage = "概要説明"
	HelpToolAboutTitle  = "ツールについて"
	HelpToolAbout       = "ツールについて説明"

	LabelTranslateTab = "名称置換"
	LabelCsvOutputTab = "Csv出力"
	LabelCsvAppendTab = "Csv追加"
	LabelSave         = "保存"

	LabelTranslateTabTip = "名称置換使い方"
	LabelCsvOutputTabTip = "Csv出力使い方"
	LabelCsvAppendTabTip = "Csv追加使い方"

	LabelOriginalModel    = "置換対象モデル(Pmx)"
	LabelOriginalModelTip = "置換対象モデルPmxファイルを選択してください"
	HelpOriginalModel     = "置換対象モデルの使い方"

	LabelDictionaryCsv    = "置換辞書データ(Csv)"
	LabelDictionaryCsvTip = "置換辞書データファイルを選択してください"
	HelpDictionaryCsv     = "置換辞書データの使い方"

	LabelOutputModel    = "出力モデル(Pmx)"
	LabelOutputModelTip = "出力モデル(Pmx)ファイルパスを指定してください"
	HelpOutputModel     = "出力モデルの使い方"

	LabelOutputCsv    = "出力Csv"
	LabelOutputCsvTip = "出力Csvファイルパスを指定してください"
	HelpOutputCsv     = "出力Csvファイルパスの使い方"

	LabelAppendSourceCsv    = "追加元Csvデータ"
	LabelAppendSourceCsvTip = "追加元Csvデータファイルを選択してください"
	HelpAppendSourceCsv     = "追加元Csvデータの使い方"

	LabelAppendTargetCsv    = "追加対象Csvデータ"
	LabelAppendTargetCsvTip = "追加対象Csvデータファイルを選択してください"
	HelpAppendTargetCsv     = "追加対象Csvデータの使い方"

	HelpTranslateTableTitle = "名称置換テーブル"
	HelpTranslateTable      = "名称置換テーブルの使い方"
	HelpCsvOutputTableTitle = "Csv出力候補テーブル"
	HelpCsvOutputTable      = "Csv出力候補テーブルの使い方"
	HelpCsvAppendTableTitle = "Csv追加テーブル"
	HelpCsvAppendTable      = "Csv追加テーブルの使い方"

	HelpTranslateSaveTitle = "保存ボタン(名称置換)"
	HelpTranslateSave      = "保存ボタン(名称置換)の使い方"
	HelpCsvOutputSaveTitle = "保存ボタン(Csv出力)"
	HelpCsvOutputSave      = "保存ボタン(Csv出力)の使い方"
	HelpCsvAppendSaveTitle = "保存ボタン(Csv追加)"
	HelpCsvAppendSave      = "保存ボタン(Csv追加)の使い方"

	HelpAppendOutputCsvTitle = "出力Csv(追加結果)"
	HelpAppendOutputCsv      = "出力Csv(追加結果)の使い方"

	HelpOpenButtonTitle    = "開くボタン"
	HelpOpenButton         = "開くボタンの使い方"
	HelpHistoryButtonTitle = "履歴ボタン"
	HelpHistoryButton      = "履歴ボタンの使い方"

	LabelTableType         = "種類"
	LabelTableIndex        = "インデックス"
	LabelTableSegmented    = "分割"
	LabelTableSourceName   = "元名称"
	LabelTableJapaneseName = "日本語名称"
	LabelTableEnglishName  = "英語名称"
	LabelNameEditDialog    = "名称変更"
	LabelOK                = "OK"
	LabelCancel            = "キャンセル"
	MessageTextRequired    = "文字列未入力"

	TypePath        = "パス"
	TypeModel       = "モデル"
	TypeMaterial    = "材質"
	TypeTexture     = "テクスチャ"
	TypeBone        = "ボーン"
	TypeMorph       = "モーフ"
	TypeDisplaySlot = "表示枠"
	TypeRigidBody   = "剛体"
	TypeJoint       = "ジョイント"

	MessageLoadFailed   = "読み込み失敗"
	MessageOutputFailed = "出力失敗"
	MessageBuildFailed  = "生成失敗"
	MessageOutputDone   = "出力成功"

	ErrorCsvReaderMissing             = "error_csv_reader_missing"
	ErrorCsvModelConvertFailed        = "error_csv_model_convert_failed"
	ErrorCsvModelNil                  = "error_csv_model_nil"
	ErrorCsvHeaderNotFound            = "error_csv_header_not_found"
	ErrorTranslationCsvColumnsInvalid = "error_translation_csv_columns_invalid"
	ErrorTranslateTargetModelRequired = "error_translate_target_model_required"
	ErrorOutputModelPathRequired      = "error_output_model_path_required"
	ErrorModelWriterMissing           = "error_model_writer_missing"
	ErrorOutputDirCreateFailed        = "error_output_dir_create_failed"
	ErrorCsvTargetModelRequired       = "error_csv_target_model_required"
	ErrorOutputCsvPathRequired        = "error_output_csv_path_required"
	ErrorCsvWriterMissing             = "error_csv_writer_missing"
	ErrorTextureOutputDirCreateFailed = "error_texture_output_dir_create_failed"
	ErrorTextureSaveFailed            = "error_texture_save_failed"
	ErrorPrerequisiteMissing          = "error_prerequisite_missing"
)
