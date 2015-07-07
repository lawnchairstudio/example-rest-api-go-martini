package main

import (
  "flag"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/sessions"
  "database/sql"
  "github.com/auth0/go-jwt-middleware"
  "encoding/base64"
  "github.com/dgrijalva/jwt-go"
  "net/http"
  _ "github.com/lib/pq"
)

/**
 * Access database and return a connection.
 */
func SetupDatabase() *sql.DB {
  db, err := sql.Open("postgres", "dbname=segment-go-test sslmode=disable")
  handleError(err)
  return db
}

func handleError(err error) {
  if err != nil {
    panic(err)
  }
}

func main() {

  jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      decoded, err := base64.URLEncoding.DecodeString("fo5kctDR_R7ys1wwl6WpSs1YmWNgGoG7VD1spH5pgdXzZT8dVHcX0W1FLKb_DaQj")
      if err != nil {
        return nil, err
      }
      return decoded, nil
    },
  })

  m := martini.Classic()

  port := flag.String("port", "8000", "HTTP Port")

  flag.Parse()

  m.Map(SetupDatabase())

  store := sessions.NewCookieStore([]byte("secret123"))

  m.Use(sessions.Sessions("segment_session", store))

  m.Use(render.Renderer(render.Options{
    IndentJSON: true, // Output human readable JSON
  }))

  // for auth, just pass jwtMiddleware.CheckJWT first
  // eg: m.Get("/", jwtMiddleware.CheckJWT, Index)

  m.Get("/", Index)
  m.Get("/session/new", Session)
  m.Post("/user/new", NewUser)
  m.Get("/users/:id", Users)
  m.Get("/segments/:id", Segments)
  m.Get("/authtest", jwtMiddleware.CheckJWT, func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
    res.WriteHeader(200) // HTTP 200
  })

  // Get the PORT from the environment. Necessary for Heroku.
  m.RunOnAddr(":" + *port)

}

/**
 * @func Index
 * @description Index route controller.
 */
func Index(r render.Render) {
  r.JSON(200, map[string]interface{}{
    "create_session": "https://api.segment.social/session/new",
    "segments": "https://api.segment.social/segments/{segment_id}",
    "users": "https://api.segment.social/users/{user_id}",
    "create_user": "https://api.segment.social/user/new",
  })
}

func Session(session sessions.Session) string {
  session.Set("hello", "world")
  return "OK"
}

func NewUser(r render.Render, params martini.Params, db *sql.DB) {
  rows, err := db.Query("INSERT INTO users (name, password, email) VALUES ($1, $2, $3)", "a", "b", "c");
  handleError(err)
  defer rows.Close()
  r.JSON(200, map[string]interface{}{
    "name": "a",
    "password": "b",
    "email": "c",
  })
}

func Users(r render.Render, params martini.Params, s sessions.Session, db *sql.DB) {
  key := s.Get("hello")
  if key == nil {
    r.JSON(400, map[string]interface{}{
      "session": "false",
    })
  } else {
    r.JSON(200, map[string]interface{}{
      "id": params["id"],
    })
  }
}

func Segments(r render.Render, params martini.Params, s sessions.Session, db *sql.DB) {
  r.JSON(200, map[string]interface{}{
    "id": params["id"],
  })
}
