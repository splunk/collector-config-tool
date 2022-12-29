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
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"
	"net/http"
)

type OttoExtension struct {
	logger     *zap.Logger
	config     *Config
	settings   extension.CreateSettings
	httpServer *http.Server
}

func (o *OttoExtension) Start(_ context.Context, host component.Host) error {
	mux, err := createMux(o.logger, host, o.config.Receivers, o.config.Processors, o.config.Exporters)
	if err != nil {
		return err
	}
	o.logger.Info("starting otto http server", zap.String("endpoint", o.config.HTTPServerSettings.Endpoint))
	server, err := o.config.ToServer(host, o.settings.TelemetrySettings, mux)
	if err != nil {
		return err
	}
	o.httpServer = server
	o.httpServer.Addr = o.config.HTTPServerSettings.Endpoint
	go func() {
		err2 := o.httpServer.ListenAndServe()
		if err2 != http.ErrServerClosed {
			o.logger.Error("error starting otto server", zap.Error(err2))
		}
	}()
	return nil
}

func (o *OttoExtension) Shutdown(ctx context.Context) error {
	if o.httpServer == nil {
		return nil
	}
	return o.httpServer.Shutdown(ctx)
}
