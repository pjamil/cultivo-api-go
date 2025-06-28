### O que é `golang-migrate/migrate` e por que usá-lo?

`golang-migrate/migrate` é uma ferramenta de linha de comando e biblioteca Go para gerenciar migrações de banco de dados. Em termos simples, ela permite que você controle as alterações no esquema do seu banco de dados (criação de tabelas, adição de colunas, etc.) de forma versionada.

**Por que é importante?**
*   **Controle de Versão:** Assim como você versiona seu código, você versiona seu banco de dados. Isso garante que todos os desenvolvedores estejam trabalhando com o mesmo esquema de banco de dados e que o ambiente de produção possa ser atualizado de forma controlada.
*   **Consistência:** Garante que o esquema do banco de dados seja o mesmo em todos os ambientes (desenvolvimento, staging, produção).
*   **Reversão Segura:** Permite reverter alterações no esquema do banco de dados de forma controlada, caso algo dê errado.

### Como usar no seu projeto `cultivo-api-go`

Integramos o `golang-migrate/migrate` de duas formas principais:

1.  **Via `Makefile` (para desenvolvimento local):** Comandos para criar, aplicar e reverter migrações manualmente.
2.  **Via `docker-compose` (para desenvolvimento e produção):** As migrações são aplicadas automaticamente quando o serviço `app` é iniciado.

Vamos detalhar cada um:

---
#### 1. Gerenciando Migrações no Desenvolvimento Local (via `Makefile`)

No seu ambiente de desenvolvimento local, você usará os comandos do `Makefile` para criar e aplicar migrações.

*   **`make migrate-create` (Criar uma Nova Migração)**
    *   **Propósito:** Este comando é usado para gerar um novo par de arquivos de migração: um arquivo `.up.sql` (para aplicar a alteração) e um arquivo `.down.sql` (para reverter a alteração).
    *   **Como usar:**
        ```bash
        make migrate-create
        ```
        Ao executar, o terminal irá te perguntar: `Enter migration name:`. Digite um nome descritivo para a sua migração (ex: `add_column_to_planta_table`, `create_new_user_settings_table`).
    *   **O que acontece:** A ferramenta `migrate` criará dois arquivos na pasta `internal/infrastructure/database/migrations/` com um timestamp no início (ex: `000002_add_new_feature.up.sql` e `000002_add_new_feature.down.sql`). Você então editará esses arquivos com o SQL necessário.
        *   No arquivo `.up.sql`, você escreverá o SQL para *aplicar* a mudança (ex: `ALTER TABLE plantas ADD COLUMN cor VARCHAR(50);`).
        *   No arquivo `.down.sql`, você escreverá o SQL para *reverter* a mudança (ex: `ALTER TABLE plantas DROP COLUMN cor;`).

*   **`make migrate-up` (Aplicar Migrações Pendentes)**
    *   **Propósito:** Este comando executa todas as migrações `.up.sql` que ainda não foram aplicadas ao seu banco de dados local.
    *   **Como usar:**
        ```bash
        make migrate-up
        ```
    *   **Quando usar:**
        *   Depois de criar novos arquivos de migração e preenchê-los com SQL.
        *   Quando você puxa alterações do repositório que incluem novas migrações de outros desenvolvedores.
    *   **O que acontece:** O `migrate` verifica o estado do seu banco de dados, identifica quais migrações `.up.sql` ainda não foram executadas e as aplica em ordem cronológica.

*   **`make migrate-down` (Reverter a Última Migração Aplicada)**
    *   **Propósito:** Este comando reverte a *última* migração que foi aplicada ao seu banco de dados local, executando o SQL do arquivo `.down.sql` correspondente.
    *   **Como usar:**
        ```bash
        make migrate-down
        ```
    *   **Quando usar:** Principalmente durante o desenvolvimento, se você cometer um erro em uma migração recém-aplicada e precisar desfazê-la rapidamente.
    *   **Cuidado:** Use com cautela, pois reverter migrações pode levar à perda de dados se você estiver revertendo uma migração que removeu colunas ou tabelas.

---
#### 2. Migrações Automáticas com `docker-compose`

Para simplificar o desenvolvimento e garantir que o banco de dados esteja sempre pronto, configuramos o `docker-compose` para executar as migrações automaticamente quando o serviço `app` é iniciado.

