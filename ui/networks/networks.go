package networks

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/networks/netdialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
)

const (
	viewNetworkNameColIndex = 0 + iota
	viewNetworkVersionColIndex
	viewNetworkPluginColIndex
)

var (
	errNoNetworkRemove     = errors.New("there is no network to remove")
	errNoNetworkInspect    = errors.New("there is no network to display inspect")
	errNoNetworkDisconnect = errors.New("there is no network to disconnect")
	errNoNetworkConnect    = errors.New("there is no network to connect")
)

// Networks implemnents the Networks page primitive.
type Networks struct {
	*tview.Box

	title            string
	headers          []string
	table            *tview.Table
	errorDialog      *dialogs.ErrorDialog
	progressDialog   *dialogs.ProgressDialog
	confirmDialog    *dialogs.ConfirmDialog
	cmdDialog        *dialogs.CommandDialog
	messageDialog    *dialogs.MessageDialog
	sortDialog       *dialogs.SortDialog
	createDialog     *netdialogs.NetworkCreateDialog
	connectDialog    *netdialogs.NetworkConnectDialog
	disconnectDialog *netdialogs.NetworkDisconnectDialog
	networkList      networkListReport
	selectedID       string
	confirmData      string
	appFocusHandler  func()
}

type networkListReport struct {
	mu        sync.Mutex
	report    []types.Network
	sortBy    string
	ascending bool
}

// NewNetworks returns nets page view.
func NewNetworks() *Networks {
	nets := &Networks{
		Box:              tview.NewBox(),
		title:            "networks",
		headers:          []string{i18n.T("id"), i18n.T("name"), i18n.T("driver")},
		errorDialog:      dialogs.NewErrorDialog(),
		progressDialog:   dialogs.NewProgressDialog(),
		confirmDialog:    dialogs.NewConfirmDialog(),
		messageDialog:    dialogs.NewMessageDialog(""),
		sortDialog:       dialogs.NewSortDialog([]string{"name", "driver"}, 0),
		createDialog:     netdialogs.NewNetworkCreateDialog(),
		connectDialog:    netdialogs.NewNetworkConnectDialog(),
		disconnectDialog: netdialogs.NewNetworkDisconnectDialog(),
		networkList:      networkListReport{sortBy: "name", ascending: true},
	}

	// Build command dialog with translations
	nets.buildCommandDialog()

	nets.table = tview.NewTable()
	nets.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(i18n.T(nets.title))))
	nets.table.SetBorderColor(style.BorderColor)
	nets.table.SetBackgroundColor(style.BgColor)
	nets.table.SetTitleColor(style.FgColor)
	nets.table.SetBorder(true)

	for i := range nets.headers {
		nets.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(nets.headers[i]))). //nolint:perfsprint
													SetExpansion(1).
													SetBackgroundColor(style.PageHeaderBgColor).
													SetTextColor(style.PageHeaderFgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	nets.table.SetFixed(1, 1)
	nets.table.SetSelectable(true, false)

	// set command dialog functions
	nets.cmdDialog.SetSelectedFunc(func() {
		nets.cmdDialog.Hide()
		nets.runCommand(nets.cmdDialog.GetSelectedItem())
	})
	nets.cmdDialog.SetCancelFunc(func() {
		nets.cmdDialog.Hide()
	})

	// set message dialog functions
	nets.messageDialog.SetCancelFunc(func() {
		nets.messageDialog.Hide()
	})

	// set confirm dialogs functions
	nets.confirmDialog.SetSelectedFunc(func() {
		nets.confirmDialog.Hide()

		switch nets.confirmData {
		case utils.PruneCommandLabel:
			nets.prune()
		case "rm":
			nets.remove()
		}
	})

	nets.confirmDialog.SetCancelFunc(func() {
		nets.confirmDialog.Hide()
	})

	// set create dialog functions
	nets.createDialog.SetCancelFunc(func() {
		nets.createDialog.Hide()
	})

	nets.createDialog.SetCreateFunc(func() {
		nets.createDialog.Hide()
		nets.create()
	})

	// set connect dialog functions
	nets.connectDialog.SetCancelFunc(nets.connectDialog.Hide)
	nets.connectDialog.SetConnectFunc(nets.connect)

	// set disconnect dialog functions
	nets.disconnectDialog.SetCancelFunc(nets.disconnectDialog.Hide)
	nets.disconnectDialog.SetDisconnectFunc(nets.disconnect)

	// set sort dialog functions
	nets.sortDialog.SetCancelFunc(nets.sortDialog.Hide)
	nets.sortDialog.SetSelectFunc(nets.SortView)

	return nets
}

// SetAppFocusHandler sets application focus handler.
func (nets *Networks) SetAppFocusHandler(handler func()) {
	nets.appFocusHandler = handler
}

// GetTitle returns primitive title.
func (nets *Networks) GetTitle() string {
	return nets.title
}

