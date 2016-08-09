# Restful API demo

About this project.

# Developing Environment Setup

    git clone git@github.com:dreampuf/restful_api_demo.git
    cd restful_api_demo
    export GOPATH="$PWD"
    export PATH="$GOPATH/bin:$PATH"
    go get github.com/constabulary/gb/...
    gb vendor restore
    
    gb build cmd/...
    
    gb test
    
## With Docker

    ## Setup a database
    docker run --name rest_api_database -p 5432:5432 -e POSTGRES_PASSWORD=api -e POSTGRES_USER=api -e POSTGRES_DB=api -d postgres
    
    ## DB migration
    bin/dbmigrate init
    bin/dbmigrate up
    
    ## Start the application
    bin/restapi
    
    
    
# Testing with Docker Compose

*TODO*
