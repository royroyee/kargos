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

	// Store data in the DB every 30 second
	insertTicker := time.NewTicker(30 * time.Second) // test
	go func() {
		for range insertTicker.C {
			kh.storeNodeInDB()
			kh.storeControllerInDB()
			kh.storePersistentVolumeInDB()
		}
	}()

	// Delete old data(node) from DB every 25 hours
	//deleteTicker := time.NewTicker(25 * time.Hour) // test
	//go func() {
	//	for range deleteTicker.C {
	//		kh.deleteNodeFromDB()
	//	}
	//}()

	// Delte old data(pod) from DB every 5 minutes
	//deleteTicker = time.NewTicker(5 * time.Minute) // test
	//go func() {
	//	for range deleteTicker.C {
	//		kh.deletePodFromDB()
	//	}
	//}()

	// Delete old data(event) from DB every 24 hours
	deleteTicker := time.NewTicker(24 * time.Hour) // test
	go func() {
		for range deleteTicker.C {
			kh.deleteEventFromDB()
		}
	}()

	// Wait indefinitely
	select {}
}

// Create MongoDB Session
func GetDBSession() *mgo.Session {
	log.Println("Create DB Session .. ")
	session, err := mgo.Dial("mongodb://db-service:27017") // db-service is name of mongodb service(kubernetes)
	// session, err := mgo.Dial("mongodb://localhost:27017")

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
	nodeList, err := kh.GetNodeList() // TODO change node struct
	if err != nil {
		return
	}

	// Delete values that should not be in db before saving node data.
	kh.deleteNodeFromDB(nodeList)

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := cloneSession.DB("kargos").C("node")

	bulk := collection.Bulk()
	for _, node := range nodeList {
		bulk.Upsert(bson.M{"name": node.Name}, node) // duplicate processing : name of node
	}
	_, err = bulk.Run()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Node Data stored successfully")
}

