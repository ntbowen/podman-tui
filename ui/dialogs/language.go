package dialogs

import (
	"fmt"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	langWidthOffset = 6
)

const (
	langTableFocus = 0 + iota
	langFormFocus
)

// LanguageDialog is a language selection dialog.
type LanguageDialog struct {
	*tview.Box

	layout        *tview.Flex
	table         *tview.Table
	form          *tview.Form
	display       bool
	languages     []string
	languageNames []string
	width         int
	height        int
	focusElement  int
	selectedStyle tcell.Style
	cancelHandler func()
	selectHandler func()
}

// NewLanguageDialog returns a language selection dialog primitive.
func NewLanguageDialog() *LanguageDialog {
	// Get supported languages
	languages := i18n.GetSupportedLanguages()
	languageNames := make([]string, len(languages))
	
	// Initialize with language names
	for i, code := range languages {
		languageNames[i] = i18n.GetLanguageName(code)
	}

	form := tview.NewForm().
		AddButton(i18n.T("Cancel"), nil).
		SetButtonsAlign(tview.AlignRight)

	form.SetBackgroundColor(style.DialogBgColor)
	form.SetButtonBackgroundColor(style.ButtonBgColor)

	langTable := tview.NewTable()
	langTable.SetBackgroundColor(style.DialogBgColor)

	// Language table header
	langTable.SetCell(0, 0,
		tview.NewTableCell(fmt.Sprintf("[%s::b]%s", style.GetColorHex(style.TableHeaderFgColor), i18n.T("LANGUAGE CODE"))).
			SetExpansion(1).
			SetBackgroundColor(style.TableHeaderBgColor).
			SetTextColor(style.TableHeaderFgColor).
			SetAlign(tview.AlignLeft).
			SetSelectable(false))

	langTable.SetCell(0, 1,
		tview.NewTableCell(fmt.Sprintf("[%s::b]%s", style.GetColorHex(style.TableHeaderFgColor), i18n.T("LANGUAGE NAME"))).
			SetExpansion(1).
			SetBackgroundColor(style.TableHeaderBgColor).
			SetTextColor(style.TableHeaderFgColor).
			SetAlign(tview.AlignLeft).
			SetSelectable(false))

	// Calculate max widths using widthfix for proper display
	col1Width := i18n.GetDisplayWidth(i18n.T("LANGUAGE CODE"))
	col2Width := i18n.GetDisplayWidth(i18n.T("LANGUAGE NAME"))

	for i := range languages {
		// Add language code
		langTable.SetCell(i+1, 0,
			tview.NewTableCell(languages[i]).
				SetAlign(tview.AlignLeft).
				SetSelectable(true).
				SetTextColor(style.DialogFgColor))
		
		// Add language name (will be updated with mark in Display())
		langTable.SetCell(i+1, 1,
			tview.NewTableCell(languageNames[i]).
				SetAlign(tview.AlignLeft).
				SetSelectable(true).
				SetTextColor(style.DialogFgColor))

		codeWidth := i18n.GetDisplayWidth(languages[i])
		if codeWidth > col1Width {
			col1Width = codeWidth
		}

		// Reserve space for the mark (✅ + space)
		nameWidth := i18n.GetDisplayWidth(languageNames[i]) + 4
		if nameWidth > col2Width {
			col2Width = nameWidth
		}
	}

	langWidth := col1Width + col2Width + 6 //nolint:mnd

	langTable.SetFixed(1, 1)
	langTable.SetSelectable(true, false)
	langTable.SetBackgroundColor(style.DialogBgColor)

	langLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	langLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	langLayout.AddItem(langTable, 0, 1, true)
	langLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	// Layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(langLayout, 0, 1, true)
	layout.AddItem(form, DialogFormHeight, 0, true)
	layout.SetBorder(true)
	layout.SetTitle(i18n.T("[ SELECT LANGUAGE ]"))
	layout.SetBorderColor(style.DialogBorderColor)
	layout.SetBackgroundColor(style.DialogBgColor)

	return &LanguageDialog{
		Box:          tview.NewBox().SetBorder(false),
		layout:       layout,
		table:        langTable,
		form:         form,
		display:      false,
		focusElement: langTableFocus,
		selectedStyle: tcell.StyleDefault.
			Background(style.DialogFgColor).
			Foreground(style.DialogBgColor),
		languages:     languages,
		languageNames: languageNames,
		width:         langWidth + langWidthOffset,
		height:        len(languages) + TableHeightOffset + DialogFormHeight,
	}
}

