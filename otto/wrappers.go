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
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processortest"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

type metricsReceiverWrapper struct {
	receiver.Metrics
	repeater *metricsRepeater
}

func newMetricsReceiverWrapper(logger *zap.Logger, ws *websocket.Conn, cfg component.Config, factory receiver.Factory) (metricsReceiverWrapper, error) {
	repeater := newMetricsRepeater(logger, ws)
	r, err := factory.CreateMetricsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		cfg,
		repeater,
	)
	return metricsReceiverWrapper{
		Metrics:  r,
		repeater: repeater,
	}, err
}

func (w metricsReceiverWrapper) setNextMetricsConsumer(next consumer.Metrics) {
	w.repeater.setNext(next)
}

func (w metricsReceiverWrapper) waitForStopMessage() {
	w.repeater.waitForStopMessage()
}

func (w metricsReceiverWrapper) Start(ctx context.Context, host component.Host) error {
	return w.Metrics.Start(ctx, host)
}

func (w metricsReceiverWrapper) Shutdown(ctx context.Context) error {
	return w.Metrics.Shutdown(ctx)
}

type logsReceiverWrapper struct {
	receiver.Logs
	repeater *logsRepeater
}

func newLogsReceiverWrapper(logger *zap.Logger, ws *websocket.Conn, cfg component.Config, factory receiver.Factory) (logsReceiverWrapper, error) {
	repeater := newLogsRepeater(logger, ws)
	r, err := factory.CreateLogsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		cfg,
		repeater,
	)
	return logsReceiverWrapper{
		Logs:     r,
		repeater: repeater,
	}, err
}

func (w logsReceiverWrapper) setNextLogsConsumer(next consumer.Logs) {
	w.repeater.setNext(next)
}

func (w logsReceiverWrapper) waitForStopMessage() {
	w.repeater.waitForStopMessage()
}

func (w logsReceiverWrapper) Start(ctx context.Context, host component.Host) error {
	return w.Logs.Start(ctx, host)
}

func (w logsReceiverWrapper) Shutdown(ctx context.Context) error {
	return w.Logs.Shutdown(ctx)
}

type tracesReceiverWrapper struct {
	receiver.Traces
	repeater *tracesRepeater
}

func newTracesReceiverWrapper(logger *zap.Logger, ws *websocket.Conn, cfg component.Config, factory receiver.Factory) (tracesReceiverWrapper, error) {
	repeater := newTracesRepeater(logger, ws)
	r, err := factory.CreateTracesReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		cfg,
		repeater,
	)
	return tracesReceiverWrapper{
		Traces:   r,
		repeater: repeater,
	}, err
}

func (w tracesReceiverWrapper) setNextTracesConsumer(next consumer.Traces) {
	w.repeater.setNext(next)
}

func (w tracesReceiverWrapper) waitForStopMessage() {
	w.repeater.waitForStopMessage()
}

func (w tracesReceiverWrapper) Start(ctx context.Context, host component.Host) error {
	return w.Traces.Start(ctx, host)
}

func (w tracesReceiverWrapper) Shutdown(ctx context.Context) error {
	return w.Traces.Shutdown(ctx)
}

type metricsProcessorWrapper struct {
	processor.Metrics
	repeater *metricsRepeater
}

func newMetricsProcessorWrapper(logger *zap.Logger, ws *websocket.Conn, cfg component.Config, factory processor.Factory) (metricsProcessorWrapper, error) {
	repeater := newMetricsRepeater(logger, ws)
	p, err := factory.CreateMetricsProcessor(
		context.Background(),
		processortest.NewNopCreateSettings(),
		cfg,
		repeater,
	)
	return metricsProcessorWrapper{
		Metrics:  p,
		repeater: repeater,
	}, err
}

func (w metricsProcessorWrapper) setNextMetricsConsumer(next consumer.Metrics) {
	w.repeater.setNext(next)
}

func (w metricsProcessorWrapper) waitForStopMessage() {
	w.repeater.waitForStopMessage()
}

func (w metricsProcessorWrapper) Start(ctx context.Context, host component.Host) error {
	return w.Metrics.Start(ctx, host)
}

func (w metricsProcessorWrapper) Shutdown(ctx context.Context) error {
	return w.Metrics.Shutdown(ctx)
}

type logsProcessorWrapper struct {
	processor.Logs
	repeater *logsRepeater
}

func newLogsProcessorWrapper(logger *zap.Logger, ws *websocket.Conn, cfg component.Config, factory processor.Factory) (logsProcessorWrapper, error) {
	repeater := newLogsRepeater(logger, ws)
	p, err := factory.CreateLogsProcessor(
		context.Background(),
		processortest.NewNopCreateSettings(),
		cfg,
		repeater,
	)
	return logsProcessorWrapper{
		Logs:     p,
		repeater: repeater,
	}, err
}

