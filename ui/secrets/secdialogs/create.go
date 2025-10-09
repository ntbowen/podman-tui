package secdialogs

import (
	"strings"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/pdcs/secrets"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	secretCreateDialogMaxWidth  = 80
	secretCreateDialogMaxHeight = 15
	labelWidth                  = 13
)

const (
	secretNameFocus = 0 + iota
	secretReplaceFocus
	secretFileFocus
	secretTextFocus
	secretLabelsFocus
	secretDriverFocus
	secretDriverOptionsFocus
	secretFormFocus
)

// SecretCreateDialog implements secret create dialog.
type SecretCreateDialog struct {
	*tview.Box

	layout              *tview.Flex
	form                *tview.Form
	secretName          *tview.InputField
	secretFile          *tview.InputField
	secretText          *tview.InputField
	secretLabels        *tview.InputField
	secretDriver        *tview.DropDown
	secretDriverOptions *tview.InputField
	secretReplace       *tview.Checkbox
	display             bool
	focusElement        int
	createHandler       func()
	cancelHandler       func()
}

// NewSecretCreateDialog returns new secret create dialog primitive.
func NewSecretCreateDialog() *SecretCreateDialog {
	createDialog := &SecretCreateDialog{
		Box:                 tview.NewBox(),
		layout:              tview.NewFlex().SetDirection(tview.FlexColumn),
		form:                tview.NewForm(),
		secretName:          tview.NewInputField(),
		secretFile:          tview.NewInputField(),
		secretText:          tview.NewInputField(),
		secretLabels:        tview.NewInputField(),
		secretDriver:        tview.NewDropDown(),
		secretDriverOptions: tview.NewInputField(),
		secretReplace:       tview.NewCheckbox(),
		display:             false,
		focusElement:        secretNameFocus,
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected

	// secret name field
	createDialog.secretName.SetBackgroundColor(bgColor)
	createDialog.secretName.SetLabel(utils.StringToInputLabel(i18n.T("name:"), labelWidth))
	createDialog.secretName.SetFieldStyle(style.InputFieldStyle)
	createDialog.secretName.SetLabelStyle(style.InputLabelStyle)

	// secret file field
	createDialog.secretFile.SetBackgroundColor(bgColor)
	createDialog.secretFile.SetLabel(utils.StringToInputLabel(i18n.T("secret file:"), labelWidth))
	createDialog.secretFile.SetFieldStyle(style.InputFieldStyle)
	createDialog.secretFile.SetLabelStyle(style.InputLabelStyle)

	// secret text field
	createDialog.secretText.SetBackgroundColor(bgColor)
	createDialog.secretText.SetLabel(utils.StringToInputLabel(i18n.T("secret text:"), labelWidth))
	createDialog.secretText.SetFieldStyle(style.InputFieldStyle)
	createDialog.secretText.SetLabelStyle(style.InputLabelStyle)

	// secret labels field
	createDialog.secretLabels.SetBackgroundColor(bgColor)
	createDialog.secretLabels.SetLabel(utils.StringToInputLabel(i18n.T("labels:"), labelWidth))
	createDialog.secretLabels.SetFieldStyle(style.InputFieldStyle)
	createDialog.secretLabels.SetLabelStyle(style.InputLabelStyle)

	// secret replace
	replaceLabel := i18n.T("replace") + " "
	createDialog.secretReplace.SetLabel(replaceLabel)
	createDialog.secretReplace.SetChecked(false)
	createDialog.secretReplace.SetBackgroundColor(bgColor)
	createDialog.secretReplace.SetLabelColor(fgColor)
	createDialog.secretReplace.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// secret driver
	createDialog.secretDriver.SetBackgroundColor(bgColor)
	createDialog.secretDriver.SetLabelColor(fgColor)
	createDialog.secretDriver.SetLabel(i18n.T("driver:"))
	createDialog.secretDriver.SetLabelWidth(labelWidth)
	createDialog.secretDriver.SetOptions([]string{i18n.T("file"), i18n.T("pass"), i18n.T("shell")}, nil)
	createDialog.secretDriver.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	createDialog.secretDriver.SetFocusedStyle(style.DropDownFocused)
	createDialog.secretDriver.SetCurrentOption(0)
	createDialog.secretDriver.SetFieldStyle(style.InputFieldStyle)

	// secret driver options field
	createDialog.secretDriverOptions.SetBackgroundColor(bgColor)
	createDialog.secretDriverOptions.SetLabel(i18n.T("driver options:") + " ")
	createDialog.secretDriverOptions.SetFieldStyle(style.InputFieldStyle)
	createDialog.secretDriverOptions.SetLabelStyle(style.InputLabelStyle)

	// form
	createDialog.form.AddButton(i18n.T("Cancel"), nil)
	createDialog.form.AddButton(i18n.T("Create"), nil)
	createDialog.form.SetButtonsAlign(tview.AlignRight)
	createDialog.form.SetBackgroundColor(bgColor)
	createDialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	optionsLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	nameAndReplaceRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	nameAndReplaceRow.SetBackgroundColor(bgColor)
	nameAndReplaceRow.AddItem(createDialog.secretName, 0, 1, true)
	nameAndReplaceRow.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	nameAndReplaceRow.AddItem(createDialog.secretReplace, len(replaceLabel)+1, 0, true)

	driverRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	driverRow.SetBackgroundColor(bgColor)
	driverRow.AddItem(createDialog.secretDriver, labelWidth+6, 0, true) //nolint:mnd
	driverRow.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	driverRow.AddItem(createDialog.secretDriverOptions, 0, 1, true)

	optionsLayout.SetBackgroundColor(bgColor)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(nameAndReplaceRow, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(createDialog.secretFile, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(createDialog.secretText, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(createDialog.secretLabels, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(driverRow, 1, 0, true)
	optionsLayout.AddItem(createDialog.form, dialogs.DialogFormHeight, 0, true)

	createDialog.layout.SetBackgroundColor(bgColor)
	createDialog.layout.SetBorder(true)
	createDialog.layout.SetBorderColor(style.DialogBorderColor)
	createDialog.layout.SetTitle(i18n.T("PODMAN SECRET CREATE"))
	createDialog.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	createDialog.layout.AddItem(optionsLayout, 0, 1, true)
	createDialog.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	return createDialog
}

// Display displays this primitive.
func (d *SecretCreateDialog) Display() {
	d.display = true
	d.focusElement = secretNameFocus

	d.secretName.SetText("")
	d.secretFile.SetText("")
	d.secretText.SetText("")
	d.secretLabels.SetText("")
	d.secretDriver.SetCurrentOption(0)
	d.secretDriverOptions.SetText("")
	d.secretReplace.SetChecked(false)
}

// IsDisplay returns true if primitive is shown.
func (d *SecretCreateDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *SecretCreateDialog) Hide() {
	d.display = false
}

// SetRect set rects for this primitive.
func (d *SecretCreateDialog) SetRect(x, y, width, height int) {
	if width > secretCreateDialogMaxWidth {
		emptySpace := (width - secretCreateDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = secretCreateDialogMaxWidth
	}

	if height > secretCreateDialogMaxHeight {
		emptySpace := (height - secretCreateDialogMaxHeight) / 2 //nolint:mnd
		y += emptySpace
		height = secretCreateDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// HasFocus returns whether or not this primitive has focus.
func (d *SecretCreateDialog) HasFocus() bool {
	if d.layout.HasFocus() || d.form.HasFocus() {
		return true
	}

	if d.secretName.HasFocus() || d.secretFile.HasFocus() {
		return true
	}

	if d.secretText.HasFocus() || d.secretLabels.HasFocus() {
		return true
	}

	if d.secretDriver.HasFocus() || d.secretDriverOptions.HasFocus() {
		return true
	}

	return d.secretReplace.HasFocus() || d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *SecretCreateDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case secretNameFocus:
		delegate(d.secretName)
	case secretReplaceFocus:
		delegate(d.secretReplace)
	case secretFileFocus:
		delegate(d.secretFile)
	case secretTextFocus:
		delegate(d.secretText)
	case secretLabelsFocus:
		delegate(d.secretLabels)
	case secretDriverFocus:
		delegate(d.secretDriver)
	case secretDriverOptionsFocus:
		delegate(d.secretDriverOptions)
	case secretFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.Key {
				d.focusElement = secretNameFocus

				d.Focus(delegate)
				d.form.SetFocus(0)

				return nil
			}

			return event
		})

		delegate(d.form)
	}
}

// Draw draws this primitive into the screen.
func (d *SecretCreateDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)

	x, y, width, height := d.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// InputHandler returns input handler function for this primitive.
func (d *SecretCreateDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("secret create dialog: event %v received", event)

		if event.Key() == utils.CloseDialogKey.Key {
			if !d.secretDriver.HasFocus() {
				d.cancelHandler()

				return
			}
		}

		if event.Key() == utils.SwitchFocusKey.Key {
			d.setFocusElement()
		}

		if d.secretName.HasFocus() {
			if nameHandler := d.secretName.InputHandler(); nameHandler != nil {
				nameHandler(event, setFocus)

				return
			}
		}

		if d.secretReplace.HasFocus() {
			if replaceHandler := d.secretReplace.InputHandler(); replaceHandler != nil {
				replaceHandler(event, setFocus)

				return
			}
		}

		if d.secretFile.HasFocus() {
			if fileHandler := d.secretFile.InputHandler(); fileHandler != nil {
				fileHandler(event, setFocus)

				return
			}
		}

		if d.secretText.HasFocus() {
			if textHandler := d.secretText.InputHandler(); textHandler != nil {
				textHandler(event, setFocus)

				return
			}
		}

		if d.secretLabels.HasFocus() {
			if labelsHandler := d.secretLabels.InputHandler(); labelsHandler != nil {
				labelsHandler(event, setFocus)

				return
			}
		}

		if d.secretDriver.HasFocus() {
			event = utils.ParseKeyEventKey(event)
			if driverHandler := d.secretDriver.InputHandler(); driverHandler != nil {
				driverHandler(event, setFocus)

				return
			}
		}

		if d.secretDriverOptions.HasFocus() {
			if driverOptionsHandler := d.secretDriverOptions.InputHandler(); driverOptionsHandler != nil {
				driverOptionsHandler(event, setFocus)

				return
			}
		}

		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)

				return
			}
		}
	})
}

