package networks

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rs/zerolog/log"
)

func (nets *Networks) runCommand(cmd string) {
	switch cmd {
	case "connect":
		nets.cconnect()
	case "create":
		nets.createDialog.Display()
	case "disconnect":
		nets.cdisconnect()
	case "inspect":
		nets.inspect()
	case utils.PruneCommandLabel:
		nets.cprune()
	case "rm":
		nets.rm()
	}
}

func (nets *Networks) displayError(title string, err error) {
	log.Error().Msgf("%s: %v", strings.ToLower(title), err)
	nets.errorDialog.SetTitle(title)
	nets.errorDialog.SetText(fmt.Sprintf("%v", err))
	nets.errorDialog.Display()
}

func (nets *Networks) cconnect() {
	if nets.selectedID == "" {
		nets.displayError("", errNoNetworkConnect)

		return
	}

	initData := func() {
		nets.progressDialog.SetTitle(i18n.T("podman network connect"))
		nets.progressDialog.Display()

		cntListReport, err := containers.List()
		if err != nil {
			nets.progressDialog.Hide()
			nets.displayError(i18n.T("NETWORK CONNECT ERROR"), err)
			nets.appFocusHandler()

			return
		}

		netID, netName := nets.getSelectedItem()

		nets.connectDialog.SetNetworkInfo(netID, netName)
		nets.connectDialog.SetContainers(cntListReport)
		nets.progressDialog.Hide()
		nets.connectDialog.Display()
		nets.appFocusHandler()
	}

	go initData()
}

func (nets *Networks) connect() {
	connectOptions := nets.connectDialog.GetConnectOptions()

	connect := func() {
		nets.connectDialog.Hide()
		nets.progressDialog.SetTitle(i18n.T("podman network connect"))
		nets.progressDialog.Display()

		err := networks.Connect(connectOptions)

		nets.progressDialog.Hide()

		if err != nil {
			nets.displayError(i18n.T("NETWORK CONNECT ERROR"), err)
		}

		nets.appFocusHandler()
	}

	go connect()
}

func (nets *Networks) cdisconnect() {
	if nets.selectedID == "" {
		nets.displayError("", errNoNetworkDisconnect)

		return
	}

	initData := func() {
		nets.progressDialog.SetTitle(i18n.T("podman network disconnect"))
		nets.progressDialog.Display()

		cntListReport, err := containers.List()

		nets.progressDialog.Hide()

		if err != nil {
			nets.displayError(i18n.T("NETWORK DISCONNECT ERROR"), err)
			nets.appFocusHandler()

			return
		}

		netID, netName := nets.getSelectedItem()

		nets.disconnectDialog.SetNetworkInfo(netID, netName)
		nets.disconnectDialog.SetContainers(cntListReport)
		nets.disconnectDialog.Display()
		nets.appFocusHandler()
	}

	go initData()
}

func (nets *Networks) disconnect() {
	disconnect := func() {
		networkName, containerID := nets.disconnectDialog.GetDisconnectOptions()

		nets.disconnectDialog.Hide()
		nets.progressDialog.SetTitle(i18n.T("podman network disconnect"))
		nets.progressDialog.Display()

		err := networks.Disconnect(networkName, containerID)

		nets.progressDialog.Hide()

		if err != nil {
			nets.displayError(i18n.T("NETWORK DISCONNECT ERROR"), err)
		}

		nets.appFocusHandler()
	}

	go disconnect()
}

func (nets *Networks) create() {
	createOpts := nets.createDialog.NetworkCreateOptions()

	_, err := networks.Create(createOpts)
	if err != nil {
		nets.displayError(i18n.T("NETWORK CREATE ERROR"), err)

		return
	}

	nets.UpdateData()
}

func (nets *Networks) inspect() {
	netID, netName := nets.getSelectedItem()
	if netID == "" {
		nets.displayError("", errNoNetworkInspect)

		return
	}

	data, err := networks.Inspect(netID)
	if err != nil {
		title := fmt.Sprintf("NETWORK (%s) INSPECT ERROR", netID)
		nets.displayError(title, err)

		return
	}

	headerLabel := fmt.Sprintf("%s (%s)", netID, netName)

	nets.messageDialog.SetTitle(i18n.T("podman network inspect"))
	nets.messageDialog.SetText(dialogs.MessageNetworkInfo, headerLabel, data)
	nets.messageDialog.DisplayFullSize()
}

func (nets *Networks) cprune() {
	nets.confirmDialog.SetTitle(i18n.T("podman network prune"))

	nets.confirmData = utils.PruneCommandLabel

	nets.confirmDialog.SetText(i18n.T("Are you sure you want to remove all un used network ?"))
	nets.confirmDialog.Display()
}

func (nets *Networks) prune() {
	nets.progressDialog.SetTitle(i18n.T("network prune in progress"))
	nets.progressDialog.Display()

	prune := func() {
		err := networks.Prune()

		nets.progressDialog.Hide()

		if err != nil {
			nets.displayError(i18n.T("NETWORK PRUNE ERROR"), err)
		}

		nets.appFocusHandler()
		nets.UpdateData()
	}

	go prune()
}

func (nets *Networks) rm() {
	netID, netName := nets.getSelectedItem()
	if netID == "" {
		nets.displayError("", errNoNetworkRemove)

		return
	}

	nets.confirmDialog.SetTitle(i18n.T("podman network remove"))
	nets.confirmData = "rm"

	bgColor := style.GetColorHex(style.DialogBorderColor)
	fgColor := style.GetColorHex(style.DialogFgColor)
	networkItem := fmt.Sprintf("[%s:%s:b]%s[:-:-] %s (%s)", fgColor, bgColor, i18n.T("NETWORK ID:"), netID, netName)

	description := fmt.Sprintf("%s\n\n%s", networkItem, i18n.T("Are you sure you want to remove the selected network?"))
	nets.confirmDialog.SetText(description)
	nets.confirmDialog.Display()
}

func (nets *Networks) remove() {
	nets.progressDialog.SetTitle(i18n.T("network remove in progress"))
	nets.progressDialog.Display()

	remove := func(id string) {
		err := networks.Remove(id)

		nets.progressDialog.Hide()

		if err != nil {
			title := fmt.Sprintf(i18n.T("NETWORK (%s) REMOVE ERROR"), nets.selectedID)

			nets.displayError(title, err)
		}

		nets.appFocusHandler()
		nets.UpdateData()
	}

	go remove(nets.selectedID)
}
