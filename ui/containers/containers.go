package containers

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/ui/containers/cntdialogs"
	"github.com/containers/podman-tui/ui/containers/cntdialogs/vterm"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/rivo/tview"
)

const (
	viewContainersIDColIndex = 0 + iota
	viewContainersImageColIndex
	viewContainersPodColIndex
	viewContainersCreatedAtColIndex
	viewContainersStatusColIndex
	viewContainersNamesColIndex
	viewContainersPortsColIndex
)

var (
	errNoContainerAttach       = errors.New("there is no container to attach")
	errNoContainerHealthCheck  = errors.New("there is no container to perform healthcheck")
	errNoContainerCommit       = errors.New("there is no container to commit")
	errNoContainerStat         = errors.New("there is no container to display stats")
	errNoContainerCheckpoint   = errors.New("there is no container to perform checkpoint")
	errNoContainerExec         = errors.New("there is no container to perform exec")
	errNoContainerDiff         = errors.New("there is no container to display diff")
	errNoContainerInspect      = errors.New("there is no container to inspect")
	errNoContainerKill         = errors.New("there is no container to kill")
	errNoContainerLogs         = errors.New("there is no container to display logs")
	errNoContainerPause        = errors.New("there is no container to pause")
	errNoContainerUnpause      = errors.New("there is no container to unpause")
	errNoContainerPorts        = errors.New("there is no container to display ports")
	errNoContainerRename       = errors.New("there is no container to rename")
	errNoContainerRemove       = errors.New("there is no container to remove")
	errNoContainerStart        = errors.New("there is no container to start")
	errNoContainerStop         = errors.New("there is no container to stop")
	errNoContainerTop          = errors.New("there is no container to display top")
	errEmptyContainerImageName = errors.New("empty container image name")
)

// getTranslatedError returns the translated error message
func getTranslatedError(err error) string {
	switch err {
	case errNoContainerAttach:
		return i18n.T("there is no container to attach")
	case errNoContainerHealthCheck:
		return i18n.T("there is no container to perform healthcheck")
	case errNoContainerCommit:
		return i18n.T("there is no container to commit")
	case errNoContainerStat:
		return i18n.T("there is no container to display stats")
	case errNoContainerCheckpoint:
		return i18n.T("there is no container to perform checkpoint")
	case errNoContainerExec:
		return i18n.T("there is no container to perform exec")
	case errNoContainerDiff:
		return i18n.T("there is no container to display diff")
	case errNoContainerInspect:
		return i18n.T("there is no container to inspect")
	case errNoContainerKill:
		return i18n.T("there is no container to kill")
	case errNoContainerLogs:
		return i18n.T("there is no container to display logs")
	case errNoContainerPause:
		return i18n.T("there is no container to pause")
	case errNoContainerUnpause:
		return i18n.T("there is no container to unpause")
	case errNoContainerPorts:
		return i18n.T("there is no container to display ports")
	case errNoContainerRename:
		return i18n.T("there is no container to rename")
	case errNoContainerRemove:
		return i18n.T("there is no container to remove")
	case errNoContainerStart:
		return i18n.T("there is no container to start")
	case errNoContainerStop:
		return i18n.T("there is no container to stop")
	case errNoContainerTop:
		return i18n.T("there is no container to display top")
	case errEmptyContainerImageName:
		return i18n.T("empty container image name")
	default:
		// Keep Podman's dynamic error messages as-is (no translation)
		return err.Error()
	}
}

