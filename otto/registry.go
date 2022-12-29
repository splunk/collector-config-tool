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
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"sort"
)

type componentRegistry struct {
	logger              *zap.Logger
	supportedReceivers  []component.ID
	supportedProcessors []component.ID
	supportedExporters  []component.ID
	host                component.Host
}

func (h *componentRegistry) factoriesToComponentTypes() componentTypes {
	out := componentTypes{}
	for _, name := range h.supportedReceivers {
		factory := h.host.GetFactory(component.KindReceiver, name.Type())
		if factory == nil {
			h.logger.Warn("unknown receiver", zap.Any("factory", name.Type()))
			continue
		}
		sp := receiverSupportedPipelines(factory)
		if sp.metrics {
			out.Metrics.Receivers = append(out.Metrics.Receivers, name.String())
		}
		if sp.traces {
			out.Traces.Receivers = append(out.Traces.Receivers, name.String())
		}
		if sp.logs {
			out.Logs.Receivers = append(out.Logs.Receivers, name.String())
		}
	}
	sort.Strings(out.Metrics.Receivers)
	sort.Strings(out.Traces.Receivers)
	sort.Strings(out.Logs.Receivers)

	for _, name := range h.supportedProcessors {
		factory := h.host.GetFactory(component.KindProcessor, name.Type())
		if factory == nil {
			h.logger.Warn("unknown receiver", zap.Any("factory", name.Type()))
			continue
		}
		sp := processorSupportedPipelines(factory)
		if sp.metrics {
			out.Metrics.Processors = append(out.Metrics.Processors, name.String())
		}
		if sp.traces {
			out.Traces.Processors = append(out.Traces.Processors, name.String())
		}
		if sp.logs {
			out.Logs.Processors = append(out.Logs.Processors, name.String())
		}
	}
	sort.Strings(out.Metrics.Processors)
	sort.Strings(out.Traces.Processors)
	sort.Strings(out.Logs.Processors)

	for _, name := range h.supportedExporters {
		factory := h.host.GetFactory(component.KindExporter, name.Type())
		if factory == nil {
			h.logger.Warn("unknown receiver", zap.Any("factory", name.Type()))
			continue
		}
		sp := exporterSupportedPipelines(factory)
		if sp.metrics {
			out.Metrics.Exporters = append(out.Metrics.Exporters, name.String())
		}
		if sp.traces {
			out.Traces.Exporters = append(out.Traces.Exporters, name.String())
		}
		if sp.logs {
			out.Logs.Exporters = append(out.Logs.Exporters, name.String())
		}
	}
	sort.Strings(out.Metrics.Exporters)
	sort.Strings(out.Traces.Exporters)
	sort.Strings(out.Logs.Exporters)

	return out
}

func (h *componentRegistry) receivers() map[component.ID]receiver.Factory {
	components := make(map[component.ID]receiver.Factory, len(h.supportedReceivers))
	for _, name := range h.supportedReceivers {
		components[name] = h.host.GetFactory(component.KindReceiver, name.Type()).(receiver.Factory)
	}
	return components
}

func (h *componentRegistry) processors() map[component.ID]processor.Factory {
	components := make(map[component.ID]processor.Factory, len(h.supportedProcessors))
	for _, name := range h.supportedProcessors {
		components[name] = h.host.GetFactory(component.KindProcessor, name.Type()).(processor.Factory)
	}
	return components
}

func (h *componentRegistry) exporters() map[component.ID]exporter.Factory {
	components := make(map[component.ID]exporter.Factory, len(h.supportedExporters))
	for _, name := range h.supportedExporters {
		components[name] = h.host.GetFactory(component.KindExporter, name.Type()).(exporter.Factory)
	}
	return components
}
