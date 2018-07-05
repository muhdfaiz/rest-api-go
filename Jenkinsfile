#!/bin/bash
node {
    def workspace
    
    stage('Build') {

        try {
            //slackSend (color: '#FFFF00', message: "STARTED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")

            // Clean Up Directory
            deleteDir()
            
            // Set Go Environment Variables
            workspace = pwd()
            env.GOPATH="${workspace}/"
            env.GOROOT="/usr/local/go"
            env.PATH="${env.PATH}:/usr/local/go/bin:${workspace}/bin"
    
            // Display all environment variables
            sh "printenv"
            
            // Create directory src, bin and pkg
            sh 'mkdir bin'
            sh 'mkdir src'
            sh 'mkdir pkg'
    
            // Checkout shoppermate api repo in subfolder src/bitbucket.org/cliqers/shoppermate-api'
            checkout([$class: 'GitSCM', branches: scm.branches, doGenerateSubmoduleConfigurations: false, extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'src/bitbucket.org/cliqers/shoppermate-api']], submoduleCfg: [], userRemoteConfigs: [[credentialsId: 'ccbd457c-f55b-48a7-88d2-f78639575da6', url: 'https://muhdfaiz@bitbucket.org/cliqers/shoppermate-api.git']]])
            
            dir('src/bitbucket.org/cliqers/shoppermate-api') {
                sh 'mkdir storages'
                sh 'rm -rf glide.lock'
                sh 'glide install'
                sh 'touch .env'
                sh 'echo "ENVIRONMENT:local" >> .env'
                sh 'echo "DEBUG:true" >> .env'
                sh 'echo "DEBUG_DATABASE:false" >> .env'
                sh 'echo "ENABLE_HTTPS:false" >> .env'
                sh 'echo "SHOPPERMATE_EMAIL_API_URL:http://api.example.com:5000/" >> .env'
                sh 'echo "SEND_EMAIL_EVENT=false" >> .env'
                sh 'echo "SEND_SMS=false" >> .env'
                sh 'echo "FACEBOOK_APP_ID:" >> .env'
                sh 'echo "FACEBOOK_APP_SECRET:" >> .env'
                sh 'echo "DEBUG_FACEBOOK_APP_ID:" >> .env'
                sh 'echo "DEBUG_FACEBOOK_APP_SECRET:" >> .env'
                sh 'echo "AWS_ACCESS_KEY_ID:" >> .env'
                sh 'echo "AWS_SECRET_ACCESS_KEY:" >> .env'
                sh 'echo "AWS_S3_BUCKET_NAME:" >> .env'
                sh 'echo "AWS_S3_REGION_NAME:" >> .env'
                sh 'echo "AWS_S3_URL:https://s3-ap-southeast-1.amazonaws.com/" >> .env'
                sh 'echo "JWT_TOKEN_SECRET:" >> .env'
                sh 'echo "STORAGE_PATH:" >> .env'
                sh 'echo "TEST_DB_HOST:localhost" >> .env'
                sh 'echo "TEST_DB_PORT:3306" >> .env'
                sh 'echo "TEST_DB_NAME:" >> .env'
                sh 'echo "TEST_DB_USERNAME:" >> .env'
                sh 'echo "TEST_DB_PASSWORD:" >> .env'
                sh 'echo "MOCEAN_SMS_URL:" >> .env'
                sh 'echo "MOCEAN_SMS_USERNAME:" >> .env'
                sh 'echo "MOCEAN_SMS_PASSWORD:" >> .env'
                sh 'echo "MOCEAN_SMS_CODING:" >> .env'
                sh 'echo "MCOEAN_SMS_SUBJECT:" >> .env'
                sh 'echo "MAX_DEAL_RADIUS_IN_KM=10" >> .env'
                sh 'echo "UTC_TIMEZONE=8" >> .env'
                sh 'printenv'
            }

            dir('src/bitbucket.org/cliqers/shoppermate-api/vendor/github.com/jstemmer/go-junit-report') {
                sh "go build -o ${workspace}/bin/go-junit-report"
            }

            dir('bin') {
                sh 'chmod 777 go-junit-report'
            }

        } catch (e) {
            currentBuild.result = "FAILED"
            //slackSend (color: '#FF0000', message: "FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")
            throw e
        }
    }
    
    stage('Test') {
        try {
            dir('src/bitbucket.org/cliqers/shoppermate-api') {
                sh 'go test -v $(go list ./... | grep -v /vendor/) | go-junit-report > report.xml'
            }
            //slackSend (color: '#00FF00', message: "SUCCESSFUL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")
        } catch (e) {
            currentBuild.result = "FAILED"
            //slackSend (color: '#FF0000', message: "FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")
            throw e
        }

    }

    stage('Result') {
        dir('src/bitbucket.org/cliqers/shoppermate-api') {
            junit '*.xml' 
        }
        step([$class: 'JUnitResultArchiver', testResults: 'src/bitbucket.org/cliqers/shoppermate-api/*.xml'])
        
        if (currentBuild.result == 'UNSTABLE')
            currentBuild.result = 'FAILURE'
    }
}