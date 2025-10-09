package system

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/sysinfo"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rs/zerolog/log"
)

func (sys *System) runCommand(cmd string) {
	switch cmd {
	case "add connection", i18n.T("add connection"):
		sys.connAddDialog.Display()
	case "connect", i18n.T("connect"):
		sys.connect()
	case "disconnect", i18n.T("disconnect"):
		sys.disconnect()
	case "disk usage", i18n.T("disk usage"):
		sys.df()
	case "events", i18n.T("events"):
		sys.events()
	case "info", i18n.T("info"):
		sys.info()
	case utils.PruneCommandLabel, i18n.T("prune"):
		sys.cprune()
	case "remove connection", i18n.T("remove connection"):
		sys.cremove()
	case "language", i18n.T("language"):
		sys.showLanguageDialog()
	case "set default", i18n.T("set default"):
		sys.setDefault()
	default:
		log.Warn().Msgf("unknown system command: %s", cmd)
	}
}

func (sys *System) showLanguageDialog() {
	sys.languageDialog.Display()
}

func (sys *System) changeLanguage() {
	langCode := sys.languageDialog.GetSelectedLanguage()
	if langCode == "" {
		return
	}
	i18n.SetLanguage(langCode)
	if err := i18n.SaveLanguageConfig(); err != nil {
		sys.displayError("LANGUAGE SWITCH ERROR", err)
		return
	}
	sys.buildCommandDialog()
	sys.updateTableHeaders()
	sys.rebuildDialogs()
	sys.languageDialog.Hide()

	sys.messageDialog.SetTitle("LANGUAGE CHANGED")
	sys.messageDialog.SetText(dialogs.MessageSystemInfo, i18n.GetLanguageName(langCode), i18n.T("The interface has been updated"))
	sys.messageDialog.Display()

	if sys.refreshUIHandler != nil {
		sys.refreshUIHandler()
	}
}


func (sys *System) displayError(title string, err error) {
	log.Error().Msgf("%s: %v", strings.ToLower(title), err)
	sys.errorDialog.SetTitle(i18n.T(title))
	
	// For error messages, keep the original English text
	// Complex nested errors are better left untranslated for debugging
	errMsg := fmt.Sprintf("%v", err)
	
	sys.errorDialog.SetText(errMsg)
	sys.errorDialog.Display()
}

func (sys *System) addConnection() {
	sys.connAddDialog.Hide()
	name, uri, identity := sys.connAddDialog.GetItems()
	sys.progressDialog.SetTitle(i18n.T("adding new connection"))
	sys.progressDialog.Display()

	go func() {
		err := sys.connectionAddFunc(name, uri, identity)
		sys.progressDialog.Hide()
		sys.UpdateData()

		if err != nil {
			sys.displayError("ADD NEW CONNECTION ERROR", err)
		}

		sys.appFocusHandler()
	}()
}

func (sys *System) connect() {
	selectedItem := sys.getSelectedItem()
	// empty table
	if selectedItem.name == "" {
		return
	}

	dest := registry.Connection{
		Name:     selectedItem.name,
		URI:      selectedItem.uri,
		Identity: selectedItem.identity,
	}

	sys.eventDialog.SetText("")
	sys.connectionConnectFunc(dest)
	sys.UpdateData()
}

func (sys *System) disconnect() {
	sys.connectionDisconnectFunc()
	sys.eventDialog.SetText("")
	sys.UpdateData()
}

func (sys *System) df() {
	if !sys.destIsSet() {
		return
	}

	sys.progressDialog.SetTitle(i18n.T("podman disk usage in progress"))
	sys.progressDialog.Display()

	diskUsage := func() {
		response, err := sysinfo.DiskUsage()

		sys.progressDialog.Hide()

		if err != nil {
			sys.displayError("SYSTEM DISK USAGE ERROR", err)
			sys.appFocusHandler()

			return
		}

		connName := registry.ConnectionName()
		sys.dfDialog.SetServiceName(connName)
		sys.dfDialog.UpdateDiskSummary(response)
		sys.dfDialog.Display()
		sys.appFocusHandler()
	}

	go diskUsage()
}

func (sys *System) events() {
	if !sys.destIsSet() {
		return
	}

	connName := registry.ConnectionName()
	sys.eventDialog.SetServiceName(connName)
	sys.eventDialog.Display()
}

