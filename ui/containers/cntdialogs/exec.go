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
	// maxheight = button height + total input widgets row + 11.
	execDialogMaxHeight    = dialogs.DialogFormHeight + 7 + 9
	execDialogMaxWidth     = 80
	execDialogLabelWidth   = 14
	execDialogLabelPadding = 1
)

const (
	execCommandFieldFocus = 0 + iota
	execInteractiveFieldFocus
	execTtyFieldFocus
	execPrivilegedFieldFocus
	execWorkingDirFieldFocus
	execEnvVariablesFieldFocus
	execEnvFileFieldFocus
	execUserFieldFocus
	execDetachFieldFocus
	execFormFieldFocus
)

// ContainerExecDialog represents container exec dialog primitive.
type ContainerExecDialog struct {
	*tview.Box

	layout        *tview.Flex
	cntInfo       *tview.InputField
	command       *tview.InputField
	interactive   *tview.Checkbox
	tty           *tview.Checkbox
	privileged    *tview.Checkbox
	detach        *tview.Checkbox
	workingDir    *tview.InputField
	envVariables  *tview.InputField
	envFile       *tview.InputField
	user          *tview.InputField
	form          *tview.Form
	display       bool
	containerID   string
	focusElement  int
	execHandler   func()
	cancelHandler func()
}

