-- 000001_create_initial_schema.up.sql

-- Cria a tabela ambientes
CREATE TABLE IF NOT EXISTS ambientes (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    tipo VARCHAR(50) NOT NULL,
    comprimento NUMERIC NOT NULL,
    altura NUMERIC NOT NULL,
    largura NUMERIC NOT NULL,
    tempo_exposicao INTEGER NOT NULL,
    orientacao VARCHAR(20)
);

-- Cria a tabela fotos
CREATE TABLE IF NOT EXISTS fotos (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    owner_id INTEGER,
    owner_type VARCHAR(50),
    url VARCHAR(255) NOT NULL,
    descricao TEXT,
    usuario_id INTEGER,
    ambiente_id INTEGER REFERENCES ambientes(id) ON DELETE SET NULL
);

-- Cria a tabela micro_climas
CREATE TABLE IF NOT EXISTS micro_climas (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    ambiente_id INTEGER REFERENCES ambientes(id) ON DELETE CASCADE,
    data_medicao TIMESTAMP WITH TIME ZONE,
    temperatura NUMERIC,
    umidade NUMERIC,
    luminosidade NUMERIC
);

-- Cria a tabela usuarios
CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    senha_hash VARCHAR(255) NOT NULL,
    preferencias JSON
);

-- Cria a tabela geneticas
CREATE TABLE IF NOT EXISTS geneticas (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    tipo_genetica VARCHAR(50) NOT NULL,
    tipo_especie VARCHAR(50) NOT NULL,
    tempo_floracao INTEGER NOT NULL,
    origem VARCHAR(100) NOT NULL,
    caracteristicas TEXT
);

-- Cria a tabela meio_cultivos
CREATE TABLE IF NOT EXISTS meio_cultivos (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    tipo VARCHAR(100) NOT NULL,
    descricao TEXT
);

-- Cria a tabela plantas
CREATE TABLE IF NOT EXISTS plantas (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    nome VARCHAR(255) NOT NULL,
    comecando_de VARCHAR(100) NOT NULL,
    especie VARCHAR(100) NOT NULL,
    data_plantio TIMESTAMP WITH TIME ZONE NOT NULL,
    data_colheita TIMESTAMP WITH TIME ZONE,
    status VARCHAR(100) NOT NULL,
    notas TEXT,
    foto_capa_id INTEGER REFERENCES fotos(id) ON DELETE SET NULL,
    genetica_id INTEGER REFERENCES geneticas(id) ON DELETE CASCADE,
    meio_cultivo_id INTEGER REFERENCES meio_cultivos(id) ON DELETE CASCADE,
    ambiente_id INTEGER REFERENCES ambientes(id) ON DELETE CASCADE,
    planta_mae_id INTEGER REFERENCES plantas(id) ON DELETE SET NULL,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE
);