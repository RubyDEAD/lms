DEPENDENCIES:
-    RABBITMQ
-    REACT (NPM INSTALL)
-    GO LANG

TO RUN THE SERVICES (RUN THEM SIMULTANEOUSLY):
-    Book-Service
          cd bookService
          go run server.go
-    Patron-Service
          cd patron-service
          go run server.go
-    Borrowing-Service
          cd borrowing-service
          go run server.go
-    Fine-Service
          cd fine_service
          go run main.go

(NOTE: AFTER "go run server.go/main.go" ON EACH SERVICES, IT SHOULD GIVE YOU AN GRAPHQL ENDPOINT)


ONLY ADMIN CAN ADD A BOOK TO TEST THIS FUNCTION USE THIS ACCOUNT:
admin@gmail.com
admin123
     
