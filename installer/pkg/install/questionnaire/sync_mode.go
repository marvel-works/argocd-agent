package questionnaire

import (
	"encoding/base64"
	"encoding/json"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/util"
	"github.com/elliotchance/orderedmap"
)

func AskAboutSyncOptions(installOptions *install.InstallCmdOptions) {
	syncModes := orderedmap.NewOrderedMap()
	syncModes.Set("Import existing Argo applications to Codefresh and auto-import all new ones in the future", "CONTINUE_SYNC")
	syncModes.Set("Import all existing Argo applications to Codefresh", "ONE_TIME_SYNC")
	syncModes.Set("Select specific Argo applications to import", "SELECT")
	syncModes.Set("Do not import anything from Argo to Codefresh", "NONE")

	_, autoSyncMode := prompt.Select(util.ConvertIntToStringArray(syncModes.Keys()), "Select argocd sync behavior please")

	syncMode, _ := syncModes.Get(autoSyncMode)

	if syncMode == "SELECT" {

		argoToken := installOptions.Argo.Token

		if installOptions.Argo.Username != "" {
			argoToken, _ = argo.GetToken(installOptions.Argo.Username, installOptions.Argo.Password, installOptions.Argo.Host)
		}

		applications, _ := argo.GetApplications(argoToken, installOptions.Argo.Host)

		applicationNames := make([]string, 0)

		for _, prj := range applications {
			applicationNames = append(applicationNames, prj.Metadata.Name)
		}

		_, applicationsForSync := prompt.Multiselect(applicationNames, "Please select application for sync")

		applicationsAsJson, _ := json.Marshal(applicationsForSync)

		installOptions.Codefresh.ApplicationsForSync = base64.StdEncoding.EncodeToString(applicationsAsJson)
	}

	installOptions.Codefresh.SyncMode = syncMode.(string)
}