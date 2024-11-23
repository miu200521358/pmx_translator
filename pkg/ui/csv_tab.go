package ui

import (
	"regexp"
	"runtime"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
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
	// 拡張子変更(大文字小文字無視)
	pmxToCsvRegex := regexp.MustCompile(`(?i)\.pmx$`)

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
			toolState.SetEnabled(false)

			if canLoad, err := toolState.OriginalCsvPmxPicker.CanLoad(); !canLoad {
				if err != nil {
					mlog.ET(mi18n.T("読み込み失敗"), err.Error())
				}
				return
			}

			resultChan := make(chan loadPmxResult, 1)
			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()

				var loadResult loadPmxResult
				rep := repository.NewPmxRepository()
				if data, err := rep.Load(path); err != nil {
					loadResult.model = nil
					loadResult.err = err
					resultChan <- loadResult
					return
				} else {
					loadResult.model = data.(*pmx.PmxModel)
					loadResult.err = nil
					resultChan <- loadResult
				}
			}()

			// 非同期で結果を受け取る
			go func() {
				wg.Wait()
				close(resultChan)

				result := <-resultChan

				if result.err != nil {
					mlog.ET(mi18n.T("読み込み失敗"), err.Error())
				} else if result.model != nil {
					toolState.OriginalCsvPmxPicker.SetCache(result.model)

					// CsvTableView
					toolState.CsvTableView.ResetModel(result.model)
				}

				go func() {
					runtime.GC() // 読み込み時のメモリ解放
				}()

				defer toolState.ControlWindow.Synchronize(func() {
					// 出力パス設定
					outputPath := mutils.CreateOutputPath(path, "")
					outputPath = pmxToCsvRegex.ReplaceAllString(outputPath, ".csv")
					toolState.OutputCsvPicker.SetPath(outputPath)
					// 画面活性化
					toolState.SetEnabled(true)
				})
			}()
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
