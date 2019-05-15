package com.aireno.gradle;

import org.gradle.api.DefaultTask;
import org.gradle.api.tasks.TaskAction;

class DCRedeploy extends DefaultTask {
     String serviceName 
     String composePath = "."
     String projectName = ""
     String fileName = ""

     @TaskAction
     def run() {
        def str = new ArrayList <String> ();
        str.add('docker-compose')
        if (projectName != ""){
            str.add('-p')
            str.add(projectName)
         }
        if (fileName != ""){
            str.add('-f')
            str.add(fileName)
        }
        println '==== Redeploy ' + serviceName
        project.exec {
            workingDir composePath
            commandLine str + 'stop' + serviceName
        }
        project.exec {
            workingDir composePath
            commandLine str + 'rm' + '-f' + serviceName
        }
        project.exec {
            workingDir composePath
            commandLine str + 'build' + serviceName
        }
        project.exec {
            workingDir composePath
            commandLine str + 'up' + '-d' + serviceName
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