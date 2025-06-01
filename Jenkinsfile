pipeline {
    agent any

    environment {
        DOCKERHUB_CREDENTIALS = credentials('docker-hub-creds')
        // VPS_SSH_CREDENTIALS = credentials('vps-ssh-key')
        APP_NAME = 'cultivo-api-go'
        DOCKER_IMAGE = "pjamil/${APP_NAME}:${env.BUILD_NUMBER}"
        DOCKER_IMAGE_LATEST = "pjamil/${APP_NAME}:latest"
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', 
                url: 'git@gitea.paulojamil.dev.br:paulojamil.dev.br/cultivo-api-go.git',
                credentialsId: 'gitea-token'
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    docker.build(DOCKER_IMAGE)
                }
            }
        }

        stage('Push to Docker Hub') {
            steps {
                script {
                    docker.withRegistry('https://registry.hub.docker.com', DOCKERHUB_CREDENTIALS) {
                        docker.image(DOCKER_IMAGE).push()
                        docker.image(DOCKER_IMAGE).push('latest')
                    }
                }
            }
        }

        // stage('Deploy to VPS') {
        //     steps {
        //         sshagent(['vps-ssh-key']) {
        //             sh """
        //                 ssh -o StrictHostKeyChecking=no seu-usuario@seu-ip-vps << EOF
        //                     docker pull ${DOCKER_IMAGE}
        //                     cd /opt/${APP_NAME}
        //                     docker-compose down
        //                     docker-compose up -d
        //                 EOF
        //             """
        //         }
        //     }
        // }
    }

    post {
        success {
            slackSend color: 'good', message: "Build ${env.BUILD_NUMBER} deployed successfully!"
        }
        failure {
            slackSend color: 'danger', message: "Build ${env.BUILD_NUMBER} failed!"
        }
    }
}