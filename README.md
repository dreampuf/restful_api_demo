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
    
    bin/restapi -help
    Usage of ./bin/restapi:
      -dbhost string
           	Connecting DB Host (default "127.0.0.1")
      -dbname string
           	Connecting DB Name (default "api")
      -dbpassword string
           	Connecting DB Password (default "api")
      -dbport uint
           	Connecting DB Port (default 5432)
      -dbuser string
           	Connecting DB User (default "api")
      -host string
           	WebService Host (default "127.0.0.1")
      -port uint
           	WebService Port (default 8080)
    
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
