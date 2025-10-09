package system

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/system/sysdialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
)

var ErrConnectionInprogres = errors.New("connection is in progress, need to disconnect")

// System implemnents the system information page primitive.
type System struct {
	*tview.Box

	title                    string
	connTable                *tview.Table
	connTableHeaders         []string
	cmdDialog                *dialogs.CommandDialog
	confirmDialog            *dialogs.ConfirmDialog
	messageDialog            *dialogs.MessageDialog
	progressDialog           *dialogs.ProgressDialog
	errorDialog              *dialogs.ErrorDialog
	sortDialog               *dialogs.SortDialog
	languageDialog           *dialogs.LanguageDialog
	eventDialog              *sysdialogs.EventsDialog
	dfDialog                 *sysdialogs.DfDialog
	connPrgDialog            *sysdialogs.ConnectDialog
	connAddDialog            *sysdialogs.AddConnectionDialog
	confirmData              string
	connectionList           connectionListReport
	connectionListFunc       func() []registry.Connection
	connectionAddFunc        func(string, string, string) error
	connectionRemoveFunc     func(string) error
	connectionSetDefaultFunc func(string) error
	connectionConnectFunc    func(registry.Connection)
	connectionDisconnectFunc func()
	appFocusHandler          func()
	refreshUIHandler         func()
}

type connectionListReport struct {
	mu        sync.Mutex
	report    []registry.Connection
	sortBy    string
	ascending bool
}

type sysSelectedItem struct {
	name     string
	status   string
	uri      string
	identity string
}

