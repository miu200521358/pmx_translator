package ui

import (
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/mcsv"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/pmx_translator/pkg/domain"
	"github.com/miu200521358/pmx_translator/pkg/usecase"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewTranslatePage(mWidgets *controller.MWidgets) declarative.TabPage {
	var translateTab *walk.TabPage

	translateState := domain.NewTranslateState()

	var translateTableView *walk.TableView

	pmxSavePicker := widget.NewPmxSaveFilePicker(
		mi18n.T("出力モデル(Pmx)"),
		mi18n.T("出力モデル(Pmx)ファイルパスを指定してください"),
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

				translateState.Model = model
				translateState.LoadData()

				pmxSavePicker.SetPath(translateState.OutputPath)
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	csvLoadPicker := widget.NewCsvLoadFilePicker(
		"LangCsv",
		mi18n.T("置換辞書データ(Csv)"),
		mi18n.T("置換辞書データファイルを選択してください"),
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				csvData := data.(*mcsv.CsvModel)

				translateState.CsvData = csvData
				translateState.LoadData()

				pmxSavePicker.SetPath(translateState.OutputPath)
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	mWidgets.Widgets = append(mWidgets.Widgets, pmxLoadPicker, csvLoadPicker, pmxSavePicker)

	return declarative.TabPage{

		Title:    mi18n.T("名称置換"),
		AssignTo: &translateTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					pmxLoadPicker.Widgets(),
					csvLoadPicker.Widgets(),
					declarative.TableView{
						AssignTo:         &translateTableView,
						AlternatingRowBG: true,
						CheckBoxes:       true,
						ColumnsOrderable: true,
						MultiSelection:   true,
						Model:            translateState.NameModel,
						MinSize:          declarative.Size{Width: 400, Height: 250},
						Columns: []declarative.TableViewColumn{
							{Title: "#", Width: 50},
							{Title: "No.", Width: 50},
							{Title: mi18n.T("種類"), Width: 80},
							{Title: mi18n.T("インデックス"), Width: 40},
							{Title: mi18n.T("元名称"), Width: 150},
							{Title: mi18n.T("日本語名称"), Width: 150},
							{Title: mi18n.T("英語名称"), Width: 150},
						},
						StyleCell: func(style *walk.CellStyle) {
							if translateState.NameModel.Checked(style.Row()) {
								style.BackgroundColor = walk.RGB(159, 255, 243)
							} else {
								style.BackgroundColor = walk.RGB(255, 255, 255)
							}
						},
						OnSelectedIndexesChanged: func() {
							if dlg := newTranslateTextChangeDialog(
								translateState,
								translateTableView.CurrentIndex(),
								&walk.Point{X: mWidgets.Position.X + 100, Y: mWidgets.Position.Y + 100},
							); dlg != nil {
								if cmd, err := dlg.Run(nil); err == nil {
									if cmd == walk.DlgCmdOK {
										translateState.NameModel.Records[translateTableView.CurrentIndex()].Checked = true
										translateState.NameModel.PublishRowsReset()
									}
								} else {
									panic(err)
								}
							}
						},
					},
					pmxSavePicker.Widgets(),
					declarative.VSeparator{},
					declarative.PushButton{
						Text: mi18n.T("保存"),
						OnClicked: func() {
							if err := usecase.SavePmx(
								translateState.Model,
								translateState.NameModel.Records,
								translateState.OutputPath,
							); err == nil {
								mlog.IT(mi18n.T("出力成功"), mi18n.T("出力成功メッセージ", map[string]interface{}{"Path": translateState.OutputPath}))
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

func newTranslateTextChangeDialog(translateState *domain.TranslateState, recordIndex int, position *walk.Point) *declarative.Dialog {
	var okBtn *walk.PushButton
	var cancelBtn *walk.PushButton
	var db *walk.DataBinder
	var jpTxt *walk.TextEdit
	var enTxt *walk.TextEdit

	if len(translateState.NameModel.Records) == 0 || recordIndex < 0 {
		return nil
	}

	return newTextChangeDialog(translateState.TextChangeDialog, okBtn, cancelBtn, db,
		translateState.NameModel.Records[recordIndex], jpTxt, enTxt, position)
}
