// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otto

import (
	"encoding/json"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"net/http"
)

type componentHandler struct {
	logger   *zap.Logger
	registry *componentRegistry
}

func (h *componentHandler) ServeHTTP(resp http.ResponseWriter, _ *http.Request) {
	jsn, err := h.factoriesToComponentTypeJSON()
	if err != nil {
		h.logger.Info("componentHandler: ServeHTTP: error getting components", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = resp.Write(jsn)
	if err != nil {
		h.logger.Info("componentHandler: ServeHTTP: error writing response", zap.Error(err))
	}
}

func (h *componentHandler) factoriesToComponentTypeJSON() ([]byte, error) {
	cmp := h.registry.factoriesToComponentTypes()
	return json.Marshal(cmp)
}

type componentTypes struct {
	Metrics rpe `json:"metrics"`
	Logs    rpe `json:"logs"`
	Traces  rpe `json:"traces"`
}

type rpe struct {
	Receivers  []string `json:"receivers"`
	Processors []string `json:"processors"`
	Exporters  []string `json:"exporters"`
}

type supportedPipelines struct {
	metrics bool
	traces  bool
	logs    bool
}

func receiverSupportedPipelines(fact component.Factory) supportedPipelines {
	out := supportedPipelines{}
	if f, ok := fact.(receiver.Factory); ok {
		if f.MetricsReceiverStability() != component.StabilityLevelUndefined {
			out.metrics = true
		}
		if f.TracesReceiverStability() != component.StabilityLevelUndefined {
			out.traces = true
		}
		if f.LogsReceiverStability() != component.StabilityLevelUndefined {
			out.logs = true
		}
	}
	return out
}

func processorSupportedPipelines(fact component.Factory) supportedPipelines {
	out := supportedPipelines{}
	if f, ok := fact.(processor.Factory); ok {
		if f.MetricsProcessorStability() != component.StabilityLevelUndefined {
			out.metrics = true
		}
		if f.TracesProcessorStability() != component.StabilityLevelUndefined {
			out.traces = true
		}
		if f.LogsProcessorStability() != component.StabilityLevelUndefined {
			out.logs = true
		}
	}
	return out
}

func exporterSupportedPipelines(fact component.Factory) supportedPipelines {
	out := supportedPipelines{}
	if f, ok := fact.(exporter.Factory); ok {
		if f.MetricsExporterStability() != component.StabilityLevelUndefined {
			out.metrics = true
		}
		if f.TracesExporterStability() != component.StabilityLevelUndefined {
			out.traces = true
		}
		if f.LogsExporterStability() != component.StabilityLevelUndefined {
			out.logs = true
		}
	}
	return out
}
