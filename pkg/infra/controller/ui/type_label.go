//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"

	"github.com/miu200521358/pmx_translator/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/pmx_translator/pkg/domain"
)

// typeLabelByKey は名称種別キーを表示文言へ変換する。
func typeLabelByKey(translator i18n.II18n, typeKey string) string {
	messageKey := typeMessageKey(typeKey)
	if messageKey == "" {
		return typeKey
	}
	return i18n.TranslateOrMark(translator, messageKey)
}

// typeMessageKey は名称種別キーに対応するメッセージキーを返す。
func typeMessageKey(typeKey string) string {
	switch typeKey {
	case domain.NameTypePath:
		return messages.TypePath
	case domain.NameTypeModel:
		return messages.TypeModel
	case domain.NameTypeMaterial:
		return messages.TypeMaterial
	case domain.NameTypeTexture:
		return messages.TypeTexture
	case domain.NameTypeBone:
		return messages.TypeBone
	case domain.NameTypeMorph:
		return messages.TypeMorph
	case domain.NameTypeDisplaySlot:
		return messages.TypeDisplaySlot
	case domain.NameTypeRigidBody:
		return messages.TypeRigidBody
	case domain.NameTypeJoint:
		return messages.TypeJoint
	default:
		return ""
	}
}
