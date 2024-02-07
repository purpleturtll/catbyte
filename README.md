docker-compose up -d

To run services run these in separate terminal windows:
go run .\cmd\api\main.go
go run .\cmd\message_processor\main.go
go run .\cmd\reporter\main.go

Example requests in tests.http file, can be run using "REST Client" VSCode extension or used as reference.
