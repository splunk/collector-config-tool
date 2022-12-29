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
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

type metricsRepeater struct {
	logger    *zap.Logger
	ws        *websocket.Conn
	marshaler pmetric.Marshaler
	next      consumer.Metrics
	stop      chan struct{}
}

func newMetricsRepeater(logger *zap.Logger, ws *websocket.Conn) *metricsRepeater {
	return &metricsRepeater{
		logger:    logger,
		ws:        ws,
		marshaler: &pmetric.JSONMarshaler{},
		stop:      make(chan struct{}),
	}
}

func (r *metricsRepeater) setNext(next consumer.Metrics) {
	r.next = next
}

func (*metricsRepeater) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{}
}

func (r *metricsRepeater) ConsumeMetrics(ctx context.Context, pmetrics pmetric.Metrics) error {
	err := doWritePayload(r.ws, metrics(pmetrics), r.stop)
	if err != nil {
		r.logger.Error("metricsRepeater", zap.Error(err))
		return nil
	}
	if r.next == nil {
		return nil
	}
	return r.next.ConsumeMetrics(ctx, pmetrics)
}

func (r *metricsRepeater) waitForStopMessage() {
	<-r.stop
}

type logsRepeater struct {
	logger    *zap.Logger
	ws        *websocket.Conn
	marshaler plog.Marshaler
	next      consumer.Logs
	stop      chan struct{}
}

func newLogsRepeater(logger *zap.Logger, ws *websocket.Conn) *logsRepeater {
	return &logsRepeater{
		logger:    logger,
		ws:        ws,
		marshaler: &plog.JSONMarshaler{},
		stop:      make(chan struct{}),
	}
}

func (*logsRepeater) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{}
}

func (r *logsRepeater) waitForStopMessage() {
	<-r.stop
}

func (r *logsRepeater) ConsumeLogs(ctx context.Context, plogs plog.Logs) error {
	err := doWritePayload(r.ws, logs(plogs), r.stop)
	if err != nil {
		r.logger.Error("logsRepeater", zap.Error(err))
		return nil
	}
	if r.next == nil {
		return nil
	}
	return r.next.ConsumeLogs(ctx, plogs)
}

func (r *logsRepeater) setNext(next consumer.Logs) {
	r.next = next
}

type tracesRepeater struct {
	logger    *zap.Logger
	ws        *websocket.Conn
	marshaler ptrace.Marshaler
	next      consumer.Traces
	stop      chan struct{}
}

func newTracesRepeater(logger *zap.Logger, ws *websocket.Conn) *tracesRepeater {
	return &tracesRepeater{
		logger:    logger,
		ws:        ws,
		marshaler: &ptrace.JSONMarshaler{},
		stop:      make(chan struct{}),
	}
}

func (r *tracesRepeater) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{}
}

func (r *tracesRepeater) ConsumeTraces(ctx context.Context, ptraces ptrace.Traces) error {
	err := doWritePayload(r.ws, traces(ptraces), r.stop)
	if err != nil {
		r.logger.Error("tracesRepeater", zap.Error(err))
		return nil
	}
	if r.next == nil {
		return nil
	}
	return r.next.ConsumeTraces(ctx, ptraces)
}

func (r *tracesRepeater) setNext(next consumer.Traces) {
	r.next = next
}

func (r *tracesRepeater) waitForStopMessage() {
	<-r.stop
}

func doWritePayload(ws *websocket.Conn, payload json.Marshaler, stop chan struct{}) error {
	err := writePayload(ws, payload)
	if err != nil {
		stop <- struct{}{}
	}
	return err
}

func writePayload(ws *websocket.Conn, payload json.Marshaler) error {
	envelopeJson, err := json.Marshal(wsMessageEnvelope{
		Payload: payload,
	})
	if err != nil {
		return fmt.Errorf("error marshaling envelope: %w", err)
	}
	_, err = ws.Write(envelopeJson)
	if err != nil {
		return fmt.Errorf("error writing envelope json to websocket: %v\n", err)
	}
	return nil
}
