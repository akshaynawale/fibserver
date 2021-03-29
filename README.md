# fibserver
Fibonacci Server with GUI for learning Docker Compose, Kubernetes and GRPC

This project has three main component 
1. worker - (runs as a GRPC server and responds to requests with calculating a fibonacci number at a giben index)
2. webServer - ( runs a  http server for client interaction, tries to get fibonacci number from redis if not found gets it from worker with GRPC)
3. Redis Server - (stores the fibonacci number calculated at a index and acts as a cache)

# :man_technologist:
