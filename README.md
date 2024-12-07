# magic-survey

**magic-survey** is a powerful :)

## Setup Environment
- [Bootup Documention](./doc/bootup.md)


## Docs
- [Git Strategy](./doc/git-strategy.md)
- [architecture.png](./doc/architecture/architecture.png)


## Patterns used
1. **Service Layer Pattern** Usage: for example In my code, the `userService` is used in multiple handler functions (UserCreate, Verify2FACode, Login). The service layer provides a way to separate business logic from the controller (handlers). Each function in the handler invokes the service layer to perform the actual business logic, such as creating a user, verifying a 2FA code, or logging in a user.
   Benefit: This separation improves maintainability and testability, as the business logic is decoupled from the HTTP layer.
2. **DTO (Data Transfer Object) Pattern** Usage: DTOs like `CreateUserDTO`, `UpdateUserDTO`, `LoginRequest`, and `Verify2FACodeRequest` are used to transfer data between the client and server. They are typically used to encapsulate the data required for operations, often including validation.

3. **Middleware Pattern**
Usage: `The WithAuthMiddleware` function is a classic example of a middleware pattern. Middleware functions are used to process HTTP requests before they reach the actual handler. In this case, the middleware validates the JWT token and checks the validity of the user before processing the request.

4. **Observer Pattern (Signal Handling)**
   Usage: The use of `os.Signal` and the `signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)` function is an implementation of the Observer Pattern. The program listens for specific operating system signals (like `SIGINT`or `SIGTERM`) and reacts to them by shutting down the server gracefully.
5. **Context Propagation**: Usage: The use of `context.WithCancel` and `context.WithTimeout` is a common practice in Go for propagating cancellation signals and deadlines. The context is used to manage the lifecycle of operations, particularly during server shutdown.
6. **Go-Routine Pattern (Concurrency)**: Usage: The server is run in a separate goroutine using go func() { ... }. This allows the main function to remain responsive to the termination signal while the server continues to run concurrently.