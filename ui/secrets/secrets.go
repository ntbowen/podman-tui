package secrets

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/secrets/secdialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/rivo/tview"
)

const (
	viewSecretsIDColIndex = 0 + iota
	viewSecretsNameColIndex
	viewSecretsDriverColIndex
	viewSecretsCreatedColIndex
	viewSecretsUpdatedColIndex
)

var (
	errNoSecretRemove        = errors.New(i18n.T("there is no secret to remove"))
	errNoSecretInspect       = errors.New(i18n.T("there is no secret to display inspect"))
	errSecretFileAndText     = errors.New(i18n.T("cannot select secret file and secret text together"))
	errEmptySecretFileOrText = errors.New(i18n.T("secret content not provided"))
)

// Secrets implements the secrets page primitive.
type Secrets struct {
	*tview.Box

	title           string
	headers         []string
	table           *tview.Table
	cmdDialog       *dialogs.CommandDialog
	messageDialog   *dialogs.MessageDialog
	errorDialog     *dialogs.ErrorDialog
	progressDialog  *dialogs.ProgressDialog
	confirmDialog   *dialogs.ConfirmDialog
	sortDialog      *dialogs.SortDialog
	createDialog    *secdialogs.SecretCreateDialog
	secretList      secretListReport
	appFocusHandler func()
}

type secretListReport struct {
	mu        sync.Mutex
	report    []*types.SecretInfoReport
	sortBy    string
	ascending bool
}

// NewSecrets returns secrets page view.
func NewSecrets() *Secrets {
	secrets := &Secrets{
		Box:            tview.NewBox(),
		title:          "secrets",
		headers:        []string{i18n.T("id"), i18n.T("name"), i18n.T("driver"), i18n.T("created"), i18n.T("updated")},
		table:          tview.NewTable(),
		messageDialog:  dialogs.NewMessageDialog(""),
		errorDialog:    dialogs.NewErrorDialog(),
		progressDialog: dialogs.NewProgressDialog(),
		confirmDialog:  dialogs.NewConfirmDialog(),
		sortDialog:     dialogs.NewSortDialog([]string{"name", "driver", "created", "updated"}, 2), //nolint:mnd
		createDialog:   secdialogs.NewSecretCreateDialog(),
	}

	// Build command dialog with translations
	secrets.buildCommandDialog()

	secrets.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(i18n.T(secrets.title))))
	secrets.table.SetBorderColor(style.BorderColor)
	secrets.table.SetBackgroundColor(style.BgColor)
	secrets.table.SetTitleColor(style.FgColor)
	secrets.table.SetBorder(true)

	secrets.table.SetFixed(1, 1)
	secrets.table.SetSelectable(true, false)

	// set command dialog functions
	secrets.cmdDialog.SetSelectedFunc(func() {
		secrets.cmdDialog.Hide()
		secrets.runCommand(secrets.cmdDialog.GetSelectedItem())
	})

	secrets.cmdDialog.SetCancelFunc(func() {
		secrets.cmdDialog.Hide()
	})

	// set message dialog function
	secrets.messageDialog.SetCancelFunc(func() {
		secrets.messageDialog.Hide()
	})

	// set confirm dialog functions
	secrets.confirmDialog.SetSelectedFunc(func() {
		secrets.confirmDialog.Hide()

		secrets.remove()
	})

	secrets.confirmDialog.SetCancelFunc(func() {
		secrets.confirmDialog.Hide()
	})

	// set create dialog function
	secrets.createDialog.SetCancelFunc(func() {
		secrets.createDialog.Hide()
	})

	secrets.createDialog.SetCreateFunc(func() {
		secrets.createDialog.Hide()
		secrets.create()
	})

	// set sort dialog function
	secrets.sortDialog.SetCancelFunc(secrets.sortDialog.Hide)
	secrets.sortDialog.SetSelectFunc(secrets.SortView)

	return secrets
}

// SetAppFocusHandler sets application focus handler.
func (s *Secrets) SetAppFocusHandler(handler func()) {
	s.appFocusHandler = handler
}

// GetTitle returns primitive title.
func (s *Secrets) GetTitle() string {
	return s.title
}

