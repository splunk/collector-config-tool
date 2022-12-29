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
	"embed"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
)

//go:embed web/static
var staticFiles embed.FS

func createMux(logger *zap.Logger, host component.Host, supportedReceivers []component.ID, supportedProcessors []component.ID, supportedExporters []component.ID) (http.Handler, error) {
	mux := http.NewServeMux()
	static, err := fs.Sub(staticFiles, "web/static")
	if err != nil {
		return nil, err
	}
	staticFS := http.FS(static)
	mux.Handle("/", http.FileServer(staticFS))

	registry := &componentRegistry{
		host:                host,
		supportedExporters:  supportedExporters,
		supportedProcessors: supportedProcessors,
		supportedReceivers:  supportedReceivers,
	}

	mux.Handle("/components", &componentHandler{
		logger:   logger,
		registry: registry,
	})

	ottoPipeline := &pipeline{}
	mux.Handle("/cfg-metadata/", &cfgschemaHandler{
		logger:   logger,
		pipeline: ottoPipeline,
	})

	mux.Handle("/jsonToYAML", &jsonToYAMLHandler{
		logger: logger,
	})

	wsHandlers := map[string]wsHandler{}
	registerReceiverHandlers(logger, registry, host, wsHandlers, ottoPipeline)
	registerProcessorHandlers(logger, registry, host, wsHandlers, ottoPipeline)
	registerExporterHandlers(logger, registry, host, wsHandlers, ottoPipeline)
	mux.Handle("/ws/", httpWsHandler{handlers: wsHandlers})

	return mux, nil
}

func registerReceiverHandlers(logger *zap.Logger, registry *componentRegistry, host component.Host, handlers map[string]wsHandler, ppln *pipeline) {
	for componentName, factory := range registry.receivers() {
		const componentType = "receiver"
		path := "/ws/" + componentType + "/" + componentName.String()
		handlers[path] = receiverSocketHandler{
			host:            host,
			logger:          logger,
			pipeline:        ppln,
			receiverFactory: factory,
		}
	}
}

func registerProcessorHandlers(logger *zap.Logger, registry *componentRegistry, host component.Host, handlers map[string]wsHandler, ppln *pipeline) {
	for componentName, factory := range registry.processors() {
		const componentType = "processor"
		path := "/ws/" + componentType + "/" + componentName.String()
		handlers[path] = processorSocketHandler{
			host:             host,
			logger:           logger,
			pipeline:         ppln,
			processorFactory: factory,
		}
	}
}

func registerExporterHandlers(logger *zap.Logger, registry *componentRegistry, host component.Host, handlers map[string]wsHandler, ppln *pipeline) {
	for componentName, factory := range registry.exporters() {
		const componentType = "exporter"
		path := "/ws/" + componentType + "/" + componentName.String()
		handlers[path] = exporterSocketHandler{
			host:            host,
			logger:          logger,
			pipeline:        ppln,
			exporterFactory: factory,
		}
	}
}
