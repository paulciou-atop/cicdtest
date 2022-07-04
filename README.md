
# ATOP NMS

All NMS related materials go under here. Each sub directory from here should have a README.md file.

Each sub directory also should have a build.bat and build.sh files.  These are build scripts that can be run without any arguments on Windows or Linux/Mac OS.  They serve as a tool and documentation for building everything under the sub directory.  The script should recursively build everything under the sub directory.




## Project Definition Draft

https://docs.google.com/document/d/1pRh-eG8kp5kSrpZhoqz-xixzDGHV45GzBKPFGOZxOBM/edit?usp=sharing

## Project Management and Worklogs

https://docs.google.com/document/d/10Z7Q4gAzrcHc9xCRwlTINYAlUvfIi2c43wehRvIgNtc/edit?usp=sharing


## Build
1. Generate gRPC code first
   Go to /api directory and run build script
   ```
   $ cd /api
   $ ./build.sh
   ```
