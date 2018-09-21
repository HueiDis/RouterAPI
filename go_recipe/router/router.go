package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/golang/protobuf/proto"

	"testRouterAPI/go_recipe/logging"
	"testRouterAPI/go_recipe/recipe"

	"github.com/gorilla/mux"
	"github.com/rs/xid"
)

func GetRecipe() *recipe.Recipe {
	return &recipe.Recipe{
		Id:      "111",
		Factory: recipe.FactoryID(2),
		//Factory: (recipe.FactoryID_value["KY"]),
		SpecID: "U0000001",
	}
}

// NewRouter returns a new router for the server.
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(logging.Middleware) // Log all requests.

	// Create router for api endpoints.
	handleAPIRoutes("/api/", r)

	// Create router for web endpoints, including web backend calls that require auth.
	handleWebRoutes("/web/", r)

	// Other endpoints without auth.

	// Serve static files (css, js, img, ...)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))).Methods("GET")

	// [Toy example] Handle index with inline html + served javascript file.
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		io.WriteString(w, `
		<!doctype html>
		<head>
			<meta charset="utf-8">
			<title>Index</title>
			<script type="text/javascript" src="/static/js/example.js"></script>
		</head>
		<body>
			<p>Index!</p>
		</body>
		</html>
		`)
	}).Methods("GET")

	// Server status.
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// TODO: add backend health status.
		io.WriteString(w, `{"alive": true}`)
	}).Methods("GET")

	// Catch not-founds.
	r.PathPrefix("/").Handler(http.NotFoundHandler())

	return r
}

func handleAPIRoutes(pathPrefix string, parentRouter *mux.Router) {
	r := mux.NewRouter()
	// TODO: add api auth middleware.

	r.HandleFunc(
		path.Join(pathPrefix, "spec_example"),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if r.URL.Query().Get("format") == "json" {
				bs, err := json.Marshal(ExampleSpec())
				if err != nil {
					logging.WithRequestID("rid", r).Fatalln(err)
				}
				w.Write(bs)
			} else {
				bs, err := proto.Marshal(ExampleSpec())
				if err != nil {
					logging.WithRequestID("rid", r).Fatalln(err)
				}
				w.Write(bs)
			}
		},
	)

	// TODO: add api endpoints.

	parentRouter.PathPrefix(pathPrefix).Handler(r).Methods("GET") // Only GET endpoints for now.
}

func handleWebRoutes(pathPrefix string, parentRouter *mux.Router) {
	r := mux.NewRouter()
	// TODO: add web auth middleware.

	// [Toy example] Template usage.
	r.HandleFunc(
		path.Join(pathPrefix, "simple_example.html"),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)

			t, err := template.ParseFiles("tpl/simple_example.html")
			if err != nil {
				logging.WithRequestID("rid", r).Fatalln(err)
			}
			_, err = t.ParseGlob("tpl/public/*.html")
			if err != nil {
				logging.WithRequestID("rid", r).Fatalln(err)
			}

			data := struct {
				Title string
				Rows  []string
			}{
				Title: "Template example",
				Rows:  []string{"Row A", "Row B", "Row C"},
			}

			err = t.Execute(w, data)
			if err != nil {
				logging.WithRequestID("rid", r).Fatalln(err)
			}
		})

	r.HandleFunc(
		path.Join(pathPrefix, "spec_example.html"),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)

			t, err := template.ParseFiles("tpl/spec_example.html")
			if err != nil {
				logging.WithRequestID("rid", r).Fatalln(err)
			}
			_, err = t.ParseGlob("tpl/public/*.html")
			if err != nil {
				logging.WithRequestID("rid", r).Fatalln(err)
			}
			_, err = t.ParseGlob("tpl/spec_tables/*.html")
			if err != nil {
				logging.WithRequestID("rid", r).Fatalln(err)
			}

			err = t.Execute(w, exampleSpecView())
			if err != nil {
				logging.WithRequestID("rid", r).Fatalln(err)
			}
		})

	// TODO: add web endpoints.

	parentRouter.PathPrefix(pathPrefix).Handler(r)
}

// Example spec object.

func ExampleSpec() *recipe.Spec {
	return &recipe.Spec{
		Id:        xid.New().String(),
		ProductID: "U0123456",
		VersionHistory: []*recipe.Version{
			&recipe.Version{
				Status:    recipe.Version_TEST,
				Major:     "1",
				Minor:     "7",
				UpdatedAt: time.Now().AddDate(0, -2, 0).UnixNano(), // 2 months ago
				UpdatedBy: "testbot",
				Note:      "let's test",
			},
			&recipe.Version{
				Status:    recipe.Version_PRODUCTION,
				Major:     "1",
				Minor:     "3",
				UpdatedAt: time.Now().UnixNano(),
				UpdatedBy: "prodbot",
				Note:      "go to prduction",
			},
		},
		ParentSpecID: "", // No parent.
		Materials: []*recipe.Material{
			&recipe.Material{
				Type: "Polymer",
				Name: "Material One",
				Param: &recipe.Param{
					Unit:  "phr",
					Value: 90.0,
				},
			},
			&recipe.Material{
				Type: "Filler",
				Name: "Material Two",
				Param: &recipe.Param{
					Unit:  "phr",
					Value: 10.0,
				},
			},
			&recipe.Material{
				Type: "Activator",
				Name: "Material Three",
				Param: &recipe.Param{
					Unit:  "phr",
					Value: 2.0,
					Aux:   1.0,
				},
			},
		},
		Tools: []*recipe.Tool{}, // No tools.
	}
}

type materialView struct {
	Type, Name, Value string
}

func toMaterialViews(ms []*recipe.Material) (r []*materialView) {
	for _, m := range ms {
		value := ""
		p := m.GetParam()
		if p == nil {
			value = "N/A"
		} else {
			aux := p.GetAux()
			if aux == 0.0 {
				value += fmt.Sprintf("%f", p.GetValue())
			} else {
				min := p.GetValue()
				max := min + aux
				value += fmt.Sprintf("%f ~ %f", min, max)
			}

			errorMargin := p.GetErrorMargin()
			if errorMargin == 0.0 {
				value += fmt.Sprintf(" (%s)", p.GetUnit())
			} else {
				value += fmt.Sprintf(" +- %f (%s)", errorMargin, p.GetUnit())
			}
		}

		r = append(r, &materialView{
			Type:  m.GetType(),
			Name:  m.GetName(),
			Value: value,
		})
	}
	return
}

type specView struct {
	Title          string
	VersionHistory []*recipe.Version
	Materials      []*materialView
}

func exampleSpecView() *specView {
	s := ExampleSpec()
	verHist := s.GetVersionHistory()
	title := s.ProductID
	if verHist != nil && len(verHist) > 0 { // This should always be true in our case.
		ver := verHist[len(verHist)-1]
		title = title + fmt.Sprintf("-%s-%s-%s",
			ver.Status.String(),
			ver.Major,
			ver.Minor)
	}
	return &specView{
		Title:          title,
		VersionHistory: verHist,
		Materials:      toMaterialViews(s.GetMaterials()),
	}
}
