package comparator

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/git"
	"testing"
)

func TestEnvironmentComparatorWithSameEnv(t *testing.T) {

	envComparator := EnvComparator{}

	env1 := codefresh.Environment{
		Gitops:       git.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   nil,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	env2 := codefresh.Environment{
		Gitops:       git.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   nil,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	if !envComparator.Compare(&env1, &env2) {
		t.Errorf("'EnvComparator' comparation failed")
	}

}

func TestEnvironmentComparatorWithSameEnvAndActivities(t *testing.T) {

	envComparator := EnvComparator{}

	act1 := codefresh.EnvironmentActivity{
		Name:         "test",
		TargetImages: nil,
		Status:       "test",
		LiveImages:   nil,
	}

	act2 := codefresh.EnvironmentActivity{
		Name:         "test2",
		TargetImages: nil,
		Status:       "test2",
		LiveImages:   nil,
	}

	activities1 := make([]codefresh.EnvironmentActivity, 0)

	activities2 := make([]codefresh.EnvironmentActivity, 0)

	activities1 = append(activities1, act1)
	activities1 = append(activities1, act2)

	activities2 = append(activities2, act2)
	activities2 = append(activities2, act1)

	env1 := codefresh.Environment{
		Gitops:       git.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   activities1,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	env2 := codefresh.Environment{
		Gitops:       git.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   activities2,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	if !envComparator.Compare(&env1, &env2) {
		t.Errorf("'EnvComparator' comparation failed")
	}

}

func TestEnvironmentComparatorWithSameEnvAndDifferentActivities(t *testing.T) {

	envComparator := EnvComparator{}

	act1 := codefresh.EnvironmentActivity{
		Name:         "test",
		TargetImages: nil,
		Status:       "test",
		LiveImages:   nil,
	}

	act2 := codefresh.EnvironmentActivity{
		Name:         "test",
		TargetImages: nil,
		Status:       "test4",
		LiveImages:   nil,
	}

	activities1 := make([]codefresh.EnvironmentActivity, 0)

	activities2 := make([]codefresh.EnvironmentActivity, 0)

	activities1 = append(activities1, act1)
	activities1 = append(activities1, act2)

	activities2 = append(activities2, act1)

	env1 := codefresh.Environment{
		Gitops:       git.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   activities1,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	env2 := codefresh.Environment{
		Gitops:       git.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   activities2,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	if envComparator.Compare(&env1, &env2) {
		t.Errorf("'EnvComparator' comparation failed")
	}

}

func TestEnvironmentComparatorWithDiffEnv(t *testing.T) {

	envComparator := EnvComparator{}

	env1 := codefresh.Environment{
		Gitops:       git.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    123,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   nil,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	env2 := codefresh.Environment{
		Gitops:       git.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   nil,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	if envComparator.Compare(&env1, &env2) {
		t.Errorf("'EnvComparator' comparation failed")
	}

}
