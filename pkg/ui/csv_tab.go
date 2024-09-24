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
	toolState.CsvTab = widget.NewMTabPage(mi18n.T("Csv出力"))
	controlWindow.AddTabPage(toolState.CsvTab.TabPage)

	toolState.CsvTab.SetLayout(walk.NewVBoxLayout())

	var err error

	{
		label, err := walk.NewTextLabel(toolState.CsvTab)
		if err != nil {
			widget.RaiseError(err)
		}
		label.SetText(mi18n.T("CsvTabLabel"))
	}

	walk.NewVSeparator(toolState.CsvTab)

	{
		toolState.OriginalCsvPmxPicker = widget.NewPmxReadFilePicker(
			controlWindow,
			toolState.CsvTab,
			"OriginalPmx",
			mi18n.T("置換対象モデル(Pmx)"),
			mi18n.T("置換対象モデルPmxファイルを選択してください"),
			mi18n.T("置換対象モデルの使い方"))

		toolState.OriginalCsvPmxPicker.SetOnPathChanged(func(path string) {
			if data, err := toolState.OriginalCsvPmxPicker.Load(); err == nil {
				if data == nil {
					return
				}

				// 出力パス設定
				outputPath := mutils.CreateOutputPath(path, "")
				outputPath = strings.ReplaceAll(outputPath, ".pmx", ".csv")
				toolState.OutputCsvPicker.SetPath(outputPath)

				// CsvTableView
				toolState.CsvTableView.ResetModel(data.(*pmx.PmxModel))
			} else {
				mlog.E(mi18n.T("読み込み失敗"), err)
			}
		})
	}

	{
		toolState.OutputCsvPicker = widget.NewCsvSaveFilePicker(
			controlWindow,
			toolState.CsvTab,
			mi18n.T("出力Csv"),
			mi18n.T("出力Csvファイルパスを指定してください"),
			mi18n.T("出力Csvファイルパスの使い方"))
	}

	walk.NewVSeparator(toolState.CsvTab)

	// CsvTableView
	toolState.CsvTableView = NewCsvTableView(toolState.CsvTab, nil)

	walk.NewVSpacer(toolState.CsvTab)

	// OKボタン
	{
		toolState.SaveButton, err = walk.NewPushButton(toolState.CsvTab)
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

	// No.でソート
	if err := toolState.CsvTableView.Model.Sort(1, walk.SortAscending); err != nil {
		mlog.ET(mi18n.T("出力失敗"), mi18n.T("Csv出力失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
		return
	}

	if err := usecase.CsvSave(
		toolState.OriginalCsvPmxPicker.GetCache().(*pmx.PmxModel),
		toolState.CsvTableView.Model.CheckedNames(),
		toolState.OutputCsvPicker.GetPath()); err != nil {
		mlog.ET(mi18n.T("出力失敗"), mi18n.T("Csv出力失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
		return
	}

	widget.Beep()
}
