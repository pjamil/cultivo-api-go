-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

DROP TABLE IF EXISTS registro_crescimento;
DROP TABLE IF EXISTS registro_planta;
DROP TABLE IF EXISTS planta_foto_ids;

-- Adicionar coluna DiarioCultivoID na tabela tarefas, se ainda não existir
ALTER TABLE tarefas ADD COLUMN IF NOT EXISTS diario_cultivo_id BIGINT;

-- Adicionar coluna AmbienteID na tabela fotos, se ainda não existir
ALTER TABLE fotos ADD COLUMN IF NOT EXISTS ambiente_id BIGINT;

-- Remover a coluna Estagios da tabela plantas, se existir
ALTER TABLE plantas DROP COLUMN IF EXISTS estagios;

-- Remover a coluna Tarefas da tabela plantas, se existir
ALTER TABLE plantas DROP COLUMN IF EXISTS tarefas;

-- Alterar o tipo da coluna tipo na tabela registro_diario para o novo ENUM, se necessário
-- ATENÇÃO: Se você estiver usando um ENUM real no PostgreSQL, precisará de um ALTER TYPE.
-- Se for apenas uma string com validação na aplicação, esta parte não é necessária.
-- Exemplo para ENUM:
-- ALTER TYPE registro_tipo ADD VALUE 'crescimento';

-- Atualizar a coluna ReferenciaTipo na tabela registro_diario para o novo ENUM, se necessário
-- ATENÇÃO: Se você estiver usando um ENUM real no PostgreSQL, precisará de um ALTER TYPE.
-- Se for apenas uma string com validação na aplicação, esta parte não é necessária.
-- Exemplo para ENUM:
-- ALTER TYPE referencia_tipo ADD VALUE 'planta';
