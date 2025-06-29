package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// EspacoCultivo represents a row from 'public.espaco_cultivo'.
type EspacoCultivo struct {
	ID             int64          `json:"id"`              // id
	Nome           string         `json:"nome"`            // nome
	AreaCultivo    string         `json:"area_cultivo"`    // area_cultivo
	TipoIluminacao string         `json:"tipo_iluminacao"` // tipo_iluminacao
	QuantidadeLuz  sql.NullString `json:"quantidade_luz"`  // quantidade_luz
	Ativo          bool           `json:"ativo"`           // ativo
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [EspacoCultivo] exists in the database.
func (ec *EspacoCultivo) Exists() bool {
	return ec._exists
}

// Deleted returns true when the [EspacoCultivo] has been marked for deletion
// from the database.
func (ec *EspacoCultivo) Deleted() bool {
	return ec._deleted
}

// Insert inserts the [EspacoCultivo] to the database.
func (ec *EspacoCultivo) Insert(ctx context.Context, db DB) error {
	switch {
	case ec._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case ec._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO public.espaco_cultivo (` +
		`nome, area_cultivo, tipo_iluminacao, quantidade_luz, ativo` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) RETURNING id`
	// run
	logf(sqlstr, ec.Nome, ec.AreaCultivo, ec.TipoIluminacao, ec.QuantidadeLuz, ec.Ativo)
	if err := db.QueryRowContext(ctx, sqlstr, ec.Nome, ec.AreaCultivo, ec.TipoIluminacao, ec.QuantidadeLuz, ec.Ativo).Scan(&ec.ID); err != nil {
		return logerror(err)
	}
	// set exists
	ec._exists = true
	return nil
}

// Update updates a [EspacoCultivo] in the database.
func (ec *EspacoCultivo) Update(ctx context.Context, db DB) error {
	switch {
	case !ec._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case ec._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with composite primary key
	const sqlstr = `UPDATE public.espaco_cultivo SET ` +
		`nome = $1, area_cultivo = $2, tipo_iluminacao = $3, quantidade_luz = $4, ativo = $5 ` +
		`WHERE id = $6`
	// run
	logf(sqlstr, ec.Nome, ec.AreaCultivo, ec.TipoIluminacao, ec.QuantidadeLuz, ec.Ativo, ec.ID)
	if _, err := db.ExecContext(ctx, sqlstr, ec.Nome, ec.AreaCultivo, ec.TipoIluminacao, ec.QuantidadeLuz, ec.Ativo, ec.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [EspacoCultivo] to the database.
func (ec *EspacoCultivo) Save(ctx context.Context, db DB) error {
	if ec.Exists() {
		return ec.Update(ctx, db)
	}
	return ec.Insert(ctx, db)
}

// Upsert performs an upsert for [EspacoCultivo].
func (ec *EspacoCultivo) Upsert(ctx context.Context, db DB) error {
	switch {
	case ec._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO public.espaco_cultivo (` +
		`id, nome, area_cultivo, tipo_iluminacao, quantidade_luz, ativo` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6` +
		`)` +
		` ON CONFLICT (id) DO ` +
		`UPDATE SET ` +
		`nome = EXCLUDED.nome, area_cultivo = EXCLUDED.area_cultivo, tipo_iluminacao = EXCLUDED.tipo_iluminacao, quantidade_luz = EXCLUDED.quantidade_luz, ativo = EXCLUDED.ativo `
	// run
	logf(sqlstr, ec.ID, ec.Nome, ec.AreaCultivo, ec.TipoIluminacao, ec.QuantidadeLuz, ec.Ativo)
	if _, err := db.ExecContext(ctx, sqlstr, ec.ID, ec.Nome, ec.AreaCultivo, ec.TipoIluminacao, ec.QuantidadeLuz, ec.Ativo); err != nil {
		return logerror(err)
	}
	// set exists
	ec._exists = true
	return nil
}

// Delete deletes the [EspacoCultivo] from the database.
func (ec *EspacoCultivo) Delete(ctx context.Context, db DB) error {
	switch {
	case !ec._exists: // doesn't exist
		return nil
	case ec._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM public.espaco_cultivo ` +
		`WHERE id = $1`
	// run
	logf(sqlstr, ec.ID)
	if _, err := db.ExecContext(ctx, sqlstr, ec.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	ec._deleted = true
	return nil
}

// EspacoCultivoByID retrieves a row from 'public.espaco_cultivo' as a [EspacoCultivo].
//
// Generated from index 'espaco_cultivo_pk'.
func EspacoCultivoByID(ctx context.Context, db DB, id int64) (*EspacoCultivo, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, nome, area_cultivo, tipo_iluminacao, quantidade_luz, ativo ` +
		`FROM public.espaco_cultivo ` +
		`WHERE id = $1`
	// run
	logf(sqlstr, id)
	ec := EspacoCultivo{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, id).Scan(&ec.ID, &ec.Nome, &ec.AreaCultivo, &ec.TipoIluminacao, &ec.QuantidadeLuz, &ec.Ativo); err != nil {
		return nil, logerror(err)
	}
	return &ec, nil
}

// EspacoCultivoByNome retrieves a row from 'public.espaco_cultivo' as a [EspacoCultivo].
//
// Generated from index 'espaco_cultivo_unique'.
func EspacoCultivoByNome(ctx context.Context, db DB, nome string) (*EspacoCultivo, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, nome, area_cultivo, tipo_iluminacao, quantidade_luz, ativo ` +
		`FROM public.espaco_cultivo ` +
		`WHERE nome = $1`
	// run
	logf(sqlstr, nome)
	ec := EspacoCultivo{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, nome).Scan(&ec.ID, &ec.Nome, &ec.AreaCultivo, &ec.TipoIluminacao, &ec.QuantidadeLuz, &ec.Ativo); err != nil {
		return nil, logerror(err)
	}
	return &ec, nil
}
