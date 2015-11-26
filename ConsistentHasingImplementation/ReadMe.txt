Purpose:
This program helps to understand how consistent hashing works and to implement a simple RESTful key-value cache data store.

Steps:
PART 1 - Server Side Cache Data Store
  You will be implementing a simple RESTful key-value data store with the following features:

PUT http://localhost:3000/keys/{key_id}/{value}
E.g. http://localhost:3000/keys/1/foobar
Response: 200

GET http://localhost:3000/keys/{key_id}
E.g. http://localhost:3000/keys/1
Response: {
                     “key” : 1,
                     “value” : “foobar”
                   }
        
GET http://localhost:3000/keys
E.g. http://localhost:3000/keys
Response: [
          {
                     “key” : 1,
                     “value” : “foobar”
           },
                       {
                                 “key” : 2,
                                “value” : “b”
                       }
            ]
Run three server instances using ports 3000, 3001, and 3002.
A: http://localhost:3000
B: http://localhost:3001
C: http://localhost:3002

Part -2 Consistent Hashing on Client Side
Client side implements a consistent hashing client in GO to support PUT and GET [/keys/{key_id}] operations. 
Consistent hashing algorithm is used to hash server hostnames and keys.
{Key => Value}
1 => a
2 => b
3 => c
4 => d
5 => e
6 => f
7 => g
8 => h
9 => i
10 => j
