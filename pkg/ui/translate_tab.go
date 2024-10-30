package ui

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/pmx_translator/pkg/usecase"
	"github.com/miu200521358/walk/pkg/walk"
)

func newTranslateTab(controlWindow *controller.ControlWindow, toolState *ToolState) {
	toolState.TranslateTab = widget.NewMTabPage(mi18n.T("名称置換"))
	controlWindow.AddTabPage(toolState.TranslateTab.TabPage)

	toolState.TranslateTab.SetLayout(walk.NewVBoxLayout())

	var err error

	{
		label, err := walk.NewTextLabel(toolState.TranslateTab)
		if err != nil {
			widget.RaiseError(err)
		}
		label.SetText(mi18n.T("TranslateTabLabel"))
	}

	walk.NewVSeparator(toolState.TranslateTab)

	{
		toolState.OriginalPmxPicker = widget.NewPmxReadFilePicker(
			controlWindow,
			toolState.TranslateTab,
			"OriginalPmx",
			mi18n.T("置換対象モデル(Pmx)"),
			mi18n.T("置換対象モデルPmxファイルを選択してください"),
			mi18n.T("置換対象モデルの使い方"))

		toolState.OriginalPmxPicker.SetOnPathChanged(func(path string) {
			if data, err := toolState.OriginalPmxPicker.Load(path); err == nil {
				if data == nil {
					return
				}
				toolState.TranslateModel.Model = data.(*pmx.PmxModel)

				if toolState.TranslateModel.LangCsv != nil {
					toolState.TranslateTableView.ResetModel(
						toolState.TranslateModel.Model, toolState.TranslateModel.LangCsv)

					// 出力パス設定
					outputPath := mutils.CreateOutputPath(
						toolState.TranslateTableView.NameModel.Records[0].JapaneseNameText, "")
					toolState.OutputPmxPicker.SetPath(outputPath)
				}
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		})
	}

	{
		toolState.LangCsvPicker = widget.NewCsvReadFilePicker(
			controlWindow,
			toolState.TranslateTab,
			"LangCsv",
			mi18n.T("置換辞書データ(Csv)"),
			mi18n.T("置換辞書データファイルを選択してください"),
			mi18n.T("置換辞書データの使い方"))

		toolState.LangCsvPicker.SetOnPathChanged(func(path string) {
			if data, err := toolState.LangCsvPicker.Load(path); err == nil {
				if data == nil {
					return
				}
				toolState.TranslateModel.LangCsv = data.(*core.CsvModel)

				if toolState.TranslateModel.Model != nil {
					toolState.TranslateTableView.ResetModel(
						toolState.TranslateModel.Model, toolState.TranslateModel.LangCsv)

					// 出力パス設定
					outputPath := mutils.CreateOutputPath(
						toolState.TranslateTableView.NameModel.Records[0].JapaneseNameText, "")
					toolState.OutputPmxPicker.SetPath(outputPath)
				}
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		})
	}

	toolState.TranslateTableView = NewTranslateTableView(toolState.TranslateTab, nil, nil)

	{
		toolState.OutputPmxPicker = widget.NewPmxSaveFilePicker(
			controlWindow,
			toolState.TranslateTab,
			mi18n.T("出力モデル(Pmx)"),
			mi18n.T("出力モデル(Pmx)ファイルパスを指定してください"),
			mi18n.T("出力モデルの使い方"))
	}

	walk.NewVSpacer(toolState.TranslateTab)

	// OKボタン
	{
		toolState.SaveButton, err = walk.NewPushButton(toolState.TranslateTab)
		if err != nil {
			widget.RaiseError(err)
		}
		toolState.SaveButton.SetText(mi18n.T("保存"))
		toolState.SaveButton.Clicked().Attach(toolState.onClickSave)
	}

	toolState.App.SetFuncGetModels(
		func() [][]*pmx.PmxModel {
			return [][]*pmx.PmxModel{
				{toolState.TranslateModel.Model},
			}
		},
	)

	toolState.App.SetFuncGetMotions(
		func() [][]*vmd.VmdMotion {
			return [][]*vmd.VmdMotion{
				{toolState.TranslateModel.Motion},
			}
		},
	)
}

func (toolState *ToolState) onClickSave() {
	if !toolState.OriginalPmxPicker.Exists() {
		mlog.ILT("生成失敗", "生成失敗メッセージ")
		return
	}

	if err := usecase.Save(
		toolState.OriginalPmxPicker.GetCache().(*pmx.PmxModel),
		toolState.TranslateTableView.NameModel.Records,
		toolState.OutputPmxPicker.GetPath()); err != nil {
		mlog.ET(mi18n.T("出力失敗"), mi18n.T("出力失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
		return
	}

	widget.Beep()
}