*   **Como funciona:** No arquivo `entrypoint.sh` (que é executado antes do seu aplicativo Go), adicionamos a linha:
    ```bash
    /go/bin/migrate -path /app/internal/infrastructure/database/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up
    ```
    Esta linha instrui o `migrate` a aplicar todas as migrações pendentes.
*   **Como usar:**
    ```bash
    docker-compose up --build
    ```
    Quando você executa este comando (ou `docker-compose up` se a imagem já estiver construída), o contêiner do seu aplicativo Go será iniciado, e as migrações serão executadas automaticamente antes que o aplicativo comece a rodar. Isso é muito conveniente para garantir que seu banco de dados de desenvolvimento esteja sempre atualizado.

---
#### 3. Gerenciando Migrações em Produção/VPS (via `Makefile.prod`)

Para ambientes de produção ou VPS, você usará um `Makefile` separado (`Makefile.prod`) que aponta para o banco de dados correto e espera que as variáveis de ambiente do banco de dados estejam configuradas.

*   **`make -f Makefile.prod migrate-up` (Aplicar Migrações em Produção)**
    *   **Propósito:** Aplica todas as migrações pendentes ao seu banco de dados de produção.
    *   **Como usar:**
        ```bash
        # Certifique-se de que as variáveis de ambiente do banco de dados de produção estejam configuradas
        # Ex: export DB_HOST=your_prod_db_host
        #     export DB_USER=your_prod_db_user
        #     ...
        make -f Makefile.prod migrate-up
        ```
    *   **Quando usar:** Quando você implanta uma nova versão do seu aplicativo que inclui alterações no esquema do banco de dados.
    *   **Importante:** Este comando deve ser executado no servidor de produção ou em um ambiente que tenha acesso ao banco de dados de produção e as credenciais corretas.

*   **`make -f Makefile.prod migrate-down` (Reverter Migrações em Produção)**
    *   **Propósito:** Reverte a última migração aplicada no banco de dados de produção.
    *   **Como usar:**
        ```bash
        # Certifique-se de que as variáveis de ambiente do banco de dados de produção estejam configuradas
        make -f Makefile.prod migrate-down
        ```
    *   **CUIDADO EXTREMO:** **NUNCA** use este comando em um ambiente de produção a menos que você saiba exatamente o que está fazendo e tenha um backup completo do seu banco de dados. Reverter migrações em produção pode levar à perda de dados e instabilidade do sistema. Geralmente, em produção, você só aplica migrações (`migrate-up`). Se houver um problema, a estratégia comum é aplicar uma nova migração para corrigir o problema, não reverter.

---
### Exemplo de Fluxo de Trabalho Típico

Vamos imaginar que você precisa adicionar uma nova coluna `cor` à tabela `plantas`.

1.  **Crie a migração:**
    ```bash
    make migrate-create
    # Digite: add_cor_to_plantas_table
    ```
    Isso criará `00000X_add_cor_to_plantas_table.up.sql` e `00000X_add_cor_to_plantas_table.down.sql`.

2.  **Edite os arquivos SQL:**
    *   Abra `00000X_add_cor_to_plantas_table.up.sql` e adicione:
        ```sql
        ALTER TABLE plantas ADD COLUMN cor VARCHAR(50);
        ```
    *   Abra `00000X_add_cor_to_plantas_table.down.sql` e adicione:
        ```sql
        ALTER TABLE plantas DROP COLUMN cor;
        ```

3.  **Aplique a migração no desenvolvimento:**
    ```bash
    make migrate-up
    ```
    Ou, se você estiver usando Docker Compose para o desenvolvimento, basta reiniciar seus contêineres:
    ```bash
    docker-compose down && docker-compose up --build
    ```
    O `entrypoint.sh` cuidará de aplicar a nova migração.

4.  **Teste suas alterações:** Verifique se a nova coluna está presente e se seu aplicativo funciona corretamente com ela.

5.  **Commit e Push:**
    ```bash
    git add .
    git commit -m "feat: add 'cor' column to plantas table"
    git push
    ```

6.  **Implantação em Produção:**
    Quando você implantar seu código em produção (por exemplo, via CI/CD ou manualmente), você executará:
    ```bash
    # No servidor de produção, após atualizar o código
    make -f Makefile.prod migrate-up
    ```
    Isso aplicará a nova coluna `cor` à tabela `plantas` no seu banco de dados de produção.
