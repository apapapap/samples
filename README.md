# User Management 

## Basic user management using GoLang protocol buffers.

### Pre-requisities
- Go 1.17 or higher
- Protoc 3.6.1 or higher

### Run the application as GRPC server
To run this application use the `Makefile`, you can use the following commands:
- **make clean**
  - This command will clean the existing auto-generated `.pb.go` files from the `pb/` directory.
 
- **make gen**
  - This command will auto-generate the `.pb.go` files in the `pb/` directory.  

- **make grpc**
  - This command will start the GRPC server on 0.0.0.0 at port 8000

- **make client**
  - This command will start the GRPC client and dial on 0.0.0.0 at port 8000
  - The client performs the following in the respective order:
    - Add a Role
    - List Roles
    - Fetch a role
    - Add a User with above fetched role
    - List Users
    - Fetch a User

### Run the application as REST server via GRPC gateway
To run this application use the `Makefile`, you can use the following commands:
- **make clean**
  - This command will clean the existing auto-generated `.pb.go` files from the `pb/` directory.
 
- **make gen**
  - This command will auto-generate the `.pb.go` files in the `pb/` directory.  

- **make rest**
  - This command will start the REST server on 0.0.0.0 at port 8001

- **Use client(Postman/curl) of your choice to run the REST endpoint**

### Run the application as both gRPC server and REST server via GRPC gateway
To run this application use the `Makefile`, you can use the following commands:
- **make clean**
  - This command will clean the existing auto-generated `.pb.go` files from the `pb/` directory.
 
- **make gen**
  - This command will auto-generate the `.pb.go` files in the `pb/` directory.  

- **make both**
  - This command will start the GRPC server on 0.0.0.0 at port 8000 and the REST server on 0.0.0.0 at port 8001

- **Use client(Postman/curl) of your choice to run the REST endpoint and/or you can also run `make client`**
