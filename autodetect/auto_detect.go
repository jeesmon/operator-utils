/*
SPDX-License-Identifier: Apache-2.0
*/

package autodetect

import (
	"fmt"
	"os"
	"time"

	"github.com/jeesmon/operator-utils/actions"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	DefaultAutoDetectTick = 2 * time.Minute
)

var log = logf.Log.WithName("autodetect")

type DetectConfig struct {
	GroupVersionKinds []schema.GroupVersionKind
	Delay             *time.Duration
	ExitOnChange      bool
}

// Background represents a procedure that runs in the background, periodically auto-detecting features
type Background struct {
	config DetectConfig
	dc     discovery.ServerResourcesInterface
	ticker *time.Ticker
}

// New creates a new auto-detect runner
func NewAutoDetect(mgr manager.Manager, config DetectConfig) (*Background, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}

	bg := &Background{
		config: config,
		dc:     dc,
	}

	// periodically attempts to auto detect all the capabilities for this operator
	if config.Delay == nil {
		bg.ticker = time.NewTicker(DefaultAutoDetectTick)
	} else {
		bg.ticker = time.NewTicker(*config.Delay)
	}

	return bg, nil
}

// Start initializes the auto-detection process that runs in the background
func (b *Background) Start() {
	go func() {
		b.autoDetectCapabilities()

		for range b.ticker.C {
			b.autoDetectCapabilities()
		}
	}()
}

// Stop causes the background process to stop auto detecting capabilities
func (b *Background) Stop() {
	b.ticker.Stop()
}

func (b *Background) autoDetectCapabilities() {
	previousState := make(map[string]bool)
	for _, gvk := range b.config.GroupVersionKinds {
		previousState[gvk.String()] = IsResourceAvailable(gvk)
	}

	b.DetectCapabilities()

	for _, gvk := range b.config.GroupVersionKinds {
		before := previousState[gvk.String()]
		after := IsResourceAvailable(gvk)

		if !before && after {
			if b.config.ExitOnChange {
				log.Info(fmt.Sprintf("%s is deployed in cluster. Restarting operator to enable all APIs ....", gvk.String()))
				os.Exit(1)
			} else {
				log.Info(fmt.Sprintf("%s is deployed in cluster", gvk.String()))
			}
		} else if !before && !after {
			log.Info(fmt.Sprintf("%s is not deployed in cluster", gvk.String()))
		} else if before && !after {
			if b.config.ExitOnChange {
				log.Info(fmt.Sprintf("%s is undeployed. Restarting operator to disable some APIs ....", gvk.String()))
				os.Exit(1)
			} else {
				log.Info(fmt.Sprintf("%s is undeployed in cluster", gvk.String()))
			}
		}
	}
}

// DetectCapabilities populates state manager
func (b *Background) DetectCapabilities() {
	_, apiLists, err := b.dc.ServerGroupsAndResources()
	if err != nil {
		log.Error(err, "Failed to get API List")
		return
	}

	stateManager := actions.GetStateManager()
	for _, gvk := range b.config.GroupVersionKinds {
		exists := false
		for _, apiList := range apiLists {
			if apiList.GroupVersion == gvk.GroupVersion().String() {
				for _, r := range apiList.APIResources {
					if r.Kind == gvk.Kind {
						exists = true
						break
					}
				}
			}
		}
		stateManager.SetState(gvk.String(), exists)
	}
}

// IsResourceAvailable gets the state from state manager
func IsResourceAvailable(gvk schema.GroupVersionKind) bool {
	stateManager := actions.GetStateManager()
	available, _ := stateManager.GetState(gvk.String()).(bool)

	return available
}
