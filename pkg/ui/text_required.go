package ui

import (
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/walk/pkg/walk"
)

type textRequired struct {
	title string
}

func (tr textRequired) Create() (walk.Validator, error) {
	return &textRequiredValidator{title: tr.title}, nil
}

type textRequiredValidator struct {
	title string
}

func TextRequiredValidator(title string) walk.Validator {
	return &textRequiredValidator{title: title}
}

func (tv textRequiredValidator) Validate(v interface{}) error {
	if v == nil || v == "" {
		return walk.NewValidationError(
			mi18n.T("文字列未入力"),
			mi18n.T("文字列未入力テキスト", map[string]interface{}{"Title": tv.title}),
		)
	}

	return nil
}
