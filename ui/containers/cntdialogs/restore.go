package cntdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	cntRestoreContainersFocus = 0 + iota
	cntRestorePodFocus
	cntRestoreNameFocus
	cntRestorePublishPortsFocus
	cntRestoreImportArchiveFocus
	cntRestoreKeepFocus
	cntRestoreIgnoreStaticIPFocus
	cntRestoreIgnoreStaticMACFocus
	cntRestoreFileLocksFocus
	cntRestorePrintStatsFocus
	cntRestoreTCPEstablishedFocus
	cntRestoreIgnroeVolumesFocus
	cntRestoreIgnoreRootFsFocus
	cntRestoreFormFocus
)

const (
	cntRestoreDialogLabelWidth            = 14
	cntRestoreDialogPadding               = 2
	cntRestoreDialogChkGroupColTwoWidth   = 18
	cntRestoreDialogChkGroupColThreeWidth = 20
	cntRestoreDialogChkGroupColFourWidth  = 16
	cntRestoreDialogMaxWidth              = 88 // Optimized width for all languages after translation improvements
	cntRestoreDialogMaxHeight             = 17
	cntRestoreDialogSingleFieldWidth      = cntRestoreDialogMaxWidth -
		cntRestoreDialogLabelWidth - (2 * cntRestoreDialogPadding) //nolint:mnd
)

// ContainerRestoreDialog implements container restore dialog primitive.
type ContainerRestoreDialog struct {
	*tview.Box

	layout          *tview.Flex
	containers      *tview.DropDown
	pods            *tview.DropDown
	name            *tview.InputField
	publishPorts    *tview.InputField
	importArchive   *tview.InputField
	ignoreRootFS    *tview.Checkbox
	ignoreVolumes   *tview.Checkbox
	ignoreStaticIP  *tview.Checkbox
	ignoreStaticMAC *tview.Checkbox
	keep            *tview.Checkbox
	tcpEstablished  *tview.Checkbox
	fileLocks       *tview.Checkbox
	printStats      *tview.Checkbox
	form            *tview.Form
	row5            *tview.Flex
	row6            *tview.Flex
	display         bool
	focusElement    int
	restoreHandler  func()
	cancelHandler   func()
}

