-- 000001_create_initial_schema.down.sql
-- Drop junction tables first
DROP TABLE IF EXISTS diario_plantas;
DROP TABLE IF EXISTS diario_ambientes;
DROP TABLE IF EXISTS tarefa_plantas;

-- Drop tables that depend on others
DROP TABLE IF EXISTS anotacoes;
DROP TABLE IF EXISTS estagio_crescimentos;
DROP TABLE IF EXISTS registro_crescimentos;
DROP TABLE IF EXISTS registro_plantas;
DROP TABLE IF EXISTS registro_diarios;
DROP TABLE IF EXISTS midias; -- Depends on colecao_midias and usuarios

-- Drop main tables
DROP TABLE IF EXISTS tarefas; -- Depends on plantas, ambientes, usuarios, diario_cultivos
DROP TABLE IF EXISTS plantas; -- Depends on fotos, geneticas, meio_cultivos, ambientes, usuarios

-- Drop tables that are referenced by others, or have no FKs
DROP TABLE IF EXISTS colecao_midias; -- Referenced by midias
DROP TABLE IF EXISTS diario_cultivos; -- Referenced by colecao_midias, registro_diarios, tarefas, diario_plantas, diario_ambientes
DROP TABLE IF EXISTS fotos; -- Referenced by plantas
DROP TABLE IF EXISTS micro_climas;
DROP TABLE IF EXISTS tipo_tarefas;
DROP TABLE IF EXISTS substratos;
DROP TABLE IF EXISTS vasos;
DROP TABLE IF EXISTS geneticas;
DROP TABLE IF EXISTS meio_cultivos;
DROP TABLE IF EXISTS usuarios;
DROP TABLE IF EXISTS ambientes;
DROP TABLE IF EXISTS tarefa_templates;