// NewContainerExecDialog returns new container exec dialog.
func NewContainerExecDialog() *ContainerExecDialog {
	dialog := &ContainerExecDialog{
		Box:          tview.NewBox(),
		cntInfo:      tview.NewInputField(),
		command:      tview.NewInputField(),
		tty:          tview.NewCheckbox(),
		interactive:  tview.NewCheckbox(),
		privileged:   tview.NewCheckbox(),
		detach:       tview.NewCheckbox(),
		workingDir:   tview.NewInputField(),
		envVariables: tview.NewInputField(),
		envFile:      tview.NewInputField(),
		user:         tview.NewInputField(),
		display:      false,
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor

	// label (container ID and Name)
	dialog.cntInfo.SetBackgroundColor(style.DialogBgColor)
	dialog.cntInfo.SetLabel("[::b]" + i18n.T("CONTAINER ID:"))
	dialog.cntInfo.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.cntInfo.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// Calculate label width
	labelWidth := i18n.CalcMaxWidth(
		i18n.T("command:"),
		i18n.T("interactive:"),
		i18n.T("user:"),
		i18n.T("working dir:"),
		i18n.T("env vars:"),
		i18n.T("env file:"),
	) + 1

	// command
	dialog.command.SetBackgroundColor(style.DialogBgColor)
	dialog.command.SetLabel(i18n.PadToWidth(i18n.T("command:"), labelWidth))
	dialog.command.SetFieldStyle(style.InputFieldStyle)
	dialog.command.SetLabelStyle(style.InputLabelStyle)

	// interactive
	dialog.interactive.SetBackgroundColor(bgColor)
	dialog.interactive.SetBorder(false)
	dialog.interactive.SetLabel(i18n.T("interactive:"))
	dialog.interactive.SetLabelColor(fgColor)
	dialog.interactive.SetLabelWidth(labelWidth)
	dialog.interactive.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// tty
	ttyLabelWidth := i18n.GetDisplayWidth(i18n.T("tty:")) + 1

	dialog.tty.SetBackgroundColor(bgColor)
	dialog.tty.SetBorder(false)
	dialog.tty.SetLabel(i18n.T("tty:"))
	dialog.tty.SetLabelColor(fgColor)
	dialog.tty.SetLabelWidth(ttyLabelWidth)
	dialog.tty.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// privileged
	privilegedLabelWidth := i18n.GetDisplayWidth(i18n.T("privileged:")) + 1

	dialog.privileged.SetBackgroundColor(bgColor)
	dialog.privileged.SetBorder(false)
	dialog.privileged.SetLabel(i18n.T("privileged:"))
	dialog.privileged.SetLabelColor(fgColor)
	dialog.privileged.SetLabelWidth(privilegedLabelWidth)
	dialog.privileged.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// detach
	detachLabelWidth := i18n.GetDisplayWidth(i18n.T("detach:")) + 1

	dialog.detach.SetBackgroundColor(bgColor)
	dialog.detach.SetBorder(false)
	dialog.detach.SetLabel(i18n.T("detach:"))
	dialog.detach.SetLabelColor(fgColor)
	dialog.detach.SetLabelWidth(detachLabelWidth)
	dialog.detach.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// working dir
	dialog.workingDir.SetBackgroundColor(style.DialogBgColor)
	dialog.workingDir.SetLabel(i18n.PadToWidth(i18n.T("working dir:"), labelWidth))
	dialog.workingDir.SetFieldStyle(style.InputFieldStyle)
	dialog.workingDir.SetLabelStyle(style.InputLabelStyle)

	// env variables
	dialog.envVariables.SetBackgroundColor(style.DialogBgColor)
	dialog.envVariables.SetLabel(i18n.PadToWidth(i18n.T("env vars:"), labelWidth))
	dialog.envVariables.SetFieldStyle(style.InputFieldStyle)
	dialog.envVariables.SetLabelStyle(style.InputLabelStyle)

	// env file
	dialog.envFile.SetBackgroundColor(style.DialogBgColor)
	dialog.envFile.SetLabel(i18n.PadToWidth(i18n.T("env file:"), labelWidth))
	dialog.envFile.SetFieldStyle(style.InputFieldStyle)
	dialog.envFile.SetLabelStyle(style.InputLabelStyle)

	// user
	dialog.user.SetBackgroundColor(style.DialogBgColor)
	dialog.user.SetLabel(i18n.PadToWidth(i18n.T("user:"), labelWidth))
	dialog.user.SetFieldStyle(style.InputFieldStyle)
	dialog.user.SetLabelStyle(style.InputLabelStyle)

	// form fields
	dialog.form = tview.NewForm().
		AddButton(i18n.T("Cancel"), nil).
		AddButton(i18n.T("Execute"), nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// main dialog layout
	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetTitle(i18n.T("PODMAN CONTAINER EXEC"))

	mLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	// label
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.cntInfo, 1, 0, true)
	// command
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.command, 1, 0, true)

	// interactive, tty, privileged and detach
	checkBoxWidth := labelWidth + 4 //nolint:mnd
	labelPaddings := 5
	checkBoxLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	checkBoxLayout.SetBackgroundColor(bgColor)
	checkBoxLayout.AddItem(dialog.interactive, checkBoxWidth, 0, false)
	checkBoxLayout.AddItem(dialog.tty, ttyLabelWidth+labelPaddings, 0, false)
	checkBoxLayout.AddItem(dialog.privileged, privilegedLabelWidth+labelPaddings, 0, false)
	checkBoxLayout.AddItem(dialog.detach, detachLabelWidth+labelPaddings, 0, false)
	checkBoxLayout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(checkBoxLayout, 1, 0, true)

	// user
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.user, 1, 0, true)
	// working dir
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.workingDir, 1, 0, true)
	// env variables
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.envVariables, 1, 0, true)
	// env file
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.envFile, 1, 0, true)

	// main layout
	mainLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mainLayout.AddItem(mLayout, 0, 1, true)
	mainLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)

	dialog.layout.AddItem(mainLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive.
func (d *ContainerExecDialog) Display() {
	d.focusElement = execCommandFieldFocus

	d.command.SetText("")
	d.tty.SetChecked(true)
	d.interactive.SetChecked(false)
	d.privileged.SetChecked(false)
	d.detach.SetChecked(false)
	d.workingDir.SetText("")
	d.envVariables.SetText("")
	d.envFile.SetText("")
	d.user.SetText("")

	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *ContainerExecDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ContainerExecDialog) Hide() {
	d.SetContainerID("", "")
	d.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (d *ContainerExecDialog) HasFocus() bool { //nolint:cyclop
	if d.command.HasFocus() || d.tty.HasFocus() {
		return true
	}

	if d.interactive.HasFocus() || d.privileged.HasFocus() {
		return true
	}

	if d.workingDir.HasFocus() || d.envVariables.HasFocus() {
		return true
	}

	if d.envFile.HasFocus() || d.user.HasFocus() {
		return true
	}

	if d.detach.HasFocus() || d.form.HasFocus() {
		return true
	}

	return d.Box.HasFocus() || d.layout.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *ContainerExecDialog) Focus(delegate func(p tview.Primitive)) { //nolint:gocognit,cyclop
	switch d.focusElement {
	// command field focus
	case execCommandFieldFocus:
		d.command.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execInteractiveFieldFocus
				d.Focus(delegate)

				return nil
			}

			if event.Key() == tcell.KeyEnter {
				d.execHandler()

				return nil
			}

			return event
		})

		delegate(d.command)

		return
	// interactive field focus
	case execInteractiveFieldFocus:
		d.interactive.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execTtyFieldFocus

				d.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(d.interactive)

		return
	// tty field focus
	case execTtyFieldFocus:
		d.tty.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execPrivilegedFieldFocus

				d.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(d.tty)

		return
	// privileged field focus
	case execPrivilegedFieldFocus:
		d.privileged.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execDetachFieldFocus

				d.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(d.privileged)

		return
	// detach field focus
	case execDetachFieldFocus:
		d.detach.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execUserFieldFocus

				d.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(d.detach)

		return
	// user field focus
	case execUserFieldFocus:
		d.user.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execWorkingDirFieldFocus

				d.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(d.user)

		return
	// working directory field focus
	case execWorkingDirFieldFocus:
		d.workingDir.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execEnvVariablesFieldFocus

				d.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(d.workingDir)

		return
	// env variable field focus
	case execEnvVariablesFieldFocus:
		d.envVariables.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execEnvFileFieldFocus

				d.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(d.envVariables)

		return
	// env file field focus
	case execEnvFileFieldFocus:
		d.envFile.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execFormFieldFocus

				d.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(d.envFile)

		return
	// form field focus
	case execFormFieldFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execCommandFieldFocus

				d.Focus(delegate)
				d.form.SetFocus(execCommandFieldFocus)

				return nil
			}

			if event.Key() == tcell.KeyEnter {
				d.execHandler()

				return nil
			}

			return event
		})

		delegate(d.form)
	}
}

