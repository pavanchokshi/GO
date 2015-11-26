Description:
The purpose of this program is to enter the name of places and then plan a trip.

Steps:
1. User enters the name of various places, the application fetches its latitude-longitude using GoogleMaps API and store it into MongoDB.
2. User can edit/update/delete place names
3. User enters the list of locations to plan a trip. Application finds the shortest/cheapest route using UBER APIs and gives the best possible route using "low_estimate" parameter of UBER API.
4. User can request a cab from the application and once the destination is reached the subsequent call will request a cab for next location in the route.

The assignment folder structure is as below:

FALL15CMPE273ASSIGNMENT3
  |
  |___components
  |       |__app_config.go : contains the config details like server and authorization token adn URL of UBER
  |       |__google.go : contains the functions for getting the locations using GoogleMaps API
  |       |__models.go : contains all the struct declarations related to JSON used across the program
  |       |__mongo.go : contains the functions for making connection to MongoLab
  |       |__uber.go : contains the functions that calls the UBER API
  |
  |___interactions
  |       |__location.go : handles the GET,POST,PUT,DELETE request related to "Location"
  |       |__trip.go : handles the GET,POST,PUT request related to "Trip"
  |
  |___server.go : function for creating a server that will handle the http requests