// Info : Events other than the warning cirtical type
func (kh K8sHandler) GetNodeOverview(page int, perPage int) ([]cm.Node, error) {
	var result []cm.Node
	collection := kh.session.DB("kargos").C("node")

	skip := (page - 1) * perPage
	limit := perPage

	err := collection.Find(bson.M{}).Skip(skip).Limit(limit).Sort("name").All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

// Delete all Node data older than 25 hours
func (kh K8sHandler) deleteNodeFromDB(nodeList []cm.Node) {
	//collection := kh.session.DB("kargos").C("node")
	//
	//cutoff := time.Now().Add(-25 * time.Hour)
	//_, err := collection.RemoveAll(bson.M{"timestamp": bson.M{"$lte": cutoff}})
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := kh.session.DB("kargos").C("node")

	nodeNames := make([]string, 0)
	for _, node := range nodeList {
		nodeNames = append(nodeNames, node.Name)
	}

	// Delete the node from the database if it's not in the nodeNames list
	_, err := collection.RemoveAll(bson.M{"name": bson.M{"$nin": nodeNames}})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Node Data deleted successfully")
}

//func (kh K8sHandler) GetRecordOfNode(nodeName string) (cm.RecordOfNode, cm.RecordOfNode, cm.RecordOfNode) {
//	var hours24, hours12, hours6 cm.RecordOfNode
//
//	collection := kh.session.DB("kargos").C("node")
//
//	// last 24 hours
//	filter := bson.M{
//		"$and": []bson.M{
//			{"name": nodeName},
//			{"timestamp": bson.M{"$lte": time.Now().Add(-20 * time.Hour)}},
//		},
//	}
//
//	err := collection.Find(filter).One(&hours24)
//	if err != nil {
//		log.Println(err)
//		return hours24, hours12, hours6
//	}
//
//	// last 12 hours
//	filter = bson.M{
//		"$and": []bson.M{
//			{"name": nodeName},
//			{"timestamp": bson.M{"$lte": time.Now().Add(-10 * time.Hour)}},
//			{"timestamp": bson.M{"$gte": time.Now().Add(-15 * time.Hour)}},
//		},
//	}
//
//	err = collection.Find(filter).One(&hours12)
//	if err != nil {
//		log.Println(err)
//		return hours24, hours12, hours6
//	}
//
//	// last 6 hours
//	filter = bson.M{
//		"$and": []bson.M{
//			{"name": nodeName},
//			{"timestamp": bson.M{"$lte": time.Now().Add(-4 * time.Hour)}},
//			{"timestamp": bson.M{"$gte": time.Now().Add(-9 * time.Hour)}},
//		},
//	}
//
//	err = collection.Find(filter).One(&hours6)
//	if err != nil {
//		log.Println(err)
//		return hours24, hours12, hours6
//	}
//
//	return hours24, hours12, hours6
//}

// Store Pod Data into DB when kargos agents send container data to gRPC Server (container.go)
// default : 60 second
func (kh K8sHandler) StorePodInDB(podList []cm.Pod) {

	// Delete values that should not be in db before saving pod data.
	kh.deletePodFromDB(podList)

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := cloneSession.DB("kargos").C("pod")

	bulk := collection.Bulk()
	for _, pod := range podList {
		bulk.Upsert(bson.M{"name": pod.Name}, pod) // duplicate processing : name of pod
	}
	_, err := bulk.Run()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Pod Data stored successfully")
}

// Delete all Pod data older than 5 Minutes
func (kh K8sHandler) deletePodFromDB(podList []cm.Pod) {
	//collection := kh.session.DB("kargos").C("pod")
	//
	//cutoff := time.Now().Add(-5 * time.Minute)
	//_, err := collection.RemoveAll(bson.M{"timestamp": bson.M{"$lte": cutoff}})
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Println("Old data of pods deleted successfully")

	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := kh.session.DB("kargos").C("pod")

	podNames := make([]string, 0)
	for _, pod := range podList {
		podNames = append(podNames, pod.Name)
	}

	// Delete the pod from the database if it's not in the podNames list
	_, err := collection.RemoveAll(bson.M{"name": bson.M{"$nin": podNames}})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Pod Data deleted successfully")
}

// Info : Events other than the warning cirtical type
func (kh K8sHandler) GetPodOverview(page int, perPage int) ([]cm.PodOverview, error) {
	var result []cm.PodOverview
	collection := kh.session.DB("kargos").C("pod")

	skip := (page - 1) * perPage
	limit := perPage

	err := collection.Find(bson.M{}).Skip(skip).Limit(limit).Sort("namespace").All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
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

func (kh K8sHandler) StoreEvents(event string) {
	cloneSession := kh.session.Clone()

	collection := cloneSession.DB("kargos").C("event")

	err := collection.Insert(event)
	if err != nil {
		log.Println(err)
	}

	log.Println("Event stored successfully")
}

// Only events of Warning , Critical type
func (kh K8sHandler) GetAlerts(page int, perPage int) ([]cm.Event, error) {
	var result []cm.Event
	collection := kh.session.DB("kargos").C("event")

	skip := (page - 1) * perPage
	limit := perPage

	filter := bson.M{
		"$or": []bson.M{
			{"type": "Warning"},
			{"type": "Critical"},
		},
	}

	err := collection.Find(filter).Skip(skip).Limit(limit).Sort("-created").All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

// Info : Events other than the warning cirtical type
func (kh K8sHandler) GetInfo(page int, perPage int) ([]cm.Event, error) {
	var result []cm.Event
	collection := kh.session.DB("kargos").C("event")

	skip := (page - 1) * perPage
	limit := perPage

	filter := bson.M{
		"$nor": []bson.M{
			{"type": "Warning"},
			{"type": "Critical"},
		},
	}

	err := collection.Find(filter).Skip(skip).Limit(limit).Sort("-created").All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) StoreEventInDB(event cm.Event) {

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()

	collection := cloneSession.DB("kargos").C("event")

	err := collection.Insert(event)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Event Data stored successfully")

}

// Delete all event data older than 24 hours
func (kh K8sHandler) deleteEventFromDB() {
	collection := kh.session.DB("kargos").C("event")

	cutoff := time.Now().Add(-24 * time.Minute)
	_, err := collection.RemoveAll(bson.M{"timestamp": bson.M{"$lte": cutoff}})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Old data of events deleted successfully")
}

func (kh K8sHandler) storeControllerInDB() {
	controllerList, err := kh.Controller()
	if err != nil {
		return
	}

	kh.deleteControllerFromDB(controllerList)

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := cloneSession.DB("kargos").C("controller")

	bulk := collection.Bulk()
	for _, controller := range controllerList {
		bulk.Upsert(bson.M{"name": controller.Name, "namespace": controller.Namespace}, controller)
	}

	result, err := bulk.Run()
	if err != nil {
		log.Println(err)
	}

	log.Println("Controller Data stored successfully : ", result)
}

func (kh K8sHandler) deleteControllerFromDB(controllerList []cm.Controller) {

	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := cloneSession.DB("kargos").C("controller")

	// Get the list of controller names from the controllerList
	controllerNames := make([]string, 0)
	for _, controller := range controllerList {
		controllerNames = append(controllerNames, controller.Name)
	}

	// Delete the controller from the database if it's not in the controllerNames list
	_, err := collection.RemoveAll(bson.M{"name": bson.M{"$nin": controllerNames}})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Controller Data deleted successfully")
}

func (kh K8sHandler) GetControllers(page int, perPage int) ([]cm.Controller, error) {
	var result []cm.Controller
	collection := kh.session.DB("kargos").C("controller")

	skip := (page - 1) * perPage
	limit := perPage

	err := collection.Find(bson.M{}).Skip(skip).Limit(limit).Sort("namespace").All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) storePersistentVolumeInDB() {
	pvList, err := kh.PersistentVolume()
	if err != nil {
		return
	}

	kh.deletePersistentVolumeFromDB(pvList)

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := cloneSession.DB("kargos").C("persistentvolume")

	bulk := collection.Bulk()
	for _, pv := range pvList {
		bulk.Upsert(bson.M{"name": pv.Name, "claim": pv.Claim}, pv)
	}

	result, err := bulk.Run()
	if err != nil {
		log.Println(err)
	}

	log.Println("Persistent Volume Data stored successfully : ", result)
}

func (kh K8sHandler) deletePersistentVolumeFromDB(pvList []cm.PersistentVolume) {

	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := cloneSession.DB("kargos").C("controller")

	// Get the list of persistent volume names from the pvList
	pvNames := make([]string, 0)
	for _, pv := range pvList {
		pvNames = append(pvNames, pv.Name)
	}

	// Delete the controller from the database if it's not in the controllerNames list
	_, err := collection.RemoveAll(bson.M{"name": bson.M{"$nin": pvNames}})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Persistent Volume Data deleted successfully")
}

func (kh K8sHandler) GetPersistentVolume(page int, perPage int) ([]cm.PersistentVolume, error) {
	var result []cm.PersistentVolume
	collection := kh.session.DB("kargos").C("persistentvolume")

	skip := (page - 1) * perPage
	limit := perPage

	err := collection.Find(bson.M{}).Skip(skip).Limit(limit).Sort("claim").All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}
