**DEPENDENCIES:**
-    RABBITMQ
-    REACT (NPM INSTALL)
-    GO LANG

**TO RUN THE SERVICES (RUN THEM SIMULTANEOUSLY):**
-    Book-Service:
          cd bookService then
          go run server.go
-    Patron-Service:
          cd patron-service then
          go run server.go
-    Borrowing-Service:
          cd borrowing-service then
          go run server.go
-    Fine-Service:
          cd fine_service then
          go run main.go
-    API-Gateway:
          cd API-gateway then
          go run server.go

(NOTE: AFTER "go run server.go/main.go" ON EACH SERVICES, IT SHOULD GIVE YOU A GRAPHQL ENDPOINT)


**TO RUN THE WEBPAGE (RUN THE SERVICES FIRST BEFORE THE WEBPAGE):
          cd lms then npm run start**

**ONLY ADMIN CAN ADD A BOOK TO TEST THIS FUNCTION USE THIS ACCOUNT:
**EMAIL: admin@gmail.com
PASSWORD: admin123****
     
