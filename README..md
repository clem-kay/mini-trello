# 📝 Mini Trello

A lightweight Trello-inspired task management app built with [Go](https://golang.org/) and [Fiber](https://gofiber.io/).  
The app lets you create boards, lists, and cards to organize tasks in a simple Kanban-style workflow.

---

## 🚀 Features
- Create, view, update, and delete **boards**
- Each board contains **lists** (e.g., "To Do", "In Progress", "Done")
- Add **cards** (tasks) under lists
- Move cards between lists
- RESTful API design with Fiber
- In-memory storage (no database required)

---

## 🛠 Tech Stack
- **Backend**: Go 1.23+, Fiber v2
- **Frontend**: None (use Postman / curl for testing)
- **Database**: In-memory (can be extended later)

---

## ⚡ Getting Started

### 1. Clone & Run
```bash
git clone https://github.com/your-username/mini-trello.git
cd mini-trello
go mod init github.com/your-username/mini-trello
go get github.com/gofiber/fiber/v2
go run main.go
