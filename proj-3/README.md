# golangExpert

To run migrate up, run bellow command inside project folder:
migrate -path=internal/infra/database/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up