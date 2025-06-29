-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

-- Reverter a adição da coluna DiarioCultivoID na tabela tarefas
ALTER TABLE tarefas DROP COLUMN IF EXISTS diario_cultivo_id;

-- Reverter a adição da coluna AmbienteID na tabela fotos
ALTER TABLE fotos DROP COLUMN IF EXISTS ambiente_id;

-- Reverter a remoção da coluna Estagios da tabela plantas (se necessário, adicione o tipo original)
-- ALTER TABLE plantas ADD COLUMN IF NOT EXISTS estagios <TIPO_ORIGINAL>;

-- Reverter a remoção da coluna Tarefas da tabela plantas (se necessário, adicione o tipo original)
-- ALTER TABLE plantas ADD COLUMN IF NOT EXISTS tarefas <TIPO_ORIGINAL>;

-- ATENÇÃO: Recriar tabelas removidas é complexo e requer o esquema original.
-- As instruções abaixo são apenas placeholders e podem não funcionar sem o esquema completo.

-- CREATE TABLE IF NOT EXISTS registro_crescimento (
--    id SERIAL PRIMARY KEY,
--    planta_id BIGINT,
--    altura NUMERIC,
--    ...
-- );

-- CREATE TABLE IF NOT EXISTS registro_planta (
--    id SERIAL PRIMARY KEY,
--    planta_id BIGINT,
--    data_registro TIMESTAMP,
--    ...
-- );

-- CREATE TABLE IF NOT EXISTS planta_foto_ids (
--    planta_id BIGINT,
--    foto_id BIGINT,
--    ...
-- );
