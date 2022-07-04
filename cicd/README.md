## CI/CD related things

There are three directories in CICD now. 

(1) docker - This directory contains all the Docker files and docker comppose files that use in the workflows in .github/workflows/ directory. 

- Dockerfile_base_linux: This file is to build the docker base image with all utilties, packages and libaries for building NMS code. 
- Dockerfile_build_linux: This file is to build the docker build iamge based on the docker base image and the current snapshot of the NMS code. 
- docker-compose-cicdetest.yml: This compose file will run containers for all services (one container per service). It will create the container (name ends with "_cicdnms") that will act like a client for testing, which will have all grpcurl, the CLI commands for all services and the iputils. 

(2) build_script - This directory contains the build scripts to run locally for the developers if they want to run manually in their own testbed, not through CICD workflow. 
- build_base_image.sh - to generate the base image
- build_nms.sh - to generate the whole NMS. However, this is to pull the base image from appleatop/nms_base_image:latest (please get the repository token from me)

(3) test_script - This directory contains the shellscript for the simulation and integration test, which will be linked to worksflows in .github/workflows/ directory. 


## To Do list 
- Add CICD for other platforms
- Generate simulation test scripts and testbed system integration test scripts