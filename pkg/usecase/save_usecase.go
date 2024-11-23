package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/pmx_translator/pkg/domain"
)

func Translate(text, enText string, langDict *core.CsvModel, modelName string) (string, string) {
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

func getTranslatedNames(
	number int, jpText, enText string, nameItems []*domain.NameItem,
) (string, string) {
	var item *domain.NameItem
	for _, it := range nameItems {
		if it.Number == number {
			item = it
			break
		}
	}
	if item.Checked {
		return item.JapaneseNameText, item.EnglishNameText
	}
	return jpText, enText
}

func Save(model *pmx.PmxModel, nameItems []*domain.NameItem, outputJpPath string) error {
	number := 2

	{
		jpName, enName := getTranslatedNames(number, model.Name(), model.EnglishName(), nameItems)
		model.SetName(jpName)
		model.SetEnglishName(enName)
		number++
	}

	for _, mat := range model.Materials.Data {
		jpName, enName := getTranslatedNames(number, mat.Name(), mat.EnglishName(), nameItems)
		mat.SetName(jpName)
		mat.SetEnglishName(enName)
		number++
	}

	jpDir, _, _ := mutils.SplitPath(outputJpPath)

	for _, tex := range model.Textures.Data {
		if tex.Name() == "" {
			number++
			continue
		}

		orgName := tex.Name()
		jpPath, _ := getTranslatedNames(number, orgName, "", nameItems)
		tex.SetName(jpPath)
		number++

		dir, _, _ := mutils.SplitPath(model.Path())
		if !mutils.CanSave(outputJpPath) {
			if err := os.MkdirAll(jpDir, 0755); err != nil {
				mlog.E("ディレクトリ作成失敗: %s", err)
				return err
			}
		}

		orgTexPath := filepath.Join(dir, orgName)
		jpTexPath := filepath.Join(jpDir, jpPath)
		if orgTexPath != jpTexPath {
			copyTex(orgTexPath, jpTexPath)
		}
	}

	for _, bone := range model.Bones.Data {
		jpName, enName := getTranslatedNames(number, bone.Name(), bone.EnglishName(), nameItems)
		bone.SetName(jpName)
		bone.SetEnglishName(enName)
		number++
	}

	for _, morph := range model.Morphs.Data {
		jpName, enName := getTranslatedNames(number, morph.Name(), morph.EnglishName(), nameItems)
		morph.SetName(jpName)
		morph.SetEnglishName(enName)
		number++
	}

	for _, disp := range model.DisplaySlots.Data {
		jpName, enName := getTranslatedNames(number, disp.Name(), disp.EnglishName(), nameItems)
		disp.SetName(jpName)
		disp.SetEnglishName(enName)
		number++
	}

	for _, rb := range model.RigidBodies.Data {
		jpName, enName := getTranslatedNames(number, rb.Name(), rb.EnglishName(), nameItems)
		rb.SetName(jpName)
		rb.SetEnglishName(enName)
		number++
	}

	for _, joint := range model.Joints.Data {
		jpName, enName := getTranslatedNames(number, joint.Name(), joint.EnglishName(), nameItems)
		joint.SetName(jpName)
		joint.SetEnglishName(enName)
		number++
	}

	if !mutils.CanSave(outputJpPath) {
		jpDir, _ := filepath.Split(outputJpPath)
		if err := os.MkdirAll(jpDir, 0755); err != nil {
			mlog.E("ディレクトリ作成失敗: %s", err)
			return err
		}
	}

	pmxRep := repository.NewPmxRepository()
	if err := pmxRep.Save(outputJpPath, model, false); err != nil {
		return err
	}

	mlog.IT(mi18n.T("出力成功"), mi18n.T("出力成功メッセージ", map[string]interface{}{"Path": outputJpPath}))

	return nil
}

func copyTex(texPath string, copyTexPath string) error {
	texFile, err := os.ReadFile(texPath)
	if err != nil {
		mlog.E(fmt.Sprintf("Failed to read original pmx tex file: %s", texPath), err)
		return err
	}

	err = os.MkdirAll(filepath.Dir(copyTexPath), 0755)
	if err != nil {
		mlog.E(fmt.Sprintf("Failed to create original pmx tex tmp directory: %s", copyTexPath), err)
		return err
	}

	err = os.WriteFile(copyTexPath, texFile, 0644)
	if err != nil {
		mlog.E(fmt.Sprintf("Failed to write original pmx tex tmp file: %s", copyTexPath), err)
		return err
	}

	return nil
}
