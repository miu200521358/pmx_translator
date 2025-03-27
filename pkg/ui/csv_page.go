package ui

import (
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/pmx_translator/pkg/domain"
	"github.com/miu200521358/pmx_translator/pkg/usecase"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewCsvPage(mWidgets *controller.MWidgets) declarative.TabPage {
	var csvTab *walk.TabPage

	csvState := domain.NewCsvState()

	var csvTableView *walk.TableView

	csvSavePicker := widget.NewCsvSaveFilePicker(
		mi18n.T("出力Csv"),
		mi18n.T("出力Csvファイルパスを指定してください"),
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
		},
	)

	pmxLoadPicker := widget.NewPmxLoadFilePicker(
		"OriginalPmx",
		mi18n.T("置換対象モデル(Pmx)"),
		mi18n.T("置換対象モデルPmxファイルを選択してください"),
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if path == "" {
				cw.StoreModel(0, 0, nil)
				return
			}

			if data, err := rep.Load(path); err == nil {
				model := data.(*pmx.PmxModel)
				cw.StoreModel(0, 0, model)

				csvState.Model = model
				csvState.LoadData()

				csvSavePicker.SetPath(csvState.OutputPath)
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	mWidgets.Widgets = append(mWidgets.Widgets, pmxLoadPicker, csvSavePicker)

	return declarative.TabPage{

		Title:    mi18n.T("Csv出力"),
		AssignTo: &csvTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					pmxLoadPicker.Widgets(),
					declarative.TableView{
						AssignTo:         &csvTableView,
						AlternatingRowBG: true,
						CheckBoxes:       true,
						ColumnsOrderable: true,
						MultiSelection:   true,
						Model:            csvState.NameModel,
						MinSize:          declarative.Size{Width: 400, Height: 250},
						Columns: []declarative.TableViewColumn{
							{Title: "#", Width: 50},
							{Title: "No.", Width: 50},
							{Title: mi18n.T("種類"), Width: 80},
							{Title: mi18n.T("分割"), Width: 50},
							{Title: mi18n.T("日本語名称"), Width: 150},
							{Title: mi18n.T("英語名称"), Width: 150},
						},
						StyleCell: func(style *walk.CellStyle) {
							if csvState.NameModel.Checked(style.Row()) {
								style.BackgroundColor = walk.RGB(159, 255, 243)
							} else {
								style.BackgroundColor = walk.RGB(255, 255, 255)
							}
						},
						OnSelectedIndexesChanged: func() {
							if dlg := newCsvTextChangeDialog(
								csvState,
								csvTableView.CurrentIndex(),
								&walk.Point{X: mWidgets.Position.X + 100, Y: mWidgets.Position.Y + 100},
							); dlg != nil {
								if cmd, err := dlg.Run(nil); err == nil {
									if cmd == walk.DlgCmdOK {
										csvState.NameModel.Records[csvTableView.CurrentIndex()].Checked = true
										csvState.NameModel.PublishRowsReset()
									}
								} else {
									panic(err)
								}
							}
						},
					},
					csvSavePicker.Widgets(),
					declarative.VSeparator{},
					declarative.PushButton{
						Text: mi18n.T("保存"),
						OnClicked: func() {
							if err := usecase.SaveCsv(csvState); err == nil {
								mlog.IT(mi18n.T("出力成功"), mi18n.T("出力成功メッセージ", map[string]any{"Path": csvState.OutputPath}))
							}

							controller.Beep()
						},
					},
					declarative.VSpacer{},
				},
			},
		},
	}
}

func newCsvTextChangeDialog(csvState *domain.CsvState, recordIndex int, position *walk.Point) *declarative.Dialog {
	var okBtn *walk.PushButton
	var cancelBtn *walk.PushButton
	var db *walk.DataBinder
	var jpTxt *walk.TextEdit
	var enTxt *walk.TextEdit

	if len(csvState.NameModel.Records) == 0 || recordIndex < 0 {
		return nil
	}

	return newTextChangeDialog(csvState.TextChangeDialog, okBtn, cancelBtn, db,
		csvState.NameModel.Records[recordIndex], jpTxt, enTxt, position)
}
