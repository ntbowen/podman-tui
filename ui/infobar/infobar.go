package infobar

import (
	"fmt"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InfoBarViewHeight info bar height.
const (
	InfoBarViewHeight = 5
	connectionCellRow = 0
	hostnameCellRow   = 1
	osCellRow         = 2
	memCellRow        = 3
	swapCellRow       = 4
	dataCol1Index     = 1
	dataCol2Index     = 2
	dataCol3Index     = 4
	dataCol4Index     = 5
	totalRows         = 3
	defaultPerc       = 0.00
)

// InfoBar implements the info bar primitive.
type InfoBar struct {
	*tview.Box

	table      *tview.Table
	title      string
	connStatus registry.ConnStatus
}

// NewInfoBar returns info bar view.
func NewInfoBar() *InfoBar {
	table := tview.NewTable()
	headerColor := style.GetColorHex(style.InfoBarItemFgColor)
	emptyCell := func() *tview.TableCell {
		return tview.NewTableCell("")
	}

	// empty column
	for i := range 5 {
		table.SetCell(i, 0, emptyCell())
	}

	// valueColor := Styles.InfoBar.ValueFgColor
	table.SetCell(
		connectionCellRow,
		dataCol1Index,
		tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Connection:"))),
	)

	disconnectStatus := fmt.Sprintf("%s %s", style.HeavyRedCrossMark, i18n.T("DISCONNECTED"))
	table.SetCell(connectionCellRow, dataCol2Index, tview.NewTableCell(disconnectStatus))

	table.SetCell(
		hostnameCellRow,
		dataCol1Index,
		tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Hostname:"))),
	)

	table.SetCell(hostnameCellRow, dataCol2Index, emptyCell())

	table.SetCell(
		osCellRow,
		dataCol1Index,
		tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("OS type:"))),
	)

	table.SetCell(osCellRow, dataCol2Index, emptyCell())

	table.SetCell(
		memCellRow,
		dataCol1Index,
		tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Memory usage:"))),
	)

	table.SetCell(memCellRow, dataCol2Index, tview.NewTableCell(utils.ProgressUsageString(defaultPerc)))
	table.SetCell(swapCellRow, dataCol1Index, tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Swap usage:"))))
	table.SetCell(swapCellRow, dataCol2Index, tview.NewTableCell(utils.ProgressUsageString(defaultPerc)))

	// empty column
	for i := range dataCol4Index {
		table.SetCell(i, totalRows, emptyCell())
	}

	table.SetCell(
		connectionCellRow,
		dataCol3Index,
		tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Kernel version:"))),
	)

	table.SetCell(connectionCellRow, dataCol4Index, emptyCell())

	table.SetCell(
		hostnameCellRow,
		dataCol3Index,
		tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("API version:"))),
	)

	table.SetCell(hostnameCellRow, dataCol4Index, emptyCell())

	table.SetCell(
		osCellRow,
		dataCol3Index,
		tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("OCI runtime:"))),
	)

	table.SetCell(osCellRow, dataCol4Index, emptyCell())

	table.SetCell(memCellRow,
		dataCol3Index,
		tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Conmon version:"))),
	)

	table.SetCell(memCellRow, dataCol4Index, emptyCell())

	table.SetCell(swapCellRow,
		dataCol3Index,
		tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Buildah version:"))),
	)

	table.SetCell(swapCellRow, dataCol4Index, emptyCell())

	// infobar
	infoBar := &InfoBar{
		Box:        tview.NewBox(),
		title:      "infobar",
		table:      table,
		connStatus: registry.ConnectionStatusDisconnected,
	}

	return infoBar
}

// UpdatePodmanInfo updates api, conmon and buildah version values.
func (info *InfoBar) UpdatePodmanInfo(apiVers string, ociRuntime string, conmonVers string, buildahVers string) {
	info.table.GetCell(hostnameCellRow, dataCol4Index).SetText(apiVers)
	info.table.GetCell(osCellRow, dataCol4Index).SetText(ociRuntime)
	info.table.GetCell(memCellRow, dataCol4Index).SetText(conmonVers)
	info.table.GetCell(swapCellRow, dataCol4Index).SetText(buildahVers)
}

// UpdateBasicInfo updates hostname, kernel and os type values.
func (info *InfoBar) UpdateBasicInfo(hostname string, kernel string, ostype string) {
	info.table.GetCell(hostnameCellRow, dataCol2Index).SetText(hostname)
	info.table.GetCell(osCellRow, dataCol2Index).SetText(ostype)
	info.table.GetCell(connectionCellRow, dataCol4Index).SetText(kernel)
}

// UpdateSystemUsageInfo updates memory and swap values.
func (info *InfoBar) UpdateSystemUsageInfo(memUsage float64, swapUsage float64) {
	memUsageText := utils.ProgressUsageString(memUsage)
	swapUsageText := utils.ProgressUsageString(swapUsage)

	info.table.GetCell(memCellRow, dataCol2Index).SetText(memUsageText)
	info.table.GetCell(swapCellRow, dataCol2Index).SetText(swapUsageText)
}

// UpdateConnStatus updates connection status value.
func (info *InfoBar) UpdateConnStatus(status registry.ConnStatus) {
	var connStatus string

	info.connStatus = status

	switch info.connStatus {
	case registry.ConnectionStatusConnected:
		connStatus = fmt.Sprintf("%s %s", style.HeavyGreenCheckMark, i18n.T("STATUS_OK"))
	case registry.ConnectionStatusConnectionError:
		connStatus = fmt.Sprintf("%s %s", style.HeavyRedCrossMark, i18n.T("STATUS_ERROR"))
	default:
		connStatus = fmt.Sprintf("%s %s", style.HeavyRedCrossMark, i18n.T("DISCONNECTED"))
	}

	info.table.GetCell(connectionCellRow, dataCol2Index).SetText(connStatus)
}

// RefreshLabels refreshes all labels with current language translations.
func (info *InfoBar) RefreshLabels() {
	headerColor := style.GetColorHex(style.InfoBarItemFgColor)
	
	// Update left column labels
	info.table.GetCell(connectionCellRow, dataCol1Index).SetText(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Connection:")))
	info.table.GetCell(hostnameCellRow, dataCol1Index).SetText(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Hostname:")))
	info.table.GetCell(osCellRow, dataCol1Index).SetText(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("OS type:")))
	info.table.GetCell(memCellRow, dataCol1Index).SetText(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Memory usage:")))
	info.table.GetCell(swapCellRow, dataCol1Index).SetText(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Swap usage:")))
	
	// Update right column labels
	info.table.GetCell(connectionCellRow, dataCol3Index).SetText(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Kernel version:")))
	info.table.GetCell(hostnameCellRow, dataCol3Index).SetText(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("API version:")))
	info.table.GetCell(osCellRow, dataCol3Index).SetText(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("OCI runtime:")))
	info.table.GetCell(memCellRow, dataCol3Index).SetText(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Conmon version:")))
	info.table.GetCell(swapCellRow, dataCol3Index).SetText(fmt.Sprintf("[%s::]%s", headerColor, i18n.T("Buildah version:")))
	
	// Update connection status with current language
	info.UpdateConnStatus(info.connStatus)
}

// Draw draws this primitive onto the screen.
func (info *InfoBar) Draw(screen tcell.Screen) {
	info.DrawForSubclass(screen, info)
	info.SetBorder(false)

	x, y, width, height := info.GetInnerRect()

	info.table.SetRect(x, y, width, height)
	info.table.SetBorder(false)
	info.table.Draw(screen)
}
