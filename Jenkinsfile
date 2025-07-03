pipeline {
    agent any

    environment {
        // Variáveis de ambiente para o banco de dados de teste
        DB_HOST_TEST = 'db_test'
        DB_PORT_TEST = '5432' // Porta interna do container do DB
        DB_USER_TEST = 'testuser'
        DB_PASSWORD_TEST = 'testpassword'
        DB_NAME_TEST = 'cultivo_test_db'
        APP_ENV = 'test' // Para que o Go carregue a config de teste
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Backend (Go)') {
            steps {
                dir('cultivo-api-go') {
                    sh 'go mod tidy'
                    sh 'go build -o bin/cultivo-api ./cmd/cultivo-api-go'
                }
            }
        }

        stage('Test Backend (Go)') {
            steps {
                dir('cultivo-api-go') {
                    // Iniciar o ambiente de teste do DB
                    sh 'docker-compose -f docker-compose.test.yml up -d'
                    // Esperar o DB ficar pronto (pode ser necessário um script de espera mais robusto)
                    sh 'sleep 10' // Apenas um delay simples, idealmente usar um healthcheck
                    // Executar testes de integração
                    sh 'go test -v ./internal/integration_test/...'
                    // Executar testes unitários
                    sh 'go test -v ./...'
                    // Parar e remover o ambiente de teste do DB
                    sh 'docker-compose -f docker-compose.test.yml down -v'
                }
            }
        }

        stage('Build Frontend (Angular)') {
            steps {
                dir('cultivo-web') {
                    sh 'npm install'
                    sh 'npm run build -- --configuration production'
                }
            }
        }

        stage('Test Frontend (Angular E2E)') {
            steps {
                dir('cultivo-web') {
                    // Iniciar o backend Go para os testes E2E do frontend
                    // Assumindo que o backend pode ser iniciado separadamente para E2E
                    // ou que o ambiente de teste do backend já está rodando
                    // Para simplificar, vamos apenas rodar os testes Cypress
                    // Em um cenário real, você iniciaria o backend e o frontend aqui
                    sh 'npx cypress run'
                }
            }
        }

        stage('Push Docker Images') {
            steps {
                script {
                    // Login no Docker Hub (ou outro registry)
                    // As credenciais 'docker-hub-credentials' devem ser configuradas no Jenkins
                    // com o ID 'docker-hub-credentials'
                    withCredentials([usernamePassword(credentialsId: 'docker-hub-credentials', passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
                        sh "echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin"
                    }

                    // Push da imagem do Backend
                    dir('cultivo-api-go') {
                        sh 'docker build -t your-dockerhub-username/cultivo-api-go:latest .'
                        sh 'docker push your-dockerhub-username/cultivo-api-go:latest'
                    }

                    // Push da imagem do Frontend
                    dir('cultivo-web') {
                        sh 'docker build -t your-dockerhub-username/cultivo-web:latest .'
                        sh 'docker push your-dockerhub-username/cultivo-web:latest'
                    }

                    sh 'docker logout'
                }
            }
        }
    }
}