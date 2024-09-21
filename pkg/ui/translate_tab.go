package ui

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/pmx_renamer/pkg/usecase"
	"github.com/miu200521358/walk/pkg/walk"
)

func newTab(controlWindow *controller.ControlWindow, toolState *ToolState) {
	toolState.Tab = widget.NewMTabPage("翻訳")
	controlWindow.AddTabPage(toolState.Tab.TabPage)

	toolState.Tab.SetLayout(walk.NewVBoxLayout())

	var err error

	{
		// Step1. ファイル選択文言
		label, err := walk.NewTextLabel(toolState.Tab)
		if err != nil {
			widget.RaiseError(err)
		}
		label.SetText("Step1Label")
	}

	walk.NewVSeparator(toolState.Tab)

	{
		toolState.OriginalPmxPicker = widget.NewPmxReadFilePicker(
			controlWindow,
			toolState.Tab,
			"OriginalPmx",
			"日本語化対象モデル(Pmx)",
			"日本語化対象モデルPmxファイルを選択してください",
			"日本語化対象モデルの使い方")

		toolState.OriginalPmxPicker.SetOnPathChanged(func(path string) {
			if data, err := toolState.OriginalPmxPicker.Load(); err == nil {
				// 出力パス設定
				outputPath := mutils.CreateOutputPath(path, "jp")
				toolState.OutputPmxPicker.SetPath(outputPath)

				toolState.TranslateModel.Model = data.(*pmx.PmxModel)
			} else {
				mlog.E(fmt.Sprintf("読み込み失敗: %s", err))
			}
		})
	}

	{
		toolState.LangCsvPicker = widget.NewCsvReadFilePicker(
			controlWindow,
			toolState.Tab,
			"LangCsv",
			"日本誤翻訳辞書(csv)",
			"日本誤翻訳辞書Csvファイルを選択してください",
			"日本誤翻訳辞書の使い方")

		toolState.LangCsvPicker.SetOnPathChanged(func(path string) {
			if data, err := toolState.LangCsvPicker.Load(); err == nil {
				toolState.TranslateModel.LangCsv = data.(*core.CsvModel)
			} else {
				mlog.E(fmt.Sprintf("読み込み失敗: %s", err))
			}
		})
	}

	{
		toolState.OutputPmxPicker = widget.NewPmxSaveFilePicker(
			controlWindow,
			toolState.Tab,
			"出力モデル(Pmx)",
			"出力モデル(Pmx)ファイルパスを指定してください",
			"出力モデルの使い方")
	}

	walk.NewVSpacer(toolState.Tab)

	// OKボタン
	{
		toolState.SaveButton, err = walk.NewPushButton(toolState.Tab)
		if err != nil {
			widget.RaiseError(err)
		}
		toolState.SaveButton.SetText("保存")
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
		mlog.ILT("設定失敗", "Step1失敗")
		return
	}

	if err := usecase.Save(
		toolState.OriginalPmxPicker.GetCache().(*pmx.PmxModel),
		toolState.LangCsvPicker.GetCache().(*core.CsvModel),
		toolState.OutputPmxPicker.GetPath()); err != nil {
		mlog.E(fmt.Sprintf("保存失敗: %s", err))
		return
	}

	widget.Beep()
}
