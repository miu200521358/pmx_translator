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

func Translate(text string, langDict *core.CsvModel, langIndex int, modelName string) string {
	// 先にモデル名一致の翻訳を行う
	for _, row := range langDict.Records() {
		if row[0] == modelName && row[1] != "" {
			text = strings.ReplaceAll(text, row[1], row[langIndex])
		}
	}

	// モデル名を問わない翻訳
	for _, row := range langDict.Records() {
		if row[0] == "" && row[1] != "" {
			text = strings.ReplaceAll(text, row[1], row[langIndex])
		}
	}

	return text
}

func getTranslatedNames(
	jpText, enText string, nameItems []*domain.NameItem,
) (string, string) {
	newJpText := jpText
	newEnText := enText

	for _, item := range nameItems {
		if item.Checked {
			newJpText = strings.ReplaceAll(newJpText, item.NameText, item.JapaneseNameText)
			newEnText = strings.ReplaceAll(newEnText, item.NameText, item.EnglishNameText)
		}
	}
	return newJpText, newEnText
}

func TranslateOutputPath(model *pmx.PmxModel, nameItems []*domain.NameItem) string {
	{
		jpName, enName := getTranslatedNames(model.Name(), model.EnglishName(), nameItems)
		model.SetName(jpName)
		model.SetEnglishName(enName)
	}

	outputJpPath := model.Path()
	{
		path, fileName, ext := mutils.SplitPath(outputJpPath)
		jpFileName, _ := getTranslatedNames(fileName, "", nameItems)

		paths := strings.Split(path, string(filepath.Separator))
		for i, p := range paths {
			if i < 2 {
				paths[i] = p
				continue
			}
			if p == "" {
				continue
			}
			paths[i], _ = getTranslatedNames(p, "", nameItems)
		}

		outputJpPath = fmt.Sprintf(
			"%s%s%s", paths[0], string(filepath.Separator), filepath.Join(append(paths[1:], jpFileName+ext)...))
	}

	return mutils.CreateOutputPath(outputJpPath, "")
}

func Save(model *pmx.PmxModel, nameItems []*domain.NameItem, outputJpPath string) error {
	{
		jpName, enName := getTranslatedNames(model.Name(), model.EnglishName(), nameItems)
		model.SetName(jpName)
		model.SetEnglishName(enName)
	}

	for _, mat := range model.Materials.Data {
		jpName, enName := getTranslatedNames(mat.Name(), mat.EnglishName(), nameItems)
		mat.SetName(jpName)
		mat.SetEnglishName(enName)
	}

	jpDir, _, _ := mutils.SplitPath(outputJpPath)

	for _, tex := range model.Textures.Data {
		if tex.Name() == "" {
			continue
		}

		orgName := tex.Name()
		jpPath, _ := getTranslatedNames(orgName, "", nameItems)
		tex.SetName(jpPath)

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
		jpName, enName := getTranslatedNames(bone.Name(), bone.EnglishName(), nameItems)
		bone.SetName(jpName)
		bone.SetEnglishName(enName)
	}

	for _, morph := range model.Morphs.Data {
		jpName, enName := getTranslatedNames(morph.Name(), morph.EnglishName(), nameItems)
		morph.SetName(jpName)
		morph.SetEnglishName(enName)
	}

	for _, disp := range model.DisplaySlots.Data {
		jpName, enName := getTranslatedNames(disp.Name(), disp.EnglishName(), nameItems)
		disp.SetName(jpName)
		disp.SetEnglishName(enName)
	}

	for _, rb := range model.RigidBodies.Data {
		jpName, enName := getTranslatedNames(rb.Name(), rb.EnglishName(), nameItems)
		rb.SetName(jpName)
		rb.SetEnglishName(enName)
	}

	for _, joint := range model.Joints.Data {
		jpName, enName := getTranslatedNames(joint.Name(), joint.EnglishName(), nameItems)
		joint.SetName(jpName)
		joint.SetEnglishName(enName)
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
