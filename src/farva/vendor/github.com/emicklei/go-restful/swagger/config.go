/*
Copyright 2016 Planet Labs 

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package swagger

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

// PostBuildDeclarationMapFunc can be used to modify the api declaration map.
type PostBuildDeclarationMapFunc func(apiDeclarationMap *ApiDeclarationList)

type Config struct {
	// url where the services are available, e.g. http://localhost:8080
	// if left empty then the basePath of Swagger is taken from the actual request
	WebServicesUrl string
	// path where the JSON api is avaiable , e.g. /apidocs
	ApiPath string
	// [optional] path where the swagger UI will be served, e.g. /swagger
	SwaggerPath string
	// [optional] location of folder containing Swagger HTML5 application index.html
	SwaggerFilePath string
	// api listing is constructed from this list of restful WebServices.
	WebServices []*restful.WebService
	// will serve all static content (scripts,pages,images)
	StaticHandler http.Handler
	// [optional] on default CORS (Cross-Origin-Resource-Sharing) is enabled.
	DisableCORS bool
	// Top-level API version. Is reflected in the resource listing.
	ApiVersion string
	// If set then call this handler after building the complete ApiDeclaration Map
	PostBuildHandler PostBuildDeclarationMapFunc
	// Swagger global info struct
	Info Info
}
