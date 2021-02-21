package agent

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ophum/humstack-redeployment/pkg/api"
	"github.com/ophum/humstack-redeployment/pkg/client"
	"github.com/ophum/humstack/pkg/api/system"
	hsClient "github.com/ophum/humstack/pkg/client"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type RedeploymentAgent struct {
	client         *client.RedeploymentClient
	humstackClient *hsClient.Clients
	config         *Config
	logger         *zap.Logger
	nodeName       string
}

func NewRedeploymentAgent(client *client.RedeploymentClient, humstackClient *hsClient.Clients, config *Config, logger *zap.Logger) *RedeploymentAgent {
	nodeName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return &RedeploymentAgent{
		client:         client,
		humstackClient: humstackClient,
		config:         config,
		logger:         logger,
		nodeName:       nodeName,
	}
}

func (a *RedeploymentAgent) Run() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rdList, err := a.client.List()
			if err != nil {
				a.logger.Error("get redeployment list", zap.String("msg", err.Error()), zap.Time("time", time.Now()))
				continue
			}

			for _, rd := range rdList {
				// 再展開したいVMがどのノードに展開されているか調べる

				if err := a.sync(rd); err != nil {
					a.logger.Error(
						fmt.Sprintf("sync redeployment `%s`", rd.ID),
						zap.String("msg", err.Error()),
						zap.Time("time", time.Now()),
					)
					continue
				}
			}
		}
	}
}

