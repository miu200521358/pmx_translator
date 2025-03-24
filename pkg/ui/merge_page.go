package ui

import (
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/mcsv"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/pmx_translator/pkg/domain"
	"github.com/miu200521358/pmx_translator/pkg/usecase"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewMergePage(mWidgets *controller.MWidgets) declarative.TabPage {
	var mergeTab *walk.TabPage

	mergeState := domain.NewMergeState()

	var mergeTableView *walk.TableView

	mergedCsvSavePicker := widget.NewCsvSaveFilePicker(
		mi18n.T("出力Csv"),
		mi18n.T("出力Csvファイルパスを指定してください"),
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
		},
	)

	originalCsvLoadPicker := widget.NewCsvLoadFilePicker(
		"AppendOriginalCsv",
		mi18n.T("追加元Csvデータ"),
		mi18n.T("追加元Csvデータファイルを選択してください"),
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				csvData := data.(*mcsv.CsvModel)

				mergeState.OriginalCsvModel = csvData
				mergeState.LoadData()

				mergedCsvSavePicker.SetPath(mergeState.OutputPath)
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	mergedCsvLoadPicker := widget.NewCsvLoadFilePicker(
		"AppendCsv",
		mi18n.T("追加対象Csvデータ"),
		mi18n.T("追加対象Csvデータファイルを選択してください"),
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				csvData := data.(*mcsv.CsvModel)

				mergeState.MergedCsvModel = csvData
				mergeState.LoadData()

				mergedCsvSavePicker.SetPath(mergeState.OutputPath)
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	mWidgets.Widgets = append(mWidgets.Widgets, originalCsvLoadPicker, mergedCsvSavePicker, mergedCsvLoadPicker)

	return declarative.TabPage{

		Title:    mi18n.T("Csv追加"),
		AssignTo: &mergeTab,
		Layout:   declarative.VBox{},
		Background: declarative.SystemColorBrush{
			Color: walk.SysColorInactiveCaption,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					originalCsvLoadPicker.Widgets(),
					mergedCsvLoadPicker.Widgets(),
					declarative.TableView{
						AssignTo:         &mergeTableView,
						AlternatingRowBG: true,
						CheckBoxes:       true,
						ColumnsOrderable: true,
						MultiSelection:   true,
						Model:            mergeState.NameModel,
						MinSize:          declarative.Size{Width: 400, Height: 250},
						Columns: []declarative.TableViewColumn{
							{Title: "#", Width: 50},
							{Title: "No.", Width: 50},
							{Title: mi18n.T("元名称"), Width: 150},
							{Title: mi18n.T("日本語名称"), Width: 150},
							{Title: mi18n.T("英語名称"), Width: 150},
						},
						StyleCell: func(style *walk.CellStyle) {
							if mergeState.NameModel.Checked(style.Row()) {
								if mergeState.NameModel.IsOriginal(style.Row()) {
									// 黄色(元データ)
									style.BackgroundColor = walk.RGB(239, 255, 160)
								} else {
									// 水色(追加データで追加対象)
									style.BackgroundColor = walk.RGB(159, 255, 243)
								}
							} else {
								style.BackgroundColor = walk.RGB(255, 255, 255)
							}
						},
						OnSelectedIndexesChanged: func() {
							if dlg := newMergeTextChangeDialog(
								mergeState,
								mergeTableView.CurrentIndex(),
								&walk.Point{X: mWidgets.Position.X + 100, Y: mWidgets.Position.Y + 100},
							); dlg != nil {
								if cmd, err := dlg.Run(nil); err == nil {
									if cmd == walk.DlgCmdOK {
										mergeState.NameModel.Records[mergeTableView.CurrentIndex()].Checked = true
										mergeState.NameModel.PublishRowsReset()
									}
								} else {
									panic(err)
								}
							}
						},
					},
					mergedCsvSavePicker.Widgets(),
					declarative.VSeparator{},
					declarative.PushButton{
						Text: mi18n.T("保存"),
						OnClicked: func() {
							if err := usecase.MergeCsv(mergeState); err == nil {
								mlog.IT(mi18n.T("出力成功"), mi18n.T("出力成功メッセージ", map[string]any{"Path": mergeState.OutputPath}))
							}

							app.Beep()
						},
					},
					declarative.VSpacer{},
				},
			},
		},
	}
}

func newMergeTextChangeDialog(mergeState *domain.MergeState, recordIndex int, position *walk.Point) *declarative.Dialog {
	var okBtn *walk.PushButton
	var cancelBtn *walk.PushButton
	var db *walk.DataBinder
	var jpTxt *walk.TextEdit
	var enTxt *walk.TextEdit

	if len(mergeState.NameModel.Records) == 0 || recordIndex < 0 {
		return nil
	}

	return newTextChangeDialog(mergeState.TextChangeDialog, okBtn, cancelBtn, db,
		mergeState.NameModel.Records[recordIndex], jpTxt, enTxt, position)
}
