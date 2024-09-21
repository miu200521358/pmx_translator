package ui

import (
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/pmx_translator/pkg/usecase"
	"github.com/miu200521358/walk/pkg/walk"
)

func newCsvTab(controlWindow *controller.ControlWindow, toolState *ToolState) {
	toolState.Tab = widget.NewMTabPage(mi18n.T("Csv出力"))
	controlWindow.AddTabPage(toolState.Tab.TabPage)

	toolState.Tab.SetLayout(walk.NewVBoxLayout())

	var err error

	{
		// Step1. ファイル選択文言
		label, err := walk.NewTextLabel(toolState.Tab)
		if err != nil {
			widget.RaiseError(err)
		}
		label.SetText(mi18n.T("CsvTabLabel"))
	}

	walk.NewVSeparator(toolState.Tab)

	{
		toolState.OriginalCsvPmxPicker = widget.NewPmxReadFilePicker(
			controlWindow,
			toolState.Tab,
			"OriginalPmx",
			mi18n.T("置換対象モデル(Pmx)"),
			mi18n.T("置換対象モデルPmxファイルを選択してください"),
			mi18n.T("置換対象モデルの使い方"))

		toolState.OriginalCsvPmxPicker.SetOnPathChanged(func(path string) {
			if _, err := toolState.OriginalCsvPmxPicker.Load(); err == nil {
				// 出力パス設定
				outputPath := mutils.CreateOutputPath(path, "")
				outputPath = strings.ReplaceAll(outputPath, ".pmx", ".csv")
				toolState.OutputCsvPicker.SetPath(outputPath)
			} else {
				mlog.E(mi18n.T("読み込み失敗"), err)
			}
		})
	}

	{
		toolState.OutputCsvPicker = widget.NewCsvSaveFilePicker(
			controlWindow,
			toolState.Tab,
			mi18n.T("出力Csv"),
			mi18n.T("出力Csvファイルパスを指定してください"),
			mi18n.T("出力Csvファイルパスの使い方"))
	}

	walk.NewVSpacer(toolState.Tab)

	// OKボタン
	{
		toolState.SaveButton, err = walk.NewPushButton(toolState.Tab)
		if err != nil {
			widget.RaiseError(err)
		}
		toolState.SaveButton.SetText(mi18n.T("保存"))
		toolState.SaveButton.Clicked().Attach(toolState.onClickCsvSave)
	}
}

func (toolState *ToolState) onClickCsvSave() {
	if !toolState.OriginalCsvPmxPicker.Exists() {
		mlog.ILT("生成失敗", "生成失敗メッセージ")
		return
	}

	if err := usecase.CsvSave(
		toolState.OriginalCsvPmxPicker.GetCache().(*pmx.PmxModel),
		toolState.OutputCsvPicker.GetPath()); err != nil {
		mlog.ET(mi18n.T("出力失敗"), mi18n.T("Csv出力失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
		return
	}

	widget.Beep()
}
