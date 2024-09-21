package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

func translate(text string, langDict *core.CsvModel, langIndex int, modelName string) string {
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

func Save(model *pmx.PmxModel, langDict *core.CsvModel, outputPath string) error {
	modelName := model.Name()
	model.SetName(translate(modelName, langDict, 2, modelName))
	model.SetEnglishName(translate(modelName, langDict, 3, modelName))

	for _, mat := range model.Materials.Data {
		chTxt := mat.Name()
		mat.SetName(translate(chTxt, langDict, 2, modelName))
		mat.SetEnglishName(translate(chTxt, langDict, 3, modelName))
	}

	for _, bone := range model.Bones.Data {
		chTxt := bone.Name()
		bone.SetName(translate(chTxt, langDict, 2, modelName))
		bone.SetEnglishName(translate(chTxt, langDict, 3, modelName))
	}

	for _, morph := range model.Morphs.Data {
		chTxt := morph.Name()
		morph.SetName(translate(chTxt, langDict, 2, modelName))
		morph.SetEnglishName(translate(chTxt, langDict, 3, modelName))
	}

	for _, disp := range model.DisplaySlots.Data {
		chTxt := disp.Name()
		disp.SetName(translate(chTxt, langDict, 2, modelName))
		disp.SetEnglishName(translate(chTxt, langDict, 3, modelName))
	}

	for _, rb := range model.RigidBodies.Data {
		chTxt := rb.Name()
		rb.SetName(translate(chTxt, langDict, 2, modelName))
		rb.SetEnglishName(translate(chTxt, langDict, 3, modelName))
	}

	for _, joint := range model.Joints.Data {
		chTxt := joint.Name()
		joint.SetName(translate(chTxt, langDict, 2, modelName))
		joint.SetEnglishName(translate(chTxt, langDict, 3, modelName))
	}

	outputJpPath := translate(outputPath, langDict, 2, modelName)

	if outputJpPath != outputPath {
		chDir, _ := filepath.Split(outputPath)

		jpDir, _ := filepath.Split(outputJpPath)
		if err := os.MkdirAll(jpDir, 0755); err != nil {
			mlog.E("ディレクトリ作成失敗: %s", err)
			return err
		}

		for _, tex := range model.Textures.Data {
			chPath := tex.Name()
			jpPath := translate(chPath, langDict, 2, modelName)
			tex.SetName(jpPath)

			copyTex(filepath.Join(chDir, chPath), filepath.Join(jpDir, jpPath))
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

	// 仮パスのフォルダ構成を作成する
	err = os.MkdirAll(filepath.Dir(copyTexPath), 0755)
	if err != nil {
		mlog.E(fmt.Sprintf("Failed to create original pmx tex tmp directory: %s", copyTexPath), err)
		return err
	}

	// 作業フォルダにファイルを書き込む
	err = os.WriteFile(copyTexPath, texFile, 0644)
	if err != nil {
		mlog.E(fmt.Sprintf("Failed to write original pmx tex tmp file: %s", copyTexPath), err)
		return err
	}

	return nil
}