func (a *RedeploymentAgent) sync(rd *api.Redeployment) error {
	switch rd.Status.State {
	case "":
		/**
		 * StateをRequestedにする
		 */
		rd.Status.State = api.RedeploymentStateRequested
		if _, err := a.client.Update(rd); err != nil {
			return errors.Wrap(err, "update state")
		}
	case api.RedeploymentStateRequested:
		/**
		 * 1. StateをStoppingVMにする
		 * 2. VMのActionStateをPowerOffにする
		 */
		rd.Status.State = api.RedeploymentStateStoppingVM
		if _, err := a.client.Update(rd); err != nil {
			return errors.Wrap(err, "update state")
		}

		vmList, err := a.getVirtualMachines(rd.Spec.Group, rd.Spec.Namespace, rd.Spec.VMIDPrefix)
		if err != nil {
			return errors.Wrap(err, "get filtered vm list")
		}

		for _, vm := range vmList {
			vm.Spec.ActionState = system.VirtualMachineActionStatePowerOff
			if _, err := a.humstackClient.SystemV0().VirtualMachine().Update(vm); err != nil {
				return errors.Wrap(err, "change vm action state to powerOff")
			}
		}
	case api.RedeploymentStateStoppingVM:
		/**
		 * 再展開対象のすべてのVMのStateがStoppedの場合
		 * StateをStoppedVMにする
		 */
		vmList, err := a.getVirtualMachines(rd.Spec.Group, rd.Spec.Namespace, rd.Spec.VMIDPrefix)
		if err != nil {
			return errors.Wrap(err, "get filtered vm list")
		}

		for _, vm := range vmList {
			if vm.Status.State != system.VirtualMachineStateStopped {
				return nil
			}
		}

		rd.Status.State = api.RedeploymentStateStoppedVM
		if _, err := a.client.Update(rd); err != nil {
			return errors.Wrap(err, "update state")
		}
	case api.RedeploymentStateStoppedVM:
		/**
		* 1. StateをDeletingBSにする
		* 2. 再展開対象のVMのBSを削除する
			* BSのStateをErrorにすると既存の物が削除される
		* 3. StateをDeletedBSにする
		*/
		rd.Status.State = api.RedeploymentStateDeletingBS
		if _, err := a.client.Update(rd); err != nil {
			return errors.Wrap(err, "update state")
		}

		vmList, err := a.getVirtualMachines(rd.Spec.Group, rd.Spec.Namespace, rd.Spec.VMIDPrefix)
		if err != nil {
			return errors.Wrap(err, "get filtered vm list")
		}

		for _, vm := range vmList {
			for _, bsID := range vm.Spec.BlockStorageIDs {
				bs, err := a.humstackClient.SystemV0().BlockStorage().Get(vm.Group, vm.Namespace, bsID)
				if err != nil {
					return err
				}

				// bsのstateをerrorにするとbsAgentが作り直してくれる
				bs.Status.State = system.BlockStorageStateError
				if _, err := a.humstackClient.SystemV0().BlockStorage().Update(bs); err != nil {
					return err
				}
			}
		}

	case api.RedeploymentStateDeletingBS:
		vmList, err := a.getVirtualMachines(rd.Spec.Group, rd.Spec.Namespace, rd.Spec.VMIDPrefix)
		if err != nil {
			return errors.Wrap(err, "get filtered vm list")
		}
		for _, vm := range vmList {
			for _, bsID := range vm.Spec.BlockStorageIDs {
				bs, err := a.humstackClient.SystemV0().BlockStorage().Get(vm.Group, vm.Namespace, bsID)
				if err != nil {
					return err
				}
				if bs.Status.State != system.BlockStorageStateActive {
					return nil
				}
			}
		}

		rd.Status.State = api.RedeploymentStateDeletedBS
		if _, err := a.client.Update(rd); err != nil {
			return errors.Wrap(err, "update state")
		}
	case api.RedeploymentStateDeletedBS:
		/**
		 * StateをPendingStartVMにする
		 */
		rd.Status.State = api.RedeploymentStatePendingStartVM
		if _, err := a.client.Update(rd); err != nil {
			return errors.Wrap(err, "update state")
		}
	case api.RedeploymentStatePendingStartVM:
		/**
		 * 現在時刻 >= rd.Spec.RestartTimeの場合
		 * 1. VMのActionStateをPowerOnにする
		 * 2. StateをRestartedVMにする
		 */
		if time.Now().Before(rd.Spec.RestartTime) {
			return nil
		}

		rd.Status.State = api.RedeploymentStateRestartingVM
		if _, err := a.client.Update(rd); err != nil {
			return errors.Wrap(err, "update state")
		}

		vmList, err := a.getVirtualMachines(rd.Spec.Group, rd.Spec.Namespace, rd.Spec.VMIDPrefix)
		if err != nil {
			return errors.Wrap(err, "get filtered vm list")
		}

		for _, vm := range vmList {
			vm.Spec.ActionState = system.VirtualMachineActionStatePowerOn
			if _, err := a.humstackClient.SystemV0().VirtualMachine().Update(vm); err != nil {
				return errors.Wrap(err, "change vm action state to powerOn")
			}
		}

		rd.Status.State = api.RedeploymentStateRestartedVM
		if _, err := a.client.Update(rd); err != nil {
			return errors.Wrap(err, "update state")
		}
	case api.RedeploymentStateRestartedVM:
		/**
		 * 再展開対象のすべてのVMのStateがRunningの場合
		 * StateをDoneにする
		 */
		vmList, err := a.getVirtualMachines(rd.Spec.Group, rd.Spec.Namespace, rd.Spec.VMIDPrefix)
		if err != nil {
			return errors.Wrap(err, "get filtered vm list")
		}

		for _, vm := range vmList {
			if vm.Status.State != system.VirtualMachineStateRunning {
				return nil
			}
		}

		rd.Status.State = api.RedeploymentStateDone
		if _, err := a.client.Update(rd); err != nil {
			return errors.Wrap(err, "update state")
		}
	}
	return nil
}

func (a *RedeploymentAgent) getVirtualMachines(groupID, namespaceID, idPrefix string) ([]*system.VirtualMachine, error) {
	vmList, err := a.humstackClient.SystemV0().VirtualMachine().List(groupID, namespaceID)
	if err != nil {
		return []*system.VirtualMachine{}, errors.Wrap(err, "get vm list from humstack")
	}

	filtered := []*system.VirtualMachine{}
	for _, vm := range vmList {
		if strings.HasPrefix(vm.ID, idPrefix) {
			filtered = append(filtered, vm)
		}
	}
	return filtered, nil
}

func fileIsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
