# Library Management System
## BE WARY OF PORTS BEING USED
### API-GATEWAY
#### Port Number used: 8081

### Book Service
#### Port Number used: 8080

##### query notes:
  getBooks:[Book!]! Returns  List of Book
  
  getAuthor:[Author!]! Returns List of Author
  
  getBookCopies:[Book_copies!]!  Returns all list of Book_copies
 
  getFilteredBooks(filter: BookFilter): [Book!]! # Filter books based on criteria

  getBookById(id: String!): Book  Returns a Book based on the inputed id(uuid)
  getBookCopiesById(id: String!):[Book_copies!]! Returns a List of Book_copies based on the inputted id(book_id(uuid) ni siya kapoy ilis)

  searchBooks(query: String!): [Book!]! # Search books by title or author

  getAvailbleBookCopyByID(id:String!):Book_copies! Returns the first Available Book_copy based on the inputted id(book_id(uuid) ni siya)
 ##### mutation notes:
    addBook(title: String!, author_name: String!, datePublished: String!,description:String!): Book!     Adds a book. "!" means required na naay input

  addBCopy(book_id:String!): Book_copies!  adds Book_copy need book_id input(uuid)
  addAuthor(author_name:String!): Author!  adds Author 

  updateBook(id: ID!, title: String, author_name:String!, datePublished: String, description:String): Book! updates a Book

  updateBookCopyStatus(id: ID!, book_status:String):Book_copies! updates the status of a BookCopy 
  dapat ang input sa book_status kay naa ra ani nila : Available , Not Available, In Use , Borrwed. 

  deleteBook(id: ID!): Boolean!
  deleteBCopy(id: ID!): Boolean!
  deleteAuthor(id: ID!): Boolean!

### Fine Service
#### Port Number used: 

### Borrowing Service
#### Port Number used: 

### Patron Service
#### Port Number used: 8069
To do:
1. Usernames and emails (if need be)

<br>
Notes for Front-End:
Inig himo sa patron-creation page, applyi og regex ang phone number na field <br>
Mao ni ako gi gamit : '^[0-9]{10,15}$'



