package interactions

import (
 
    "fmt"
	"components"
    "encoding/json"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "net/http"
   	"github.com/julienschmidt/httprouter"
)

type MongodbHandler struct	{
			session *mgo.Session
}

/*
const(
    timeout = time.Duration(time.Second*100)
)

func GetGoogleLocation(address string) (float64, float64)  {

    gmapsDetails := components.GoogleLocation{}
	address = strings.Replace(address," ","%20",-1)
    url := fmt.Sprintf("http://maps.google.com/maps/api/geocode/json?address=%s",address)
	client := http.Client{Timeout: timeout}

    res, err := client.Get(url)
    if err != nil {
        fmt.Errorf("Cannot read Google API: %v", err)
    }
    defer res.Body.Close()
    decoder := json.NewDecoder(res.Body)

    err = decoder.Decode(&gmapsDetails)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the Google JSON: %v", err)
    }
    return gmapsDetails.Results[0].Geometry.Location.Lat, gmapsDetails.Results[0].Geometry.Location.Lng
}
*/

func CreateMongodbHandler (sess *mgo.Session) *MongodbHandler {
	return &MongodbHandler{sess}
}

func (sess MongodbHandler) GetLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	
	locationDetails := components.Location{}
	id := bson.ObjectIdHex(p.ByName("locationID"))
	err := sess.session.DB("test_db").C("test").FindId(id).One(&locationDetails)
  	if err != nil {
    	fmt.Printf("got an error finding a doc %v\n")
    } 	
	if(locationDetails.ID == ""){
		fmt.Fprintf(rw,"Invalid LocationID")
	}else{
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
    	rw.WriteHeader(200)
    	json.NewEncoder(rw).Encode(locationDetails)
	}
}

func (sess MongodbHandler) AddLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	locationDetails := components.Location{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&locationDetails)
	if(err!=nil)    {
		fmt.Errorf("Error in decoding the Input JSON: %v", err)
	}
	address := locationDetails.Address+"%20"+locationDetails.City+"%20"+locationDetails.State+"%20"+locationDetails.Zip
	locationDetails.Coordinate.Lat,locationDetails.Coordinate.Lng = components.GetGoogleLocation(address)
	locationDetails.ID = bson.NewObjectId()
	err = sess.session.DB("test_db").C("test").Insert(locationDetails)
	if err != nil {
		fmt.Printf("Can't insert document: %v\n", err)
	}
	err = sess.session.DB("test_db").C("test").FindId(locationDetails.ID).One(&locationDetails)
	if err != nil {
		fmt.Printf("got an error finding a doc %v\n")
	} 	
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(locationDetails)	
}

func (sess MongodbHandler) UpdateLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	locationDetails := components.Location{}
	tempLocationDetails := components.Location{}
	//searching for existing location details from db
	id := bson.ObjectIdHex(p.ByName("locationID"))
	err := sess.session.DB("test_db").C("test").FindId(id).One(&tempLocationDetails)
  	if err != nil {
    	fmt.Printf("got an error finding a doc %v\n")
    } 
	//decoding the input JSON from request.Body
	decoder := json.NewDecoder(req.Body)
    err = decoder.Decode(&locationDetails)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the Input JSON: %v", err)
    }

	address := locationDetails.Address+" "+locationDetails.City+" "+locationDetails.State+" "+locationDetails.Zip
	locationDetails.Coordinate.Lat,locationDetails.Coordinate.Lng = components.GetGoogleLocation(address)
	locationDetails.ID = id
    if locationDetails.Name == "" {
      locationDetails.Name = tempLocationDetails.Name
    }
	err = sess.session.DB("test_db").C("test").UpdateId(id,locationDetails)
  	if err != nil {
    	fmt.Printf("got an error updating a doc %v\n")
    } 

	err = sess.session.DB("test_db").C("test").FindId(id).One(&locationDetails)
  	if err != nil {
    	fmt.Printf("got an error finding a doc %v\n")
    }
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
    rw.WriteHeader(201)
    json.NewEncoder(rw).Encode(locationDetails)
}

func (sess MongodbHandler) DeleteLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	id := bson.ObjectIdHex(p.ByName("locationID"))
	err := sess.session.DB("test_db").C("test").RemoveId(id)
  	if err != nil {
    	fmt.Printf("got an error deleting a doc %v\n")
    }
	rw.WriteHeader(200)
}

func (sess MongodbHandler) GetOptions(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Print("Inside Options")
	rw.Header().Set("Content-Type", "application/x-www-form-urlencoded;application/json; charset=UTF-8")
	rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")
	rw.Header().Set("Origin", "localhost")
	rw.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, X-HTTP-Method-Override, Content-Type, Accept")
	rw.Header().Set("Access-Control-Allow-Methods","POST,GET,PUT,DELETE")
	rw.WriteHeader(200)
}





