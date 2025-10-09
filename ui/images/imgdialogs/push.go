package imgdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	imagePushDialogMaxWidth     = 90
	imagePushDialogMaxHeight    = 15
	imagePushDialogLabelPadding = 4
)

const (
	imagePushDesitnationFocus = 0 + iota
	imagePushCompressFocus
	imagePushFormatFocus
	imagePushSkipTLSVerifyFocus
	imagePushUsernameFocus
	imagePushPasswordFocus
	imagePushAuthFileFocus
	imagePushFormFocus
)

// ImagePushDialog represents image push dialog primitive.
type ImagePushDialog struct {
	*tview.Box

	layout        *tview.Flex
	contentLayout *tview.Flex // holds the main input layout
	imageInfo     *tview.InputField
	destination   *tview.InputField
	compress      *tview.Checkbox
	format        *tview.DropDown
	skipTLSVerify *tview.Checkbox
	authFile      *tview.InputField
	username      *tview.InputField
	password      *tview.InputField
	form          *tview.Form
	display       bool
	pushHandler   func()
	cancelHandler func()
	focusElement  int
}

// NewImagePushDialog returns a new image push dialog primitive.
func NewImagePushDialog() *ImagePushDialog {
	dialog := &ImagePushDialog{
		Box:           tview.NewBox(),
		layout:        tview.NewFlex(),
		imageInfo:     tview.NewInputField(),
		destination:   tview.NewInputField(),
		compress:      tview.NewCheckbox(),
		format:        tview.NewDropDown(),
		skipTLSVerify: tview.NewCheckbox(),
		authFile:      tview.NewInputField(),
		username:      tview.NewInputField(),
		password:      tview.NewInputField(),
		form:          tview.NewForm(),
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	
	// Calculate label width dynamically based on translated text
	labels := []string{
		i18n.T("destination:"),
		i18n.T("compress:"),
		i18n.T("authfile:"),
		i18n.T("username:"),
	}
	labelWidth := 0
	for _, label := range labels {
		if width := i18n.GetDisplayWidth(label); width > labelWidth {
			labelWidth = width
		}
	}

	// image info field
	imageIDLabel := i18n.T("IMAGE ID:")
	dialog.imageInfo.SetBackgroundColor(style.DialogBgColor)
	dialog.imageInfo.SetLabel("[::b]" + imageIDLabel)
	dialog.imageInfo.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.imageInfo.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// destination input field
	dialog.destination.SetBackgroundColor(bgColor)
	dialog.destination.SetLabel(utils.StringToInputLabel(i18n.T("destination:"), labelWidth))
	dialog.destination.SetFieldStyle(style.InputFieldStyle)
	dialog.destination.SetLabelStyle(style.InputLabelStyle)

	// compress checkbox
	compressLabel := i18n.T("compress:")
	dialog.compress.SetBackgroundColor(bgColor)
	dialog.compress.SetLabelColor(fgColor)
	dialog.compress.SetLabel(compressLabel)
	dialog.compress.SetLabelWidth(i18n.GetDisplayWidth(compressLabel) + 1)
	dialog.compress.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// format dropdown
	formatLabel := i18n.T("format:")
	dialog.format.SetLabel(formatLabel)
	dialog.format.SetTitleAlign(tview.AlignRight)
	dialog.format.SetLabelColor(fgColor)
	dialog.format.SetLabelWidth(i18n.GetDisplayWidth(formatLabel) + 1)
	dialog.format.SetBackgroundColor(bgColor)
	dialog.format.SetOptions([]string{
		"oci",
		"v2v2",
		"v2v1",
	},
		nil)
	dialog.format.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	dialog.format.SetFocusedStyle(style.DropDownFocused)
	dialog.format.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// skipTLSVerify checkbox
	skipTLSVerifyLabel := i18n.T("skip tls verify:")

	dialog.skipTLSVerify.SetBackgroundColor(bgColor)
	dialog.skipTLSVerify.SetLabelColor(fgColor)
	dialog.skipTLSVerify.SetLabel(skipTLSVerifyLabel)
	dialog.skipTLSVerify.SetLabelWidth(i18n.GetDisplayWidth(skipTLSVerifyLabel) + 1)
	dialog.skipTLSVerify.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// authfile input field
	dialog.authFile.SetBackgroundColor(bgColor)
	dialog.authFile.SetLabel(utils.StringToInputLabel(i18n.T("authfile:"), labelWidth))
	dialog.authFile.SetFieldStyle(style.InputFieldStyle)
	dialog.authFile.SetLabelStyle(style.InputLabelStyle)

	// username input field
	dialog.username.SetBackgroundColor(bgColor)
	dialog.username.SetLabel(utils.StringToInputLabel(i18n.T("username:"), labelWidth))
	dialog.username.SetFieldStyle(style.InputFieldStyle)
	dialog.username.SetLabelStyle(style.InputLabelStyle)

	// password input field
	passwordLabel := i18n.T("password:")

	dialog.password.SetBackgroundColor(bgColor)
	dialog.password.SetLabel(utils.StringToInputLabel(passwordLabel, i18n.GetDisplayWidth(passwordLabel)+1))
	dialog.password.SetFieldStyle(style.InputFieldStyle)
	dialog.password.SetLabelStyle(style.InputLabelStyle)
	dialog.password.SetMaskCharacter('*')

	// form
	dialog.form.AddButton(i18n.T("Cancel"), nil)
	dialog.form.AddButton(i18n.T("Push"), nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	// dropdowns and checkbox row layour
	compressWidth := i18n.GetDisplayWidth(i18n.T("compress:")) + 5  //nolint:mnd
	formatWidth := i18n.GetDisplayWidth(formatLabel) + 10           //nolint:mnd
	skipTLSWidth := i18n.GetDisplayWidth(skipTLSVerifyLabel) + 5    //nolint:mnd
	
	dcLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	dcLayout.AddItem(dialog.compress, compressWidth, 0, true)
	dcLayout.AddItem(utils.EmptyBoxSpace(bgColor), 2, 0, false)  //nolint:mnd
	dcLayout.AddItem(dialog.format, formatWidth, 0, true)
	dcLayout.AddItem(utils.EmptyBoxSpace(bgColor), 2, 0, false)  //nolint:mnd
	dcLayout.AddItem(dialog.skipTLSVerify, skipTLSWidth, 0, true)
	dcLayout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)

	// username and password row layout
	userPassLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	userPassLayout.AddItem(dialog.username, 0, 1, true)
	userPassLayout.AddItem(utils.EmptyBoxSpace(bgColor), 3, 0, false) //nolint:mnd
	userPassLayout.AddItem(dialog.password, 0, 1, true)

	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(dialog.imageInfo, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(dialog.destination, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(dcLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(userPassLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(dialog.authFile, 0, 1, true)

	dialog.contentLayout = tview.NewFlex().SetDirection(tview.FlexColumn)
	dialog.contentLayout.SetBackgroundColor(bgColor)
	dialog.contentLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	dialog.contentLayout.AddItem(layout, 0, 1, true)
	dialog.contentLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetTitle(i18n.T("PODMAN IMAGE PUSH"))
	dialog.layout.AddItem(dialog.contentLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	dialog.Hide()

	return dialog
}

// Display displays this primitive.
func (d *ImagePushDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *ImagePushDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ImagePushDialog) Hide() {
	d.display = false
	d.focusElement = imagePushDesitnationFocus

	d.destination.SetText("")
	d.compress.SetChecked(false)
	d.format.SetCurrentOption(0)
	d.skipTLSVerify.SetChecked(false)
	d.authFile.SetText("")
	d.username.SetText("")
	d.password.SetText("")
}

// HasFocus returns whether or not this primitive has focus.
func (d *ImagePushDialog) HasFocus() bool { //nolint:cyclop
	if d.destination.HasFocus() || d.compress.HasFocus() {
		return true
	}

	if d.format.HasFocus() || d.skipTLSVerify.HasFocus() {
		return true
	}

	if d.username.HasFocus() || d.password.HasFocus() {
		return true
	}

	if d.authFile.HasFocus() || d.form.HasFocus() {
		return true
	}

	if d.layout.HasFocus() || d.Box.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *ImagePushDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case imagePushDesitnationFocus:
		delegate(d.destination)
	case imagePushCompressFocus:
		delegate(d.compress)
	case imagePushFormatFocus:
		delegate(d.format)
	case imagePushSkipTLSVerifyFocus:
		delegate(d.skipTLSVerify)
	case imagePushAuthFileFocus:
		delegate(d.authFile)
	case imagePushUsernameFocus:
		delegate(d.username)
	case imagePushPasswordFocus:
		delegate(d.password)
	case imagePushFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.Key {
				d.focusElement = imagePushDesitnationFocus
				d.Focus(delegate)
				d.form.SetFocus(0)

				return nil
			}

			return event
		})

		delegate(d.form)
	}
}

// InputHandler returns input handler function for this primitive.
func (d *ImagePushDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("image push dialog: event %v received", event)

		if event.Key() == utils.SwitchFocusKey.Key {
			d.setFocusElement()
		}

		if event.Key() == tcell.KeyEsc && !d.dropdownHasFocus() {
			d.cancelHandler()

			return
		}

		if d.destination.HasFocus() {
			if destinationHandler := d.destination.InputHandler(); destinationHandler != nil {
				destinationHandler(event, setFocus)

				return
			}
		}

		if d.compress.HasFocus() {
			if compressHandler := d.compress.InputHandler(); compressHandler != nil {
				compressHandler(event, setFocus)

				return
			}
		}

		if d.format.HasFocus() {
			if formatHandler := d.format.InputHandler(); formatHandler != nil {
				event = utils.ParseKeyEventKey(event)
				formatHandler(event, setFocus)

				return
			}
		}

		if d.skipTLSVerify.HasFocus() {
			if skipTLSVerifyHandler := d.skipTLSVerify.InputHandler(); skipTLSVerifyHandler != nil {
				skipTLSVerifyHandler(event, setFocus)

				return
			}
		}

		if d.authFile.HasFocus() {
			if authFileHandler := d.authFile.InputHandler(); authFileHandler != nil {
				authFileHandler(event, setFocus)

				return
			}
		}

		if d.username.HasFocus() {
			if usernameHandler := d.username.InputHandler(); usernameHandler != nil {
				usernameHandler(event, setFocus)

				return
			}
		}

		if d.password.HasFocus() {
			if passwordHandler := d.password.InputHandler(); passwordHandler != nil {
				passwordHandler(event, setFocus)

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

// SetRect set rects for this primitive.
func (d *ImagePushDialog) SetRect(x, y, width, height int) {
	if width > imagePushDialogMaxWidth {
		emptySpace := (width - imagePushDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = imagePushDialogMaxWidth
	}

	if height > imagePushDialogMaxHeight {
		emptySpace := (height - imagePushDialogMaxHeight) / 2 //nolint:mnd
		y += emptySpace
		height = imagePushDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ImagePushDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	x, y, width, height := d.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetPushFunc sets form push button selected function.
func (d *ImagePushDialog) SetPushFunc(handler func()) *ImagePushDialog {
	d.pushHandler = handler
	pushButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	pushButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *ImagePushDialog) SetCancelFunc(handler func()) *ImagePushDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd
	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetImageInfo sets selected image ID and name in push dialog.
func (d *ImagePushDialog) SetImageInfo(id string, name string) {
	containerInfo := fmt.Sprintf("%12s (%s)", id, name)
	containerInfo = utils.LabelWidthLeftPadding(containerInfo, imagePushDialogLabelPadding)

	d.imageInfo.SetText(containerInfo)
}

// GetImagePushOptions returns image push options based on user inputs.
func (d *ImagePushDialog) GetImagePushOptions() images.ImagePushOptions {
	var opts images.ImagePushOptions

	opts.Destination = strings.TrimSpace(d.destination.GetText())
	_, format := d.format.GetCurrentOption()
	format = strings.TrimSpace(format)
	opts.Format = format
	opts.Compress = d.compress.IsChecked()
	opts.SkipTLSVerify = d.skipTLSVerify.IsChecked()
	opts.Username = strings.TrimSpace(d.username.GetText())
	opts.Password = strings.TrimSpace(d.password.GetText())
	opts.AuthFile = strings.TrimSpace(d.authFile.GetText())

	return opts
}

// dropdownHasFocus returns true if image push dialog dropdown primitives.
// has focus.
func (d *ImagePushDialog) dropdownHasFocus() bool {
	return d.format.HasFocus()
}

func (d *ImagePushDialog) setFocusElement() {
	switch d.focusElement {
	case imagePushDesitnationFocus:
		d.focusElement = imagePushCompressFocus
	case imagePushCompressFocus:
		d.focusElement = imagePushFormatFocus
	case imagePushFormatFocus:
		d.focusElement = imagePushSkipTLSVerifyFocus
	case imagePushSkipTLSVerifyFocus:
		d.focusElement = imagePushUsernameFocus
	case imagePushUsernameFocus:
		d.focusElement = imagePushPasswordFocus
	case imagePushPasswordFocus:
		d.focusElement = imagePushAuthFileFocus
	case imagePushAuthFileFocus:
		d.focusElement = imagePushFormFocus
	}
}

// UpdateLanguage updates all translatable text in the dialog.
func (d *ImagePushDialog) UpdateLanguage() {
	bgColor := style.DialogBgColor
	
	// Update dialog title
	d.layout.SetTitle(i18n.T("PODMAN IMAGE PUSH"))
	
	// Calculate label width dynamically
	labels := []string{
		i18n.T("destination:"),
		i18n.T("compress:"),
		i18n.T("authfile:"),
		i18n.T("username:"),
	}
	labelWidth := 0
	for _, label := range labels {
		if width := i18n.GetDisplayWidth(label); width > labelWidth {
			labelWidth = width
		}
	}
	
	// Update image ID label
	imageIDLabel := i18n.T("IMAGE ID:")
	d.imageInfo.SetLabel("[::b]" + imageIDLabel)
	
	// Update field labels
	d.destination.SetLabel(utils.StringToInputLabel(i18n.T("destination:"), labelWidth))
	
	compressLabel := i18n.T("compress:")
	d.compress.SetLabel(compressLabel)
	d.compress.SetLabelWidth(i18n.GetDisplayWidth(compressLabel) + 1)
	
	formatLabel := i18n.T("format:")
	d.format.SetLabel(formatLabel)
	d.format.SetLabelWidth(i18n.GetDisplayWidth(formatLabel) + 1)
	
	skipTLSVerifyLabel := i18n.T("skip tls verify:")
	d.skipTLSVerify.SetLabel(skipTLSVerifyLabel)
	d.skipTLSVerify.SetLabelWidth(i18n.GetDisplayWidth(skipTLSVerifyLabel) + 1)
	
	d.authFile.SetLabel(utils.StringToInputLabel(i18n.T("authfile:"), labelWidth))
	d.username.SetLabel(utils.StringToInputLabel(i18n.T("username:"), labelWidth))
	
	passwordLabel := i18n.T("password:")
	d.password.SetLabel(utils.StringToInputLabel(passwordLabel, i18n.GetDisplayWidth(passwordLabel)+1))
	
	// Rebuild the checkbox/dropdown row with correct widths
	compressWidth := i18n.GetDisplayWidth(i18n.T("compress:")) + 5  //nolint:mnd
	formatWidth := i18n.GetDisplayWidth(formatLabel) + 10           //nolint:mnd
	skipTLSWidth := i18n.GetDisplayWidth(skipTLSVerifyLabel) + 5    //nolint:mnd
	
	dcLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	dcLayout.AddItem(d.compress, compressWidth, 0, true)
	dcLayout.AddItem(utils.EmptyBoxSpace(bgColor), 2, 0, false)  //nolint:mnd
	dcLayout.AddItem(d.format, formatWidth, 0, true)
	dcLayout.AddItem(utils.EmptyBoxSpace(bgColor), 2, 0, false)  //nolint:mnd
	dcLayout.AddItem(d.skipTLSVerify, skipTLSWidth, 0, true)
	dcLayout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	
	// username and password row layout
	userPassLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	userPassLayout.AddItem(d.username, 0, 1, true)
	userPassLayout.AddItem(utils.EmptyBoxSpace(bgColor), 3, 0, false) //nolint:mnd
	userPassLayout.AddItem(d.password, 0, 1, true)
	
	// Rebuild main layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(d.imageInfo, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(d.destination, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(dcLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(userPassLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(d.authFile, 0, 1, true)
	
	// Rebuild content layout
	d.contentLayout.Clear()
	d.contentLayout.SetBackgroundColor(bgColor)
	d.contentLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	d.contentLayout.AddItem(layout, 0, 1, true)
	d.contentLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	
	// Update form buttons
	d.form.ClearButtons()
	d.form.AddButton(i18n.T("Cancel"), nil)
	d.form.AddButton(i18n.T("Push"), nil)
	
	// Re-set button handlers
	if d.cancelHandler != nil {
		cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd
		cancelButton.SetSelectedFunc(d.cancelHandler)
	}
	
	if d.pushHandler != nil {
		pushButton := d.form.GetButton(d.form.GetButtonCount() - 1)
		pushButton.SetSelectedFunc(d.pushHandler)
	}
}
