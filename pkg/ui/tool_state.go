package ui

import (
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/pmx_renamer/pkg/model"
	"github.com/miu200521358/walk/pkg/walk"
)

type ToolState struct {
	App               *app.MApp
	ControlWindow     *controller.ControlWindow
	TranslateModel    *model.TranslateModel
	Tab               *widget.MTabPage
	OriginalPmxPicker *widget.FilePicker
	LangCsvPicker     *widget.FilePicker
	OutputPmxPicker   *widget.FilePicker
	SaveButton        *walk.PushButton
}

func NewToolState(app *app.MApp, controlWindow *controller.ControlWindow) *ToolState {
	toolState := &ToolState{
		App:            app,
		ControlWindow:  controlWindow,
		TranslateModel: model.NewTranslateModel(),
	}

	newTab(controlWindow, toolState)

	return toolState
}
