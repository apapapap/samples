# User Management 

## Basic user management using GoLang protocol buffers.

### Pre-requisities
- Go 1.17 or higher
- Protoc 3.6.1 or higher

To run this application use the `Makefile`, you can use the following commands:
- make clean
  - This command will clean the existing auto-generated `.pb.go` files from the `pb/` directory.
 
- make gen
  - This command will auto-generate the `.pb.go` files in the `pb/` directory.  

- make server
  - This command will start the GRPC server on 0.0.0.0 at port 8000

- make client
  - This command will start the GRPC client and dial on 0.0.0.0 at port 8000
  - The client performs the following in the respective order:
    - Add a User
    - List Users
    - Fetch a User
