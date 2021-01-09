package api

import (
	"time"

	"github.com/ophum/humstack/pkg/api/meta"
)

const (
	APITypeRedeploymentV0 meta.APIType = "customv0/redeployment"
)

type RedeploymentSpec struct {
	VMNamePrefix string    `json:"vmNamePrefix" yaml:"vmNamePrefix"`
	RestartTime  time.Time `json:"restartTime" yaml:"restartTime"`
}

type RedeploymentState string

const (
	RedeploymentStateRequested      RedeploymentState = "Requested"
	RedeploymentStateStoppingVM     RedeploymentState = "StoppingVM"
	RedeploymentStateDeletingBS     RedeploymentState = "DeletingBS"
	RedeploymentStatePendingStartVM RedeploymentState = "PendingRestartVM"
	RedeploymentStateRestartingVM   RedeploymentState = "RestartingVM"
	RedeploymentStateDone           RedeploymentState = "Done"
	RedeploymentStateError          RedeploymentState = "Error"
)

type RedeploymentStatus struct {
	State RedeploymentState `json:"state" yaml:"state"`
}

type Redeployment struct {
	meta.Meta `json:"meta" yaml:"meta"`
	Spec      RedeploymentSpec   `json:"spec" yaml:"spec"`
	Status    RedeploymentStatus `json:"status" yaml:"status"`
}