// SetCreateFunc sets form create button selected function.
func (d *SecretCreateDialog) SetCreateFunc(handler func()) *SecretCreateDialog {
	d.createHandler = handler
	createButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	createButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *SecretCreateDialog) SetCancelFunc(handler func()) *SecretCreateDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd

	cancelButton.SetSelectedFunc(handler)

	return d
}

// GetCreateOptions returns secret create options.
func (d *SecretCreateDialog) GetCreateOptions() *secrets.SecretCreateOptions {
	var createOptions secrets.SecretCreateOptions

	createOptions.Name = strings.TrimSpace(d.secretName.GetText())
	createOptions.File = strings.TrimSpace(d.secretFile.GetText())
	createOptions.Text = strings.TrimSpace(d.secretText.GetText())
	createOptions.Labels = strings.Split(strings.TrimSpace(d.secretLabels.GetText()), " ")
	createOptions.Replace = d.secretReplace.IsChecked()
	_, translatedDriver := d.secretDriver.GetCurrentOption()
	createOptions.Driver = d.getEnglishDriver(translatedDriver)
	createOptions.DriverOptions = strings.Split(strings.TrimSpace(d.secretDriverOptions.GetText()), " ")

	return &createOptions
}

// getEnglishDriver maps translated driver name back to English
func (d *SecretCreateDialog) getEnglishDriver(translatedDriver string) string {
	driverMap := map[string]string{
		i18n.T("file"):  "file",
		i18n.T("pass"):  "pass",
		i18n.T("shell"): "shell",
	}

	if englishDriver, ok := driverMap[translatedDriver]; ok {
		return englishDriver
	}

	// Fallback to original if not found
	return translatedDriver
}

