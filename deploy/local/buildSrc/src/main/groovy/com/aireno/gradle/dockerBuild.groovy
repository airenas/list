package com.aireno.gradle;

import org.gradle.api.DefaultTask;
import org.gradle.api.tasks.TaskAction;

class DCRedeploy extends DefaultTask {
     String serviceName 
     String composePath = "."

     @TaskAction
     def run() {

        println '==== Redeploy ' + serviceName
        project.exec {
            workingDir composePath
            commandLine 'docker-compose', 'stop', serviceName
        }
        project.exec {
            workingDir composePath
            commandLine 'docker-compose', 'rm', '-f', serviceName
        }
        project.exec {
            workingDir composePath
            commandLine 'docker-compose', 'build', serviceName
        }
        project.exec {
            workingDir composePath
            commandLine 'docker-compose', 'up', '-d', serviceName
        }
        println '==== Redeployed ' + serviceName
     }
 }

 class DBuild extends DefaultTask {
     String tag 
     String dir = "."

     @TaskAction
     def run() {

        println '==== Building ' + tag
        project.exec {
            commandLine 'docker', 'build', '-t', tag, dir
        }
        println '==== Built ' + tag
     }
 }
 
 class DPush extends DefaultTask {
     String tag 
     @TaskAction
     def run() {
        println '==== Pushing ' + tag
        project.exec {
            commandLine 'docker', 'push', tag
        }
        println '==== Pushed ' + tag
     }
 }