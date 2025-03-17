package ui

import (
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

func newTextChangeDialog(
	dlg *walk.Dialog,
	okBtn, cancelBtn *walk.PushButton,
	db *walk.DataBinder,
	dataSource any,
	jpTxt, enTxt *walk.TextEdit,
	position *walk.Point,
) *declarative.Dialog {
	return &declarative.Dialog{
		AssignTo:      &dlg,
		DefaultButton: &okBtn,
		CancelButton:  &cancelBtn,
		Title:         mi18n.T("名称変更"),
		Layout:        declarative.VBox{},
		MinSize:       declarative.Size{Width: 400, Height: 200},
		Position:      position,
		DataBinder: declarative.DataBinder{
			AssignTo:   &db,
			DataSource: dataSource,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.Grid{Columns: 2},
				Children: []declarative.Widget{
					declarative.Label{
						Text: mi18n.T("種類"),
					},
					declarative.Label{
						Text: declarative.Bind("TypeText"),
					},
					declarative.Label{
						Text: mi18n.T("元名称"),
					},
					declarative.Label{
						Text: declarative.Bind("NameText"),
					},
					declarative.Label{
						Text: mi18n.T("日本語名称"),
					},
					declarative.TextEdit{
						AssignTo: &jpTxt,
						Text:     declarative.Bind("JapaneseNameText", textRequired{title: mi18n.T("日本語名称")}),
					},
					declarative.Label{
						Text: mi18n.T("英語名称"),
					},
					declarative.TextEdit{
						AssignTo: &enTxt,
						Text:     declarative.Bind("EnglishNameText"),
					},
				},
			},
			declarative.Composite{
				Layout: declarative.HBox{
					Alignment: declarative.AlignHFarVCenter,
				},
				Children: []declarative.Widget{
					declarative.PushButton{
						AssignTo: &okBtn,
						Text:     mi18n.T("OK"),
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								panic(err)
							}
							dlg.Accept()
						},
					},
					declarative.PushButton{
						AssignTo: &cancelBtn,
						Text:     mi18n.T("キャンセル"),
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}
}
