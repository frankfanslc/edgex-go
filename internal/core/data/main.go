/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2019 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package data

import (
	"context"
	"os"

	"github.com/edgexfoundry/edgex-go"
	"github.com/edgexfoundry/edgex-go/internal"
	"github.com/edgexfoundry/edgex-go/internal/core/data/config"
	"github.com/edgexfoundry/edgex-go/internal/core/data/container"
	"github.com/edgexfoundry/edgex-go/internal/core/data/messaging"
	pkgHandlers "github.com/edgexfoundry/edgex-go/internal/pkg/bootstrap/handlers"
	"github.com/edgexfoundry/edgex-go/internal/pkg/telemetry"

	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/flags"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/handlers"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/interfaces"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"

	"github.com/gorilla/mux"
)

func Main(ctx context.Context, cancel context.CancelFunc, router *mux.Router) {
	startupTimer := startup.NewStartUpTimer(common.CoreDataServiceKey)

	// All common command-line flags have been moved to DefaultCommonFlags. Service specific flags can be add here,
	// by inserting service specific flag prior to call to commonFlags.Parse().
	// Example:
	// 		flags.FlagSet.StringVar(&myvar, "m", "", "Specify a ....")
	//      ....
	//      flags.Parse(os.Args[1:])
	//
	f := flags.New()
	f.Parse(os.Args[1:])

	configuration := &config.ConfigurationStruct{}
	dic := di.NewContainer(di.ServiceConstructorMap{
		container.ConfigurationName: func(get di.Get) interface{} {
			return configuration
		},
	})

	httpServer := handlers.NewHttpServer(router, true)

	bootstrap.Run(
		ctx,
		cancel,
		f,
		common.CoreDataServiceKey,
		internal.ConfigStemCore,
		configuration,
		startupTimer,
		dic,
		true,
		[]interfaces.BootstrapHandler{
			pkgHandlers.NewDatabase(httpServer, configuration, container.DBClientInterfaceName).BootstrapHandler, // add v2 db client bootstrap handler
			messaging.BootstrapHandler,
			NewBootstrap(router, common.CoreDataServiceKey).BootstrapHandler,
			telemetry.BootstrapHandler,
			httpServer.BootstrapHandler,
			handlers.NewStartMessage(common.CoreDataServiceKey, edgex.Version).BootstrapHandler,
		},
	)
}