// translateErrorTitle translates error title like "CONTAINER (id) ACTION ERROR"
func translateErrorTitle(title string) string {
	if title == "" {
		return ""
	}

	// Check if title contains a container ID in format "CONTAINER (id) ... ERROR"
	if strings.HasPrefix(title, "CONTAINER (") && strings.Contains(title, ") ") {
		// Extract the container ID
		idStart := len("CONTAINER (")
		idEnd := strings.Index(title, ") ")
		if idEnd > idStart {
			containerID := title[idStart:idEnd]
			// Get the error action part (after ") ")
			actionPart := title[idEnd+2:]
			
			// Build the template key: "CONTAINER (%s) ACTION ERROR"
			templateKey := fmt.Sprintf("CONTAINER (%%s) %s", actionPart)
			
			// Translate the template
			translatedTemplate := i18n.T(templateKey)
			
			// Format with the actual container ID
			return fmt.Sprintf(translatedTemplate, containerID)
		}
	}
	
	// For titles without container ID, translate directly
	return i18n.T(title)
}

// Containers implements the containers page primitive.
type Containers struct {
	*tview.Box

	title            string
	headers          []string
	table            *tview.Table
	errorDialog      *dialogs.ErrorDialog
	cmdDialog        *dialogs.CommandDialog
	cmdInputDialog   *dialogs.SimpleInputDialog
	confirmDialog    *dialogs.ConfirmDialog
	messageDialog    *dialogs.MessageDialog
	progressDialog   *dialogs.ProgressDialog
	sortDialog       *dialogs.SortDialog
	topDialog        *dialogs.TopDialog
	createDialog     *cntdialogs.ContainerCreateDialog
	runDialog        *cntdialogs.ContainerCreateDialog
	execDialog       *cntdialogs.ContainerExecDialog
	statsDialog      *cntdialogs.ContainerStatsDialog
	commitDialog     *cntdialogs.ContainerCommitDialog
	checkpointDialog *cntdialogs.ContainerCheckpointDialog
	restoreDialog    *cntdialogs.ContainerRestoreDialog
	terminalDialog   *vterm.VtermDialog
	containersList   containerListReport
	selectedID       string
	selectedName     string
	confirmData      string
	fastRefreshChan  chan bool
	appFocusHandler  func()
}

type containerListReport struct {
	mu        sync.Mutex
	report    []entities.ListContainer
	sortBy    string
	ascending bool
}

