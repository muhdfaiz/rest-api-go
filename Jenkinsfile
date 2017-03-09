#!/bin/bash
node {
    def workspace
    
    stage('Build') {

        try {
            slackSend (color: '#FFFF00', message: "STARTED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")

            // Clean Up Directory
            deleteDir()
            
            // Set Go Environment Variables
            workspace = pwd()
            env.GOPATH="${workspace}/"
            env.GOROOT="/usr/local/go"
            env.PATH="${env.PATH}:/usr/local/go/bin:${workspace}/bin"
            
            // Reload shell to take effect latest environment variables
            sh "source ~/.profile"
    
            // Display all environment variables
            sh "printenv"
            
            // Create directory src, bin and pkg
            sh 'mkdir bin'
            sh 'mkdir src'
            sh 'mkdir pkg'
    
            // Checkout shoppermate api repo in subfolder src/bitbucket.org/cliqers/shoppermate-api'
            checkout([$class: 'GitSCM', branches: [[name: '**']], doGenerateSubmoduleConfigurations: false, extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'src/bitbucket.org/cliqers/shoppermate-api']], submoduleCfg: [], userRemoteConfigs: [[credentialsId: '26ce324d-eab2-4d6f-b59b-ffa8100c6920', url: 'https://muhdfaiz@bitbucket.org/cliqers/shoppermate-api.git']]])
            
            dir('src/bitbucket.org/cliqers/shoppermate-api') {
                sh 'mkdir storages'
                sh 'glide install'
                sh 'touch .env'
                sh 'echo -e "ENVIRONMENT:local" >> .env'
                sh 'echo -e "DEBUG:false" >> .env'
                sh 'echo -e "DEBUG_DATABASE:false" >> .env'
                sh 'echo -e "ENABLE_HTTPS:false" >> .env'
                sh 'echo -e "SHOPPERMATE_EMAIL_API_URL:http://api.shoppermate.com:5000/" >> .env'
                sh 'echo -e "SEND_EMAIL_EVENT=false" >> .env'
                sh 'echo -e "SEND_SMS=false" >> .env'
                sh 'echo -e "FACEBOOK_APP_ID:1390260387947574" >> .env'
                sh 'echo -e "FACEBOOK_APP_SECRET:6646a7a1057f9dd9c6a3e18f3615b081" >> .env'
                sh 'echo -e "DEBUG_FACEBOOK_APP_ID:1572196483090312" >> .env'
                sh 'echo -e "DEBUG_FACEBOOK_APP_SECRET:c013593814f1b4b3383c19f76eb038b1" >> .env'
                sh 'echo -e "AWS_ACCESS_KEY_ID:AKIAISDRKR2IUDLU44FA" >> .env'
                sh 'echo -e "AWS_SECRET_ACCESS_KEY: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" >> .env'
                sh 'echo -e "AWS_S3_BUCKET_NAME:shoppermate-local" >> .env'
                sh 'echo -e "AWS_S3_REGION_NAME:ap-southeast-1" >> .env'
                sh 'echo -e "AWS_S3_URL:https://s3-ap-southeast-1.amazonaws.com/" >> .env'
                sh 'echo -e "JWT_TOKEN_SECRET:gN2T5znLzeSTBvdeKPGZBAUdFb6fSrjK" >> .env'
                sh 'echo -e "STORAGE_PATH:src/bitbucket.org/cliqers/shoppermate-api/storages/" >> .env'
                sh 'echo -e "TEST_DB_HOST:localhost" >> .env'
                sh 'echo -e "TEST_DB_PORT:3306" >> .env'
                sh 'echo -e "TEST_DB_NAME:shoppermate_test" >> .env'
                sh 'echo -e "TEST_DB_USERNAME:root" >> .env'
                sh 'echo -e "TEST_DB_PASSWORD:123456" >> .env'
                sh 'echo -e "MOCEAN_SMS_URL:http://183.81.161.84:13016/cgi-bin/sendsms" >> .env'
                sh 'echo -e "MOCEAN_SMS_USERNAME:shoppermate" >> .env'
                sh 'echo -e "MOCEAN_SMS_PASSWORD:s28Dua3p" >> .env'
                sh 'echo -e "MOCEAN_SMS_CODING:1" >> .env'
                sh 'echo -e "MCOEAN_SMS_SUBJECT:ShopperMate" >> .env'
                sh 'echo -e "MAX_DEAL_RADIUS_IN_KM=10" >> .env'
                sh 'echo -e "UTC_TIMEZONE=8" >> .env'
                sh 'mysql -u root "-p123456" shoppermate_test -e "show tables" | grep -v Tables_in | grep -v "+" | gawk \'{print "drop table " $1 ";"}\' | mysql -u root "-p123456" shoppermate_test'
                sh 'mysql -u root "-p123456" shoppermate_test < shoppermate_test.sql'
                sh 'printenv'
            }

            dir('src/bitbucket.org/cliqers/shoppermate-api/vendor/github.com/tebeka/go2xunit') {
                sh 'go build -o ../../../../go2xunit'
            }

        } catch (e) {
            currentBuild.result = "FAILED"
            slackSend (color: '#FF0000', message: "FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")
            throw e
        }
    }
    
    stage('Test') {
        try {
            dir('src/bitbucket.org/cliqers/shoppermate-api/application/v1_1') {
                sh 'go test -v | tee ../../test_result.out'
            }

            dir('src/bitbucket.org/cliqers/shoppermate-api/') {
                sh './go2xunit -fail -input test_result.out -output tests.xml'
            }
            
            slackSend (color: '#00FF00', message: "SUCCESSFUL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")
        } catch (e) {
            currentBuild.result = "FAILED"
            slackSend (color: '#FF0000', message: "FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")
            throw e
        }

    }

    stage('Result') {
        junit '*.xml'
    }
}
