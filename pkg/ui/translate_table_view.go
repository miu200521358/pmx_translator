package ui

import (
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/pmx_translator/pkg/domain"
	"github.com/miu200521358/pmx_translator/pkg/usecase"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type TranslateTableView struct {
	*declarative.TableView
	NameModel *TranslateNameModel
}

type TranslateNameModel struct {
	walk.TableModelBase
	walk.SorterBase
	tv         *walk.TableView
	sortColumn int
	sortOrder  walk.SortOrder
	Records    []*domain.NameItem
}

func NewTranslateNameModel(model *pmx.PmxModel, charaCsv *core.CsvModel) *TranslateNameModel {
	m := new(TranslateNameModel)
	m.ResetRows(model, charaCsv)
	return m
}

func (m *TranslateNameModel) RowCount() int {
	return len(m.Records)
}

func (m *TranslateNameModel) SetParent(parent *walk.TableView) {
	m.tv = parent
}

func (m *TranslateNameModel) Value(row, col int) interface{} {
	item := m.Records[row]

	switch col {
	case 0:
		return item.Checked
	case 1:
		return item.Number
	case 2:
		return item.TypeText
	case 3:
		return item.Index
	case 4:
		return item.NameText
	case 5:
		return item.JapaneseNameText
	case 6:
		return item.EnglishNameText
	}

	panic("unexpected col")
}

func (m *TranslateNameModel) Checked(row int) bool {
	return m.Records[row].Checked
}

func (m *TranslateNameModel) SetChecked(row int, checked bool) error {
	m.Records[row].Checked = checked

	return nil
}

func (m *TranslateNameModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.Records, func(i, j int) bool {
		a, b := m.Records[i], m.Records[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			av := 0
			if a.Checked {
				av = 1
			}
			bv := 0
			if b.Checked {
				bv = 1
			}
			return c(av < bv)
		case 1:
			return c(a.Number < b.Number)
		case 2:
			return c(a.TypeText < b.TypeText)
		case 3:
			return c(a.Index < b.Index)
		case 4:
			return c(a.NameText < b.NameText)
		case 5:
			return c(a.JapaneseNameText < b.JapaneseNameText)
		case 6:
			return c(a.EnglishNameText < b.EnglishNameText)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *TranslateNameModel) AddRecord(
	fileName string, index int, jpTxt, enTxt, fieldKey string, charaCsv *core.CsvModel,
) {
	jpTransTxt, enTransTxt := usecase.Translate(jpTxt, enTxt, charaCsv, fileName)

	item := &domain.NameItem{
		Checked:          jpTxt != jpTransTxt || enTxt != enTransTxt,
		Number:           len(m.Records) + 1,
		TypeText:         mi18n.T(fieldKey),
		Index:            index,
		NameText:         jpTxt,
		JapaneseNameText: jpTransTxt,
		EnglishNameText:  enTransTxt,
	}
	m.Records = append(m.Records, item)

}

func (m *TranslateNameModel) ResetRows(model *pmx.PmxModel, charaCsv *core.CsvModel) {
	m.Records = make([]*domain.NameItem, 0)

	m.PublishRowsReset()

	if model == nil || charaCsv == nil {
		return
	}

	_, fileName, _ := mutils.SplitPath(model.Path())

	m.AddRecord(fileName, 0, model.Path(), "", "パス", charaCsv)
	m.AddRecord(fileName, 0, model.Name(), model.EnglishName(), "モデル", charaCsv)

	for _, mat := range model.Materials.Data {
		m.AddRecord(fileName, mat.Index(), mat.Name(), mat.EnglishName(), "材質", charaCsv)
	}

	for _, tex := range model.Textures.Data {
		m.AddRecord(fileName, tex.Index(), tex.Name(), "", "テクスチャ", charaCsv)
	}

	for _, bone := range model.Bones.Data {
		m.AddRecord(fileName, bone.Index(), bone.Name(), bone.EnglishName(), "ボーン", charaCsv)
	}

	for _, morph := range model.Morphs.Data {
		m.AddRecord(fileName, morph.Index(), morph.Name(), morph.EnglishName(), "モーフ", charaCsv)
	}

	for _, disp := range model.DisplaySlots.Data {
		m.AddRecord(fileName, disp.Index(), disp.Name(), disp.EnglishName(), "表示枠", charaCsv)
	}

	for _, rb := range model.RigidBodies.Data {
		m.AddRecord(fileName, rb.Index(), rb.Name(), rb.EnglishName(), "剛体", charaCsv)
	}

	for _, joint := range model.Joints.Data {
		m.AddRecord(fileName, joint.Index(), joint.Name(), joint.EnglishName(), "ジョイント", charaCsv)
	}

	m.PublishRowsReset()
}

func NewTranslateTableView(parent walk.Container, model *pmx.PmxModel, charaCsv *core.CsvModel) *TranslateTableView {
	nameModel := NewTranslateNameModel(model, charaCsv)

	var tv *walk.TableView
	builder := declarative.NewBuilder(parent)

	dTableView := &declarative.TableView{
		AssignTo:         &tv,
		AlternatingRowBG: true,
		CheckBoxes:       true,
		ColumnsOrderable: true,
		MultiSelection:   true,
		Model:            nameModel,
		Columns: []declarative.TableViewColumn{
			{Title: "#", Width: 50},
			{Title: "No.", Width: 50},
			{Title: mi18n.T("種類"), Width: 80},
			{Title: mi18n.T("インデックス"), Width: 40},
			{Title: mi18n.T("元名称"), Width: 200},
			{Title: mi18n.T("日本語名称"), Width: 200},
			{Title: mi18n.T("英語名称"), Width: 200},
		},
		StyleCell: func(style *walk.CellStyle) {
			if nameModel.Checked(style.Row()) {
				style.BackgroundColor = walk.RGB(159, 255, 243)
			} else {
				style.BackgroundColor = walk.RGB(255, 255, 255)
			}
		},
		OnSelectedIndexesChanged: func() {
			var dlg *walk.Dialog
			var cancelBtn *walk.PushButton
			var okBtn *walk.PushButton
			var db *walk.DataBinder
			var jpTxt *walk.TextEdit
			var enTxt *walk.TextEdit

			textChangeDialog := &declarative.Dialog{
				AssignTo:      &dlg,
				CancelButton:  &cancelBtn,
				DefaultButton: &okBtn,
				Title:         mi18n.T("名称変更"),
				Layout:        declarative.VBox{},
				MinSize:       declarative.Size{Width: 400, Height: 120},
				DataBinder: declarative.DataBinder{
					AssignTo:   &db,
					DataSource: nameModel.Records[tv.CurrentIndex()],
				},
				Children: []declarative.Widget{
					declarative.Composite{
						Layout: declarative.Grid{Columns: 2},
						Children: []declarative.Widget{
							declarative.Label{
								Text: mi18n.T("種類"),
							},
							declarative.Label{
								Text: declarative.Bind("TypeText"),
							},
							declarative.Label{
								Text: mi18n.T("元名称"),
							},
							declarative.Label{
								Text: declarative.Bind("NameText"),
							},
							declarative.Label{
								Text: mi18n.T("日本語名称"),
							},
							declarative.TextEdit{
								AssignTo: &jpTxt,
								Text:     declarative.Bind("JapaneseNameText", textRequired{title: mi18n.T("日本語名称")}),
							},
							declarative.Label{
								Text: mi18n.T("英語名称"),
							},
							declarative.TextEdit{
								AssignTo: &enTxt,
								Text:     declarative.Bind("EnglishNameText"),
							},
						},
					},
					declarative.Composite{
						Layout: declarative.HBox{
							Alignment: declarative.AlignHFarVCenter,
						},
						Children: []declarative.Widget{
							declarative.PushButton{
								AssignTo: &okBtn,
								Text:     mi18n.T("OK"),
								OnClicked: func() {
									if err := db.Submit(); err != nil {
										widget.RaiseError(err)
										return
									}
									dlg.Accept()
								},
							},
							declarative.PushButton{
								AssignTo: &cancelBtn,
								Text:     mi18n.T("キャンセル"),
								OnClicked: func() {
									dlg.Cancel()
								},
							},
						},
					},
				},
			}

			if cmd, err := textChangeDialog.Run(builder.Parent().Form()); err != nil {
				widget.RaiseError(err)
			} else if cmd == walk.DlgCmdOK {
				nameModel.Records[tv.CurrentIndex()].Checked = true
				// nameModel.Records[tv.CurrentIndex()].JapaneseNameText = jpTxt.Text()
				// nameModel.Records[tv.CurrentIndex()].EnglishNameText = enTxt.Text()
			}
		},
	}

	if err := dTableView.Create(builder); err != nil {
		widget.RaiseError(err)
	}

	nameModel.SetParent(tv)
	nameTableView := &TranslateTableView{
		TableView: dTableView,
		NameModel: nameModel,
	}

	return nameTableView
}

func (tv *TranslateTableView) ResetModel(model *pmx.PmxModel, charaCsv *core.CsvModel) {
	tv.NameModel.ResetRows(model, charaCsv)
}

// -----------------------

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