// GetSelectedLanguage returns selected language code.
func (lang *LanguageDialog) GetSelectedLanguage() string {
	row, _ := lang.table.GetSelection()
	if row > 0 && row <= len(lang.languages) {
		return lang.languages[row-1]
	}
	return ""
}

// Display displays this primitive.
func (lang *LanguageDialog) Display() {
	// Update language names with current marks BEFORE showing
	lang.updateLanguageMarks()
	
	// Select current language
	currentLang := i18n.GetCurrentLanguage()
	selectedRow := 1
	
	for i, code := range lang.languages {
		if code == currentLang {
			selectedRow = i + 1
			break
		}
	}
	
	lang.table.Select(selectedRow, 0)
	lang.form.SetFocus(1)
	lang.display = true
}

// updateLanguageMarks updates the language names in the table with current marks
func (lang *LanguageDialog) updateLanguageMarks() {
	for i, code := range lang.languages {
		langName := i18n.GetLanguageNameWithMark(code)
		lang.languageNames[i] = langName
		lang.table.GetCell(i+1, 1).SetText(langName)
	}
}

// IsDisplay returns true if primitive is shown.
func (lang *LanguageDialog) IsDisplay() bool {
	return lang.display
}

// Hide stops displaying this primitive.
func (lang *LanguageDialog) Hide() {
	lang.display = false
	lang.focusElement = langTableFocus
	lang.table.SetSelectedStyle(lang.selectedStyle)
}

// HasFocus returns whether or not this primitive has focus.
func (lang *LanguageDialog) HasFocus() bool {
	return lang.table.HasFocus() || lang.form.HasFocus()
}

// Focus is called when this primitive receives focus.
func (lang *LanguageDialog) Focus(delegate func(p tview.Primitive)) {
	if lang.focusElement == langTableFocus {
		delegate(lang.table)
		return
	}

	button := lang.form.GetButton(lang.form.GetButtonCount() - 1)
	button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == utils.SwitchFocusKey.Key {
			lang.focusElement = langTableFocus
			lang.Focus(delegate)
			lang.form.SetFocus(0)
			return nil
		}
		return event
	})

	delegate(lang.form)
}

// InputHandler returns input handler function for this primitive.
func (lang *LanguageDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return lang.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("language dialog: event %v received", event)

		if event.Key() == utils.CloseDialogKey.Key {
			lang.cancelHandler()
			return
		}

		if event.Key() == utils.SwitchFocusKey.Key {
			lang.setFocusElement()
		}

		if lang.form.HasFocus() {
			if formHandler := lang.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)
				return
			}
		}

		// Language table handler
		if lang.table.HasFocus() {
			if event.Key() == tcell.KeyEnter {
				lang.selectHandler()
				return
			}

			if tableHandler := lang.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
				return
			}
		}
	})
}

// SetSelectedFunc sets form enter button selected function.
func (lang *LanguageDialog) SetSelectedFunc(handler func()) *LanguageDialog {
	lang.selectHandler = handler
	return lang
}

// SetCancelFunc sets form cancel button selected function.
func (lang *LanguageDialog) SetCancelFunc(handler func()) *LanguageDialog {
	lang.cancelHandler = handler
	cancelButton := lang.form.GetButton(lang.form.GetButtonCount() - 1)
	cancelButton.SetSelectedFunc(handler)
	return lang
}

// SetRect set rects for this primitive.
func (lang *LanguageDialog) SetRect(x, y, width, height int) {
	ws := (width - lang.width) / 2     //nolint:mnd
	hs := ((height - lang.height) / 2) //nolint:mnd
	dy := y + hs
	bWidth := lang.width

	if lang.width > width {
		ws = 0
		bWidth = width - 1
	}

	bHeight := lang.height

	if lang.height >= height {
		dy = y + 1
		bHeight = height - 1
	}

	lang.Box.SetRect(x+ws, dy, bWidth, bHeight)

	x, y, width, height = lang.GetInnerRect()
	lang.layout.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (lang *LanguageDialog) Draw(screen tcell.Screen) {
	if !lang.display {
		return
	}

	lang.DrawForSubclass(screen, lang)
	lang.layout.Draw(screen)
}

func (lang *LanguageDialog) setFocusElement() {
	if lang.focusElement == langTableFocus {
		lang.focusElement = langFormFocus
		lang.table.SetSelectedStyle(tcell.StyleDefault.
			Background(style.DialogBgColor).
			Foreground(style.DialogFgColor))
	} else {
		lang.focusElement = langTableFocus
		lang.table.SetSelectedStyle(lang.selectedStyle)
	}
}
