//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/audio_api"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/infra/file/mfile"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/pmx_translator/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/pmx_translator/pkg/domain"
	"github.com/miu200521358/pmx_translator/pkg/usecase/minteractor"
)

const (
	previewWindowIndex = 0
	previewModelIndex  = 0
)

// NewTabPages は pmx_translator のタブページ群を生成する。
func NewTabPages(mWidgets *controller.MWidgets, baseServices base.IBaseServices, initialModelPath string, _ audio_api.IAudioPlayer, viewerUsecase *minteractor.PmxTranslatorUsecase) []declarative.TabPage {
	var translateTab *walk.TabPage
	var csvOutputTab *walk.TabPage
	var csvAppendTab *walk.TabPage

	var translator i18n.II18n
	var logger logging.ILogger
	var userConfig config.IUserConfig
	if baseServices != nil {
		translator = baseServices.I18n()
		logger = baseServices.Logger()
		if cfg := baseServices.Config(); cfg != nil {
			userConfig = cfg.UserConfig()
		}
	}
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	if viewerUsecase == nil {
		viewerUsecase = minteractor.NewPmxTranslatorUsecase(minteractor.PmxTranslatorUsecaseDeps{})
	}

	previewMotion := motion.NewVmdMotion("")

	translateTable := NewTranslateTableView(translator)
	csvTable := NewCsvTableView(translator)
	appendTable := NewAppendTableView(translator)

	var translateModel *model.PmxModel
	translateCsvRows := []domain.TranslationCsvRecord{}
	translateCsvLoaded := false
	translateOutputPath := ""
	var translateOutputPicker *widget.FilePicker

	refreshTranslateRows := func() {
		if translateTable == nil {
			return
		}
		if translateModel == nil || !translateCsvLoaded {
			translateTable.ResetRows([]domain.TranslateNameItem{})
			if translateOutputPicker != nil && translateModel != nil {
				defaultPath := minteractor.BuildTranslationOutputPath(nil, translateModel.Path())
				translateOutputPath = defaultPath
				translateOutputPicker.SetPath(defaultPath)
			}
			return
		}
		items := viewerUsecase.BuildTranslateNameItems(translateModel, translateCsvRows)
		translateTable.ResetRows(items)
		defaultPath := minteractor.BuildTranslationOutputPath(items, translateModel.Path())
		translateOutputPath = defaultPath
		if translateOutputPicker != nil && strings.TrimSpace(defaultPath) != "" {
			translateOutputPicker.SetPath(defaultPath)
		}
	}

	translateModelPicker := widget.NewPmxLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyPmxHistory,
		i18n.TranslateOrMark(translator, messages.LabelOriginalModel),
		i18n.TranslateOrMark(translator, messages.LabelOriginalModelTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			if strings.TrimSpace(path) == "" {
				translateModel = nil
				refreshTranslateRows()
				if cw != nil {
					cw.SetModel(previewWindowIndex, previewModelIndex, nil)
					cw.SetMotion(previewWindowIndex, previewModelIndex, nil)
				}
				return
			}

			modelData, err := viewerUsecase.LoadModel(rep, path)
			if err != nil {
				logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageLoadFailed), err)
				translateModel = nil
				refreshTranslateRows()
				if cw != nil {
					cw.SetModel(previewWindowIndex, previewModelIndex, nil)
					cw.SetMotion(previewWindowIndex, previewModelIndex, nil)
				}
				return
			}
			if modelData == nil {
				logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageLoadFailed), nil)
				translateModel = nil
				refreshTranslateRows()
				if cw != nil {
					cw.SetModel(previewWindowIndex, previewModelIndex, nil)
					cw.SetMotion(previewWindowIndex, previewModelIndex, nil)
				}
				return
			}

			translateModel = modelData
			refreshTranslateRows()
			if cw != nil {
				cw.SetModel(previewWindowIndex, previewModelIndex, modelData)
				cw.SetMotion(previewWindowIndex, previewModelIndex, previewMotion)
			}
		},
	)

	translateCsvPicker := widget.NewCsvLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyCsvHistory,
		i18n.TranslateOrMark(translator, messages.LabelDictionaryCsv),
		i18n.TranslateOrMark(translator, messages.LabelDictionaryCsvTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			_ = cw
			if strings.TrimSpace(path) == "" {
				translateCsvRows = []domain.TranslationCsvRecord{}
				translateCsvLoaded = false
				refreshTranslateRows()
				return
			}
			rows, err := viewerUsecase.LoadTranslationCsv(rep, path)
			if err != nil {
				logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageLoadFailed), err)
				translateCsvRows = []domain.TranslationCsvRecord{}
				translateCsvLoaded = false
				refreshTranslateRows()
				return
			}
			translateCsvRows = rows
			translateCsvLoaded = true
			refreshTranslateRows()
		},
	)

	translateOutputPicker = widget.NewPmxSaveFilePicker(
		userConfig,
		translator,
		i18n.TranslateOrMark(translator, messages.LabelOutputModel),
		i18n.TranslateOrMark(translator, messages.LabelOutputModelTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			_ = cw
			_ = rep
			translateOutputPath = path
		},
	)

	translateSaveButton := widget.NewMPushButton()
	translateSaveButton.SetLabel(i18n.TranslateOrMark(translator, messages.LabelSave))
	translateSaveButton.SetOnClicked(func(cw *controller.ControlWindow) {
		if translateModel == nil {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageBuildFailed), nil)
			logger.Error(i18n.TranslateOrMark(translator, messages.LabelOriginalModelTip))
			return
		}
		if !translateCsvLoaded {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageBuildFailed), nil)
			logger.Error(i18n.TranslateOrMark(translator, messages.LabelDictionaryCsvTip))
			return
		}
		if strings.TrimSpace(translateOutputPath) == "" {
			translateOutputPath = minteractor.BuildTranslationOutputPath(translateTable.Rows(), translateModel.Path())
			if translateOutputPicker != nil && strings.TrimSpace(translateOutputPath) != "" {
				translateOutputPicker.SetPath(translateOutputPath)
			}
		}
		if strings.TrimSpace(translateOutputPath) == "" {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageBuildFailed), nil)
			logger.Error(i18n.TranslateOrMark(translator, messages.LabelOutputModelTip))
			return
		}

		if err := viewerUsecase.SaveTranslatedModel(
			translateOutputPath,
			translateModel,
			translateTable.Rows(),
			minteractor.SaveOptions{},
		); err != nil {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageOutputFailed), err)
			return
		}

		if cw != nil && translateModel != nil {
			cw.SetModel(previewWindowIndex, previewModelIndex, translateModel)
			cw.SetMotion(previewWindowIndex, previewModelIndex, previewMotion)
		}
		controller.Beep()
		logger.Info("%s: %s", i18n.TranslateOrMark(translator, messages.MessageOutputDone), filepath.Base(translateOutputPath))
	})

	csvOutputModel := (*model.PmxModel)(nil)
	csvOutputPath := ""
	csvOutputModelPath := ""
	var csvOutputPicker *widget.FilePicker

	refreshCsvCandidates := func() {
		if csvOutputModel == nil {
			csvTable.ResetRows([]domain.CsvCandidateItem{})
			return
		}
		candidates, err := viewerUsecase.BuildCsvCandidates(csvOutputModel)
		if err != nil {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageLoadFailed), err)
			csvTable.ResetRows([]domain.CsvCandidateItem{})
			return
		}
		csvTable.ResetRows(candidates)
		defaultPath := minteractor.BuildCsvOutputPath(csvOutputModelPath)
		csvOutputPath = defaultPath
		if csvOutputPicker != nil && strings.TrimSpace(defaultPath) != "" {
			csvOutputPicker.SetPath(defaultPath)
		}
	}

	csvOutputModelPicker := widget.NewPmxLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyPmxHistory,
		i18n.TranslateOrMark(translator, messages.LabelOriginalModel),
		i18n.TranslateOrMark(translator, messages.LabelOriginalModelTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			csvOutputModelPath = path
			if strings.TrimSpace(path) == "" {
				csvOutputModel = nil
				refreshCsvCandidates()
				if cw != nil {
					cw.SetModel(previewWindowIndex, previewModelIndex, nil)
					cw.SetMotion(previewWindowIndex, previewModelIndex, nil)
				}
				return
			}

			modelData, err := viewerUsecase.LoadModel(rep, path)
			if err != nil {
				logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageLoadFailed), err)
				csvOutputModel = nil
				refreshCsvCandidates()
				if cw != nil {
					cw.SetModel(previewWindowIndex, previewModelIndex, nil)
					cw.SetMotion(previewWindowIndex, previewModelIndex, nil)
				}
				return
			}
			if modelData == nil {
				logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageLoadFailed), nil)
				csvOutputModel = nil
				refreshCsvCandidates()
				if cw != nil {
					cw.SetModel(previewWindowIndex, previewModelIndex, nil)
					cw.SetMotion(previewWindowIndex, previewModelIndex, nil)
				}
				return
			}

			csvOutputModel = modelData
			refreshCsvCandidates()
			if cw != nil {
				cw.SetModel(previewWindowIndex, previewModelIndex, modelData)
				cw.SetMotion(previewWindowIndex, previewModelIndex, previewMotion)
			}
		},
	)

	csvOutputPicker = widget.NewCsvSaveFilePicker(
		userConfig,
		translator,
		i18n.TranslateOrMark(translator, messages.LabelOutputCsv),
		i18n.TranslateOrMark(translator, messages.LabelOutputCsvTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			_ = cw
			_ = rep
			csvOutputPath = path
		},
	)

	csvOutputSaveButton := widget.NewMPushButton()
	csvOutputSaveButton.SetLabel(i18n.TranslateOrMark(translator, messages.LabelSave))
	csvOutputSaveButton.SetOnClicked(func(cw *controller.ControlWindow) {
		if csvOutputModel == nil {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageBuildFailed), nil)
			logger.Error(i18n.TranslateOrMark(translator, messages.LabelOriginalModelTip))
			return
		}
		if strings.TrimSpace(csvOutputPath) == "" {
			csvOutputPath = minteractor.BuildCsvOutputPath(csvOutputModel.Path())
			if csvOutputPicker != nil && strings.TrimSpace(csvOutputPath) != "" {
				csvOutputPicker.SetPath(csvOutputPath)
			}
		}
		if strings.TrimSpace(csvOutputPath) == "" {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageBuildFailed), nil)
			logger.Error(i18n.TranslateOrMark(translator, messages.LabelOutputCsvTip))
			return
		}

		if err := viewerUsecase.SaveCsvDictionary(
			csvOutputModel,
			csvTable.CheckedNames(),
			csvOutputPath,
			minteractor.SaveOptions{},
		); err != nil {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageOutputFailed), err)
			return
		}

		if cw != nil {
			cw.SetModel(previewWindowIndex, previewModelIndex, csvOutputModel)
			cw.SetMotion(previewWindowIndex, previewModelIndex, previewMotion)
		}
		controller.Beep()
		logger.Info("%s: %s", i18n.TranslateOrMark(translator, messages.MessageOutputDone), filepath.Base(csvOutputPath))
	})

	appendSourceRows := []domain.TranslationCsvRecord{}
	appendTargetRows := []domain.TranslationCsvRecord{}
	appendSourceLoaded := false
	appendTargetLoaded := false
	appendSourcePath := ""
	appendOutputPath := ""
	var appendOutputPicker *widget.FilePicker

	refreshAppendRows := func() {
		if !appendSourceLoaded || !appendTargetLoaded {
			appendTable.ResetRows([]domain.AppendNameItem{})
			return
		}
		items := viewerUsecase.BuildAppendNameItems(appendSourceRows, appendTargetRows)
		appendTable.ResetRows(items)
		defaultPath := buildAppendCsvOutputPath(appendSourcePath)
		appendOutputPath = defaultPath
		if appendOutputPicker != nil && strings.TrimSpace(defaultPath) != "" {
			appendOutputPicker.SetPath(defaultPath)
		}
	}

	appendSourcePicker := widget.NewCsvLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyCsvHistory,
		i18n.TranslateOrMark(translator, messages.LabelAppendSourceCsv),
		i18n.TranslateOrMark(translator, messages.LabelAppendSourceCsvTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			_ = cw
			appendSourcePath = path
			if strings.TrimSpace(path) == "" {
				appendSourceRows = []domain.TranslationCsvRecord{}
				appendSourceLoaded = false
				refreshAppendRows()
				return
			}
			rows, err := viewerUsecase.LoadTranslationCsv(rep, path)
			if err != nil {
				logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageLoadFailed), err)
				appendSourceRows = []domain.TranslationCsvRecord{}
				appendSourceLoaded = false
				refreshAppendRows()
				return
			}
			appendSourceRows = rows
			appendSourceLoaded = true
			refreshAppendRows()
		},
	)

	appendTargetPicker := widget.NewCsvLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyCsvHistory,
		i18n.TranslateOrMark(translator, messages.LabelAppendTargetCsv),
		i18n.TranslateOrMark(translator, messages.LabelAppendTargetCsvTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			_ = cw
			if strings.TrimSpace(path) == "" {
				appendTargetRows = []domain.TranslationCsvRecord{}
				appendTargetLoaded = false
				refreshAppendRows()
				return
			}
			rows, err := viewerUsecase.LoadTranslationCsv(rep, path)
			if err != nil {
				logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageLoadFailed), err)
				appendTargetRows = []domain.TranslationCsvRecord{}
				appendTargetLoaded = false
				refreshAppendRows()
				return
			}
			appendTargetRows = rows
			appendTargetLoaded = true
			refreshAppendRows()
		},
	)

	appendOutputPicker = widget.NewCsvSaveFilePicker(
		userConfig,
		translator,
		i18n.TranslateOrMark(translator, messages.LabelOutputCsv),
		i18n.TranslateOrMark(translator, messages.LabelOutputCsvTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			_ = cw
			_ = rep
			appendOutputPath = path
		},
	)

	appendSaveButton := widget.NewMPushButton()
	appendSaveButton.SetLabel(i18n.TranslateOrMark(translator, messages.LabelSave))
	appendSaveButton.SetOnClicked(func(cw *controller.ControlWindow) {
		_ = cw
		if !appendSourceLoaded {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageBuildFailed), nil)
			logger.Error(i18n.TranslateOrMark(translator, messages.LabelAppendSourceCsvTip))
			return
		}
		if !appendTargetLoaded {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageBuildFailed), nil)
			logger.Error(i18n.TranslateOrMark(translator, messages.LabelAppendTargetCsvTip))
			return
		}
		if strings.TrimSpace(appendOutputPath) == "" {
			appendOutputPath = buildAppendCsvOutputPath(appendSourcePath)
			if appendOutputPicker != nil && strings.TrimSpace(appendOutputPath) != "" {
				appendOutputPicker.SetPath(appendOutputPath)
			}
		}
		if strings.TrimSpace(appendOutputPath) == "" {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageBuildFailed), nil)
			logger.Error(i18n.TranslateOrMark(translator, messages.LabelOutputCsvTip))
			return
		}

		if err := viewerUsecase.SaveAppendCsv(
			appendSourceRows,
			appendTargetRows,
			appendTable.Rows(),
			appendOutputPath,
			minteractor.SaveOptions{},
		); err != nil {
			logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageOutputFailed), err)
			return
		}

		controller.Beep()
		logger.Info("%s: %s", i18n.TranslateOrMark(translator, messages.MessageOutputDone), filepath.Base(appendOutputPath))
	})

	if mWidgets != nil {
		mWidgets.Widgets = append(
			mWidgets.Widgets,
			translateModelPicker,
			translateCsvPicker,
			translateOutputPicker,
			translateSaveButton,
			csvOutputModelPicker,
			csvOutputPicker,
			csvOutputSaveButton,
			appendSourcePicker,
			appendTargetPicker,
			appendOutputPicker,
			appendSaveButton,
		)
		mWidgets.SetOnLoaded(func() {
			if mWidgets == nil || mWidgets.Window() == nil {
				return
			}
			mWidgets.Window().SetOnEnabledInPlaying(func(playing bool) {
				for _, w := range mWidgets.Widgets {
					w.SetEnabledInPlaying(playing)
				}
			})
			if strings.TrimSpace(initialModelPath) != "" {
				translateModelPicker.SetPath(initialModelPath)
			}
		})
	}

	translateTabPage := declarative.TabPage{
		Title:    i18n.TranslateOrMark(translator, messages.LabelTranslateTab),
		AssignTo: &translateTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					translateModelPicker.Widgets(),
					translateCsvPicker.Widgets(),
					translateTable.Widgets(),
					translateOutputPicker.Widgets(),
					declarative.VSeparator{},
					translateSaveButton.Widgets(),
					declarative.VSpacer{},
				},
			},
		},
	}

	csvOutputTabPage := declarative.TabPage{
		Title:    i18n.TranslateOrMark(translator, messages.LabelCsvOutputTab),
		AssignTo: &csvOutputTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					csvOutputModelPicker.Widgets(),
					csvOutputPicker.Widgets(),
					declarative.VSeparator{},
					csvTable.Widgets(),
					declarative.VSeparator{},
					csvOutputSaveButton.Widgets(),
					declarative.VSpacer{},
				},
			},
		},
	}

	csvAppendTabPage := declarative.TabPage{
		Title:    i18n.TranslateOrMark(translator, messages.LabelCsvAppendTab),
		AssignTo: &csvAppendTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					appendSourcePicker.Widgets(),
					appendTargetPicker.Widgets(),
					appendOutputPicker.Widgets(),
					declarative.VSeparator{},
					appendTable.Widgets(),
					declarative.VSeparator{},
					appendSaveButton.Widgets(),
					declarative.VSpacer{},
				},
			},
		},
	}

	return []declarative.TabPage{translateTabPage, csvOutputTabPage, csvAppendTabPage}
}

// NewTabPage は先頭タブを返す。
func NewTabPage(mWidgets *controller.MWidgets, baseServices base.IBaseServices, initialModelPath string, audioPlayer audio_api.IAudioPlayer, viewerUsecase *minteractor.PmxTranslatorUsecase) declarative.TabPage {
	return NewTabPages(mWidgets, baseServices, initialModelPath, audioPlayer, viewerUsecase)[0]
}

// buildAppendCsvOutputPath はCSV追加出力パスの既定値を生成する。
func buildAppendCsvOutputPath(sourcePath string) string {
	if strings.TrimSpace(sourcePath) == "" {
		return ""
	}
	return mfile.CreateOutputPath(sourcePath, "")
}

// logErrorTitle はタイトル付きエラーを出力する。
func logErrorTitle(logger logging.ILogger, title string, err error) {
	if logger == nil {
		return
	}
	if titled, ok := logger.(interface {
		ErrorTitle(title string, err error, msg string, params ...any)
	}); ok {
		titled.ErrorTitle(title, err, "")
		return
	}
	if err == nil {
		logger.Error("%s", title)
		return
	}
	logger.Error("%s: %s", title, err.Error())
}
