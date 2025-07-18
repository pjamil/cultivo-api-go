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

-- Cria a tabela vasos
CREATE TABLE IF NOT EXISTS vasos (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    nome VARCHAR(100) NOT NULL,
    material VARCHAR(50),
    volume NUMERIC,
    diametro NUMERIC,
    cor VARCHAR(30),
    furos_drenagem INTEGER
);

-- Cria a tabela substratos
CREATE TABLE IF NOT EXISTS substratos (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    nome VARCHAR(100) NOT NULL,
    composicao TEXT,
    ph NUMERIC,
    retencao_agua NUMERIC
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
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE SET NULL,
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

-- Cria a tabela tipo_tarefas
CREATE TABLE IF NOT EXISTS tipo_tarefas (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    nome VARCHAR(100) NOT NULL UNIQUE,
    descricao TEXT
);

-- Cria a tabela diario_cultivos
CREATE TABLE IF NOT EXISTS diario_cultivos (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    nome VARCHAR(100) NOT NULL,
    data_inicio TIMESTAMP WITH TIME ZONE NOT NULL,
    data_fim TIMESTAMP WITH TIME ZONE,
    ativo BOOLEAN DEFAULT TRUE,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    privacidade VARCHAR(20) DEFAULT 'privado',
    tags VARCHAR(200)
);

-- Cria a tabela colecao_midias
CREATE TABLE IF NOT EXISTS colecao_midias (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    diario_cultivo_id INTEGER REFERENCES diario_cultivos(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    tipo VARCHAR(50) NOT NULL,
    descricao TEXT,
    capa_url VARCHAR(255)
);

-- Cria a tabela midias
CREATE TABLE IF NOT EXISTS midias (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    tipo VARCHAR(20) NOT NULL,
    url VARCHAR(255) NOT NULL,
    thumbnail_url VARCHAR(255),
    data_captura TIMESTAMP WITH TIME ZONE,
    autor_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    descricao TEXT,
    coordenadas VARCHAR(50),
    colecao_midia_id INTEGER REFERENCES colecao_midias(id) ON DELETE CASCADE
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

-- Cria a tabela tarefas
CREATE TABLE IF NOT EXISTS tarefas (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    tipo VARCHAR(50) NOT NULL,
    descricao TEXT,
    data_agendada TIMESTAMP WITH TIME ZONE,
    data_conclusao TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) NOT NULL,
    prioridade VARCHAR(20),
    planta_id INTEGER REFERENCES plantas(id) ON DELETE SET NULL,
    ambiente_id INTEGER REFERENCES ambientes(id) ON DELETE SET NULL,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    recorrente BOOLEAN,
    frequencia_dias INTEGER,
    diario_cultivo_id INTEGER REFERENCES diario_cultivos(id) ON DELETE SET NULL
);

-- Cria a tabela anotacoes
CREATE TABLE IF NOT EXISTS anotacoes (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    titulo VARCHAR(100) NOT NULL,
    conteudo TEXT NOT NULL,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    planta_id INTEGER REFERENCES plantas(id) ON DELETE CASCADE,
    ambiente_id INTEGER REFERENCES ambientes(id) ON DELETE CASCADE,
    tarefa_id INTEGER REFERENCES tarefas(id) ON DELETE CASCADE,
    tags VARCHAR(200)
);

-- Cria a tabela estagio_crescimentos
CREATE TABLE IF NOT EXISTS estagio_crescimentos (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    planta_id INTEGER REFERENCES plantas(id) ON DELETE CASCADE,
    estagio VARCHAR(100) NOT NULL,
    data_inicio TIMESTAMP WITH TIME ZONE NOT NULL,
    data_fim TIMESTAMP WITH TIME ZONE
);

-- Cria a tabela registro_crescimentos
CREATE TABLE IF NOT EXISTS registro_crescimentos (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    planta_id INTEGER REFERENCES plantas(id) ON DELETE CASCADE,
    altura NUMERIC,
    numero_folhas INTEGER,
    diametro_caule NUMERIC,
    estagio VARCHAR(50),
    observacoes TEXT
);

-- Cria a tabela registro_plantas
CREATE TABLE IF NOT EXISTS registro_plantas (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    planta_id INTEGER REFERENCES plantas(id) ON DELETE CASCADE,
    data_registro TIMESTAMP WITH TIME ZONE NOT NULL,
    tipo VARCHAR(100) NOT NULL,
    valor VARCHAR(255),
    observacao TEXT
);



-- Cria a tabela tarefa_plantas (tabela de junção para TarefaPlanta)
CREATE TABLE IF NOT EXISTS tarefa_plantas (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    planta_id INTEGER REFERENCES plantas(id) ON DELETE CASCADE,
    descricao TEXT NOT NULL,
    data_prevista TIMESTAMP WITH TIME ZONE NOT NULL,
    data_realizada TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) NOT NULL
);

-- Cria a tabela tarefa_templates
CREATE TABLE IF NOT EXISTS tarefa_templates (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    tipo VARCHAR(50) NOT NULL,
    frequencia_dias INTEGER,
    instrucoes TEXT,
    especie_id INTEGER
);

-- Cria a tabela diario_plantas (tabela de junção para many2many)
CREATE TABLE IF NOT EXISTS diario_cultivo_plantas (
    diario_cultivo_id INTEGER REFERENCES diario_cultivos(id) ON DELETE CASCADE,
    planta_id INTEGER REFERENCES plantas(id) ON DELETE CASCADE,
    PRIMARY KEY (diario_cultivo_id, planta_id)
);

-- Cria a tabela diario_ambientes (tabela de junção para many2many)
CREATE TABLE IF NOT EXISTS diario_ambientes (
    diario_cultivo_id INTEGER REFERENCES diario_cultivos(id) ON DELETE CASCADE,
    ambiente_id INTEGER REFERENCES ambientes(id) ON DELETE CASCADE,
    PRIMARY KEY (diario_cultivo_id, ambiente_id)
);