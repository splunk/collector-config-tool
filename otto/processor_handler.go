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
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

type processorSocketHandler struct {
	host             component.Host
	logger           *zap.Logger
	pipeline         *pipeline
	processorFactory processor.Factory
}

func (h processorSocketHandler) handle(ws *websocket.Conn) {
	err := h.doHandle(ws)
	if err != nil {
		sendErr(ws, h.logger, "processorSocketHandler", err)
	}
}

func (h processorSocketHandler) doHandle(ws *websocket.Conn) error {
	pipelineType, conf, err := readSocket(ws)
	if err != nil {
		return err
	}
	processorConfig := h.processorFactory.CreateDefaultConfig()
	err = unmarshalProcessorConfig(processorConfig, conf)
	if err != nil {
		return err
	}
	switch pipelineType {
	case "metrics":
		h.attachMetricsProcessor(ws, processorConfig)
	case "logs":
		h.attachLogsProcessor(ws, processorConfig)
	case "traces":
		h.attachTracesProcessor(ws, processorConfig)
	}
	return nil
}

func (h processorSocketHandler) attachMetricsProcessor(
	ws *websocket.Conn,
	cfg component.Config,
) {
	wrapper, err := newMetricsProcessorWrapper(h.logger, ws, cfg, h.processorFactory)
	if err != nil {
		sendErr(ws, h.logger, "failed to create metrics processor wrapper", err)
		return
	}
	h.pipeline.connectMetricsProcessorWrapper(wrapper)
	err = wrapper.Start(context.Background(), h.host)
	if err != nil {
		sendErr(ws, h.logger, "failed to start metrics processor", err)
		return
	}
	wrapper.waitForStopMessage()
	h.pipeline.disconnectMetricsProcessorWrapper()
}

func (h processorSocketHandler) attachLogsProcessor(
	ws *websocket.Conn,
	cfg component.Config,
) {
	wrapper, err := newLogsProcessorWrapper(h.logger, ws, cfg, h.processorFactory)
	if err != nil {
		sendErr(ws, h.logger, "failed to create logs processor wrapper", err)
		return
	}
	h.pipeline.connectLogsProcessorWrapper(wrapper)
	err = wrapper.Start(context.Background(), h.host)
	if err != nil {
		sendErr(ws, h.logger, "failed to start logs processor", err)
		return
	}
	wrapper.waitForStopMessage()
	h.pipeline.disconnectLogsProcessorWrapper()
}

func (h processorSocketHandler) attachTracesProcessor(
	ws *websocket.Conn,
	cfg component.Config,
) {
	wrapper, err := newTracesProcessorWrapper(h.logger, ws, cfg, h.processorFactory)
	if err != nil {
		sendErr(ws, h.logger, "failed to create traces processor wrapper", err)
		return
	}
	h.pipeline.connectTracesProcessorWrapper(wrapper)
	err = wrapper.Start(context.Background(), h.host)
	if err != nil {
		sendErr(ws, h.logger, "failed to start traces processor", err)
		return
	}
	wrapper.waitForStopMessage()
	h.pipeline.disconnectTracesProcessorWrapper()
}

func unmarshalProcessorConfig(processorConfig component.Config, conf *confmap.Conf) error {
	if unmarshallable, ok := processorConfig.(confmap.Unmarshaler); ok {
		return unmarshallable.Unmarshal(conf)
	}
	return conf.Unmarshal(processorConfig)
}
