package ui

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/pmx_translator/pkg/usecase"
	"github.com/miu200521358/walk/pkg/walk"
)

func newAppendTab(controlWindow *controller.ControlWindow, toolState *ToolState) {
	toolState.AppendTab = widget.NewMTabPage(mi18n.T("Csv追加"))
	controlWindow.AddTabPage(toolState.AppendTab.TabPage)

	toolState.AppendTab.SetLayout(walk.NewVBoxLayout())

	var err error

	{
		label, err := walk.NewTextLabel(toolState.AppendTab)
		if err != nil {
			widget.RaiseError(err)
		}
		label.SetText(mi18n.T("AppendTabLabel"))
	}

	walk.NewVSeparator(toolState.AppendTab)

	{
		toolState.AppendOriginalCsvPicker = widget.NewCsvReadFilePicker(
			controlWindow,
			toolState.AppendTab,
			"AppendOriginalCsv",
			mi18n.T("追加元CSVデータ"),
			mi18n.T("追加元CSVデータファイルを選択してください"),
			mi18n.T("追加元CSVデータの使い方"))

		toolState.AppendOriginalCsvPicker.SetOnPathChanged(func(path string) {
			if data, err := toolState.AppendOriginalCsvPicker.Load(path); err == nil {
				if data == nil {
					return
				}
				toolState.TranslateModel.AppendOriginalCsv = data.(*core.CsvModel)

				if toolState.TranslateModel.AppendCsv != nil {
					toolState.AppendTableView.ResetModel(
						toolState.TranslateModel.AppendOriginalCsv, toolState.TranslateModel.AppendCsv)

					// 出力パス設定
					outputPath := mutils.CreateOutputPath(
						toolState.AppendOriginalCsvPicker.GetPath(), "")
					toolState.AppendOutputPicker.SetPath(outputPath)
				}
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		})
	}

	{
		toolState.AppendCsvPicker = widget.NewCsvReadFilePicker(
			controlWindow,
			toolState.AppendTab,
			"AppendCsv",
			mi18n.T("追加対象CSVデータ"),
			mi18n.T("追加対象CSVデータファイルを選択してください"),
			mi18n.T("追加対象CSVデータの使い方"))

		toolState.AppendCsvPicker.SetOnPathChanged(func(path string) {
			if data, err := toolState.AppendCsvPicker.Load(path); err == nil {
				if data == nil {
					return
				}
				toolState.TranslateModel.AppendCsv = data.(*core.CsvModel)

				if toolState.TranslateModel.AppendOriginalCsv != nil {
					toolState.AppendTableView.ResetModel(
						toolState.TranslateModel.AppendOriginalCsv, toolState.TranslateModel.AppendCsv)

					// 出力パス設定
					outputPath := mutils.CreateOutputPath(
						toolState.AppendOriginalCsvPicker.GetPath(), "")
					toolState.AppendOutputPicker.SetPath(outputPath)
				}
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		})
	}

	{
		toolState.AppendOutputPicker = widget.NewCsvSaveFilePicker(
			controlWindow,
			toolState.AppendTab,
			mi18n.T("出力Csv"),
			mi18n.T("出力Csvファイルパスを指定してください"),
			mi18n.T("出力Csvファイルパスの使い方"))
	}

	walk.NewVSeparator(toolState.AppendTab)

	// AppendTableView
	toolState.AppendTableView = NewAppendTableView(toolState.AppendTab, nil, nil)

	walk.NewVSpacer(toolState.AppendTab)

	// OKボタン
	{
		toolState.AppendSaveButton, err = walk.NewPushButton(toolState.AppendTab)
		if err != nil {
			widget.RaiseError(err)
		}
		toolState.AppendSaveButton.SetText(mi18n.T("保存"))
		toolState.AppendSaveButton.Clicked().Attach(toolState.onClickAppendCsvSave)
	}
}

func (toolState *ToolState) onClickAppendCsvSave() {
	if !toolState.AppendOriginalCsvPicker.Exists() || !toolState.AppendCsvPicker.Exists() {
		mlog.ILT("生成失敗", "生成失敗メッセージ")
		return
	}

	// No.でソート
	if err := toolState.AppendTableView.Model.Sort(1, walk.SortAscending); err != nil {
		mlog.ET(mi18n.T("出力失敗"), mi18n.T("Csv出力失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
		return
	}

	if err := usecase.CsvAppendSave(
		toolState.AppendOriginalCsvPicker.GetCache().(*core.CsvModel),
		toolState.AppendCsvPicker.GetCache().(*core.CsvModel),
		toolState.AppendTableView.Model.CheckedItems(),
		toolState.AppendOutputPicker.GetPath()); err != nil {
		mlog.ET(mi18n.T("出力失敗"), mi18n.T("Csv出力失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
		return
	}

	widget.Beep()
}