// NewContainers returns containers page view.
func NewContainers() *Containers {
	containers := &Containers{
		Box:              tview.NewBox(),
		title:            "containers",
		headers:          []string{i18n.T("container id"), i18n.T("image"), i18n.T("pod"), i18n.T("created"), i18n.T("status"), i18n.T("names"), i18n.T("ports")},
		errorDialog:      dialogs.NewErrorDialog(),
		cmdInputDialog:   dialogs.NewSimpleInputDialog(""),
		messageDialog:    dialogs.NewMessageDialog(""),
		progressDialog:   dialogs.NewProgressDialog(),
		confirmDialog:    dialogs.NewConfirmDialog(),
		topDialog:        dialogs.NewTopDialog(),
		sortDialog:       dialogs.NewSortDialog([]string{"name", "pod", "image", "created", "status"}, 3), //nolint:mnd
		createDialog:     cntdialogs.NewContainerCreateDialog(cntdialogs.ContainerCreateOnlyDialogMode),
		runDialog:        cntdialogs.NewContainerCreateDialog(cntdialogs.ContainerCreateAndRunDialogMode),
		execDialog:       cntdialogs.NewContainerExecDialog(),
		statsDialog:      cntdialogs.NewContainerStatsDialog(),
		commitDialog:     cntdialogs.NewContainerCommitDialog(),
		checkpointDialog: cntdialogs.NewContainerCheckpointDialog(),
		restoreDialog:    cntdialogs.NewContainerRestoreDialog(),
		terminalDialog:   vterm.NewVtermDialog(),
		containersList:   containerListReport{sortBy: "created", ascending: true},
	}

	containers.topDialog.SetTitle(i18n.T("podman container top"))

	// Build command dialog with translations
	containers.buildCommandDialog()

	containers.table = tview.NewTable()
	containers.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(containers.title)))
	containers.table.SetBorderColor(style.BorderColor)
	containers.table.SetTitleColor(style.FgColor)
	containers.table.SetBackgroundColor(style.BgColor)
	containers.table.SetBorder(true)

	for i := range containers.headers {
		containers.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(containers.headers[i]))). //nolint:perfsprint
															SetExpansion(1).
															SetBackgroundColor(style.PageHeaderBgColor).
															SetTextColor(style.PageHeaderFgColor).
															SetAlign(tview.AlignLeft).
															SetSelectable(false))
	}

	containers.table.SetFixed(1, 1)
	containers.table.SetSelectable(true, false)

	// set command dialog functions
	// NOTE: cmdDialog handlers (SetSelectedFunc, SetCancelFunc) are set in buildCommandDialog()

	// set input cmd dialog functions
	containers.cmdInputDialog.SetCancelFunc(containers.cmdInputDialog.Hide)
	containers.cmdInputDialog.SetSelectedFunc(containers.cmdInputDialog.Hide)

	// set message dialog functions
	containers.messageDialog.SetCancelFunc(containers.messageDialog.Hide)

	// set container top dialog functions
	containers.topDialog.SetCancelFunc(containers.topDialog.Hide)

	// set confirm dialogs functions
	containers.confirmDialog.SetSelectedFunc(func() {
		containers.confirmDialog.Hide()

		switch containers.confirmData {
		case utils.PruneCommandLabel:
			containers.prune()
		case "rm":
			containers.remove()
		}
	})

	containers.confirmDialog.SetCancelFunc(containers.confirmDialog.Hide)

	// set create dialog functions
	containers.createDialog.SetCancelFunc(func() {
		containers.createDialog.Hide()
	})

	containers.createDialog.SetHandlerFunc(func() {
		containers.createDialog.Hide()
		containers.create()
	})

	// set run dialog functions
	containers.runDialog.SetCancelFunc(func() {
		containers.runDialog.Hide()
	})

	containers.runDialog.SetHandlerFunc(func() {
		containers.runDialog.Hide()
		containers.run()
	})

	// set exec dialog functions
	containers.execDialog.SetCancelFunc(containers.execDialog.Hide)
	containers.execDialog.SetExecFunc(containers.exec)

	// terminal dialog
	containers.terminalDialog.SetCancelFunc(containers.terminalDialog.Hide)
	containers.terminalDialog.SetFastRefreshHandler(func() {
		containers.fastRefreshChan <- true
	})

	// set stats dialogs functions
	containers.statsDialog.SetDoneFunc(containers.statsDialog.Hide)

	// set commit dialog functions
	containers.commitDialog.SetCommitFunc(containers.commit)
	containers.commitDialog.SetCancelFunc(containers.commitDialog.Hide)

	// set checkpoint dialog functions
	containers.checkpointDialog.SetCheckpointFunc(containers.checkpoint)
	containers.checkpointDialog.SetCancelFunc(containers.checkpointDialog.Hide)

	// set restore dialog functions
	containers.restoreDialog.SetRestoreFunc(containers.restore)
	containers.restoreDialog.SetCancelFunc(containers.restoreDialog.Hide)

	// set sort dialog functions
	containers.sortDialog.SetSelectFunc(containers.SortView)
	containers.sortDialog.SetCancelFunc(containers.sortDialog.Hide)

	return containers
}

// SetAppFocusHandler sets application focus handler.
func (cnt *Containers) SetAppFocusHandler(handler func()) {
	cnt.appFocusHandler = handler
}

// GetTitle returns primitive title.
func (cnt *Containers) GetTitle() string {
	return cnt.title
}

