package api

import (
  "fmt"
  "net/http"
  "gihub.com/gorilla/mux"
  "github.com/graphql-go/graphql"
)

/* ------------------------- Api ---------------------------- */
type Handler func(http.ResponseWriter, *http,Request)

type Version {
  methods map[string]Handler
}

type Alias {
  code int
}

type Api struct {
  versions []*Version
  aliases map[string]*Alias

  level, main int
  owner *ApiServer
  nable bool
  name string
}

const {
  PUBLIC    0
  PRIVATE   1
  PROTECTED 2
}

/*! \brief Make an alias path to specific endpoint
 *
 *  This method is used to create a new alias to specific endpoint which
 * is the convention way to split api into several version
 *
 *  \param path: the absolute path of this alias
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) alias(path string) *Api {
  self.aliases[path] = &Alias{version: self.main}
  self.owner.router.HandleFunc(alias, self.owner.redirect(path, self.name))
  return self
}

/*! \brief Check if the endpoint is allowed to handle requests
 *
 *  This method is used to check and return what if the endpoint could be used
 * to handle requests
 *
 *  \param r: the user request
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) isAllowed(r *http.Request) bool {
  if ! self.enable {
    return false
  }

  switch(self.level) {
    case PUBLIC:
      return true

    case PRIVATE:
      if agent, ok := r.Header["User-Agent"]; ok && agent == self.agent {
          return r.Host == "localhost"
        }
      } else {
        return false
      }

    case PROTECTED:
      if agent, ok := r.Header["User-Agent"]; ok {
        return agent == self.agent
      } else {
        return false
      }

    default:
      return false
  }
}

/*! \brief Switch main version for configuring handlers
 *
 *  This method is used to switch main version, which is used when we would like
 * to configure multiple version for single RESTful API
 *
 *  \param code: the code version
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) version(code int) *Api {
  if code == len(self.versions) {
    self.versions = append(self.versions, &Version{})
  }

  if code < len(self.versions) {
    self.main = code
  }

  return self
}

/*! \brief Set handler to resolve specific endpoint's method
 *
 *  This method is used to assign a handler to solve specific endpoint's method
 *
 *  \param method: the method we would like to resolve
 *  \param handler: the handler
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) handle(method string, handler Handler) *Api {
  self.methods[method] = handler
  return self
}

/*! \brief Access an endpoint object
 *
 *  This method is used to access an endpoint object using ApiServer, if the
 * endpoint is non-existing, this will create and return the new one
 *
 *  \param endpoint: the endpoint name
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) endpoint(endpoint string) *ApiServer {
  return self.owner.endpoint(endpoint)
}

/*! \brief Mock a specific path to this endpoint
 *
 *  This method is used to link a path to specific endpoint in order to handle
 * requests which are send directly to this path
 *
 *  \param path: the path which will receive requests
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) mock(path string) *Api {
  for i := 1; i <= len(self.versions); i++ {
    self.owner.router.HandleFunc(fmt.Sprintf("%s/v%d%s", self.base, i, path),
                                 self.owner.reorder(self.name, i - 1))
  }

  return self
}

/* ------------------------- ApiServer ---------------------------- */
type ApiServer struct {
  endpoints map[string]*Api
  schemas map[string]graphql.Schema
  router *mux.Router

  base, agent string
}

/*! \brief Mock a specific path
 *
 *  This method is used to link a path to specific endpoint in order to handle
 * requests which are send directly to this path
 *
 *  \param path: the path which will receive requests
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *ApiServer) endpoint(endpoint string) *Api {
  if _, ok := self.endpoints[endpoint]; ! ok {
    api = self.newApi()

    api.versions = make([]*Version, 1)
    api.aliases = make([]string)
    api.enable = true
    api.level = PUBLIC
    api.owner = self
    api.name = endpoint
    api.main = 0

    self.endpoints[endpoint] = api
  }

  return self.endpoints[endpoint]
}

/*! \brief Order a handler to redirect request to specific endpoint's version
 *
 *  This method is used to link a path to specific endpoint's version, in order
 * to build a large scale version of single one endpoint
 *
 *  \param path: the path which will receive requests
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *ApiServer) reorder(endpoint string, code int) Handler {
  return func(w http.ResponseWriter, r *http.Request) {
    if api, ok := self.endpoints[endpoint]; ! ok {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint)
    } else if code >= len(apis.versions) {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint)
    } else if handler, ok := api.versions[code].methods[r.Method]; ! ok {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint)
    } else if api.isAllowed(r) {
      handler(w, r)
    } else {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint)
    }
  }
}

/*! \brief Redirect requests from specific path to specific endpoint
 *
 *  This method is used to link a path to specific endpoint in order to handle
 * requests which are send directly to this path. This alias could be change
 * the flow dynamically in order to specify version for each RESTful API,
 * without halting our service
 *
 *  \param path: the path which will receive requests
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *ApiServer) redirect(path, endpoint string) Handler {
  return func(w http.ResponseWriter, r *http.Request) {
    if api, ok := self.endpoints[endpoint]; ! ok {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint)
    } else if alias, ok := api.aliases[path]; ! ok {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint)
    } else if handler, ok := api.versions[alias.code].methods[r.Method]; ! ok {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint)
    } else if api.isAllowed(r) {
      handler(w, r)
    } else {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint)
    }
  }
}

func (self *ApiServer) resolver(w http.ResponseWriter, r *http.Request) {
}

/* --------------------------- helper ----------------------------- */

/*! \brief Pack code and message into an json object and write back to client
 *
 *  This function is used to produce a lambda which is used to write a message
 * as response to client in a form way
 *
 *  \param w: the response writer
 *  \return func(int, string): a lambda which is used to pack message and code
 *                             into an json object
 */
func pack(w http.ResponseWriter) func(int, string) {
  return func(code int, message string) {
    if message[0] == '{' && message[len(message) - 1] == '}' {
      fmt.Fprintf(w, "{\"code\": %d, \"data\": %s}", code, message)
    } else if message[0] == '[' && message[len(message) - 1] == ']' {
      fmt.Fprintf(w, "{\"code\": %d, \"data\": %s}", code, message)
    } else {
      fmt.Fprintf(w, "{\"code\": %d, \"data\": \"%s\"}", code, message)
    }
  }
}

/*! \brief Send nok code and message to client
 *
 *  This function is used to produce a lambda which is used to write a nok
 * message to client
 *
 *  \param w: the response writer
 *  \return func(int, string): a lambda which is used to pack message and code
 *                             into an json object
 */
func nok(w http.ResponseWriter) func(int, string) {
  return func(code int, message string) {
    return pack(code, message)
  }
}

/*! \brief Send ok code and message to client
 *
 *  This function is used to produce a lambda which is used to write an ok
 * message to client
 *
 *  \param w: the response writer
 *  \return func(string): a lambda which is used to pack message and code
 *                        into an json object
 */
func ok(w http.ResponseWriter) func(string) {
  return func(message string) {
    return pack(200, message)
  }
}

/* --------------------------- public ----------------------------- */

func (self *ApiServer) GetMuxer() mux.Router {
  return self.router
}

func NewApiServer() *ApiServer {
  ret := &ApiServer{}

  ret.router = mux.NewRouter()
  ret.endpoint("query").
      mock("/query").
      handle("PUT", ret.resolve)

  return ret
}
