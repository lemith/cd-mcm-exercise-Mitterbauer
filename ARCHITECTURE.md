Client
  |
  |  HTTP request
  |  e.g. GET /products/1
  v
Go HTTP server
  |
  v
gorilla/mux Router
  |
  |  Matches URL + HTTP method
  |  Example:
  |  GET /products/{id:[0-9]+}
  v
Handler function
  |
  |  Example:
  |  h.GetProduct(w, r)
  |
  |  - reads URL params with mux.Vars(r)
  |  - reads JSON body with json.NewDecoder(r.Body)
  |  - validates product input where needed
  v
Store layer
  |
  |  Example:
  |  h.Store.GetByID(id)
  |  h.Store.Create(p)
  |  h.Store.Update(id, p)
  |  h.Store.Delete(id)
  v
Database / storage
  |
  |  MemoryStore: in-memory Go data structure
  |  PostgresStore: PostgreSQL database
  v
Store result
  |
  |  product, product list, or error
  v
Handler response
  |
  |  respondJSON(...)
  |  respondError(...)
  v
HTTP JSON response
  |
  v
Client

## MemoryStore vs. PostgresStore

Use MemoryStore when: 
    - Running unit or integration tests with no external dependencies
    - Running a local demo or prototype where restart-safe state is not required
    - Executing in CI pipelines where provisioning a databse adds latency

Use PostgresStore when: 
    - Data must survive a process restart or deployment
    - Multiple instances of the service run behind a load balancer
    - You need relational queries, joins, transactions or audit logging
    - You are operating in a production or staging enviroment

Key difference
    - MemoryStore = temporary storage in application memory
    - PostgresStore = persistent storage in an external database