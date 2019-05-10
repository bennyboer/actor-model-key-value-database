pipeline {
    agent none
    stages {
        stage('Build') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh '''
                    git reset --hard
                    go version
                    go env
                    cd build
                    chmod +x ./install_protoc.sh
                    . ./install_protoc.sh
                    cd ..
                    chmod +x ./build.sh
                    ./build.sh
                '''
            }
        }
        stage('Test') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'echo run tests with code coverage...'
                sh 'go test ./... -cover'
            }
        }
        stage('Lint') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }   
            steps {
                sh 'golangci-lint run --deadline 20m --enable-all --disable=goimports --disable=lll --disable=dupl --tests=false'
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