// NewContainerRestoreDialog returns new container dialog primitive.
func NewContainerRestoreDialog() *ContainerRestoreDialog {
	dialog := &ContainerRestoreDialog{
		Box:             tview.NewBox(),
		layout:          tview.NewFlex(),
		containers:      tview.NewDropDown(),
		pods:            tview.NewDropDown(),
		name:            tview.NewInputField(),
		publishPorts:    tview.NewInputField(),
		importArchive:   tview.NewInputField(),
		keep:            tview.NewCheckbox(),
		ignoreStaticIP:  tview.NewCheckbox(),
		ignoreStaticMAC: tview.NewCheckbox(),
		fileLocks:       tview.NewCheckbox(),
		printStats:      tview.NewCheckbox(),
		ignoreRootFS:    tview.NewCheckbox(),
		ignoreVolumes:   tview.NewCheckbox(),
		tcpEstablished:  tview.NewCheckbox(),
		form:            tview.NewForm(),
	}

	fgColor := style.DialogFgColor
	bgColor := style.DialogBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected

	// Calculate dynamic checkbox column widths (minimal padding)
	col1Width := i18n.CalcMaxWidth(
		i18n.T("keep:"),
		i18n.T("print Stats:"),
	)
	col2Width := i18n.CalcMaxWidth(
		i18n.T("ignore static IP:"),
		i18n.T("tcp established:"),
	)
	col3Width := i18n.CalcMaxWidth(
		i18n.T("ignore static MAC:"),
		i18n.T("ignore volumes:"),
	)
	col4Width := i18n.CalcMaxWidth(
		i18n.T("file locks:"),
		i18n.T("ignore rootfs:"),
	)

	// containers
	containersLabel := fmt.Sprintf("[:#%x:b]%s[:-:-]", style.DialogBorderColor.Hex(), i18n.T("CONTAINER ID:"))

	dialog.containers.SetLabel(containersLabel)
	dialog.containers.SetLabelWidth(cntRestoreDialogLabelWidth)
	dialog.containers.SetFieldWidth(cntRestoreDialogSingleFieldWidth)
	dialog.containers.SetBackgroundColor(bgColor)
	dialog.containers.SetLabelColor(fgColor)
	dialog.containers.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	dialog.containers.SetFocusedStyle(style.DropDownFocused)
	dialog.containers.SetFieldStyle(style.InputFieldStyle)
	dialog.SetContainers(nil)
	dialog.containers.SetCurrentOption(0)

	// pod
	dialog.pods.SetLabel(i18n.T("pod:"))
	dialog.pods.SetLabelWidth(cntRestoreDialogLabelWidth)
	dialog.pods.SetFieldWidth(cntRestoreDialogSingleFieldWidth)
	dialog.pods.SetBackgroundColor(bgColor)
	dialog.pods.SetLabelColor(fgColor)
	dialog.pods.SetLabelColor(fgColor)
	dialog.pods.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	dialog.pods.SetFocusedStyle(style.DropDownFocused)
	dialog.pods.SetFieldStyle(style.InputFieldStyle)
	dialog.SetPods(nil)
	dialog.pods.SetCurrentOption(0)

	// name
	dialog.name.SetBackgroundColor(style.DialogBgColor)
	dialog.name.SetLabel(i18n.PadToWidth(i18n.T("name:"), cntRestoreDialogLabelWidth))
	dialog.name.SetFieldStyle(style.InputFieldStyle)
	dialog.name.SetLabelStyle(style.InputLabelStyle)

	// Publish ports
	dialog.publishPorts.SetBackgroundColor(style.DialogBgColor)
	dialog.publishPorts.SetLabel(i18n.PadToWidth(i18n.T("publish:"), cntRestoreDialogLabelWidth))
	dialog.publishPorts.SetFieldStyle(style.InputFieldStyle)
	dialog.publishPorts.SetLabelStyle(style.InputLabelStyle)

	// Import
	dialog.importArchive.SetBackgroundColor(style.DialogBgColor)
	dialog.importArchive.SetLabel(i18n.PadToWidth(i18n.T("import:"), cntRestoreDialogLabelWidth))
	dialog.importArchive.SetFieldStyle(style.InputFieldStyle)
	dialog.importArchive.SetLabelStyle(style.InputLabelStyle)

	// keep
	dialog.keep.SetLabel(i18n.T("keep:"))
	dialog.keep.SetLabelWidth(col1Width)
	dialog.keep.SetChecked(false)
	dialog.keep.SetBackgroundColor(bgColor)
	dialog.keep.SetLabelColor(fgColor)
	dialog.keep.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// ignoreStaticIP
	dialog.ignoreStaticIP.SetLabel(i18n.T("ignore static IP:"))
	dialog.ignoreStaticIP.SetLabelWidth(col2Width)
	dialog.ignoreStaticIP.SetChecked(false)
	dialog.ignoreStaticIP.SetBackgroundColor(bgColor)
	dialog.ignoreStaticIP.SetLabelColor(fgColor)
	dialog.ignoreStaticIP.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// ignoreStaticMAC
	dialog.ignoreStaticMAC.SetLabel(i18n.T("ignore static MAC:"))
	dialog.ignoreStaticMAC.SetLabelWidth(col3Width)
	dialog.ignoreStaticMAC.SetChecked(false)
	dialog.ignoreStaticMAC.SetBackgroundColor(bgColor)
	dialog.ignoreStaticMAC.SetLabelColor(fgColor)
	dialog.ignoreStaticMAC.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// fileLocks
	dialog.fileLocks.SetLabel(i18n.T("file locks:"))
	dialog.fileLocks.SetLabelWidth(col4Width)
	dialog.fileLocks.SetChecked(false)
	dialog.fileLocks.SetBackgroundColor(bgColor)
	dialog.fileLocks.SetLabelColor(fgColor)
	dialog.fileLocks.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// printStats
	dialog.printStats.SetLabel(i18n.T("print Stats:"))
	dialog.printStats.SetLabelWidth(col1Width)
	dialog.printStats.SetChecked(false)
	dialog.printStats.SetBackgroundColor(bgColor)
	dialog.printStats.SetLabelColor(fgColor)
	dialog.printStats.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// tcpEstablished
	dialog.tcpEstablished.SetLabel(i18n.T("tcp established:"))
	dialog.tcpEstablished.SetLabelWidth(col2Width)
	dialog.tcpEstablished.SetChecked(false)
	dialog.tcpEstablished.SetBackgroundColor(bgColor)
	dialog.tcpEstablished.SetLabelColor(fgColor)
	dialog.tcpEstablished.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// ignoreVolumes
	dialog.ignoreVolumes.SetLabel(i18n.T("ignore volumes:"))
	dialog.ignoreVolumes.SetLabelWidth(col3Width)
	dialog.ignoreVolumes.SetChecked(false)
	dialog.ignoreVolumes.SetBackgroundColor(bgColor)
	dialog.ignoreVolumes.SetLabelColor(fgColor)
	dialog.ignoreVolumes.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// ignoreRootFS
	dialog.ignoreRootFS.SetLabel(i18n.T("ignore rootfs:"))
	dialog.ignoreRootFS.SetLabelWidth(col4Width)
	dialog.ignoreRootFS.SetChecked(false)
	dialog.ignoreRootFS.SetBackgroundColor(bgColor)
	dialog.ignoreRootFS.SetLabelColor(fgColor)
	dialog.ignoreRootFS.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// form
	dialog.form.AddButton(i18n.T("Cancel"), nil)
	dialog.form.AddButton(i18n.T("Restore"), nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	// layout row #one
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(dialog.containers, 0, 1, true)
	// layout row #two
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(dialog.pods, 0, 1, true)

	// layout row #three
	row := tview.NewFlex().SetDirection(tview.FlexColumn)
	row.AddItem(dialog.name, 0, 1, true)
	row.AddItem(utils.EmptyBoxSpace(bgColor), 2, 0, false)
	row.AddItem(dialog.publishPorts, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(row, 0, 1, true)

	// layout row #four
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(dialog.importArchive, 0, 1, true)

	// layout row #five
	dialog.row5 = tview.NewFlex().SetDirection(tview.FlexColumn)
	dialog.row5.AddItem(dialog.keep, col1Width+3, 0, true)
	dialog.row5.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	dialog.row5.AddItem(dialog.ignoreStaticIP, col2Width+3, 0, true)
	dialog.row5.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	dialog.row5.AddItem(dialog.ignoreStaticMAC, col3Width+3, 0, true)
	dialog.row5.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	dialog.row5.AddItem(dialog.fileLocks, col4Width+3, 0, true)
	dialog.row5.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(dialog.row5, 0, 1, true)

	// layout row #six
	dialog.row6 = tview.NewFlex().SetDirection(tview.FlexColumn)
	dialog.row6.AddItem(dialog.printStats, col1Width+3, 0, true)
	dialog.row6.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	dialog.row6.AddItem(dialog.tcpEstablished, col2Width+3, 0, true)
	dialog.row6.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	dialog.row6.AddItem(dialog.ignoreVolumes, col3Width+3, 0, true)
	dialog.row6.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	dialog.row6.AddItem(dialog.ignoreRootFS, col4Width+3, 0, true)
	dialog.row6.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(dialog.row6, 0, 1, true)

	mainOptsLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	mainOptsLayout.SetBackgroundColor(bgColor)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	mainOptsLayout.AddItem(layout, 0, 1, true)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetTitle(i18n.T("PODMAN CONTAINER RESTORE"))
	dialog.layout.AddItem(mainOptsLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive.
func (d *ContainerRestoreDialog) Display() {
	d.display = true
	d.focusElement = cntRestoreContainersFocus

	d.containers.SetCurrentOption(0)
	d.pods.SetCurrentOption(0)
	d.name.SetText("")
	d.publishPorts.SetText("")
	d.importArchive.SetText("")
	d.keep.SetChecked(false)
	d.ignoreStaticIP.SetChecked(false)
	d.ignoreStaticMAC.SetChecked(false)
	d.fileLocks.SetChecked(false)
	d.printStats.SetChecked(false)
	d.tcpEstablished.SetChecked(false)
	d.ignoreVolumes.SetChecked(false)
	d.ignoreRootFS.SetChecked(false)
}

// IsDisplay returns true if this primitive is shown.
func (d *ContainerRestoreDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ContainerRestoreDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (d *ContainerRestoreDialog) HasFocus() bool {
	for _, primitive := range d.getInnerPrimitives() {
		if primitive.HasFocus() {
			return true
		}
	}

	if d.layout.HasFocus() || d.Box.HasFocus() {
		return true
	}

	return false
}

// Focus is called when this primitive receives focus.
func (d *ContainerRestoreDialog) Focus(delegate func(p tview.Primitive)) {
	// all priviteves that can accept inputs
	if d.focusElement != cntRestoreFormFocus {
		primitives := d.getInnerPrimitives()
		delegate(primitives[d.focusElement])

		return
	}

	button := d.form.GetButton(d.form.GetButtonCount() - 1)

	button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == utils.SwitchFocusKey.Key {
			d.focusElement = cntRestoreContainersFocus

			d.Focus(delegate)
			d.form.SetFocus(0)

			return nil
		}

		return event
	})

	delegate(d.form)
}

// InputHandler returns input handler function for this primitive.
func (d *ContainerRestoreDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:cyclop,lll
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container restore dialog: event %v received", event)

		if event.Key() == utils.CloseDialogKey.Key {
			if !d.containers.HasFocus() && !d.pods.HasFocus() {
				d.cancelHandler()
			}
		}

		if event.Key() == utils.SwitchFocusKey.Key {
			if d.focusElement != cntRestoreFormFocus {
				d.setFocusElement()
			}
		}

		// all priviteves that can accept inputs
		for _, primitive := range d.getInnerPrimitives() {
			if primitive.HasFocus() {
				if d.containers.HasFocus() || d.pods.HasFocus() {
					event = utils.ParseKeyEventKey(event)
				}

				if handler := primitive.InputHandler(); handler != nil {
					handler(event, setFocus)

					return
				}
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *ContainerRestoreDialog) SetRect(x, y, width, height int) {
	if width > cntRestoreDialogMaxWidth {
		emptySpace := (width - cntRestoreDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = cntRestoreDialogMaxWidth
	}

	if height > cntRestoreDialogMaxHeight {
		emptySpace := (height - cntRestoreDialogMaxHeight) / 2 //nolint:mnd
		y += emptySpace
		height = cntRestoreDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive into the screen.
func (d *ContainerRestoreDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)

	x, y, width, height := d.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetRestoreFunc sets form restore button selected function.
func (d *ContainerRestoreDialog) SetRestoreFunc(handler func()) *ContainerRestoreDialog {
	d.restoreHandler = handler
	restoreButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	restoreButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *ContainerRestoreDialog) SetCancelFunc(handler func()) *ContainerRestoreDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd

	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetContainers sets containers dropdown options.
func (d *ContainerRestoreDialog) SetContainers(cnts [][]string) {
	emptyOptions := fmt.Sprintf("%*s", cntRestoreDialogSingleFieldWidth, " ")
	cntOptions := []string{emptyOptions}

	for i := range cnts {
		cntInfo := fmt.Sprintf("%s (%s)", utils.GetIDWithLimit(cnts[i][0]), cnts[i][1])
		cntInfoOption := fmt.Sprintf("%-*s", cntRestoreDialogSingleFieldWidth, cntInfo)
		cntOptions = append(cntOptions, cntInfoOption)
	}

	d.containers.SetOptions(cntOptions, nil)
}

// SetPods sets pods dropdown options.
func (d *ContainerRestoreDialog) SetPods(pods [][]string) {
	emptyOptions := fmt.Sprintf("%*s", cntRestoreDialogSingleFieldWidth, " ")
	podOptions := []string{emptyOptions}

	for i := range pods {
		podInfo := fmt.Sprintf("%s (%s)", utils.GetIDWithLimit(pods[i][0]), pods[i][1])
		podInfoOption := fmt.Sprintf("%-*s", cntRestoreDialogSingleFieldWidth, podInfo)
		podOptions = append(podOptions, podInfoOption)
	}

	d.pods.SetOptions(podOptions, nil)
}

func (d *ContainerRestoreDialog) GetRestoreOptions() containers.CntRestoreOptions {
	var opts containers.CntRestoreOptions

	_, cntInfoString := d.containers.GetCurrentOption()
	if strings.TrimSpace(cntInfoString) != "" {
		opts.ContainerID = strings.Split(cntInfoString, " ")[0]
	}

	_, podInfoString := d.pods.GetCurrentOption()
	if strings.TrimSpace(podInfoString) != "" {
		opts.PodID = strings.Split(podInfoString, " ")[0]
	}

	opts.Name = strings.TrimSpace(d.name.GetText())

	publishPortsList := strings.TrimSpace(d.publishPorts.GetText())
	opts.Publish = strings.Split(publishPortsList, " ")

	importArchive := strings.TrimSpace(d.importArchive.GetText())
	if strings.Index(importArchive, "~") == 0 {
		importArchive, _ = utils.ResolveHomeDir(importArchive)
	}

	opts.Import = importArchive

	opts.Keep = d.keep.IsChecked()
	opts.IgnoreStaticIP = d.ignoreStaticIP.IsChecked()
	opts.IgnoreStaticMAC = d.ignoreStaticMAC.IsChecked()
	opts.FileLocks = d.fileLocks.IsChecked()
	opts.PrintStats = d.printStats.IsChecked()
	opts.TCPEstablished = d.tcpEstablished.IsChecked()
	opts.IgnoreVolumes = d.ignoreVolumes.IsChecked()
	opts.IgnoreRootfs = d.ignoreRootFS.IsChecked()

	return opts
}

func (d *ContainerRestoreDialog) setFocusElement() {
	if d.focusElement < cntRestoreFormFocus {
		d.focusElement++

		return
	}

	d.focusElement = cntRestoreContainersFocus
}

func (d *ContainerRestoreDialog) getInnerPrimitives() []tview.Primitive {
	// the item sort is important to be same as focus element number
	return []tview.Primitive{
		d.containers,
		d.pods,
		d.name,
		d.publishPorts,
		d.importArchive,
		d.keep,
		d.ignoreStaticIP,
		d.ignoreStaticMAC,
		d.fileLocks,
		d.printStats,
		d.tcpEstablished,
		d.ignoreVolumes,
		d.ignoreRootFS,
		d.form,
	}
}

// UpdateLanguage updates all labels to current language.
func (d *ContainerRestoreDialog) UpdateLanguage() {
	bgColor := style.DialogBgColor
	
	// Update window title
	d.layout.SetTitle(i18n.T("PODMAN CONTAINER RESTORE"))

	// Update container ID label
	containersLabel := fmt.Sprintf("[:#%x:b]%s[:-:-]", style.DialogBorderColor.Hex(), i18n.T("CONTAINER ID:"))
	d.containers.SetLabel(containersLabel)

	// Update field labels
	d.pods.SetLabel(i18n.T("pod:"))
	d.name.SetLabel(i18n.PadToWidth(i18n.T("name:"), cntRestoreDialogLabelWidth))
	d.publishPorts.SetLabel(i18n.PadToWidth(i18n.T("publish:"), cntRestoreDialogLabelWidth))
	d.importArchive.SetLabel(i18n.PadToWidth(i18n.T("import:"), cntRestoreDialogLabelWidth))
	
	// Calculate dynamic checkbox column widths (minimal padding)
	col1Width := i18n.CalcMaxWidth(
		i18n.T("keep:"),
		i18n.T("print Stats:"),
	)
	col2Width := i18n.CalcMaxWidth(
		i18n.T("ignore static IP:"),
		i18n.T("tcp established:"),
	)
	col3Width := i18n.CalcMaxWidth(
		i18n.T("ignore static MAC:"),
		i18n.T("ignore volumes:"),
	)
	col4Width := i18n.CalcMaxWidth(
		i18n.T("file locks:"),
		i18n.T("ignore rootfs:"),
	)
	
	// Update checkbox labels
	d.keep.SetLabel(i18n.T("keep:"))
	d.keep.SetLabelWidth(col1Width)
	
	d.ignoreStaticIP.SetLabel(i18n.T("ignore static IP:"))
	d.ignoreStaticIP.SetLabelWidth(col2Width)
	
	d.ignoreStaticMAC.SetLabel(i18n.T("ignore static MAC:"))
	d.ignoreStaticMAC.SetLabelWidth(col3Width)
	
	d.fileLocks.SetLabel(i18n.T("file locks:"))
	d.fileLocks.SetLabelWidth(col4Width)
	
	d.printStats.SetLabel(i18n.T("print Stats:"))
	d.printStats.SetLabelWidth(col1Width)
	
	d.tcpEstablished.SetLabel(i18n.T("tcp established:"))
	d.tcpEstablished.SetLabelWidth(col2Width)
	
	d.ignoreVolumes.SetLabel(i18n.T("ignore volumes:"))
	d.ignoreVolumes.SetLabelWidth(col3Width)
	
	d.ignoreRootFS.SetLabel(i18n.T("ignore rootfs:"))
	d.ignoreRootFS.SetLabelWidth(col4Width)
	
	// Rebuild row5 and row6 with new column widths
	d.row5.Clear()
	d.row5.AddItem(d.keep, col1Width+3, 0, true)
	d.row5.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	d.row5.AddItem(d.ignoreStaticIP, col2Width+3, 0, true)
	d.row5.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	d.row5.AddItem(d.ignoreStaticMAC, col3Width+3, 0, true)
	d.row5.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	d.row5.AddItem(d.fileLocks, col4Width+3, 0, true)
	d.row5.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	
	d.row6.Clear()
	d.row6.AddItem(d.printStats, col1Width+3, 0, true)
	d.row6.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	d.row6.AddItem(d.tcpEstablished, col2Width+3, 0, true)
	d.row6.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	d.row6.AddItem(d.ignoreVolumes, col3Width+3, 0, true)
	d.row6.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	d.row6.AddItem(d.ignoreRootFS, col4Width+3, 0, true)
	d.row6.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	
	// Update form buttons
	d.form.ClearButtons()
	d.form.AddButton(i18n.T("Cancel"), d.cancelHandler)
	d.form.AddButton(i18n.T("Restore"), d.restoreHandler)
}
