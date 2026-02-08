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
		{TitleKey: messages.LabelCsvOutputTab, MessageKey: messages.LabelCsvOutputTabTip},
		{TitleKey: messages.LabelCsvAppendTab, MessageKey: messages.LabelCsvAppendTabTip},
		{TitleKey: controller.MenuSeparatorKey},
		{TitleKey: messages.LabelOriginalModel, MessageKey: messages.HelpOriginalModel},
		{TitleKey: messages.LabelDictionaryCsv, MessageKey: messages.HelpDictionaryCsv},
		{TitleKey: messages.LabelOutputModel, MessageKey: messages.HelpOutputModel},
		{TitleKey: messages.LabelOutputCsv, MessageKey: messages.HelpOutputCsv},
		{TitleKey: messages.LabelAppendSourceCsv, MessageKey: messages.HelpAppendSourceCsv},
		{TitleKey: messages.LabelAppendTargetCsv, MessageKey: messages.HelpAppendTargetCsv},
	})
}
