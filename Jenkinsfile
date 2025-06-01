pipeline {
    agent {
        docker {
            image 'golang:1.21'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    }
    environment {
        DOCKERHUB_CREDENTIALS = 'dockerhub-credentials'
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
        // stage('Push to Docker Hub') {
        //     steps {
        //         withCredentials([
        //             usernamePassword(
        //                 credentialsId: 'registry-paulojamil',
        //                 passwordVariable: 'dockerHubPassword',
        //                 usernameVariable: 'dockerHubUser')
        //         ]) {
        //             sh 'docker login -u $dockerHubUser -p $dockerHubPassword registry.paulojamil.dev.br'
        //             sh 'docker push registry.paulojamil.dev.br/siscompras-api:latest'
        //         }
        //     }
        // }
    }
}
