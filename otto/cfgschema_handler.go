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
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed cfg-metadata
var cfgMetadata embed.FS

type cfgschemaHandler struct {
	logger   *zap.Logger
	pipeline *pipeline
}

func (h *cfgschemaHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// e.g. "/cfg-metadata/receiver/redis"
	parts := strings.Split(req.RequestURI, "/")
	componentType := parts[2]
	componentName := parts[3]

	path := filepath.Join("cfg-metadata", componentType, componentName+".yaml")
	yamlFile, err := cfgMetadata.ReadFile(path)
	if err != nil {
		h.logger.Info("cfgschemaHandler: ServeHTTP: error reading yaml file", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	yamlMap := map[string]any{}
	err = yaml.Unmarshal(yamlFile, &yamlMap)
	if err != nil {
		h.logger.Info("cfgschemaHandler: ServeHTTP: error unmarshaling yaml", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	converted, err := convertAnyMapsToStringMaps(yamlMap)
	if err != nil {
		h.logger.Info("cfgschemaHandler: ServeHTTP: error coercing any keys to string keys", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	fjson, err := json.Marshal(converted)
	if err != nil {
		h.logger.Info("cfgschemaHandler: ServeHTTP: error marshaling fieldInfo to json", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(fjson)
	if err != nil {
		h.logger.Info("cfgschemaHandler: ServeHTTP: error writing response", zap.Error(err))
	}
}

func convertAnyMapsToStringMaps(v any) (any, error) {
	if strMap, ok := v.(map[string]any); ok {
		return convertStrMap(strMap)
	} else if anyMap, ok := v.(map[any]any); ok {
		return convertAnyMap(anyMap)
	} else if anySlice, ok := v.([]any); ok {
		return convertSlice(anySlice)
	}
	return v, nil
}

func convertStrMap(strMap map[string]any) (any, error) {
	for k, v := range strMap {
		converted, err := convertAnyMapsToStringMaps(v)
		if err != nil {
			return nil, err
		}
		strMap[k] = converted
	}
	return strMap, nil
}

func convertAnyMap(anyMap map[any]any) (any, error) {
	out := map[string]any{}
	for k, v := range anyMap {
		strK, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("unable to convert key to string: %v", k)
		}
		converted, err := convertAnyMapsToStringMaps(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert any map value: %v: %w", v, err)
		}
		out[strK] = converted
	}
	return out, nil
}

func convertSlice(anySlice []any) (any, error) {
	for i, v := range anySlice {
		converted, err := convertAnyMapsToStringMaps(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert any slice: %v: %w", v, err)
		}
		anySlice[i] = converted
	}
	return anySlice, nil
}
