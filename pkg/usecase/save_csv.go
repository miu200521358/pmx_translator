package usecase

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/domain/mcsv"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/pmx_translator/pkg/domain"
	"github.com/miu200521358/walk/pkg/walk"
)

func SaveCsv(csvState *domain.CsvState) error {
	if ok, err := mfile.ExistsFile(csvState.Model.Path()); !ok || err != nil {
		return fmt.Errorf("生成失敗メッセージ")
	}

	// No.でソート
	if err := csvState.NameModel.Sort(1, walk.SortAscending); err != nil {
		return fmt.Errorf(mi18n.T("Csv出力失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
	}

	records := make([][]string, 0)
	records = append(records, []string{
		mi18n.T("ファイル名"), mi18n.T("元名称"), mi18n.T("日本語名称"), mi18n.T("英語名称")})

	_, fileName, _ := mfile.SplitPath(csvState.Model.Path())

	for _, name := range csvState.NameModel.CheckedNames() {
		if !domain.ExistText(records, name) {
			records = append(records, []string{fileName, name, "", ""})
		}
	}

	data := mcsv.NewCsvModel(records)

	csvRep := repository.NewCsvRepository()
	if err := csvRep.Save(csvState.OutputPath, data, false); err != nil {
		return err
	}

	return nil
}
