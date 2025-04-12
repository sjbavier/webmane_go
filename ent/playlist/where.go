// Code generated by ent, DO NOT EDIT.

package playlist

import (
	"time"
	"webmane_go/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Playlist {
	return predicate.Playlist(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Playlist {
	return predicate.Playlist(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Playlist {
	return predicate.Playlist(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Playlist {
	return predicate.Playlist(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Playlist {
	return predicate.Playlist(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Playlist {
	return predicate.Playlist(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Playlist {
	return predicate.Playlist(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Playlist {
	return predicate.Playlist(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Playlist {
	return predicate.Playlist(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldEQ(FieldName, v))
}

// LastUpdate applies equality check predicate on the "last_update" field. It's identical to LastUpdateEQ.
func LastUpdate(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldEQ(FieldLastUpdate, v))
}

// LastAccessed applies equality check predicate on the "last_accessed" field. It's identical to LastAccessedEQ.
func LastAccessed(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldEQ(FieldLastAccessed, v))
}

// CoverArt applies equality check predicate on the "cover_art" field. It's identical to CoverArtEQ.
func CoverArt(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldEQ(FieldCoverArt, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Playlist {
	return predicate.Playlist(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Playlist {
	return predicate.Playlist(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldContainsFold(FieldName, v))
}

// LastUpdateEQ applies the EQ predicate on the "last_update" field.
func LastUpdateEQ(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldEQ(FieldLastUpdate, v))
}

// LastUpdateNEQ applies the NEQ predicate on the "last_update" field.
func LastUpdateNEQ(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldNEQ(FieldLastUpdate, v))
}

// LastUpdateIn applies the In predicate on the "last_update" field.
func LastUpdateIn(vs ...time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldIn(FieldLastUpdate, vs...))
}

// LastUpdateNotIn applies the NotIn predicate on the "last_update" field.
func LastUpdateNotIn(vs ...time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldNotIn(FieldLastUpdate, vs...))
}

// LastUpdateGT applies the GT predicate on the "last_update" field.
func LastUpdateGT(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldGT(FieldLastUpdate, v))
}

// LastUpdateGTE applies the GTE predicate on the "last_update" field.
func LastUpdateGTE(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldGTE(FieldLastUpdate, v))
}

// LastUpdateLT applies the LT predicate on the "last_update" field.
func LastUpdateLT(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldLT(FieldLastUpdate, v))
}

// LastUpdateLTE applies the LTE predicate on the "last_update" field.
func LastUpdateLTE(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldLTE(FieldLastUpdate, v))
}

// LastAccessedEQ applies the EQ predicate on the "last_accessed" field.
func LastAccessedEQ(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldEQ(FieldLastAccessed, v))
}

// LastAccessedNEQ applies the NEQ predicate on the "last_accessed" field.
func LastAccessedNEQ(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldNEQ(FieldLastAccessed, v))
}

// LastAccessedIn applies the In predicate on the "last_accessed" field.
func LastAccessedIn(vs ...time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldIn(FieldLastAccessed, vs...))
}

// LastAccessedNotIn applies the NotIn predicate on the "last_accessed" field.
func LastAccessedNotIn(vs ...time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldNotIn(FieldLastAccessed, vs...))
}

// LastAccessedGT applies the GT predicate on the "last_accessed" field.
func LastAccessedGT(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldGT(FieldLastAccessed, v))
}

// LastAccessedGTE applies the GTE predicate on the "last_accessed" field.
func LastAccessedGTE(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldGTE(FieldLastAccessed, v))
}

// LastAccessedLT applies the LT predicate on the "last_accessed" field.
func LastAccessedLT(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldLT(FieldLastAccessed, v))
}

// LastAccessedLTE applies the LTE predicate on the "last_accessed" field.
func LastAccessedLTE(v time.Time) predicate.Playlist {
	return predicate.Playlist(sql.FieldLTE(FieldLastAccessed, v))
}

// LastAccessedIsNil applies the IsNil predicate on the "last_accessed" field.
func LastAccessedIsNil() predicate.Playlist {
	return predicate.Playlist(sql.FieldIsNull(FieldLastAccessed))
}

// LastAccessedNotNil applies the NotNil predicate on the "last_accessed" field.
func LastAccessedNotNil() predicate.Playlist {
	return predicate.Playlist(sql.FieldNotNull(FieldLastAccessed))
}

// CoverArtEQ applies the EQ predicate on the "cover_art" field.
func CoverArtEQ(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldEQ(FieldCoverArt, v))
}

// CoverArtNEQ applies the NEQ predicate on the "cover_art" field.
func CoverArtNEQ(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldNEQ(FieldCoverArt, v))
}

// CoverArtIn applies the In predicate on the "cover_art" field.
func CoverArtIn(vs ...string) predicate.Playlist {
	return predicate.Playlist(sql.FieldIn(FieldCoverArt, vs...))
}

// CoverArtNotIn applies the NotIn predicate on the "cover_art" field.
func CoverArtNotIn(vs ...string) predicate.Playlist {
	return predicate.Playlist(sql.FieldNotIn(FieldCoverArt, vs...))
}

// CoverArtGT applies the GT predicate on the "cover_art" field.
func CoverArtGT(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldGT(FieldCoverArt, v))
}

// CoverArtGTE applies the GTE predicate on the "cover_art" field.
func CoverArtGTE(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldGTE(FieldCoverArt, v))
}

// CoverArtLT applies the LT predicate on the "cover_art" field.
func CoverArtLT(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldLT(FieldCoverArt, v))
}

// CoverArtLTE applies the LTE predicate on the "cover_art" field.
func CoverArtLTE(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldLTE(FieldCoverArt, v))
}

// CoverArtContains applies the Contains predicate on the "cover_art" field.
func CoverArtContains(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldContains(FieldCoverArt, v))
}

// CoverArtHasPrefix applies the HasPrefix predicate on the "cover_art" field.
func CoverArtHasPrefix(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldHasPrefix(FieldCoverArt, v))
}

// CoverArtHasSuffix applies the HasSuffix predicate on the "cover_art" field.
func CoverArtHasSuffix(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldHasSuffix(FieldCoverArt, v))
}

// CoverArtIsNil applies the IsNil predicate on the "cover_art" field.
func CoverArtIsNil() predicate.Playlist {
	return predicate.Playlist(sql.FieldIsNull(FieldCoverArt))
}

// CoverArtNotNil applies the NotNil predicate on the "cover_art" field.
func CoverArtNotNil() predicate.Playlist {
	return predicate.Playlist(sql.FieldNotNull(FieldCoverArt))
}

// CoverArtEqualFold applies the EqualFold predicate on the "cover_art" field.
func CoverArtEqualFold(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldEqualFold(FieldCoverArt, v))
}

// CoverArtContainsFold applies the ContainsFold predicate on the "cover_art" field.
func CoverArtContainsFold(v string) predicate.Playlist {
	return predicate.Playlist(sql.FieldContainsFold(FieldCoverArt, v))
}

// HasSongs applies the HasEdge predicate on the "songs" edge.
func HasSongs() predicate.Playlist {
	return predicate.Playlist(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, SongsTable, SongsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSongsWith applies the HasEdge predicate on the "songs" edge with a given conditions (other predicates).
func HasSongsWith(preds ...predicate.Music) predicate.Playlist {
	return predicate.Playlist(func(s *sql.Selector) {
		step := newSongsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Playlist) predicate.Playlist {
	return predicate.Playlist(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Playlist) predicate.Playlist {
	return predicate.Playlist(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Playlist) predicate.Playlist {
	return predicate.Playlist(sql.NotPredicates(p))
}
