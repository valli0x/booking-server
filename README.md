## Booking Server

This repository contains the code for a booking server application, built following the principles of clean architecture.

### Layers of the Application

* Entities: The entities layer represents the business objects in our application that encapsulate the business logic. In this case, we have the Order structure located in the models folder.

* Use Cases: The use cases layer defines how entities interact to perform specific business operations. In this case, we have the Storage interface that defines operations with orders, and the Cache interface that defines operations with the cache.

* Interface Adapters: The interface adapters layer adapts the data from a format convenient for internal use to a format that can be used for external representation or interaction with external systems. In this case, we have HTTP request handler functions that convert data from HTTP requests to Order structures and vice versa. We also have the InMemoryStorage structure that implements the Storage interface.

* Frameworks and Drivers: The frameworks and drivers layer consists of external systems and tools that we use to build our application. In this case, we use the web server and router provided by the go-chi/chi library, as well as the main function that starts our application.