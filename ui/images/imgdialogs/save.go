package imgdialogs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v5/libpod/define"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	imageSaveDialogMaxWidth     = 70
	imageSaveDialogMaxHeight    = 13
	imageSaveDialogLabelPadding = 1
)

var errSaveEmptyOuputName = errors.New("empty output name")

const (
	imageSaveOutputFocus = 0 + iota
	imageSaveCompressFocus
	imageSaveAcceptUncompressedFocus
	imageSaveFormatFocus
	imageSaveFormFocus
)

// ImageSaveDialog represents image save dialog primitive.
type ImageSaveDialog struct {
	*tview.Box

	layout                *tview.Flex
	contentLayout         *tview.Flex // holds the main options layout
	imageInfo             *tview.InputField
	output                *tview.InputField
	compress              *tview.Checkbox
	format                *tview.DropDown
	ociAcceptUncompressed *tview.Checkbox
	form                  *tview.Form
	display               bool
	saveHandler           func()
	cancelHandler         func()
	focusElement          int
}

// NewImageSaveDialog returns new image save dialog.
func NewImageSaveDialog() *ImageSaveDialog {
	dialog := &ImageSaveDialog{
		Box:                   tview.NewBox(),
		layout:                tview.NewFlex(),
		imageInfo:             tview.NewInputField(),
		output:                tview.NewInputField(),
		compress:              tview.NewCheckbox(),
		format:                tview.NewDropDown(),
		ociAcceptUncompressed: tview.NewCheckbox(),
		form:                  tview.NewForm(),
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	
	// Calculate label width dynamically
	labels := []string{
		i18n.T("output:"),
		i18n.T("compress:"),
		i18n.T("format:"),
	}
	labelWidth := 0
	for _, label := range labels {
		if width := i18n.GetDisplayWidth(label); width > labelWidth {
			labelWidth = width
		}
	}

	// image info
	imageInfoLabel := i18n.T("IMAGE ID:")

	dialog.imageInfo.SetBackgroundColor(style.DialogBgColor)
	dialog.imageInfo.SetLabel("[::b]" + imageInfoLabel)
	dialog.imageInfo.SetLabelWidth(i18n.GetDisplayWidth(imageInfoLabel))
	dialog.imageInfo.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.imageInfo.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// output
	dialog.output.SetBackgroundColor(bgColor)
	dialog.output.SetLabel(utils.StringToInputLabel(i18n.T("output:"), labelWidth))
	dialog.output.SetFieldStyle(style.InputFieldStyle)
	dialog.output.SetLabelStyle(style.InputLabelStyle)

	// compress
	compressLabel := i18n.T("compress:")
	dialog.compress.SetBackgroundColor(bgColor)
	dialog.compress.SetLabelColor(fgColor)
	dialog.compress.SetLabel(compressLabel)
	dialog.compress.SetLabelWidth(i18n.GetDisplayWidth(compressLabel) + 1)
	dialog.compress.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// format
	formatLabel := i18n.T("format:")
	dialog.format.SetBackgroundColor(bgColor)
	dialog.format.SetLabelColor(fgColor)
	dialog.format.SetLabel(formatLabel)
	dialog.format.SetLabelWidth(i18n.GetDisplayWidth(formatLabel) + 1)
	dialog.format.SetOptions([]string{
		i18n.T(define.V2s2Archive),
		i18n.T(define.V2s2ManifestDir),
		i18n.T(define.OCIArchive),
		i18n.T(define.OCIManifestDir),
	},
		nil)
	dialog.format.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	dialog.format.SetFocusedStyle(style.DropDownFocused)
	dialog.format.SetCurrentOption(0)
	dialog.format.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// OciAcceptUncompressed
	ociLabel := i18n.T("accept uncompressed (OCI images):")
	dialog.ociAcceptUncompressed.SetBackgroundColor(bgColor)
	dialog.ociAcceptUncompressed.SetLabelColor(fgColor)
	dialog.ociAcceptUncompressed.SetLabel(ociLabel)
	dialog.ociAcceptUncompressed.SetLabelWidth(i18n.GetDisplayWidth(ociLabel) + 1)
	dialog.ociAcceptUncompressed.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// form
	dialog.form.AddButton(i18n.T("Cancel"), nil)
	dialog.form.AddButton(i18n.T("Save"), nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	compressWidth := i18n.GetDisplayWidth(compressLabel) + 5  //nolint:mnd
	ociWidth := i18n.GetDisplayWidth(ociLabel) + 5            //nolint:mnd
	
	compressRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	compressRow.SetBackgroundColor(bgColor)
	compressRow.AddItem(dialog.compress, compressWidth, 0, true)
	compressRow.AddItem(utils.EmptyBoxSpace(bgColor), 2, 0, false) //nolint:mnd
	compressRow.AddItem(dialog.ociAcceptUncompressed, ociWidth, 0, true)
	compressRow.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)

	optionsLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.imageInfo, 0, 1, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.output, 0, 1, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(compressRow, 0, 1, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.format, 0, 1, true)

	dialog.contentLayout = tview.NewFlex().SetDirection(tview.FlexColumn)
	dialog.contentLayout.SetBackgroundColor(bgColor)
	dialog.contentLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	dialog.contentLayout.AddItem(optionsLayout, 0, 1, true)
	dialog.contentLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetTitle(i18n.T("PODMAN IMAGE SAVE"))
	dialog.layout.AddItem(dialog.contentLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive.
func (d *ImageSaveDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *ImageSaveDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ImageSaveDialog) Hide() {
	d.display = false

	d.SetImageInfo("", "")

	d.focusElement = imageSaveOutputFocus

	d.output.SetText("")
	d.compress.SetChecked(false)
	d.ociAcceptUncompressed.SetChecked(false)
	d.format.SetCurrentOption(0)
}

// HasFocus returns whether or not this primitive has focus.
func (d *ImageSaveDialog) HasFocus() bool {
	if d.output.HasFocus() || d.compress.HasFocus() {
		return true
	}

	if d.format.HasFocus() || d.ociAcceptUncompressed.HasFocus() {
		return true
	}

	if d.form.HasFocus() || d.layout.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *ImageSaveDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case imageSaveOutputFocus:
		delegate(d.output)
	case imageSaveCompressFocus:
		delegate(d.compress)
	case imageSaveAcceptUncompressedFocus:
		delegate(d.ociAcceptUncompressed)
	case imageSaveFormatFocus:
		delegate(d.format)
	case imageSaveFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.Key {
				d.focusElement = imageSaveOutputFocus
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
func (d *ImageSaveDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("image save dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc {
			if !d.format.HasFocus() {
				d.cancelHandler()

				return
			}
		}

		if event.Key() == utils.SwitchFocusKey.Key {
			d.setFocusElement()
		}

		// drop down event
		if d.format.HasFocus() {
			event = utils.ParseKeyEventKey(event)
			if formatHandler := d.format.InputHandler(); formatHandler != nil {
				formatHandler(event, setFocus)

				return
			}
		}

		if d.output.HasFocus() {
			if outputHandler := d.output.InputHandler(); outputHandler != nil {
				outputHandler(event, setFocus)

				return
			}
		}

		if d.compress.HasFocus() {
			if compressHandler := d.compress.InputHandler(); compressHandler != nil {
				compressHandler(event, setFocus)

				return
			}
		}

		if d.ociAcceptUncompressed.HasFocus() {
			if acceptUncompressedhandler := d.ociAcceptUncompressed.InputHandler(); acceptUncompressedhandler != nil {
				acceptUncompressedhandler(event, setFocus)

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
func (d *ImageSaveDialog) SetRect(x, y, width, height int) {
	if width > imageSaveDialogMaxWidth {
		emptySpace := (width - imageSaveDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = imageSaveDialogMaxWidth
	}

	if height > imageSaveDialogMaxHeight {
		emptySpace := (height - imageSaveDialogMaxHeight) / 2 //nolint:mnd
		y += emptySpace
		height = imageSaveDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ImageSaveDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	x, y, width, height := d.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetSaveFunc sets form save button selected function.
func (d *ImageSaveDialog) SetSaveFunc(handler func()) *ImageSaveDialog {
	d.saveHandler = handler
	saveButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	saveButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *ImageSaveDialog) SetCancelFunc(handler func()) *ImageSaveDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd

	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetImageInfo sets selected image ID and name in save dialog.
func (d *ImageSaveDialog) SetImageInfo(id string, name string) {
	nameSplited := strings.Split(name, "/")

	l := len(nameSplited)
	if l > 1 {
		name = nameSplited[l-1]
	}

	imageInfo := fmt.Sprintf("%s (%s)", id, name)
	imageInfo = utils.LabelWidthLeftPadding(imageInfo, imageSaveDialogLabelPadding)

	d.imageInfo.SetText(imageInfo)
}

// ImageSaveOptions prepare and returns image save options.
func (d *ImageSaveDialog) ImageSaveOptions() (images.ImageSaveOptions, error) {
	opts := images.ImageSaveOptions{
		Compressed:                  d.compress.IsChecked(),
		OciAcceptUncompressedLayers: d.ociAcceptUncompressed.IsChecked(),
	}

	// Get format index and map to original format constant
	formatIndex, _ := d.format.GetCurrentOption()
	formats := []string{
		define.V2s2Archive,
		define.V2s2ManifestDir,
		define.OCIArchive,
		define.OCIManifestDir,
	}
	opts.Format = formats[formatIndex]

	output := strings.TrimSpace(d.output.GetText())
	if output == "" {
		return opts, errSaveEmptyOuputName
	}

	outputPath, err := utils.ResolveHomeDir(output)
	if err != nil {
		return opts, err
	}

	opts.Output = outputPath

	return opts, nil
}

func (d *ImageSaveDialog) setFocusElement() {
	switch d.focusElement {
	case imageSaveOutputFocus:
		d.focusElement = imageSaveCompressFocus
	case imageSaveCompressFocus:
		d.focusElement = imageSaveAcceptUncompressedFocus
	case imageSaveAcceptUncompressedFocus:
		d.focusElement = imageSaveFormatFocus
	case imageSaveFormatFocus:
		d.focusElement = imageSaveFormFocus
	}
}

// UpdateLanguage updates all translatable text in the dialog.
func (d *ImageSaveDialog) UpdateLanguage() {
	bgColor := style.DialogBgColor
	
	// Update dialog title
	d.layout.SetTitle(i18n.T("PODMAN IMAGE SAVE"))
	
	// Calculate label width dynamically
	labels := []string{
		i18n.T("output:"),
		i18n.T("compress:"),
		i18n.T("format:"),
	}
	labelWidth := 0
	for _, label := range labels {
		if width := i18n.GetDisplayWidth(label); width > labelWidth {
			labelWidth = width
		}
	}
	
	// Update image ID label
	imageInfoLabel := i18n.T("IMAGE ID:")
	d.imageInfo.SetLabel("[::b]" + imageInfoLabel)
	d.imageInfo.SetLabelWidth(i18n.GetDisplayWidth(imageInfoLabel))
	
	// Update field labels
	d.output.SetLabel(utils.StringToInputLabel(i18n.T("output:"), labelWidth))
	
	compressLabel := i18n.T("compress:")
	d.compress.SetLabel(compressLabel)
	d.compress.SetLabelWidth(i18n.GetDisplayWidth(compressLabel) + 1)
	
	formatLabel := i18n.T("format:")
	d.format.SetLabel(formatLabel)
	d.format.SetLabelWidth(i18n.GetDisplayWidth(formatLabel) + 1)
	
	// Update format options
	currentFormatIndex, _ := d.format.GetCurrentOption()
	d.format.SetOptions([]string{
		i18n.T(define.V2s2Archive),
		i18n.T(define.V2s2ManifestDir),
		i18n.T(define.OCIArchive),
		i18n.T(define.OCIManifestDir),
	}, nil)
	d.format.SetCurrentOption(currentFormatIndex)
	
	ociLabel := i18n.T("accept uncompressed (OCI images):")
	d.ociAcceptUncompressed.SetLabel(ociLabel)
	d.ociAcceptUncompressed.SetLabelWidth(i18n.GetDisplayWidth(ociLabel) + 1)
	
	// Rebuild compress row with correct widths
	compressWidth := i18n.GetDisplayWidth(compressLabel) + 5  //nolint:mnd
	ociWidth := i18n.GetDisplayWidth(ociLabel) + 5            //nolint:mnd
	
	compressRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	compressRow.SetBackgroundColor(bgColor)
	compressRow.AddItem(d.compress, compressWidth, 0, true)
	compressRow.AddItem(utils.EmptyBoxSpace(bgColor), 2, 0, false) //nolint:mnd
	compressRow.AddItem(d.ociAcceptUncompressed, ociWidth, 0, true)
	compressRow.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	
	// Rebuild options layout
	optionsLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(d.imageInfo, 0, 1, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(d.output, 0, 1, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(compressRow, 0, 1, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(d.format, 0, 1, true)
	
	// Rebuild content layout
	d.contentLayout.Clear()
	d.contentLayout.SetBackgroundColor(bgColor)
	d.contentLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	d.contentLayout.AddItem(optionsLayout, 0, 1, true)
	d.contentLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	
	// Update form buttons
	d.form.ClearButtons()
	d.form.AddButton(i18n.T("Cancel"), nil)
	d.form.AddButton(i18n.T("Save"), nil)
	
	// Re-set button handlers
	if d.cancelHandler != nil {
		cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd
		cancelButton.SetSelectedFunc(d.cancelHandler)
	}
	
	if d.saveHandler != nil {
		saveButton := d.form.GetButton(d.form.GetButtonCount() - 1)
		saveButton.SetSelectedFunc(d.saveHandler)
	}
}
