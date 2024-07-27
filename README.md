Creating a Load Balancer:


Components:
1. Request sent from client.
2. Request redirected to server.
3. List maintaining all the servers.
4. Details of individual server.
5. Logs of all the calls.



SERVER:
1. IP Address of the server
2. Host name
3. HTTP Method
4. CPU Utilisation
5. Count of current Request handled
6. 



Features:
1. Request Distribution
2. Health monitoring
3. Session persistance
4. Scalable
5. Load balancing Algorithms
6. Fault tolerance
7. Logging


Design Patterns Involved or Used:

1. Singleton Pattern: Used to ensure that only one instance of the load balancer is created and shared across the system.

2. Strategy Pattern: Used to encapsulate different load balancing algorithms, allowing flexibility in selecting and switching between different strategies.

3. Observer Pattern: Used to monitor and track the health of servers, notifying the load balancer about any changes in the server states.

4. Proxy Pattern: Used to create proxies for servers, allowing the load balancer to handle requests, perform health checks, and manage session persistence.

5. Decorator Pattern: Used to add additional functionality or features, such as monitoring, logging, or rate limiting, to the load balancer without modifying its core implementation.