// NewSystem returns new system page view.
func NewSystem() *System {
	headers := []string{
		i18n.T("NAME"),
		i18n.T("DEFAULT"),
		i18n.T("STATUS"),
		i18n.T("URI"),
		i18n.T("IDENTITY"),
	}
	sys := &System{
		Box:              tview.NewBox(),
		title:            "system",
		connTable:        tview.NewTable(),
		connTableHeaders: headers,
		confirmDialog:    dialogs.NewConfirmDialog(),
		progressDialog:   dialogs.NewProgressDialog(),
		errorDialog:      dialogs.NewErrorDialog(),
		messageDialog:    dialogs.NewMessageDialog(""),
		sortDialog:       dialogs.NewSortDialog(headers, 0),
		languageDialog:   dialogs.NewLanguageDialog(),
		eventDialog:      sysdialogs.NewEventDialog(),
		dfDialog:         sysdialogs.NewDfDialog(),
		connPrgDialog:    sysdialogs.NewConnectDialog(),
		connAddDialog:    sysdialogs.NewAddConnectionDialog(),
		connectionList:   connectionListReport{sortBy: "name", ascending: true},
	}

	// connection table
	sys.connTable.SetBackgroundColor(style.BgColor)
	sys.connTable.SetBorder(true)
	sys.updateConnTableTitle(0)
	sys.connTable.SetTitleColor(style.FgColor)
	sys.connTable.SetBorderColor(style.BorderColor)
	sys.connTable.SetFixed(1, 1)
	sys.connTable.SetSelectable(true, false)

	for i := range sys.connTableHeaders {
		header := fmt.Sprintf("[::b]%s", sys.connTableHeaders[i]) //nolint:perfsprint
		sys.connTable.SetCell(0, i,
			tview.NewTableCell(header).
				SetExpansion(1).
				SetBackgroundColor(style.PageHeaderBgColor).
				SetTextColor(style.PageHeaderFgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}

	sys.buildCommandDialog()

	// set confirm dialogs functions.
	sys.confirmDialog.SetSelectedFunc(func() {
		sys.confirmDialog.Hide()

		switch sys.confirmData {
		case utils.PruneCommandLabel:
			sys.prune()
		case "remove_conn":
			sys.remove()
		}
	})

	sys.confirmDialog.SetCancelFunc(func() {
		sys.confirmDialog.Hide()
	})

	// set message dialog functions
	sys.messageDialog.SetCancelFunc(func() {
		sys.messageDialog.Hide()
	})

	// set event dialog functions
	sys.eventDialog.SetCancelFunc(func() {
		sys.eventDialog.Hide()
	})

	// set disk usage function
	sys.dfDialog.SetCancelFunc(func() {
		sys.dfDialog.Hide()
	})

	// set connection progress bar cancel function
	sys.connPrgDialog.SetCancelFunc(func() {
		sys.connPrgDialog.Hide()
		registry.UnsetConnection()
		sys.eventDialog.SetText("")
		sys.UpdateData()
	})

	// set connection create dialog functions
	sys.connAddDialog.SetCancelFunc(sys.connAddDialog.Hide)
	sys.connAddDialog.SetAddFunc(func() {
		sys.addConnection()
	})

	// set connection sort dialog functions
	sys.sortDialog.SetCancelFunc(sys.sortDialog.Hide)
	sys.sortDialog.SetSelectFunc(sys.SortView)

	// set language dialog functions
	sys.languageDialog.SetCancelFunc(sys.languageDialog.Hide)
	sys.languageDialog.SetSelectedFunc(func() {
		sys.changeLanguage()
	})

	return sys
}

// SetAppFocusHandler sets application focus handler.
func (sys *System) SetAppFocusHandler(handler func()) {
	sys.appFocusHandler = handler
}

// SetRefreshUIHandler sets UI refresh handler.
func (sys *System) SetRefreshUIHandler(handler func()) {
	sys.refreshUIHandler = handler
}

// GetTitle returns primitive title.
func (sys *System) GetTitle() string {
	return sys.title
}

// HasFocus returns whether or not this primitive has focus.
func (sys *System) HasFocus() bool {
	for _, dialog := range sys.getInnerDialogs(true) {
		if dialog.HasFocus() {
			return true
		}
	}

	if sys.connTable.HasFocus() || sys.Box.HasFocus() {
		return true
	}

	return false
}

// SubDialogHasFocus returns true if there is an active dialog
// displayed on the front screen.
func (sys *System) SubDialogHasFocus() bool {
	for _, dialog := range sys.getInnerDialogs(false) {
		if dialog.HasFocus() {
			return true
		}
	}

	return false
}

// Focus is called when this primitive receives focus.
func (sys *System) Focus(delegate func(p tview.Primitive)) {
	for _, dialog := range sys.getInnerDialogs(true) {
		if dialog.IsDisplay() {
			delegate(dialog)

			return
		}
	}

	delegate(sys.connTable)
}

// SetEventMessage appends podman events to textview.
func (sys *System) SetEventMessage(messages []string) {
	msg := strings.Join(messages, "\n")
	sys.eventDialog.SetText(msg)
}

// SetConnectionProgressMessage sets connection progressbar error message.
func (sys *System) SetConnectionProgressMessage(message string) {
	sys.connPrgDialog.SetMessage(message)
}

// SetConnectionProgressDestName sets connection
// progressbar title destination name.
func (sys *System) SetConnectionProgressDestName(name string) {
	sys.connPrgDialog.SetDestinationName(name)
}

// ConnectionProgressDisplay displays or hide the connection progress dialog.
func (sys *System) ConnectionProgressDisplay(display bool) {
	if display {
		sys.hideAllDialogs(false)
		sys.connPrgDialog.Display()

		return
	}

	sys.connPrgDialog.Hide()
}

// SetConnectionListFunc sets list destination function.
func (sys *System) SetConnectionListFunc(list func() []registry.Connection) {
	sys.connectionListFunc = list
}

// SetConnectionSetDefaultFunc sets set destination default function.
func (sys *System) SetConnectionSetDefaultFunc(setDefault func(dest string) error) {
	sys.connectionSetDefaultFunc = setDefault
}

// SetConnectionConnectFunc sets system connect function.
func (sys *System) SetConnectionConnectFunc(connect func(dest registry.Connection)) {
	sys.connectionConnectFunc = connect
}

// SetConnectionDisconnectFunc sets system disconnect function.
func (sys *System) SetConnectionDisconnectFunc(disconnect func()) {
	sys.connectionDisconnectFunc = disconnect
}

// SetConnectionAddFunc sets system add new connection function.
func (sys *System) SetConnectionAddFunc(add func(name string, uri string, identity string) error) {
	sys.connectionAddFunc = add
}

// SetConnectionRemoveFunc sets system remove connection function.
func (sys *System) SetConnectionRemoveFunc(remove func(name string) error) {
	sys.connectionRemoveFunc = remove
}

func (sys *System) getSelectedItem() *sysSelectedItem {
	selectedItem := sysSelectedItem{}

	if sys.connTable.GetRowCount() <= 1 {
		return &selectedItem
	}

	row, _ := sys.connTable.GetSelection()
	selectedItem.name = sys.connTable.GetCell(row, 0).Text
	selectedItem.status = sys.connTable.GetCell(row, 2).Text   //nolint:mnd
	selectedItem.uri = sys.connTable.GetCell(row, 3).Text      //nolint:mnd
	selectedItem.identity = sys.connTable.GetCell(row, 4).Text //nolint:mnd

	return &selectedItem
}

func (sys *System) hideAllDialogs(all bool) {
	for _, dialog := range sys.getInnerDialogs(all) {
		if dialog.IsDisplay() {
			dialog.Hide()
		}
	}
}

func (sys *System) getInnerDialogs(all bool) []utils.UIDialog {
	if all {
		return []utils.UIDialog{
			sys.progressDialog,
			sys.cmdDialog,
			sys.confirmDialog,
			sys.messageDialog,
			sys.dfDialog,
			sys.errorDialog,
			sys.connPrgDialog,
			sys.eventDialog,
			sys.connAddDialog,
			sys.sortDialog,
			sys.languageDialog,
		}
	}

	return []utils.UIDialog{
		sys.progressDialog,
		sys.cmdDialog,
		sys.confirmDialog,
		sys.messageDialog,
		sys.dfDialog,
		sys.errorDialog,
		sys.eventDialog,
		sys.connAddDialog,
		sys.sortDialog,
		sys.languageDialog,
	}
}

func (sys *System) buildCommandDialog() {
	if sys.cmdDialog != nil {
		sys.cmdDialog.Hide()
	}

	sys.cmdDialog = dialogs.NewCommandDialog([][]string{
		{i18n.T("add connection"), i18n.T("record destination for the Podman TUI service")},
		{i18n.T("connect"), i18n.T("connect to selected destination")},
		{i18n.T("disconnect"), i18n.T("disconnect from connected destination")},
		{i18n.T("disk usage"), i18n.T("display destination podman related disk usage")},
		{i18n.T("events"), i18n.T("display destination system events")},
		{i18n.T("info"), i18n.T("display destination podman system information")},
		{i18n.T("prune"), i18n.T("remove all unused pod, container, image and volume data")},
		{i18n.T("remove connection"), i18n.T("delete named destination for the Podman TUI")},
		{i18n.T("language"), i18n.T("switch application display language")},
		{i18n.T("set default"), i18n.T("set selected destination as a default service")},
	})

	sys.cmdDialog.SetSelectedFunc(func() {
		sys.cmdDialog.Hide()
		sys.runCommand(sys.cmdDialog.GetSelectedItem())
	})

	sys.cmdDialog.SetCancelFunc(func() {
		sys.cmdDialog.Hide()
	})
}

func (sys *System) updateTableHeaders() {
	// Update header texts
	sys.connTableHeaders = []string{
		i18n.T("NAME"),
		i18n.T("DEFAULT"),
		i18n.T("STATUS"),
		i18n.T("URI"),
		i18n.T("IDENTITY"),
	}

	// Update table header cells
	for i := range sys.connTableHeaders {
		header := fmt.Sprintf("[::b]%s", sys.connTableHeaders[i]) //nolint:perfsprint
		sys.connTable.GetCell(0, i).SetText(header)
	}
}

func (sys *System) rebuildDialogs() {
	// Rebuild add connection dialog
	sys.connAddDialog = sysdialogs.NewAddConnectionDialog()
	sys.connAddDialog.SetCancelFunc(sys.connAddDialog.Hide)
	sys.connAddDialog.SetAddFunc(func() {
		sys.addConnection()
	})

	// Rebuild disk usage dialog
	sys.dfDialog = sysdialogs.NewDfDialog()
	sys.dfDialog.SetCancelFunc(sys.dfDialog.Hide)

	// Rebuild events dialog
	sys.eventDialog = sysdialogs.NewEventDialog()
	sys.eventDialog.SetCancelFunc(sys.eventDialog.Hide)

	// Rebuild confirm dialog
	sys.confirmDialog = dialogs.NewConfirmDialog()
	sys.confirmDialog.SetSelectedFunc(func() {
		sys.confirmDialog.Hide()

		switch sys.confirmData {
		case utils.PruneCommandLabel:
			sys.prune()
		case "remove_conn":
			sys.remove()
		}
	})
	sys.confirmDialog.SetCancelFunc(func() {
		sys.confirmDialog.Hide()
	})

	// Rebuild language dialog
	sys.languageDialog = dialogs.NewLanguageDialog()
	sys.languageDialog.SetCancelFunc(sys.languageDialog.Hide)
	sys.languageDialog.SetSelectedFunc(func() {
		sys.changeLanguage()
	})

	// Rebuild message dialog
	sys.messageDialog = dialogs.NewMessageDialog("")
	sys.messageDialog.SetCancelFunc(sys.messageDialog.Hide)

	// Rebuild connection progress dialog
	sys.connPrgDialog = sysdialogs.NewConnectDialog()

	// Rebuild error dialog
	sys.errorDialog = dialogs.NewErrorDialog()

	// Rebuild sort dialog
	sys.sortDialog = dialogs.NewSortDialog(sys.connTableHeaders, 0)
	sys.sortDialog.SetCancelFunc(sys.sortDialog.Hide)
	sys.sortDialog.SetSelectFunc(sys.SortView)
}
