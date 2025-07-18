-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS registro_diarios (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    titulo VARCHAR(255) NOT NULL,
    conteudo TEXT NOT NULL,
    data TIMESTAMPTZ NOT NULL,
    tipo VARCHAR(50) NOT NULL,
    diario_cultivo_id BIGINT NOT NULL,
    planta_id BIGINT,

    referencia_id BIGINT,
    referencia_tipo VARCHAR(50),
    CONSTRAINT fk_diario_cultivo FOREIGN KEY(diario_cultivo_id) REFERENCES diario_cultivos(id) ON DELETE CASCADE,
    CONSTRAINT fk_planta FOREIGN KEY(planta_id) REFERENCES plantas(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_registro_diarios_diario_cultivo_id ON registro_diarios(diario_cultivo_id);
CREATE INDEX IF NOT EXISTS idx_registro_diarios_tipo ON registro_diarios(tipo);