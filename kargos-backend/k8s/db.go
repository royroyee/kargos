package k8s

import (
	cm "github.com/boanlab/kargos/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

func (kh K8sHandler) DBSession() {

	log.Println("Success to Create DB Session")
	// already created db session (main/initHandlers)
	defer kh.session.Close()

	// Store data in the DB every 3 Hours
	insertTicker := time.NewTicker(3 * time.Hour) // test
	go func() {
		for range insertTicker.C {
			kh.storeNodeInDB()
		}
	}()

	// Delete old data(node) from DB every 25 hours
	deleteTicker := time.NewTicker(25 * time.Hour) // test
	go func() {
		for range deleteTicker.C {
			kh.deleteNodeFromDB()
		}
	}()

	// Delte old data(pod) from DB every 5 minutes
	deleteTicker = time.NewTicker(5 * time.Minute) // test
	go func() {
		for range deleteTicker.C {
			kh.deletePodFromDB()
		}
	}()
	// Wait indefinitely
	select {}
}

// Create MongoDB Session
func GetDBSession() *mgo.Session {
	log.Println("Create DB Session .. ")
	session, err := mgo.Dial("mongodb://db-service:27017") // db-service is name of mongodb service(kubernetes)
	//session, err := mgo.Dial("mongodb://localhost:27017")

	//// Check environment variables for mongodb.
	//mongodbIP := os.Getenv("MONGODB_LISTEN_ADDR")
	//mongodbPort := os.Getenv("MONGODB_LISTEN_PORT")
	//if len(mongodbPort) == 0 || len(mongodbIP) == 0 {
	//	if len(os.Getenv("DB_SERVICE_PORT_27017_TCP_ADDR")) != 0 {
	//		mongodbIP = os.Getenv("DB_SERVICE_PORT_27017_TCP_ADDR")
	//	} else {
	//		log.Fatalf("invalid mongodb address: %s:%s", mongodbIP, mongodbPort)
	//	}
	//}

	if err != nil {
		panic(err)
	}
	return session
}

// Get Node Data In DB

// Store Node Data In DB every 6 hours
func (kh K8sHandler) storeNodeInDB() {
	nodeList, err := kh.GetNodeMetric() // Only metrics of node
	if err != nil {
		return
	}

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()

	collection := cloneSession.DB("kargos").C("node")

	for _, node := range nodeList {
		err = collection.Insert(node)
		if err != nil {
			log.Println(err)
			return
		}
	}
	log.Println("Node Data stored successfully")
}

// Delete all Node data older than 25 hours
func (kh K8sHandler) deleteNodeFromDB() {
	collection := kh.session.DB("kargos").C("node")

	cutoff := time.Now().Add(-25 * time.Hour)
	_, err := collection.RemoveAll(bson.M{"timestamp": bson.M{"$lte": cutoff}})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Old data of nodes deleted successfully")
}

func (kh K8sHandler) GetRecordOfNode(nodeName string) (cm.RecordOfNode, cm.RecordOfNode, cm.RecordOfNode) {
	var hours24, hours12, hours6 cm.RecordOfNode

	collection := kh.session.DB("kargos").C("node")

	// last 24 hours
	filter := bson.M{
		"$and": []bson.M{
			{"name": nodeName},
			{"timestamp": bson.M{"$lte": time.Now().Add(-20 * time.Hour)}},
		},
	}

	err := collection.Find(filter).One(&hours24)
	if err != nil {
		log.Println(err)
		return cm.RecordOfNode{}, cm.RecordOfNode{}, cm.RecordOfNode{}
	}

	// last 12 hours
	filter = bson.M{
		"$and": []bson.M{
			{"name": nodeName},
			{"timestamp": bson.M{"$lte": time.Now().Add(-10 * time.Hour)}},
			{"timestamp": bson.M{"$gte": time.Now().Add(-15 * time.Hour)}},
		},
	}

	err = collection.Find(filter).One(&hours12)
	if err != nil {
		log.Println(err)
		return cm.RecordOfNode{}, cm.RecordOfNode{}, cm.RecordOfNode{}
	}

	// last 6 hours
	filter = bson.M{
		"$and": []bson.M{
			{"name": nodeName},
			{"timestamp": bson.M{"$lte": time.Now().Add(-4 * time.Hour)}},
			{"timestamp": bson.M{"$gte": time.Now().Add(-9 * time.Hour)}},
		},
	}

	err = collection.Find(filter).One(&hours6)
	if err != nil {
		log.Println(err)
		return cm.RecordOfNode{}, cm.RecordOfNode{}, cm.RecordOfNode{}
	}

	return hours24, hours12, hours6
}

// Store Pod Data into DB when kargos agents send container data to gRPC Server (container.go)
// default : 60 second
func (kh K8sHandler) StorePodInDB(podList []cm.Pod) {

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()

	collection := cloneSession.DB("kargos").C("pod")

	for _, pod := range podList {
		err := collection.Insert(pod)
		if err != nil {
			log.Println(err)
			return
		}
	}

	log.Println("Pod Data stored successfully")
}

// Delete all Pod data older than 5 Minutes
func (kh K8sHandler) deletePodFromDB() {
	collection := kh.session.DB("kargos").C("pod")

	cutoff := time.Now().Add(-5 * time.Minute)
	_, err := collection.RemoveAll(bson.M{"timestamp": bson.M{"$lte": cutoff}})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Old data of pods deleted successfully")
}

func (kh K8sHandler) GetRecordOfPod(podName string) (cm.Pod, error) {
	var result = cm.Pod{}

	filter := bson.M{"name": podName}
	collection := kh.session.DB("kargos").C("pod")

	err := collection.Find(filter).One(&result)
	if err != nil {
		log.Println(err)
		return cm.Pod{}, err
	}

	return result, nil

}