func (sys *System) info() {
	if !sys.destIsSet() {
		return
	}

	sys.progressDialog.SetTitle(i18n.T("podman system info in progress"))
	sys.progressDialog.Display()

	go func() {
		data, err := sysinfo.Info()

		sys.progressDialog.Hide()

		if err != nil {
			sys.displayError("SYSTEM INFO ERROR", err)
			sys.appFocusHandler()

			return
		}

		connName := registry.ConnectionName()

		sys.messageDialog.SetTitle(i18n.T("SYSTEM INFORMATION"))
		sys.messageDialog.SetText(dialogs.MessageSystemInfo, connName, data)
		sys.messageDialog.DisplayFullSize()
		sys.appFocusHandler()
	}()
}

func (sys *System) cprune() {
	if !sys.destIsSet() {
		return
	}

	connName := registry.ConnectionName()

	sys.confirmDialog.SetTitle(i18n.T("podman system prune"))
	sys.confirmData = utils.PruneCommandLabel
	confirmMsg := fmt.Sprintf(
		i18n.T("Are you sure you want to remove all unused pod, container, image and volume data on %s?"),
		connName,
	)
	sys.confirmDialog.SetText(confirmMsg)
	sys.confirmDialog.Display()
}

func (sys *System) prune() {
	sys.progressDialog.SetTitle(i18n.T("system prune in progress"))

	prune := func() {
		report, err := sysinfo.Prune()

		sys.progressDialog.Hide()

		if err != nil {
			sys.displayError("SYSTEM PRUNE ERROR", err)
			sys.appFocusHandler()

			return
		}

		sys.messageDialog.SetTitle(i18n.T("PODMAN SYSTEM PRUNE"))
		sys.messageDialog.SetText(dialogs.MessageSystemInfo, registry.ConnectionName(), report)
		sys.messageDialog.Display()
		sys.appFocusHandler()
	}

	go prune()
}

func (sys *System) cremove() {
	selectedItem := sys.getSelectedItem()
	if selectedItem.status != "" {
		sys.displayError(
			"SYSTEM CONNECTION REMOVE",
			fmt.Errorf("%w %q", ErrConnectionInprogres, selectedItem.name))

		return
	}

	if selectedItem.name == "" {
		return
	}

	title := i18n.T("podman system connection remove")
	sys.confirmDialog.SetTitle(title)
	sys.confirmData = "remove_conn"
	bgColor := style.GetColorHex(style.DialogBorderColor)
	fgColor := style.GetColorHex(style.DialogFgColor)
	serviceItem := fmt.Sprintf("[%s:%s:b]%s:[:-:-] %s", fgColor, bgColor, i18n.T("SERVICE NAME"), selectedItem.name)

	confirmMsg := fmt.Sprintf("%s\n\n%s", //nolint:perfsprint,lll
		serviceItem, i18n.T("Are you sure you want to remove the selected service connection?"))
	sys.confirmDialog.SetText(confirmMsg)
	sys.confirmDialog.Display()
}

func (sys *System) remove() {
	selectedItem := sys.getSelectedItem()

	sys.progressDialog.SetTitle(i18n.T("removing connection"))
	sys.progressDialog.Display()

	go func() {
		err := sys.connectionRemoveFunc(selectedItem.name)
		sys.progressDialog.Hide()

		if err != nil {
			sys.displayError("SYSTEM CONNECTION REMOVE ERROR", err)
		}

		sys.appFocusHandler()
		sys.UpdateData()
	}()
}

func (sys *System) setDefault() {
	selectedItem := sys.getSelectedItem()
	setDefFunc := func() {
		sys.progressDialog.Hide()

		err := sys.connectionSetDefaultFunc(selectedItem.name)
		if err != nil {
			sys.displayError("SYSTEM CONNECTION SET DEFAULT ERROR", err)
		}

		sys.appFocusHandler()
	}

	sys.progressDialog.Display()

	go setDefFunc()
}

func (sys *System) destIsSet() bool {
	if !registry.ConnectionIsSet() {
		sys.errorDialog.SetText(i18n.T("not connected to any podman service"))
		sys.errorDialog.Display()

		return false
	}

	return true
}
