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
	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
	"gopkg.in/yaml.v3"
)

func sendErr(ws *websocket.Conn, logger *zap.Logger, msg string, err error) {
	envelopeJson, jsonErr := json.Marshal(wsMessageEnvelope{Error: err})
	if jsonErr != nil {
		logger.Error("failed to marshal envelope containing the error", zap.String("msg", msg),
			zap.NamedError("error", err), zap.NamedError("JSON error", jsonErr))
	}
	_, err = ws.Write(envelopeJson)
	if err != nil {
		logger.Error("error sending error", zap.Error(err))
	}
}

func readSocket(ws *websocket.Conn) (string, *confmap.Conf, error) {
	msg, err := readStartComponentMessage(ws)
	if err != nil {
		return "", nil, err
	}
	m := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(msg.ComponentYAML), &m)
	if err != nil {
		return "", nil, err
	}
	return msg.PipelineType, confmap.NewFromStringMap(m), nil
}