func (d *SecretCreateDialog) setFocusElement() {
	switch d.focusElement {
	case secretNameFocus:
		d.focusElement = secretReplaceFocus
	case secretReplaceFocus:
		d.focusElement = secretFileFocus
	case secretFileFocus:
		d.focusElement = secretTextFocus
	case secretTextFocus:
		d.focusElement = secretLabelsFocus
	case secretLabelsFocus:
		d.focusElement = secretDriverFocus
	case secretDriverFocus:
		d.focusElement = secretDriverOptionsFocus
	case secretDriverOptionsFocus:
		d.focusElement = secretFormFocus
	}
}

// UpdateLanguage updates all translatable text in the dialog.
func (d *SecretCreateDialog) UpdateLanguage() {
	// Update dialog title
	d.layout.SetTitle(i18n.T("PODMAN SECRET CREATE"))
	
	// Update field labels
	d.secretName.SetLabel(utils.StringToInputLabel(i18n.T("name:"), labelWidth))
	d.secretFile.SetLabel(utils.StringToInputLabel(i18n.T("secret file:"), labelWidth))
	d.secretText.SetLabel(utils.StringToInputLabel(i18n.T("secret text:"), labelWidth))
	d.secretLabels.SetLabel(utils.StringToInputLabel(i18n.T("labels:"), labelWidth))
	d.secretReplace.SetLabel(i18n.T("replace") + " ")
	d.secretDriver.SetLabel(i18n.T("driver:"))
	d.secretDriverOptions.SetLabel(i18n.T("driver options:") + " ")
	
	// Update driver dropdown options
	currentOption, _ := d.secretDriver.GetCurrentOption()
	d.secretDriver.SetOptions([]string{i18n.T("file"), i18n.T("pass"), i18n.T("shell")}, nil)
	d.secretDriver.SetCurrentOption(currentOption)
	
	// Update form buttons
	d.form.ClearButtons()
	d.form.AddButton(i18n.T("Cancel"), nil)
	d.form.AddButton(i18n.T("Create"), nil)
	
	// Re-set button handlers
	if d.cancelHandler != nil {
		cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd
		cancelButton.SetSelectedFunc(d.cancelHandler)
	}
	
	if d.createHandler != nil {
		createButton := d.form.GetButton(d.form.GetButtonCount() - 1)
		createButton.SetSelectedFunc(d.createHandler)
	}
}
