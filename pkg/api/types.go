package api

import (
	"time"

	"github.com/ophum/humstack/pkg/api/meta"
)

const (
	APITypeRedeploymentV0 meta.APIType = "customv0/redeployment"
)

type RedeploymentSpec struct {
	Group       string    `json:"group" yaml:"group"`
	Namespace   string    `json:"namespace" yaml:"namespace"`
	VMIDPrefix  string    `json:"vmIDPrefix" yaml:"vmIDPrefix"`
	RestartTime time.Time `json:"restartTime" yaml:"restartTime"`
}

type RedeploymentState string

const (
	RedeploymentStateRequested      RedeploymentState = "Requested"
	RedeploymentStateStoppingVM     RedeploymentState = "StoppingVM"
	RedeploymentStateStoppedVM      RedeploymentState = "StoppedVM"
	RedeploymentStateDeletingBS     RedeploymentState = "DeletingBS"
	RedeploymentStateDeletedBS      RedeploymentState = "DeletedBS"
	RedeploymentStatePendingStartVM RedeploymentState = "PendingRestartVM"
	RedeploymentStateRestartingVM   RedeploymentState = "RestartingVM"
	RedeploymentStateRestartedVM    RedeploymentState = "RestartedVM"
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
