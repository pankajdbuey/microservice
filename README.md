# microservice
This microservice is responsible for CRUD operation of Pet Information.  
- Pet Schema :-  
  type: string. e.g. "dog", "cat"  
  breed: string. e.g. "beagle", "tabby"  
  birthdate: RFC3339 time string e.g. "2019-10-12T07:20:50.52Z")  
  
# Database  
making use of mongo db for storage    
- Pull mongo db image - **docker pull mongo**  
- Run - **docker run -p27017:27017 -dit --name db_server  mongo**  
Note :- please use the command as it , because code is internally making use of **db_server** name, in future this variable will be reading from envirnment  
  
Now our mongo db server is up and running  


# Application  
  
- clone the repo  
- cd to micro  
- execute below command to generate binary in for linux  
**go env -w GOOS=linux  
go env -w CGO_ENABLED=0**  
- run go build  
it will create a binary micro (linux) 
- Give permission to generated binary, **chmod +x micro**   
- keep Dockerfile,  which we can find under docker folder and generated micro binary both in same the folder  
- run  **docker build -t myservice:1.0.**  
  
- Now our container image with name myservice and tag 1.0 is ready  
- Run below command to run our application which will be interacting with the mongo db in order to read/modify data  
- **docker run -p8000:8000 -dit --name db_client --link db_server test:1.0**  
  
Note :- Please use all the docker command as it is, do not change linked db server name  
  
# API endpoints  
- Create: endpoint -> http://<hostmachine ip address>:8000/pet Method -> POST Body -> {  
  "type" : "dog",  
  "breed": "local",  
  "birthdate": "2019-10-12T07:20:50.52Z"  
}  
- ReadAll: endpoint -> http://\<hostmachine ip address>:8000/pet Method -> GET  (get all resources)  
- Read by ID: endpoint -> http://\<hostmachine ip address>:8000/pet/{id} Method -> GET  
  example http://\<hostmachine ip address>:8000/pet/619881251544fb867ac4c01e  
- Edit: endpoint -> http://\<hostmachine ip address>:8000/pet/{id} Method -> PUT Body -> {  
  "type" : "dog",  
  "breed": "local",  
  "birthdate": "2019-10-12T07:20:50.52Z"  
}  
- Delete: endpoint -> http://\<hostmachine ip address>:8000/pet/{id} Method -> DELETE  
  
Note: id used in read by ID, Edit and delete, is pet id which we get from read or create api  
