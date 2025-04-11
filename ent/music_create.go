// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"
	"webmane_go/ent/music"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MusicCreate is the builder for creating a Music entity.
type MusicCreate struct {
	config
	mutation *MusicMutation
	hooks    []Hook
}

// SetPath sets the "path" field.
func (mc *MusicCreate) SetPath(s string) *MusicCreate {
	mc.mutation.SetPath(s)
	return mc
}

// SetLastUpdate sets the "last_update" field.
func (mc *MusicCreate) SetLastUpdate(t time.Time) *MusicCreate {
	mc.mutation.SetLastUpdate(t)
	return mc
}

// SetNillableLastUpdate sets the "last_update" field if the given value is not nil.
func (mc *MusicCreate) SetNillableLastUpdate(t *time.Time) *MusicCreate {
	if t != nil {
		mc.SetLastUpdate(*t)
	}
	return mc
}

// SetTitle sets the "title" field.
func (mc *MusicCreate) SetTitle(s string) *MusicCreate {
	mc.mutation.SetTitle(s)
	return mc
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (mc *MusicCreate) SetNillableTitle(s *string) *MusicCreate {
	if s != nil {
		mc.SetTitle(*s)
	}
	return mc
}

// SetArtist sets the "artist" field.
func (mc *MusicCreate) SetArtist(s string) *MusicCreate {
	mc.mutation.SetArtist(s)
	return mc
}

// SetNillableArtist sets the "artist" field if the given value is not nil.
func (mc *MusicCreate) SetNillableArtist(s *string) *MusicCreate {
	if s != nil {
		mc.SetArtist(*s)
	}
	return mc
}

// SetAlbum sets the "album" field.
func (mc *MusicCreate) SetAlbum(s string) *MusicCreate {
	mc.mutation.SetAlbum(s)
	return mc
}

// SetNillableAlbum sets the "album" field if the given value is not nil.
func (mc *MusicCreate) SetNillableAlbum(s *string) *MusicCreate {
	if s != nil {
		mc.SetAlbum(*s)
	}
	return mc
}

// SetGenre sets the "genre" field.
func (mc *MusicCreate) SetGenre(s string) *MusicCreate {
	mc.mutation.SetGenre(s)
	return mc
}

// SetNillableGenre sets the "genre" field if the given value is not nil.
func (mc *MusicCreate) SetNillableGenre(s *string) *MusicCreate {
	if s != nil {
		mc.SetGenre(*s)
	}
	return mc
}

// SetReleaseYear sets the "release_year" field.
func (mc *MusicCreate) SetReleaseYear(s string) *MusicCreate {
	mc.mutation.SetReleaseYear(s)
	return mc
}

// SetNillableReleaseYear sets the "release_year" field if the given value is not nil.
func (mc *MusicCreate) SetNillableReleaseYear(s *string) *MusicCreate {
	if s != nil {
		mc.SetReleaseYear(*s)
	}
	return mc
}

// SetCoverArt sets the "cover_art" field.
func (mc *MusicCreate) SetCoverArt(s string) *MusicCreate {
	mc.mutation.SetCoverArt(s)
	return mc
}

// SetNillableCoverArt sets the "cover_art" field if the given value is not nil.
func (mc *MusicCreate) SetNillableCoverArt(s *string) *MusicCreate {
	if s != nil {
		mc.SetCoverArt(*s)
	}
	return mc
}

// Mutation returns the MusicMutation object of the builder.
func (mc *MusicCreate) Mutation() *MusicMutation {
	return mc.mutation
}

// Save creates the Music in the database.
func (mc *MusicCreate) Save(ctx context.Context) (*Music, error) {
	mc.defaults()
	return withHooks(ctx, mc.sqlSave, mc.mutation, mc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mc *MusicCreate) SaveX(ctx context.Context) *Music {
	v, err := mc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mc *MusicCreate) Exec(ctx context.Context) error {
	_, err := mc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mc *MusicCreate) ExecX(ctx context.Context) {
	if err := mc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mc *MusicCreate) defaults() {
	if _, ok := mc.mutation.LastUpdate(); !ok {
		v := music.DefaultLastUpdate()
		mc.mutation.SetLastUpdate(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mc *MusicCreate) check() error {
	if _, ok := mc.mutation.Path(); !ok {
		return &ValidationError{Name: "path", err: errors.New(`ent: missing required field "Music.path"`)}
	}
	if v, ok := mc.mutation.Path(); ok {
		if err := music.PathValidator(v); err != nil {
			return &ValidationError{Name: "path", err: fmt.Errorf(`ent: validator failed for field "Music.path": %w`, err)}
		}
	}
	if _, ok := mc.mutation.LastUpdate(); !ok {
		return &ValidationError{Name: "last_update", err: errors.New(`ent: missing required field "Music.last_update"`)}
	}
	if v, ok := mc.mutation.Title(); ok {
		if err := music.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Music.title": %w`, err)}
		}
	}
	if v, ok := mc.mutation.Artist(); ok {
		if err := music.ArtistValidator(v); err != nil {
			return &ValidationError{Name: "artist", err: fmt.Errorf(`ent: validator failed for field "Music.artist": %w`, err)}
		}
	}
	if v, ok := mc.mutation.Album(); ok {
		if err := music.AlbumValidator(v); err != nil {
			return &ValidationError{Name: "album", err: fmt.Errorf(`ent: validator failed for field "Music.album": %w`, err)}
		}
	}
	if v, ok := mc.mutation.Genre(); ok {
		if err := music.GenreValidator(v); err != nil {
			return &ValidationError{Name: "genre", err: fmt.Errorf(`ent: validator failed for field "Music.genre": %w`, err)}
		}
	}
	if v, ok := mc.mutation.ReleaseYear(); ok {
		if err := music.ReleaseYearValidator(v); err != nil {
			return &ValidationError{Name: "release_year", err: fmt.Errorf(`ent: validator failed for field "Music.release_year": %w`, err)}
		}
	}
	return nil
}

func (mc *MusicCreate) sqlSave(ctx context.Context) (*Music, error) {
	if err := mc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	mc.mutation.id = &_node.ID
	mc.mutation.done = true
	return _node, nil
}

func (mc *MusicCreate) createSpec() (*Music, *sqlgraph.CreateSpec) {
	var (
		_node = &Music{config: mc.config}
		_spec = sqlgraph.NewCreateSpec(music.Table, sqlgraph.NewFieldSpec(music.FieldID, field.TypeInt))
	)
	if value, ok := mc.mutation.Path(); ok {
		_spec.SetField(music.FieldPath, field.TypeString, value)
		_node.Path = value
	}
	if value, ok := mc.mutation.LastUpdate(); ok {
		_spec.SetField(music.FieldLastUpdate, field.TypeTime, value)
		_node.LastUpdate = value
	}
	if value, ok := mc.mutation.Title(); ok {
		_spec.SetField(music.FieldTitle, field.TypeString, value)
		_node.Title = value
	}
	if value, ok := mc.mutation.Artist(); ok {
		_spec.SetField(music.FieldArtist, field.TypeString, value)
		_node.Artist = value
	}
	if value, ok := mc.mutation.Album(); ok {
		_spec.SetField(music.FieldAlbum, field.TypeString, value)
		_node.Album = value
	}
	if value, ok := mc.mutation.Genre(); ok {
		_spec.SetField(music.FieldGenre, field.TypeString, value)
		_node.Genre = value
	}
	if value, ok := mc.mutation.ReleaseYear(); ok {
		_spec.SetField(music.FieldReleaseYear, field.TypeString, value)
		_node.ReleaseYear = value
	}
	if value, ok := mc.mutation.CoverArt(); ok {
		_spec.SetField(music.FieldCoverArt, field.TypeString, value)
		_node.CoverArt = value
	}
	return _node, _spec
}

// MusicCreateBulk is the builder for creating many Music entities in bulk.
type MusicCreateBulk struct {
	config
	err      error
	builders []*MusicCreate
}

// Save creates the Music entities in the database.
func (mcb *MusicCreateBulk) Save(ctx context.Context) ([]*Music, error) {
	if mcb.err != nil {
		return nil, mcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(mcb.builders))
	nodes := make([]*Music, len(mcb.builders))
	mutators := make([]Mutator, len(mcb.builders))
	for i := range mcb.builders {
		func(i int, root context.Context) {
			builder := mcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MusicMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, mcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, mcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mcb *MusicCreateBulk) SaveX(ctx context.Context) []*Music {
	v, err := mcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mcb *MusicCreateBulk) Exec(ctx context.Context) error {
	_, err := mcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mcb *MusicCreateBulk) ExecX(ctx context.Context) {
	if err := mcb.Exec(ctx); err != nil {
		panic(err)
	}
}
