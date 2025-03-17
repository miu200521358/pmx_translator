package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/domain/mcsv"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/pmx_translator/pkg/domain"
)

func MergeCsv(mergeState *domain.MergeState) error {
	records := make([][]string, 0)
	records = append(records, []string{
		mi18n.T("ファイル名"), mi18n.T("元名称"), mi18n.T("日本語名称"), mi18n.T("英語名称")})

	number := 1
	for n, record := range mergeState.OriginalCsvModel.Records() {
		if n == 0 {
			continue
		}

		for _, item := range mergeState.NameModel.CheckedItems() {
			if item.Checked && number == item.Number && item.NameText == record[1] {
				records = append(records, record)
			}
		}

		number++
	}

	for n, record := range mergeState.MergedCsvModel.Records() {
		if n == 0 {
			continue
		}

		for _, item := range mergeState.NameModel.CheckedItems() {
			if item.Checked && number == item.Number && item.NameText == record[1] {
				records = append(records, record)
			}
		}

		number++
	}

	data := mcsv.NewCsvModel(records)

	csvRep := repository.NewCsvRepository()
	if err := csvRep.Save(mergeState.OutputPath, data, false); err != nil {
		return err
	}

	return nil
}