func (w logsProcessorWrapper) setNextLogsConsumer(next consumer.Logs) {
	w.repeater.setNext(next)
}

func (w logsProcessorWrapper) waitForStopMessage() {
	w.repeater.waitForStopMessage()
}

func (w logsProcessorWrapper) Start(ctx context.Context, host component.Host) error {
	return w.Logs.Start(ctx, host)
}

func (w logsProcessorWrapper) Shutdown(ctx context.Context) error {
	return w.Logs.Shutdown(ctx)
}

type tracesProcessorWrapper struct {
	processor.Traces
	repeater *tracesRepeater
}

func newTracesProcessorWrapper(logger *zap.Logger, ws *websocket.Conn, cfg component.Config, factory processor.Factory) (tracesProcessorWrapper, error) {
	repeater := newTracesRepeater(logger, ws)
	p, err := factory.CreateTracesProcessor(
		context.Background(),
		processortest.NewNopCreateSettings(),
		cfg,
		repeater,
	)
	return tracesProcessorWrapper{
		Traces:   p,
		repeater: repeater,
	}, err
}

func (w tracesProcessorWrapper) setNextTracesConsumer(next consumer.Traces) {
	w.repeater.setNext(next)
}

func (w tracesProcessorWrapper) waitForStopMessage() {
	w.repeater.waitForStopMessage()
}

func (w tracesProcessorWrapper) Start(ctx context.Context, host component.Host) error {
	return w.Traces.Start(ctx, host)
}

func (w tracesProcessorWrapper) Shutdown(ctx context.Context) error {
	return w.Traces.Shutdown(ctx)
}

type metricsExporterWrapper struct {
	exporter.Metrics
	repeater *metricsRepeater
}

func newMetricsExporterWrapper(logger *zap.Logger, ws *websocket.Conn, cfg component.Config, factory exporter.Factory) (metricsExporterWrapper, error) {
	repeater := newMetricsRepeater(logger, ws)
	e, err := factory.CreateMetricsExporter(
		context.Background(),
		exportertest.NewNopCreateSettings(),
		cfg,
	)
	return metricsExporterWrapper{
		Metrics:  e,
		repeater: repeater,
	}, err
}

func (w metricsExporterWrapper) waitForStopMessage() {
	w.repeater.waitForStopMessage()
}

func (w metricsExporterWrapper) Start(ctx context.Context, host component.Host) error {
	return w.Metrics.Start(ctx, host)
}

func (w metricsExporterWrapper) Shutdown(ctx context.Context) error {
	return w.Metrics.Shutdown(ctx)
}

type logsExporterWrapper struct {
	exporter.Logs
	repeater *logsRepeater
}

func newLogsExporterWrapper(logger *zap.Logger, ws *websocket.Conn, cfg component.Config, factory exporter.Factory) (logsExporterWrapper, error) {
	settings := exportertest.NewNopCreateSettings()
	settings.TelemetrySettings.Logger = logger
	e, err := factory.CreateLogsExporter(
		context.Background(),
		settings,
		cfg,
	)
	return logsExporterWrapper{
		Logs:     e,
		repeater: newLogsRepeater(logger, ws),
	}, err
}

func (w logsExporterWrapper) waitForStopMessage() {
	w.repeater.waitForStopMessage()
}

func (w logsExporterWrapper) Start(ctx context.Context, host component.Host) error {
	return w.Logs.Start(ctx, host)
}

func (w logsExporterWrapper) Shutdown(ctx context.Context) error {
	return w.Logs.Shutdown(ctx)
}

type tracesExporterWrapper struct {
	exporter.Traces
	repeater *tracesRepeater
}

func newTracesExporterWrapper(logger *zap.Logger, ws *websocket.Conn, cfg component.Config, factory exporter.Factory) (tracesExporterWrapper, error) {
	e, err := factory.CreateTracesExporter(
		context.Background(),
		exportertest.NewNopCreateSettings(),
		cfg,
	)
	return tracesExporterWrapper{
		Traces:   e,
		repeater: newTracesRepeater(logger, ws),
	}, err
}

func (w tracesExporterWrapper) waitForStopMessage() {
	w.repeater.waitForStopMessage()
}

func (w tracesExporterWrapper) Start(ctx context.Context, host component.Host) error {
	return w.Traces.Start(ctx, host)
}

func (w tracesExporterWrapper) Shutdown(ctx context.Context) error {
	return w.Traces.Shutdown(ctx)
}