// HasFocus returns whether or not this primitive has focus.
func (s *Secrets) HasFocus() bool {
	if s.SubDialogHasFocus() {
		return true
	}

	if s.table.HasFocus() || s.Box.HasFocus() {
		return true
	}

	return false
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus.
func (s *Secrets) SubDialogHasFocus() bool {
	for _, dialog := range s.getInnerDialogs() {
		if dialog.HasFocus() {
			return true
		}
	}

	return false
}

// Focus is called when this primitive receives focus.
func (s *Secrets) Focus(delegate func(p tview.Primitive)) {
	// error dialog
	if s.errorDialog.IsDisplay() {
		delegate(s.errorDialog)

		return
	}

	for _, dialog := range s.getInnerDialogs() {
		if dialog.IsDisplay() {
			delegate(dialog)

			return
		}
	}

	delegate(s.table)
}

// HideAllDialogs hides all sub dialogs.
func (s *Secrets) HideAllDialogs() {
	for _, dialog := range s.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.Hide()
		}
	}
}

func (s *Secrets) getSelectedItem() (int, string, string) {
	var (
		rowIndex int
		secID    string
		secName  string
	)

	if s.table.GetRowCount() <= 1 {
		return rowIndex, secID, secName
	}

	rowIndex, _ = s.table.GetSelection()
	secID = s.table.GetCell(rowIndex, 0).Text
	secName = s.table.GetCell(rowIndex, 1).Text

	return rowIndex, secID, secName
}

func (s *Secrets) getInnerDialogs() []utils.UIDialog {
	dialogs := []utils.UIDialog{
		s.progressDialog,
		s.errorDialog,
		s.confirmDialog,
		s.cmdDialog,
		s.createDialog,
		s.messageDialog,
		s.sortDialog,
	}

	return dialogs
}

// buildCommandDialog builds the command dialog with translated strings
func (s *Secrets) buildCommandDialog() {
	if s.cmdDialog != nil {
		s.cmdDialog.Hide()
	}

	// Translate both command names and descriptions
	s.cmdDialog = dialogs.NewCommandDialog([][]string{
		{i18n.T("create"), i18n.T("create a new secret")},
		{i18n.T("inspect"), i18n.T("inspect a secret")},
		{i18n.T("rm"), i18n.T("remove a secret")},
	})

	s.cmdDialog.SetSelectedFunc(func() {
		s.cmdDialog.Hide()
		// GetSelectedItem returns translated command, map it back to English
		translatedCmd := s.cmdDialog.GetSelectedItem()
		englishCmd := s.getEnglishCommand(translatedCmd)
		s.runCommand(englishCmd)
	})

	s.cmdDialog.SetCancelFunc(func() {
		s.cmdDialog.Hide()
	})
}

// getEnglishCommand maps translated command back to English command key
func (s *Secrets) getEnglishCommand(translatedCmd string) string {
	// Create reverse mapping from translated to English
	commandMap := map[string]string{
		i18n.T("create"):  "create",
		i18n.T("inspect"): "inspect",
		i18n.T("rm"):      "rm",
	}

	if englishCmd, ok := commandMap[translatedCmd]; ok {
		return englishCmd
	}

	// Fallback to original if not found (shouldn't happen)
	return translatedCmd
}

// UpdateLanguage updates all translatable text when language changes.
func (s *Secrets) UpdateLanguage() {
	// Update headers
	s.headers = []string{i18n.T("id"), i18n.T("name"), i18n.T("driver"), i18n.T("created"), i18n.T("updated")}
	
	// Update table header cells
	for i := range s.headers {
		header := fmt.Sprintf("[::b]%s", strings.ToUpper(s.headers[i]))
		s.table.GetCell(0, i).SetText(header)
	}
	
	// Update table title with translation
	translatedTitle := i18n.T(s.title)
	s.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(translatedTitle), s.table.GetRowCount()-1))
	
	// Rebuild command dialog with new translations
	s.buildCommandDialog()
	
	// Update dialogs
	s.messageDialog.UpdateLanguage()
	s.confirmDialog.UpdateLanguage()
	s.errorDialog.UpdateLanguage()
	s.createDialog.UpdateLanguage()
	
	// Rebuild sort dialog with translated headers
	s.sortDialog = dialogs.NewSortDialog([]string{i18n.T("name"), i18n.T("driver"), i18n.T("created"), i18n.T("updated")}, 2) //nolint:mnd
	s.sortDialog.SetCancelFunc(s.sortDialog.Hide)
	s.sortDialog.SetSelectFunc(s.SortView)
}