// InputHandler returns input handler function for this primitive.
func (d *ContainerExecDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container exec dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()

			return
		}

		// command field
		if d.command.HasFocus() {
			if commandHandler := d.command.InputHandler(); commandHandler != nil {
				commandHandler(event, setFocus)

				return
			}
		}

		// interactive field
		if d.interactive.HasFocus() {
			if interactiveHandler := d.interactive.InputHandler(); interactiveHandler != nil {
				interactiveHandler(event, setFocus)

				return
			}
		}

		// privileged field
		if d.privileged.HasFocus() {
			if privilegedHandler := d.privileged.InputHandler(); privilegedHandler != nil {
				privilegedHandler(event, setFocus)

				return
			}
		}

		// tty field
		if d.tty.HasFocus() {
			if ttyHandler := d.tty.InputHandler(); ttyHandler != nil {
				ttyHandler(event, setFocus)

				return
			}
		}

		// detach field
		if d.detach.HasFocus() {
			if detachHandler := d.detach.InputHandler(); detachHandler != nil {
				detachHandler(event, setFocus)

				return
			}
		}

		// working directory field
		if d.workingDir.HasFocus() {
			if workingDirHandler := d.workingDir.InputHandler(); workingDirHandler != nil {
				workingDirHandler(event, setFocus)

				return
			}
		}

		// env variables field
		if d.envVariables.HasFocus() {
			if envVariablesHandler := d.envVariables.InputHandler(); envVariablesHandler != nil {
				envVariablesHandler(event, setFocus)

				return
			}
		}

		// env file field
		if d.envFile.HasFocus() {
			if envFileHandler := d.envFile.InputHandler(); envFileHandler != nil {
				envFileHandler(event, setFocus)

				return
			}
		}

		// user field
		if d.user.HasFocus() {
			if userHandler := d.user.InputHandler(); userHandler != nil {
				userHandler(event, setFocus)

				return
			}
		}

		// form primitive
		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)

				return
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *ContainerExecDialog) SetRect(x, y, width, height int) {
	dWidth := width
	dX := x

	if width > execDialogMaxWidth {
		wEmptySpace := width - execDialogMaxWidth
		if wEmptySpace > 0 {
			dX = x + (wEmptySpace / 2) //nolint:mnd
		}

		dWidth = execDialogMaxWidth
	}

	dHeight := height
	dY := y

	if height > execDialogMaxHeight {
		hEmptySpace := height - execDialogMaxHeight

		if hEmptySpace > 0 {
			dY = y + (hEmptySpace / 2) //nolint:mnd
		}

		dHeight = execDialogMaxHeight
	}

	d.Box.SetRect(dX, dY, dWidth, dHeight)
}

