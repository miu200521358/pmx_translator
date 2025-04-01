package usecase

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/pmx_translator/pkg/domain"
)

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

func SavePmx(model *pmx.PmxModel, nameItems []*domain.NameItem, outputJpPath string) error {
	number := 2

	{
		jpName, enName := getTranslatedNames(number, model.Name(), model.EnglishName(), nameItems)
		model.SetName(jpName)
		model.SetEnglishName(enName)
		number++
	}

	model.Materials.ForEach(func(i int, mat *pmx.Material) bool {
		jpName, enName := getTranslatedNames(number, mat.Name(), mat.EnglishName(), nameItems)
		mat.SetName(jpName)
		mat.SetEnglishName(enName)
		number++
		return true
	})

	jpDir, _, _ := mfile.SplitPath(outputJpPath)

	model.Textures.ForEach(func(i int, tex *pmx.Texture) bool {
		if tex.Name() == "" {
			number++
			return true
		}

		orgName := tex.Name()
		jpPath, _ := getTranslatedNames(number, orgName, "", nameItems)
		tex.SetName(jpPath)
		number++

		dir, _, _ := mfile.SplitPath(model.Path())
		if !mfile.CanSave(outputJpPath) {
			if err := os.MkdirAll(jpDir, 0755); err != nil {
				mlog.E("ディレクトリ作成失敗: %s", err)
				return false
			}
		}

		orgTexPath := filepath.Join(dir, orgName)
		jpTexPath := filepath.Join(jpDir, jpPath)
		if orgTexPath != jpTexPath {
			copyTex(orgTexPath, jpTexPath)
		}

		return true
	})

	model.Bones.ForEach(func(i int, bone *pmx.Bone) bool {
		jpName, enName := getTranslatedNames(number, bone.Name(), bone.EnglishName(), nameItems)
		bone.SetName(jpName)
		bone.SetEnglishName(enName)
		number++
		return true
	})

	model.Morphs.ForEach(func(i int, morph *pmx.Morph) bool {
		jpName, enName := getTranslatedNames(number, morph.Name(), morph.EnglishName(), nameItems)
		morph.SetName(jpName)
		morph.SetEnglishName(enName)
		number++
		return true
	})

	model.DisplaySlots.ForEach(func(i int, disp *pmx.DisplaySlot) bool {
		jpName, enName := getTranslatedNames(number, disp.Name(), disp.EnglishName(), nameItems)
		disp.SetName(jpName)
		disp.SetEnglishName(enName)
		number++
		return true
	})

	model.RigidBodies.ForEach(func(i int, rb *pmx.RigidBody) bool {
		jpName, enName := getTranslatedNames(number, rb.Name(), rb.EnglishName(), nameItems)
		rb.SetName(jpName)
		rb.SetEnglishName(enName)
		number++
		return true
	})

	model.Joints.ForEach(func(i int, joint *pmx.Joint) bool {
		jpName, enName := getTranslatedNames(number, joint.Name(), joint.EnglishName(), nameItems)
		joint.SetName(jpName)
		joint.SetEnglishName(enName)
		number++
		return true
	})

	if !mfile.CanSave(outputJpPath) {
		jpDir, _ := filepath.Split(outputJpPath)
		if err := os.MkdirAll(jpDir, 0755); err != nil {
			mlog.E("ディレクトリ作成失敗: %s", err)
			return err
		}
	}

	pmxRep := repository.NewPmxRepository(true)
	if err := pmxRep.Save(outputJpPath, model, false); err != nil {
		return err
	}

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
