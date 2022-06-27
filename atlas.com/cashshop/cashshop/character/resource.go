package character

import (
	"atlas-cashshop/cashshop/character/wishlist"
	"atlas-cashshop/json"
	"atlas-cashshop/rest"
	"atlas-cashshop/rest/response"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

const (
	GetCharacter = "get_character"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
	r := router.PathPrefix("/characters").Subrouter()
	r.HandleFunc("/{id}", registerGetCharacter(l, db)).Methods(http.MethodGet).Queries("include", "{include}")
	r.HandleFunc("/{id}", registerGetCharacter(l, db)).Methods(http.MethodGet)
}

func registerGetCharacter(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(GetCharacter, func(span opentracing.Span) http.HandlerFunc {
		return ParseId(l, func(characterId uint32) http.HandlerFunc {
			return handleGetCharacter(l, db)(span)(characterId)
		})
	})
}

type IdHandler func(characterId uint32) http.HandlerFunc

func ParseId(l logrus.FieldLogger, next IdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			l.WithError(err).Errorf("Unable to properly parse id from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(uint32(id))(w, r)
	}
}

func handleGetCharacter(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				c, err := GetById(l, db)(characterId)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				result := response.NewDataContainer(true)
				result.AddData(c.CharacterId(), "characters", makeAttribute(c), makeCharacterRelationships(c))
				if strings.Contains(mux.Vars(r)["include"], "wishlist") {
					for _, m := range c.Wishlist() {
						result.AddIncluded(m.Id(), "wishlist", wishlist.MakeAttribute(m))
					}
				}

				err = json.ToJSON(result, w)
				if err != nil {
					l.WithError(err).Errorf("Encoding response.")
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}

func makeCharacterRelationships(c Model) map[string]*response.Relationship {
	result := make(map[string]*response.Relationship, 0)
	result["wishlist"] = &response.Relationship{
		ToOneType: false,
		Links: response.RelationshipLinks{
			Self:    "/ms/cashshop/characters/" + strconv.Itoa(int(c.CharacterId())) + "/relationships/wishlist",
			Related: "/ms/cashshop/characters/" + strconv.Itoa(int(c.CharacterId())) + "/wishlist",
		},
		Data: wishlist.MakeRelationshipData(c.Wishlist()),
	}
	return result
}
