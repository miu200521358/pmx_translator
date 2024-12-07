package ui

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/pmx_translator/pkg/domain"
	"github.com/miu200521358/walk/pkg/walk"
)

type ToolState struct {
	App                     *app.MApp
	ControlWindow           *controller.ControlWindow
	TranslateModel          *domain.TranslateModel
	TranslateTab            *widget.MTabPage
	OriginalPmxPicker       *widget.FilePicker
	LangCsvPicker           *widget.FilePicker
	OutputPmxPicker         *widget.FilePicker
	OriginalCsvPmxPicker    *widget.FilePicker
	OutputCsvPicker         *widget.FilePicker
	SaveButton              *walk.PushButton
	CsvTab                  *widget.MTabPage
	CsvTableView            *CsvTableView
	TranslateTableView      *TranslateTableView
	AppendTab               *widget.MTabPage
	AppendTableView         *AppendTableView
	AppendOriginalCsvPicker *widget.FilePicker
	AppendCsvPicker         *widget.FilePicker
	AppendOutputPicker      *widget.FilePicker
	AppendSaveButton        *walk.PushButton
}

func NewToolState(app *app.MApp, controlWindow *controller.ControlWindow) *ToolState {
	toolState := &ToolState{
		App:            app,
		ControlWindow:  controlWindow,
		TranslateModel: domain.NewTranslateModel(),
	}

	newTranslateTab(controlWindow, toolState)
	newCsvTab(controlWindow, toolState)
	newAppendTab(controlWindow, toolState)

	return toolState
}

func (toolState *ToolState) SetEnabled(enabled bool) {
	toolState.ControlWindow.SetEnabled(enabled)
}

type loadPmxResult struct {
	model *pmx.PmxModel
	err   error
}
