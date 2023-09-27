pipeline {
    agent any

    tools {
        go 'go-1.17'
    }

    environment {
        AWS_REGION = 'us-east-1'
        GO114MODULE = 'on'
        CGO_ENABLED = 0
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
    }

    stages {
        stage('Go Vet') {
            steps {
                sh 'go version'
                sh 'go vet'
            }
        }
        stage('Go Test') {
            steps {
                withCredentials([
                    conjurSecretCredential(credentialsId: 'data-vault-D-App-CybrCLI-Application-CyberArk-vault-infamous.privilegecloud.cyberark.cloud-Svc_CybrCLI-username', variable: 'PAS_USERNAME'),
                    conjurSecretCredential(credentialsId: 'data-vault-D-App-CybrCLI-Application-CyberArk-vault-infamous.privilegecloud.cyberark.cloud-Svc_CybrCLI-password', variable: 'PAS_PASSWORD'),
                    conjurSecretCredential(credentialsId: 'data-vault-D-App-CybrCLI-ccp-client-certificate-password', variable: 'CCP_CLIENT_CERT'),
                    conjurSecretCredential(credentialsId: 'data-vault-D-App-CybrCLI-ccp-priv-key-password', variable: 'CCP_CLIENT_PRIVATE_KEY')
                ]) {
                    sh '''
                        set +x
                        PAS_HOSTNAME=https://infamous.privilegecloud.cyberark.cloud
                        CCP_CLIENT_CERT=$(echo $CCP_CLIENT_CERT | base64 --decode)
                        CCP_CLIENT_PRIVATE_KEY=$(echo $CCP_CLIENT_PRIVATE_KEY | base64 --decode)
                        set -x
                        go test -v ./pkg/cybr/api
                    '''
                }
            }
            
        }
        stage('Go Build') {
            steps {
                echo "Making Linux x64 binary..."
                sh 'GOOS=linux GOARCH=amd64 go build -o ./bin/${BUILD_TIMESTAMP}_linux_cybr .'
                echo "Making Darwin x64 binary..."
                sh 'GOOS=darwin GOARCH=amd64 go build -o ./bin/${BUILD_TIMESTAMP}_darwin_cybr .'
                echo "Making Homebrew Darwin x64 binary..."
                sh 'GOOS=darwin GOARCH=amd64 go build -o ./bin/${BUILD_TIMESTAMP}_cybr'
                echo "Making Darwin ARM binary..."
                sh 'GOOS=darwin GOARCH=arm64 go build -o ./bin/${BUILD_TIMESTAMP}_darwin_arm64_cybr .'
                echo "Making Windows x64 binary..."
                sh 'GOOS=windows GOARCH=amd64 go build -o ./bin/${BUILD_TIMESTAMP}_windows_cybr.exe .'
                echo "Finished making - files output to ./bin/"
            }
        }
        stage('Release to AWS S3') {
            steps {
                withCredentials([
                    conjurSecretCredential(credentialsId: 'SyncVault-LOB_CI-D-App-CybrCLI-Cloud Service-AWSAccessKeys-jenkins_cybr-cli-awsaccesskeyid', variable: 'AWS_ACCESS_KEY_ID'),
                    conjurSecretCredential(credentialsId: 'SyncVault-LOB_CI-D-App-CybrCLI-Cloud Service-AWSAccessKeys-jenkins_cybr-cli-password', variable: 'AWS_SECRET_ACCESS_KEY')
                ]) {
                    sh 'aws s3 cp ./bin/${BUILD_TIMESTAMP}_linux_cybr s3://cybr-cli-releases'
                    sh 'aws s3 cp ./bin/${BUILD_TIMESTAMP}_darwin_cybr s3://cybr-cli-releases'
                    sh 'aws s3 cp ./bin/${BUILD_TIMESTAMP}_cybr s3://cybr-cli-releases'
                    sh 'aws s3 cp ./bin/${BUILD_TIMESTAMP}_darwin_arm64_cybr s3://cybr-cli-releases'
                    sh 'aws s3 cp ./bin/${BUILD_TIMESTAMP}_windows_cybr.exe s3://cybr-cli-releases'
                }
            }
        }
    }
}
