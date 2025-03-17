package domain

import (
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/mcsv"
)

func Translate(text, enText string, langDict *mcsv.CsvModel, modelName string) (string, string) {
	newJpText := text
	newEnText := enText

	// モデル名一致＆完全一致の翻訳を行う
	for n, row := range langDict.Records() {
		if n > 0 && row[0] == modelName && row[1] == newJpText {
			newJpText = strings.ReplaceAll(newJpText, row[1], row[2])
			if enText != "" && row[3] != "" {
				newEnText = row[3]
			}
		}
	}

	// モデル名不問＆完全一致の翻訳を行う
	for n, row := range langDict.Records() {
		if n > 0 && row[0] == "" && row[1] == newJpText {
			newJpText = strings.ReplaceAll(newJpText, row[1], row[2])
			if enText != "" && row[3] != "" {
				newEnText = row[3]
			}
		}
	}

	// モデル名一致＆部分一致翻訳を行う
	for n, row := range langDict.Records() {
		if n > 0 && row[0] == modelName && row[1] != "" {
			newJpText = strings.ReplaceAll(newJpText, row[1], row[2])
			if enText != "" && row[3] != "" {
				newEnText = strings.ReplaceAll(newEnText, row[1], row[3])
			}
		}
	}

	// モデル名不問＆部分一致翻訳
	for n, row := range langDict.Records() {
		if n > 0 && row[0] == "" && row[1] != "" {
			newJpText = strings.ReplaceAll(newJpText, row[1], row[2])
			if enText != "" && row[3] != "" {
				newEnText = strings.ReplaceAll(newEnText, row[1], row[3])
			}
		}
	}

	return newJpText, newEnText
}
