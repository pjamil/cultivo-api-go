# criar

Vou criar um arquivo de configuração para o REST Client do Visual Studio Code (uma extensão muito útil para testar APIs diretamente do editor). Este arquivo conterá todas as requisições necessárias para testar nossa API de cultivo de plantas.

## Arquivo de Requisições para REST Client

Crie um arquivo chamado `api-tests.http` na raiz do seu projeto (ou em uma pasta `requests/`):

```http
@baseUrl = http://localhost:8080/api/v1
@contentType = application/json

### Health Check
GET {{baseUrl}}/health HTTP/1.1

### List all plants
GET {{baseUrl}}/plants HTTP/1.1

### Get a specific plant
GET {{baseUrl}}/plants/1 HTTP/1.1

### Create a new plant
POST {{baseUrl}}/plants HTTP/1.1
Content-Type: {{contentType}}

{
  "name": "Tomate Cereja",
  "species": "tomato",
  "planting_date": "2023-05-15",
  "status": "growing",
  "growth_stage": "vegetative",
  "notes": "Planta saudável, sem pragas"
}

### Update a plant
PUT {{baseUrl}}/plants/1 HTTP/1.1
Content-Type: {{contentType}}

{
  "name": "Tomate Cereja Orgânico",
  "species": "tomato",
  "planting_date": "2023-05-15",
  "harvest_date": "2023-08-20",
  "status": "harvested",
  "growth_stage": "mature",
  "notes": "Colhido com sucesso"
}

### Delete a plant
DELETE {{baseUrl}}/plants/1 HTTP/1.1

### Get plants by species
GET {{baseUrl}}/plants?species=tomato HTTP/1.1

### Get plants by status
GET {{baseUrl}}/plants?status=growing HTTP/1.1

### Health Ready Check
GET {{baseUrl}}/health/ready HTTP/1.1

### Health Live Check
GET {{baseUrl}}/health/live HTTP/1.1
```

## Como Usar no VS Code

1. **Instale a extensão**:

   - Abra o VS Code
   - Vá para Extensions (Ctrl+Shift+X)
   - Busque por "REST Client"
   - Instale a extensão "REST Client" por Huachao Mao

2. **Execute as requisições**:

   - Abra o arquivo `api-tests.http`
   - Você verá "Send Request" acima de cada requisição
   - Clique para enviar cada requisição individualmente

3. **Visualize as respostas**:
   - As respostas aparecerão em um painel à direita
   - Você verá status code, headers e body da resposta

## Variáveis de Ambiente Adicionais

Você pode adicionar mais variáveis no topo do arquivo para facilitar os testes:

```http
@adminToken = seu-token-jwt-aqui
@userId = 123
```

E usar em requisições que precisam de autenticação:

```http
### Request with auth header
GET {{baseUrl}}/secure-endpoint HTTP/1.1
Authorization: Bearer {{adminToken}}
```

## Exemplo Completo com Autenticação

Se sua API tiver autenticação JWT, você pode adicionar:

```http
### Login to get token
POST http://localhost:8080/auth/login HTTP/1.1
Content-Type: {{contentType}}

{
  "email": "user@example.com",
  "password": "sua-senha"
}

### Save token from response to variable
@token = {{login.response.body.token}}

### Authenticated request
GET {{baseUrl}}/plants HTTP/1.1
Authorization: Bearer {{token}}
```

## Dicas Úteis

1. **Organização**:

   - Agrupe requisições relacionadas com comentários `###`
   - Use arquivos separados para diferentes contextos (auth.http, plants.http, etc.)

2. **Debugging**:

   - Inspecione o corpo da resposta e headers
   - Compare requisições/respostas

3. **Ambientes**:
   - Crie arquivos separados para dev, staging e production
   - Ou use variáveis para diferentes ambientes:

```http
@dev = http://localhost:8080
@prod = https://api.plantcultivation.com

@baseUrl = {{dev}}/api/v1
```

4. **Documentação**:
   - Adicione comentários para explicar cada endpoint
   - Documente os parâmetros esperados

Esta configuração permite testar rapidamente toda a API sem sair do VS Code, sendo muito útil durante o desenvolvimento. O REST Client também permite salvar históricos de requisições e tem suporte a vários tipos de autenticação.