// HasFocus returns whether or not this primitive has focus.
func (cnt *Containers) HasFocus() bool { //nolint:cyclop
	if cnt.table.HasFocus() || cnt.errorDialog.HasFocus() {
		return true
	}

	if cnt.cmdDialog.HasFocus() || cnt.progressDialog.HasFocus() {
		return true
	}

	if cnt.topDialog.HasFocus() || cnt.messageDialog.HasFocus() {
		return true
	}

	if cnt.confirmDialog.HasFocus() || cnt.cmdInputDialog.HasFocus() {
		return true
	}

	if cnt.createDialog.HasFocus() || cnt.execDialog.HasFocus() {
		return true
	}

	if cnt.statsDialog.HasFocus() || cnt.commitDialog.HasFocus() {
		return true
	}

	if cnt.checkpointDialog.HasFocus() || cnt.restoreDialog.HasFocus() {
		return true
	}

	if cnt.runDialog.HasFocus() || cnt.terminalDialog.HasFocus() {
		return true
	}

	if cnt.sortDialog.HasFocus() || cnt.Box.HasFocus() {
		return true
	}

	return false
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus.
func (cnt *Containers) SubDialogHasFocus() bool { //nolint:cyclop
	if cnt.statsDialog.HasFocus() || cnt.errorDialog.HasFocus() {
		return true
	}

	if cnt.cmdDialog.HasFocus() || cnt.progressDialog.HasFocus() {
		return true
	}

	if cnt.topDialog.HasFocus() || cnt.messageDialog.HasFocus() {
		return true
	}

	if cnt.confirmDialog.HasFocus() || cnt.cmdInputDialog.HasFocus() {
		return true
	}

	if cnt.createDialog.HasFocus() || cnt.execDialog.HasFocus() {
		return true
	}

	if cnt.commitDialog.HasFocus() || cnt.checkpointDialog.HasFocus() {
		return true
	}

	if cnt.restoreDialog.HasFocus() || cnt.terminalDialog.HasFocus() {
		return true
	}

	if cnt.sortDialog.HasFocus() || cnt.runDialog.HasFocus() {
		return true
	}

	return false
}

// Focus is called when this primitive receives focus.
func (cnt *Containers) Focus(delegate func(p tview.Primitive)) { //nolint:cyclop
	// error dialog
	if cnt.errorDialog.IsDisplay() {
		delegate(cnt.errorDialog)

		return
	}

	// command dialog
	if cnt.cmdDialog.IsDisplay() {
		delegate(cnt.cmdDialog)

		return
	}

	// command input dialog
	if cnt.cmdInputDialog.IsDisplay() {
		delegate(cnt.cmdInputDialog)

		return
	}

	// message dialog
	if cnt.messageDialog.IsDisplay() {
		delegate(cnt.messageDialog)

		return
	}

	// container top dialog
	if cnt.topDialog.IsDisplay() {
		delegate(cnt.topDialog)

		return
	}

	// confirm dialog
	if cnt.confirmDialog.IsDisplay() {
		delegate(cnt.confirmDialog)

		return
	}

	// create dialog
	if cnt.createDialog.IsDisplay() {
		delegate(cnt.createDialog)

		return
	}

	// run dialog
	if cnt.runDialog.IsDisplay() {
		delegate(cnt.runDialog)

		return
	}

	// exec dialog
	if cnt.execDialog.IsDisplay() {
		delegate(cnt.execDialog)

		return
	}

	// stats dialog
	if cnt.statsDialog.IsDisplay() {
		delegate(cnt.statsDialog)

		return
	}

	// commit dialog
	if cnt.commitDialog.IsDisplay() {
		delegate(cnt.commitDialog)

		return
	}

	// checkpoint dialog
	if cnt.checkpointDialog.IsDisplay() {
		delegate(cnt.checkpointDialog)

		return
	}

	// restore dialog
	if cnt.restoreDialog.IsDisplay() {
		delegate(cnt.restoreDialog)

		return
	}

	// terminal dialog
	if cnt.terminalDialog.IsDisplay() {
		delegate(cnt.terminalDialog)

		return
	}

	// sort dialog
	if cnt.sortDialog.IsDisplay() {
		delegate(cnt.sortDialog)

		return
	}

	delegate(cnt.table)
}

// SetFastRefreshChannel sets channel for fastRefresh func.
func (cnt *Containers) SetFastRefreshChannel(refresh chan bool) {
	cnt.fastRefreshChan = refresh
}

// HideAllDialogs hides all sub dialogs.
func (cnt *Containers) HideAllDialogs() { //nolint:cyclop
	if cnt.errorDialog.IsDisplay() {
		cnt.errorDialog.Hide()
	}

	if cnt.progressDialog.IsDisplay() {
		cnt.progressDialog.Hide()
	}

	if cnt.confirmDialog.IsDisplay() {
		cnt.confirmDialog.Hide()
	}

	if cnt.cmdDialog.IsDisplay() {
		cnt.cmdDialog.Hide()
	}

	if cnt.cmdInputDialog.IsDisplay() {
		cnt.cmdInputDialog.Hide()
	}

	if cnt.messageDialog.IsDisplay() {
		cnt.messageDialog.Hide()
	}

	if cnt.topDialog.IsDisplay() {
		cnt.topDialog.Hide()
	}

	if cnt.createDialog.IsDisplay() {
		cnt.createDialog.Hide()
	}

	if cnt.runDialog.IsDisplay() {
		cnt.runDialog.Hide()
	}

	if cnt.execDialog.IsDisplay() {
		cnt.execDialog.Hide()
	}

	if cnt.statsDialog.IsDisplay() {
		cnt.statsDialog.Hide()
	}

	if cnt.commitDialog.IsDisplay() {
		cnt.commitDialog.Hide()
	}

	if cnt.checkpointDialog.IsDisplay() {
		cnt.checkpointDialog.Hide()
	}

	if cnt.restoreDialog.IsDisplay() {
		cnt.restoreDialog.Hide()
	}

	if cnt.terminalDialog.IsDisplay() {
		cnt.terminalDialog.Hide()
	}

	if cnt.sortDialog.IsDisplay() {
		cnt.sortDialog.Hide()
	}
}

func (cnt *Containers) getSelectedItem() (string, string) {
	var (
		cntID   string
		cntName string
	)

	if cnt.table.GetRowCount() <= 1 {
		return cntID, cntName
	}

	row, _ := cnt.table.GetSelection()
	cntID = cnt.table.GetCell(row, viewContainersIDColIndex).Text
	cntName = cnt.table.GetCell(row, viewContainersNamesColIndex).Text

	return cntID, cntName
}

// buildCommandDialog builds the command dialog with translated strings
func (cnt *Containers) buildCommandDialog() {
	if cnt.cmdDialog != nil {
		cnt.cmdDialog.Hide()
	}

	// Translate both command names and descriptions
	cnt.cmdDialog = dialogs.NewCommandDialog([][]string{
		{i18n.T("attach"), i18n.T("attach to a running container")},
		{i18n.T("checkpoint"), i18n.T("checkpoints a running container")},
		{i18n.T("commit"), i18n.T("create an image from a container's changes")},
		{i18n.T("create"), i18n.T("create a new container but do not start")},
		{i18n.T("diff"), i18n.T("inspect changes to the selected container's file systems")},
		{i18n.T("exec"), i18n.T("execute the specified command inside a running container")},
		{i18n.T("healthcheck"), i18n.T("run the health check of a container")},
		{i18n.T("inspect"), i18n.T("display the configuration of a container")},
		{i18n.T("kill"), i18n.T("kill the selected running container with a SIGKILL signal")},
		{i18n.T("logs"), i18n.T("fetch the logs of the selected container")},
		{i18n.T("pause"), i18n.T("pause all the processes in the selected container")},
		{i18n.T("port"), i18n.T("list port mappings for the selected container")},
		{i18n.T("prune"), i18n.T("remove all non running containers")},
		{i18n.T("rename"), i18n.T("rename the selected container")},
		{i18n.T("restore"), i18n.T("restores a container from a checkpoint")},
		{i18n.T("rm"), i18n.T("remove the selected container")},
		{i18n.T("run"), i18n.T("runs a command in a new container from the given image")},
		{i18n.T("start"), i18n.T("start the selected containers")},
		{i18n.T("stats"), i18n.T("display container resource usage statistics")},
		{i18n.T("stop"), i18n.T("stop the selected containers")},
		{i18n.T("top"), i18n.T("display the running processes of the selected container")},
		{i18n.T("unpause"), i18n.T("unpause the selected container that was paused before")},
	})

	cnt.cmdDialog.SetSelectedFunc(func() {
		cnt.cmdDialog.Hide()
		// GetSelectedItem returns translated command, map it back to English
		translatedCmd := cnt.cmdDialog.GetSelectedItem()
		englishCmd := cnt.getEnglishCommand(translatedCmd)
		cnt.runCommand(englishCmd)
	})

	cnt.cmdDialog.SetCancelFunc(cnt.cmdDialog.Hide)
}

// getEnglishCommand maps translated command back to English command key
func (cnt *Containers) getEnglishCommand(translatedCmd string) string {
	// Create reverse mapping from translated to English
	commandMap := map[string]string{
		i18n.T("attach"):      "attach",
		i18n.T("checkpoint"):  "checkpoint",
		i18n.T("commit"):      "commit",
		i18n.T("create"):      "create",
		i18n.T("diff"):        "diff",
		i18n.T("exec"):        "exec",
		i18n.T("healthcheck"): "healthcheck",
		i18n.T("inspect"):     "inspect",
		i18n.T("kill"):        "kill",
		i18n.T("logs"):        "logs",
		i18n.T("pause"):       "pause",
		i18n.T("port"):        "port",
		i18n.T("prune"):       "prune",
		i18n.T("rename"):      "rename",
		i18n.T("restore"):     "restore",
		i18n.T("rm"):          "rm",
		i18n.T("run"):         "run",
		i18n.T("start"):       "start",
		i18n.T("stats"):       "stats",
		i18n.T("stop"):        "stop",
		i18n.T("top"):         "top",
		i18n.T("unpause"):     "unpause",
	}

	if englishCmd, exists := commandMap[translatedCmd]; exists {
		return englishCmd
	}

	// Fallback to original if not found (for English or unknown commands)
	return translatedCmd
}

// UpdateLanguage updates all translatable text when language changes
func (cnt *Containers) UpdateLanguage() {
	// DO NOT update cnt.title - it's used as page key in app.pages
	// Only update headers which are displayed
	cnt.headers = []string{
		i18n.T("container id"),
		i18n.T("image"),
		i18n.T("pod"),
		i18n.T("created"),
		i18n.T("status"),
		i18n.T("names"),
		i18n.T("ports"),
	}

	// Rebuild command dialog with new translations
	cnt.buildCommandDialog()

	// Update table header cells
	for i := range cnt.headers {
		header := fmt.Sprintf("[black::b]%s", strings.ToUpper(cnt.headers[i])) //nolint:perfsprint
		cnt.table.GetCell(0, i).SetText(header)
	}

	// Update table title with translation
	translatedTitle := i18n.T(cnt.title)
	cnt.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(translatedTitle), cnt.table.GetRowCount()-1))

	// Update dialogs
	cnt.checkpointDialog.UpdateLanguage()
	cnt.commitDialog.UpdateLanguage()
	cnt.createDialog.UpdateLanguage()
	cnt.runDialog.UpdateLanguage()
	cnt.execDialog.UpdateLanguage()
	cnt.statsDialog.UpdateLanguage()
	cnt.terminalDialog.UpdateLanguage()
	cnt.topDialog.UpdateLanguage()
	cnt.messageDialog.UpdateLanguage()
	cnt.confirmDialog.UpdateLanguage()
	cnt.cmdInputDialog.UpdateLanguage()
	cnt.restoreDialog.UpdateLanguage()

	// Rebuild error dialog to update OK button
	cnt.errorDialog = dialogs.NewErrorDialog()
}
