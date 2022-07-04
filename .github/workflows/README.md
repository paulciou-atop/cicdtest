## GitHub CICD Workflows for NMS

Currently, we have three workflows. 
(1) basedocker.yml - This is to create the base docker image, which will install all the utlities, go packages and libraries for building NMS. The image will be stored in the private docker repository which will be pulled by other workflows to build the NMS code. 

This workflow will be triggered by two methods. 
    - It is triggered by the GitHub GUI Workflow Dispatch (Run by choosing the action and click the "Run Workflow" button at The Repository/Actions page.)
    - It is scheduled to be run two times a month (on date 15 and 30) 

(2) pushfeature.yml - This will build the NMS code, run the unit test and run the simulation test. 
These will be built and run on GitHub Runners. 
It is triggered whenever the code is pushed to any branch. If you want to skip this process, please name your branch to start with "test...." or "tmp...."

(3) PRcheck.yml - This will build the NMS code, run the unit test and run the system integration test on the self-hosted runner in our Testbed with the real devices. 
It is triggered when the code requests for PR to main branch. 

## To Do List
- Currently, the workflow is mainly for Linux platform. We need to add workflows for windows and other platforms. 
- The testing script for both simulation and system integration test on real-devices are under development. 