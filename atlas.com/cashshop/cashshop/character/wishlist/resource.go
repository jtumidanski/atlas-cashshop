package wishlist

import (
	"atlas-cashshop/json"
	"atlas-cashshop/rest"
	"atlas-cashshop/rest/response"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

const (
	GetWishlist    = "get_wishlist"
	DeleteWishlist = "delete_wishlist"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
	r := router.PathPrefix("/characters").Subrouter()
	r.HandleFunc("/{id}/wishlist", registerGetWishlist(l, db)).Methods(http.MethodGet)
	r.HandleFunc("/{id}/wishlist", registerDeleteWishlist(l, db)).Methods(http.MethodDelete)
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

func registerDeleteWishlist(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(DeleteWishlist, func(span opentracing.Span) http.HandlerFunc {
		return ParseId(l, func(characterId uint32) http.HandlerFunc {
			return handleDeleteWishlist(l, db)(span)(characterId)
		})
	})
}

func handleDeleteWishlist(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, _ *http.Request) {
				err := DeleteForCharacter(l, db, span)(characterId)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			}
		}
	}
}

func registerGetWishlist(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(GetWishlist, func(span opentracing.Span) http.HandlerFunc {
		return ParseId(l, func(characterId uint32) http.HandlerFunc {
			return handleGetWishlist(l, db)(span)(characterId)
		})
	})
}

func handleGetWishlist(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, _ *http.Request) {
				wl, err := GetById(l, db)(characterId)
				if err != nil {
					l.WithError(err).Errorf("Unable to retrieve wishlist for character %d.", characterId)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				result := response.NewDataContainer(false)
				for _, wli := range wl {
					result.AddData(wli.Id(), "wishlist", MakeAttribute(wli), nil)
				}

				w.WriteHeader(http.StatusOK)
				err = json.ToJSON(result, w)
				if err != nil {
					l.WithError(err).Errorf("Writing response.")
				}
			}
		}
	}
}
