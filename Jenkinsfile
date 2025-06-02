pipeline {
    agent any
    environment {
        DOCKERHUB_CREDENTIALS = 'docker-hub'
        DOCKER_IMAGE = 'pjamil/cultivo-api-go'
    }

    stages {
        stage('Docker Build Image') {
            steps {
                sh 'docker build -t pjamil/cultivo-api-go:latest .'
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
    }
}
