pipeline {
    agent none
    stages {
        stage('Build') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'cat /etc/os-release'
                sh 'chmod +x ./build.sh'
                sh './build.sh'
            }
        }
        stage('Test') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'echo run tests...'
                sh 'echo CLI tests...'
                sh 'go test treecli'
                sh 'echo Service tests...'
                sh 'go test treeservice'
            }
        }
        stage('Lint') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }   
            steps {
                sh 'golangci-lint run --deadline 20m --enable-all'
            }
        }
        stage('Build Docker Image') {
            agent any
            steps {
                sh "docker-build-and-push -b ${BRANCH_NAME} -s treeservice -f treeservice.dockerfile"
                sh "docker-build-and-push -b ${BRANCH_NAME} -s treecli -f treecli.dockerfile"
            }
        }
    }
}
