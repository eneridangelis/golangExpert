# Clean Architecture

### How to run this project:

1. Run `docker-compose up` to start mysql and rabbitmq containers.
2. Run `migrate -path=internal/infra/database/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up` to run the table up migration.
3. Run `go run main.go wire_gen.go` to run the project.

### Services will be at follow ports:

1. GRPC: 50051
2. GraphQL: 8080
3. WebServer: 8080
4. MySQL: 3306
