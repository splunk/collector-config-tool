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

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type wsMessageEnvelope struct {
	Payload json.Marshaler
	Error   error
}

type metrics pmetric.Metrics

var (
	metricsMarshaler = &pmetric.JSONMarshaler{}
	tracesMarshaler  = &ptrace.JSONMarshaler{}
	logsMarshaler    = &plog.JSONMarshaler{}
)

func (m metrics) MarshalJSON() ([]byte, error) {
	return metricsMarshaler.MarshalMetrics(pmetric.Metrics(m))
}

type traces ptrace.Traces

func (t traces) MarshalJSON() ([]byte, error) {
	return tracesMarshaler.MarshalTraces(ptrace.Traces(t))
}

type logs plog.Logs

func (l logs) MarshalJSON() ([]byte, error) {
	return logsMarshaler.MarshalLogs(plog.Logs(l))
}
