package api

import (
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
)

/* ------------------------- Api ---------------------------- */
type Handler func(http.ResponseWriter, *http.Request)

type Version struct {
  methods map[string]Handler
}

type Alias struct {
  code map[string]string
}

type Api struct {
  aliases map[string]*Alias
  versions map[string]*Version
  mainlines map[string]string

  level int
  owner *ApiServer
  enable bool
  name, main string
}

const (
  PUBLIC    = 0
  PRIVATE   = 1
  PROTECTED = 2
)

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
  if _, ok := self.aliases[path]; ok {
    panic(fmt.Sprintf("redefine alias %s", path))
  }

  self.aliases[path] = &Alias{}
  self.aliases[path].code = self.mainlines

  self.owner.router.HandleFunc(path,
    func(w http.ResponseWriter, r *http.Request){
      code := self.aliases[path].code[r.Method]

      if ver, ok := self.versions[code]; ! ok {
        nok(w)(404, fmt.Sprintf("Not found %s", path))
      } else if handler, ok := ver.methods[r.Method]; ! ok {
        nok(w)(404, fmt.Sprintf("Not found %s", path))
      } else if self.isAllowed(r) {
        handler(w, r)
      } else {
        nok(w)(404, fmt.Sprintf("Not found %s", path))
      }
    })
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
      if agents, ok := r.Header["User-Agent"]; ok {
        if agents[0] == self.owner.agent {
          return r.Host == 'https://localhost' || r.Host == 'http://localhost'
        } else {
          return false
        }
      } else {
        return false
      }

    case PROTECTED:
      if agents, ok := r.Header["User-Agent"]; ok {
        return agents[0] == self.owner.agent
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
  if _, ok := self.mainlines[method]; ! ok {
    self.mainlines[method] = self.main
  }

  self.versions[self.main].methods[method] = handler
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
  for ver, obj := range self.versions {
    var dest string

    if len(self.base) > 0 {
      dest = fmt.Sprintf("/%s/%s%s", self.base, ver, path)
    } else {
      dest = fmt.Sprintf("/%s%s", ver, path)
    }

    self.owner.router.HandleFunc(dest, self.owner.reorder(self.name, ver))
  }

  if len(self.base) > 0 {
    path = fmt.Sprintf("/%s%s", self.base, path)
  }

  return alias(path)
}

/* ------------------------- ApiServer ---------------------------- */
type ApiServer struct {
  endpoints map[string]*Api
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
func (self *ApiServer) reorder(endpoint, code string) Handler {
  return func(w http.ResponseWriter, r *http.Request) {
    if api, ok := self.endpoints[endpoint]; ! ok {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    } else if ver, ok := apis.versions[code]; ! ok {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    } else if handler, ok := ver.methods[r.Method]; ! ok {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    } else if api.isAllowed(r) {
      handler(w, r)
    } else {
      nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
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

func NewApiServer(user_agent string) *ApiServer {
  ret := &ApiServer{}

  ret.router = mux.NewRouter()
  ret.agent = user_agent
  ret.endpoint("query").
      mock("/query").
      handle("PUT", ret.resolve)

  return ret
}
