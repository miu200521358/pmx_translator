//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"

	"github.com/miu200521358/pmx_translator/pkg/adapter/mpresenter/messages"
)

// NewMenuItems は pmx_translator のメニュー項目を生成する。
func NewMenuItems(translator i18n.II18n, logger logging.ILogger) []declarative.MenuItem {
	return controller.BuildMenuItemsWithMessages(translator, logger, []controller.MenuMessageItem{
		{TitleKey: messages.HelpOverviewTitle, MessageKey: messages.HelpOverviewMessage},
		{TitleKey: messages.HelpToolAboutTitle, MessageKey: messages.HelpToolAbout},
		{TitleKey: controller.MenuSeparatorKey},

		{TitleKey: messages.LabelTranslateTab, MessageKey: messages.LabelTranslateTabTip},
		{TitleKey: messages.LabelOriginalModel, MessageKey: messages.HelpOriginalModel},
		{TitleKey: messages.LabelDictionaryCsv, MessageKey: messages.HelpDictionaryCsv},
		{TitleKey: messages.HelpTranslateTableTitle, MessageKey: messages.HelpTranslateTable},
		{TitleKey: messages.LabelOutputModel, MessageKey: messages.HelpOutputModel},
		{TitleKey: messages.HelpTranslateSaveTitle, MessageKey: messages.HelpTranslateSave},
		{TitleKey: controller.MenuSeparatorKey},

		{TitleKey: messages.LabelCsvOutputTab, MessageKey: messages.LabelCsvOutputTabTip},
		{TitleKey: messages.LabelOriginalModel, MessageKey: messages.HelpOriginalModel},
		{TitleKey: messages.HelpCsvOutputTableTitle, MessageKey: messages.HelpCsvOutputTable},
		{TitleKey: messages.LabelOutputCsv, MessageKey: messages.HelpOutputCsv},
		{TitleKey: messages.HelpCsvOutputSaveTitle, MessageKey: messages.HelpCsvOutputSave},
		{TitleKey: controller.MenuSeparatorKey},

		{TitleKey: messages.LabelCsvAppendTab, MessageKey: messages.LabelCsvAppendTabTip},
		{TitleKey: messages.LabelAppendSourceCsv, MessageKey: messages.HelpAppendSourceCsv},
		{TitleKey: messages.LabelAppendTargetCsv, MessageKey: messages.HelpAppendTargetCsv},
		{TitleKey: messages.HelpCsvAppendTableTitle, MessageKey: messages.HelpCsvAppendTable},
		{TitleKey: messages.HelpAppendOutputCsvTitle, MessageKey: messages.HelpAppendOutputCsv},
		{TitleKey: messages.HelpCsvAppendSaveTitle, MessageKey: messages.HelpCsvAppendSave},
		{TitleKey: controller.MenuSeparatorKey},

		{TitleKey: messages.HelpOpenButtonTitle, MessageKey: messages.HelpOpenButton},
		{TitleKey: messages.HelpHistoryButtonTitle, MessageKey: messages.HelpHistoryButton},
	})
}
