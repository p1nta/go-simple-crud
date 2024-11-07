# Go Simple CRUD Application with Gin Framework

This is a simple CRUD (Create, Read, Update, Delete) application built in Go using the [Gin Web Framework](https://github.com/gin-gonic/gin). The application manages a list of `todos`, allowing users to create, read, update, and delete `todo` items.

## Features

- **Create a Todo**: Add a new `todo` item with a ID, description, and completion status.
- **Read Todos**: Retrieve the list of all `todos` or get details for a specific `todo`.
- **Update a Todo**: Edit an existing `todo` by updating its description and/or completion status.
- **Delete a Todo**: Remove a `todo` item from the list.

## Endpoints

- `GET /todos` - Retrieves the full list of todos.
- `POST /todos` - Creates a new todo item.
- `GET /todos/:id` - Retrieves details of a specific todo by ID.
- `PATCH /todos/:id` - Toggles the completion status of a todo.
- `PUT /todos/:id` - Updates the description and/or completion status of a specific todo.
- `DELETE /todos/:id` - Deletes a specific todo by ID.

## Quick Start

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.18 or higher)
- [Gin Web Framework](https://github.com/gin-gonic/gin)

### Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/your-username/go-simple-crud.git
   cd go-simple-crud
   ```