// Draw draws this primitive onto the screen.
func (d *ContainerExecDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)

	x, y, width, height := d.GetInnerRect()

	d.layout.SetRect(x, y, width, height)

	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function.
func (d *ContainerExecDialog) SetCancelFunc(handler func()) *ContainerExecDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd

	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetExecFunc sets form execute button selected function.
func (d *ContainerExecDialog) SetExecFunc(handler func()) *ContainerExecDialog {
	d.execHandler = handler
	execButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	execButton.SetSelectedFunc(handler)

	return d
}

// SetContainerID sets container ID label.
func (d *ContainerExecDialog) SetContainerID(id string, name string) {
	d.containerID = id
	containerInfo := fmt.Sprintf("%s (%s)", id, name)
	containerInfo = utils.LabelWidthLeftPadding(containerInfo, execDialogLabelPadding)

	d.cntInfo.SetText(containerInfo)
}

// ContainerExecOptions returns new container exec options.
func (d *ContainerExecDialog) ContainerExecOptions() containers.ExecOption {
	execOptions := containers.ExecOption{}
	cmdString := strings.TrimSpace(d.command.GetText())
	execOptions.Cmd = strings.Split(cmdString, " ")
	execOptions.Tty = d.tty.IsChecked()
	execOptions.Interactive = d.interactive.IsChecked()
	execOptions.Detach = d.detach.IsChecked()
	execOptions.Privileged = d.privileged.IsChecked()
	execOptions.WorkDir = strings.TrimSpace(d.workingDir.GetText())

	varString := strings.TrimSpace(d.envVariables.GetText())
	if varString != "" {
		execOptions.EnvVariables = strings.Split(varString, " ")
	}

	envFileString := strings.TrimSpace(d.envFile.GetText())
	if envFileString != "" {
		execOptions.EnvFile = strings.Split(envFileString, " ")
	}

	execOptions.User = strings.TrimSpace(d.user.GetText())

	return execOptions
}

// UpdateLanguage updates all labels to current language.
func (d *ContainerExecDialog) UpdateLanguage() {
	// Calculate label width
	labelWidth := i18n.CalcMaxWidth(
		i18n.T("command:"),
		i18n.T("interactive:"),
		i18n.T("user:"),
		i18n.T("working dir:"),
		i18n.T("env vars:"),
		i18n.T("env file:"),
	) + 1

	// Update window title
	d.layout.SetTitle(i18n.T("PODMAN CONTAINER EXEC"))

	// Update container ID label
	d.cntInfo.SetLabel("[::b]" + i18n.T("CONTAINER ID:"))

	// Update field labels
	d.command.SetLabel(i18n.PadToWidth(i18n.T("command:"), labelWidth))
	d.interactive.SetLabel(i18n.T("interactive:"))
	d.interactive.SetLabelWidth(labelWidth)

	ttyLabelWidth := i18n.GetDisplayWidth(i18n.T("tty:")) + 1
	d.tty.SetLabel(i18n.T("tty:"))
	d.tty.SetLabelWidth(ttyLabelWidth)

	privilegedLabelWidth := i18n.GetDisplayWidth(i18n.T("privileged:")) + 1
	d.privileged.SetLabel(i18n.T("privileged:"))
	d.privileged.SetLabelWidth(privilegedLabelWidth)

	detachLabelWidth := i18n.GetDisplayWidth(i18n.T("detach:")) + 1
	d.detach.SetLabel(i18n.T("detach:"))
	d.detach.SetLabelWidth(detachLabelWidth)

	d.user.SetLabel(i18n.PadToWidth(i18n.T("user:"), labelWidth))
	d.workingDir.SetLabel(i18n.PadToWidth(i18n.T("working dir:"), labelWidth))
	d.envVariables.SetLabel(i18n.PadToWidth(i18n.T("env vars:"), labelWidth))
	d.envFile.SetLabel(i18n.PadToWidth(i18n.T("env file:"), labelWidth))

	// Update form buttons
	d.form.ClearButtons()
	d.form.AddButton(i18n.T("Cancel"), d.cancelHandler)
	d.form.AddButton(i18n.T("Execute"), d.execHandler)
}