// HasFocus returns whether or not this primitive has focus.
func (nets *Networks) HasFocus() bool {
	if nets.SubDialogHasFocus() {
		return true
	}

	if nets.table.HasFocus() || nets.Box.HasFocus() {
		return true
	}

	return false
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus.
func (nets *Networks) SubDialogHasFocus() bool {
	for _, dialog := range nets.getInnerDialogs() {
		if dialog.HasFocus() {
			return true
		}
	}

	return false
}

// Focus is called when this primitive receives focus.
func (nets *Networks) Focus(delegate func(p tview.Primitive)) {
	if nets.errorDialog.IsDisplay() {
		delegate(nets.errorDialog)

		return
	}

	if nets.confirmDialog.IsDisplay() {
		delegate(nets.confirmDialog)

		return
	}

	for _, dialog := range nets.getInnerDialogs() {
		if dialog.IsDisplay() {
			delegate(dialog)

			return
		}
	}

	delegate(nets.table)
}

// HideAllDialogs hides all sub dialogs.
func (nets *Networks) HideAllDialogs() {
	for _, dialog := range nets.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.Hide()
		}
	}
}

func (nets *Networks) getSelectedItem() (string, string) {
	if nets.table.GetRowCount() <= 1 {
		return "", ""
	}

	row, _ := nets.table.GetSelection()
	netID := nets.table.GetCell(row, 0).Text
	netName := nets.table.GetCell(row, 1).Text

	return netID, netName
}

func (nets *Networks) getInnerDialogs() []utils.UIDialog {
	dialogs := []utils.UIDialog{
		nets.errorDialog,
		nets.progressDialog,
		nets.confirmDialog,
		nets.cmdDialog,
		nets.messageDialog,
		nets.connectDialog,
		nets.createDialog,
		nets.disconnectDialog,
		nets.sortDialog,
	}

	return dialogs
}

// buildCommandDialog builds the command dialog with translated strings
func (nets *Networks) buildCommandDialog() {
	if nets.cmdDialog != nil {
		nets.cmdDialog.Hide()
	}

	// Translate both command names and descriptions
	nets.cmdDialog = dialogs.NewCommandDialog([][]string{
		{i18n.T("connect"), i18n.T("connect a container to a network")},
		{i18n.T("create"), i18n.T("create a Podman CNI network")},
		{i18n.T("disconnect"), i18n.T("disconnect a container from a network")},
		{i18n.T("inspect"), i18n.T("displays the raw CNI network configuration")},
		{i18n.T("prune"), i18n.T("remove all unused networks")},
		{i18n.T("rm"), i18n.T("remove a CNI networks")},
	})

	nets.cmdDialog.SetSelectedFunc(func() {
		nets.cmdDialog.Hide()
		// GetSelectedItem returns translated command, map it back to English
		translatedCmd := nets.cmdDialog.GetSelectedItem()
		englishCmd := nets.getEnglishCommand(translatedCmd)
		nets.runCommand(englishCmd)
	})

	nets.cmdDialog.SetCancelFunc(func() {
		nets.cmdDialog.Hide()
	})
}

// getEnglishCommand maps translated command back to English command key
func (nets *Networks) getEnglishCommand(translatedCmd string) string {
	// Create reverse mapping from translated to English
	commandMap := map[string]string{
		i18n.T("connect"):    "connect",
		i18n.T("create"):     "create",
		i18n.T("disconnect"): "disconnect",
		i18n.T("inspect"):    "inspect",
		i18n.T("prune"):      "prune",
		i18n.T("rm"):         "rm",
	}

	if englishCmd, ok := commandMap[translatedCmd]; ok {
		return englishCmd
	}

	// Fallback to original if not found (shouldn't happen)
	return translatedCmd
}

// UpdateLanguage updates all translatable text when language changes.
func (nets *Networks) UpdateLanguage() {
	// Update headers
	nets.headers = []string{i18n.T("id"), i18n.T("name"), i18n.T("driver")}
	
	// Update table header cells
	for i := range nets.headers {
		header := fmt.Sprintf("[::b]%s", strings.ToUpper(nets.headers[i]))
		nets.table.GetCell(0, i).SetText(header)
	}
	
	// Update table title with translation
	translatedTitle := i18n.T(nets.title)
	nets.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(translatedTitle), nets.table.GetRowCount()-1))
	
	// Rebuild command dialog with new translations
	nets.buildCommandDialog()
	
	// Update dialogs
	nets.messageDialog.UpdateLanguage()
	nets.confirmDialog.UpdateLanguage()
	nets.errorDialog.UpdateLanguage()
	nets.connectDialog.UpdateLanguage()
	nets.disconnectDialog.UpdateLanguage()
	nets.createDialog.UpdateLanguage()
	
	// Rebuild sort dialog with translated headers
	nets.sortDialog = dialogs.NewSortDialog([]string{i18n.T("name"), i18n.T("driver")}, 0)
	nets.sortDialog.SetCancelFunc(nets.sortDialog.Hide)
	nets.sortDialog.SetSelectFunc(nets.SortView)
}
