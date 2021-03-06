package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/extract"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/jasonlvhit/gocron"
)

var EnvInitializer uint64 = 5

func isNewEnv(existingEnvs []store.Environment, newEnv codefresh.CFEnvironment) bool {
	for _, env := range existingEnvs {
		if env.Name == newEnv.Metadata.Name {
			return false
		}
		continue
	}
	return true
}

func handleNewApplications(applications []string) {
	for _, application := range applications {
		newApp, err := extract.ExtractNewApplication(application)
		if err != nil {
			logger.GetLogger().Errorf("Failed to handle new gitops application %v, reason: %v", application, err)
			continue
		}
		logger.GetLogger().Infof("Detect new gitops application %s, initiate initialization", application)
		_, err = codefresh.GetInstance().SendEnvironment(*newApp)
		if err != nil {
			logger.GetLogger().Errorf("Failed to send environment, reason %v", err)
		}
	}
}

func handleEnvDifference() {
	storeInst := store.GetStore()
	var newEnvs []store.Environment
	var applications []string
	envs, _ := codefresh.GetInstance().GetEnvironments()
	for _, env := range envs {
		if env.Spec.Type != "argo" {
			continue
		}

		newEnv := store.Environment{
			Name: env.Metadata.Name,
		}
		newEnvs = append(newEnvs, newEnv)

		if isNewEnv(storeInst.Environments, env) {
			applications = append(applications, env.Spec.Application)
		}
	}

	store.SetEnvironments(newEnvs)

	handleNewApplications(applications)

}

func StartEnvInitializer() {
	job := gocron.Every(EnvInitializer).Second().Do(handleEnvDifference)

	if job != nil {
		err := job.Error()

		if err != "" {
			panic("Cant run env changes job because " + err)
		}
	}

	go gocron.Start()
}
