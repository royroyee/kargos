package k8s

import (
	cm "github.com/boanlab/kargos/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"time"
)

func (kh K8sHandler) DBSession() {

	log.Println("Success to Create DB Session")
	// already created db session (main/initHandlers)
	defer kh.session.Close()

	// To create indexes
	kh.CreateNodeIndexes()
	kh.CreatePodIndexes()

	// Store data in the DB every 1 hour
	insertTicker := time.NewTicker(1 * time.Minute)
	go func() {
		for range insertTicker.C {
			kh.storeNodeInDB()
			kh.StorePodInfoInDB()
			kh.storeControllerInDB()
		}
	}()

	// Delete old data from DB every 25 hours
	deleteTicker := time.NewTicker(25 * time.Minute) // test
	go func() {
		for range deleteTicker.C {
			kh.deleteNodeFromDB()
			kh.deletePodFromDB()
			//		kh.deleteEventFromDB()
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

// Store Node Data In DB every 1 hour
func (kh K8sHandler) storeNodeInDB() {
	nodeList, err := kh.GetNodeList()
	if err != nil {
		return
	}

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := cloneSession.DB("kargos").C("node")

	bulk := collection.Bulk()
	for _, node := range nodeList {
		bulk.Insert(node)
	}
	_, err = bulk.Run()
	if err != nil {
		log.Println(err)
		return
	}
}

//Improved performance with indexing processing

func (kh K8sHandler) CreateNodeIndexes() error {
	// Get a reference to the "node" collection
	collection := kh.session.DB("kargos").C("node")

	// Create an index on the "timestamp" field
	index := mgo.Index{
		Key:        []string{"-timestamp"},
		Background: true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		log.Printf("error creating index: %s", err)
		return err
	}

	// Create any other indexes you need here

	return nil
}

func (kh K8sHandler) CreatePodIndexes() error {
	// Get a reference to the "node" collection
	collection := kh.session.DB("kargos").C("podusage")

	// Create an index on the "timestamp" field
	index := mgo.Index{
		Key:        []string{"-timestamp"},
		Background: true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		log.Printf("error creating index: %s", err)
		return err
	}

	// Create any other indexes you need here

	return nil
}

func (kh K8sHandler) GetNodeOverview(page int, perPage int) ([]cm.NodeOverview, error) {
	var result []cm.NodeOverview

	// Get a reference to the "node" collection
	collection := kh.session.DB("kargos").C("node")

	//// Define the query and projection
	//query := bson.M{}
	////	projection := bson.M{"_id": 0, "timestamp": 1}
	//
	//// Sort by "timestamp" in descending order
	//sort := "-timestamp"
	//
	//// Calculate the skip and limit values based on the requested page and items per page
	//skip := (page - 1) * perPage
	//limit := perPage

	// Define the pipeline stages
	pipeline := []bson.M{
		{"$sort": bson.M{"timestamp": -1}},
		{"$group": bson.M{
			"_id":           "$name",
			"name":          bson.M{"$first": "$name"},
			"cpuusage":      bson.M{"$first": "$cpuusage"},
			"ramusage":      bson.M{"$first": "$ramusage"},
			"diskallocated": bson.M{"$first": "$diskallocated"},
			"networkusage":  bson.M{"$first": "$networkusage"},
			"ip":            bson.M{"$first": "$ip"},
			"status":        bson.M{"$first": "$status"},
		}},
		{"$skip": (page - 1) * perPage},
		{"$limit": perPage},
	}

	// Execute the query and get the results
	//err := collection.Find(query).Sort(sort).Skip(skip).Limit(limit).All(&result)
	err := collection.Pipe(pipeline).All(&result)
	if err != nil {
		log.Printf("error querying database: %s", err)
		return result, err
	}

	return result, nil
}

func (kh K8sHandler) GetNodeUsageAvg() (cm.NodeUsage, error) {
	var result cm.NodeUsage

	collection := kh.session.DB("kargos").C("node")

	// Aggregate the average value of cpuusage and ramusage per minute
	pipeline := collection.Pipe([]bson.M{
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"minute": bson.M{"$minute": bson.M{"$toDate": "$timestamp"}},
				},
				"avgCpuUsage": bson.M{"$avg": "$cpuusage"},
				"avgRamUsage": bson.M{"$avg": "$ramusage"},
			},
		},
		{"$limit": 24},
	})

	// Extract the result
	var getUsage []bson.M
	err := pipeline.All(&getUsage)
	if err != nil {
		log.Println(err)
		return result, err
	}

	for _, usage := range getUsage {
		avgCpuUsage := int(usage["avgCpuUsage"].(float64))
		result.CpuUsage = append(result.CpuUsage, avgCpuUsage)

		avgRamUsage := int(usage["avgRamUsage"].(float64))
		result.RamUsage = append(result.RamUsage, avgRamUsage)
	}

	return result, nil
}

func (kh K8sHandler) GetNodeUsage(nodeName string) (cm.NodeUsage, error) {
	var result cm.NodeUsage

	collection := kh.session.DB("kargos").C("node")

	pipeline := collection.Pipe([]bson.M{
		{"$match": bson.M{"name": nodeName}},
		{"$limit": 24},
		{"$project": bson.M{
			"_id":          nil,
			"cpuusage":     1,
			"ramusage":     1,
			"networkusage": 1,
			// TODO Disk Usage
		}},
	})

	// Extract the result
	var getUsage []bson.M
	err := pipeline.All(&getUsage)
	if err != nil {
		log.Println(err)
		return result, err
	}
	for _, usage := range getUsage {
		CpuUsage := int(usage["cpuusage"].(float64))
		result.CpuUsage = append(result.CpuUsage, CpuUsage)

		RamUsage := int(usage["ramusage"].(float64))
		result.RamUsage = append(result.RamUsage, RamUsage)

		NetworkUsage := int(usage["networkusage"].(float64))
		result.NetworkUsage = append(result.NetworkUsage, NetworkUsage)
	}

	return result, nil
}

func (kh K8sHandler) GetPodUsageDetail(podName string) (cm.GetPodUsage, error) {
	var result cm.GetPodUsage
	collection := kh.session.DB("kargos").C("podusage")

	pipeline := collection.Pipe([]bson.M{
		{"$match": bson.M{"name": podName}},
		{"$limit": 24},
		{"$project": bson.M{
			"_id":      nil,
			"cpuusage": 1,
			"ramusage": 1,
			// "networkusage": 1,
			// TODO Disk Usage
		}},
	})

	// Extract the result
	var getUsage []bson.M
	err := pipeline.All(&getUsage)
	if err != nil {
		log.Println(err)
		return result, err
	}
	for _, usage := range getUsage {
		CpuUsage := int(usage["cpuusage"].(int64))
		result.CpuUsage = append(result.CpuUsage, CpuUsage)

		RamUsage := int(usage["ramusage"].(int64))
		result.RamUsage = append(result.RamUsage, RamUsage)

		//NetworkUsage := int(usage["networkusage"].(float64))
		//result.NetworkUsage = append(result.NetworkUsage, NetworkUsage)
	}

	return result, nil
}

func (kh K8sHandler) GetTopNode() (cm.TopNode, error) {
	var result cm.TopNode
	var name string
	var usage int

	collection := kh.session.DB("kargos").C("node")

	now := time.Now()
	// Define the cutoff time as 1 minute ago
	cutoffTime := now.Add(-1 * time.Minute).Format("2006-01-02 15:04")

	// Find the top 3 nodes with highest cpuusage and ramusage among the most recent data
	pipe := collection.Pipe([]bson.M{

		bson.M{"$match": bson.M{"timestamp": bson.M{"$gte": cutoffTime, "$lt": now.Format("2006-01-02 15:04")}}},
		bson.M{"$sort": bson.M{"cpuusage": -1}},
		bson.M{"$limit": 3},
	})

	var topCpu []cm.NodeCpuUsage
	err := pipe.All(&topCpu)
	if err != nil {
		return result, err
	}

	for _, node := range topCpu {
		name = node.Name
		usage = node.CpuUsage
		result.Cpu = append(result.Cpu, cm.NodeCpuUsage{Name: name, CpuUsage: usage})
	}

	// Find the top 3 nodes with highest ramusage among the most recent data
	pipe = collection.Pipe([]bson.M{

		bson.M{"$match": bson.M{"timestamp": bson.M{"$gte": cutoffTime, "$lt": now.Format("2006-01-02 15:04")}}},
		bson.M{"$sort": bson.M{"ramusage": -1}},
		bson.M{"$limit": 3},
	})

	var topRam []cm.NodeRamUsage
	err = pipe.All(&topRam)
	if err != nil {
		return result, err
	}

	for _, node := range topRam {
		name = node.Name
		usage = node.RamUsage
		result.Ram = append(result.Ram, cm.NodeRamUsage{Name: name, RamUsage: usage})
	}

	return result, nil
}

func (kh K8sHandler) GetTopPod() (cm.TopPod, error) {
	var result cm.TopPod
	var name string
	var usage int

	collection := kh.session.DB("kargos").C("podusage")

	// Define the cutoff time as 1 minute ago
	now := time.Now()
	cutoffTime := now.Add(-1 * time.Minute).Format("2006-01-02 15:04")

	// Find the top 3 pods with highest cpuusage and ramusage among the most recent data
	pipe := collection.Pipe([]bson.M{

		//bson.M{"$match": bson.M{"timestamp": bson.M{"$gte": cutoffTime, "$lte": now.Format("2006-01-02 15:04")}}},
		bson.M{"$match": bson.M{"timestamp": bson.M{"$gte": cutoffTime, "$lt": now.Format("2006-01-02 15:04")}}},

		bson.M{"$sort": bson.M{"cpuusage": -1}},
		bson.M{"$limit": 3},
	})

	var topCpu []cm.PodCpuUsage
	err := pipe.All(&topCpu)
	if err != nil {
		return result, err
	}

	for _, pod := range topCpu {
		name = pod.Name
		usage = pod.CpuUsage
		result.Cpu = append(result.Cpu, cm.PodCpuUsage{Name: name, CpuUsage: usage})
	}

	// Find the top 3 pods with highest ramusage among the most recent data
	pipe = collection.Pipe([]bson.M{

		bson.M{"$match": bson.M{"timestamp": bson.M{"$gte": cutoffTime, "$lt": now.Format("2006-01-02 15:04")}}},
		////	Group by name and take the first 3 groups
		//bson.M{"$group": bson.M{
		//	"_id":   "$name",
		//	"ram":   bson.M{"$first": "$ramusage"},
		//	"count": bson.M{"$sum": 1},
		//}}, -> Error (tried to block duplicate processing(name) as a group..)

		bson.M{"$sort": bson.M{"ramusage": -1}},
		bson.M{"$limit": 3},
	})

	var topRam []cm.PodRamUsage
	err = pipe.All(&topRam)
	if err != nil {
		return result, err
	}

	for _, pod := range topRam {
		name = pod.Name
		usage = pod.RamUsage
		result.Ram = append(result.Ram, cm.PodRamUsage{Name: name, RamUsage: usage})
	}
	return result, nil
}

// Delete all Node data older than 25 hours & deleted node
func (kh K8sHandler) deleteNodeFromDB() {

	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := kh.session.DB("kargos").C("node")

	// Delete the node data older than 25 hours
	cutoff := time.Now().Add(-25 * time.Hour)
	_, err := collection.RemoveAll(bson.M{"timestamp": bson.M{"$lte": cutoff}})
	if err != nil {
		log.Println(err)
		return
	}

}

// Store Pod Data into DB when kargos agents send container data to gRPC Server (container.go)
// default : 60 second
func (kh K8sHandler) StorePodUsageInDB(podList []cm.PodUsage) {

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := cloneSession.DB("kargos").C("podusage")

	bulk := collection.Bulk()
	for _, pod := range podList {
		//bulk.Upsert(bson.M{"name": pod.Name}, pod) // duplicate processing : name of pod
		bulk.Insert(pod)
	}
	_, err := bulk.Run()
	if err != nil {
		log.Println(err)
		return
	}
}

func (kh K8sHandler) StorePodInfoInDB() {

	podList, err := kh.GetPodInfoList()
	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := cloneSession.DB("kargos").C("podinfo")

	bulk := collection.Bulk()
	for _, pod := range podList {
		bulk.Upsert(bson.M{"name": pod.Name}, pod) // duplicate processing : name of pod
	}
	_, err = bulk.Run()
	if err != nil {
		log.Println(err)
		return
	}

	////Test (TODO DELETE)
	pods, err := kh.GetPodUsage()
	kh.StorePodUsageInDB(pods)
	// TEST
}

// Delete all Pod data older than 25 hours
func (kh K8sHandler) deletePodFromDB() {
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

	collection := kh.session.DB("kargos").C("podusage")

	// Delete the node data older than 25 hours
	cutoff := time.Now().Add(-25 * time.Hour)
	_, err := collection.RemoveAll(bson.M{"timestamp": bson.M{"$lte": cutoff}})
	if err != nil {
		log.Println(err)
		return
	}
}

func (kh K8sHandler) GetPodsOfController(controller string) (cm.PodsOfController, error) {
	var result cm.PodsOfController
	collection := kh.session.DB("kargos").C("controller")

	err := collection.Find(bson.M{"name": controller}).One(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) GetInfoOfPod(podName string) (cm.PodInfo, error) {
	var result = cm.PodInfo{}

	filter := bson.M{"name": podName}
	collection := kh.session.DB("kargos").C("podinfo")

	err := collection.Find(filter).One(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) GetPodUsageFromDB(podName string) (cm.PodUsage, error) {
	var result = cm.PodUsage{}

	filter := bson.M{"name": podName}
	collection := kh.session.DB("kargos").C("pod")

	err := collection.Find(filter).One(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}

	return result, nil

}

func (kh K8sHandler) GetEvents(eventLevel string, page int, perPage int) ([]cm.Event, error) {
	var result []cm.Event
	collection := kh.session.DB("kargos").C("event")

	skip := (page - 1) * perPage
	limit := perPage
	filter := bson.M{"eventlevel": strings.Title(eventLevel)}
	if eventLevel == "" {
		filter = bson.M{}
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
}

func (kh K8sHandler) storeControllerInDB() {
	controllerList := kh.GetController()

	kh.deleteControllerFromDB(controllerList)

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()
	defer cloneSession.Close()

	collection := cloneSession.DB("kargos").C("controller")

	bulk := collection.Bulk()
	for _, controller := range controllerList {
		bulk.Upsert(bson.M{"name": controller.Name, "namespace": controller.Namespace}, controller)
	}

	_, err := bulk.Run()
	if err != nil {
		log.Println(err)
	}

	//log.Println("Controller Data stored successfully : ", result)
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

}

func (kh K8sHandler) GetControllers(page int, perPage int) ([]cm.Controller, error) {
	var result []cm.Controller
	collection := kh.session.DB("kargos").C("controller")

	skip := (page - 1) * perPage
	limit := perPage

	err := collection.Find(bson.M{}).Skip(skip).Limit(limit).All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) GetControllersByFilter(namespace string, controller string, page int, perPage int) ([]cm.ControllerOverview, error) {
	var result []cm.ControllerOverview
	collection := kh.session.DB("kargos").C("controller")

	skip := (page - 1) * perPage
	limit := perPage
	var filter bson.M
	if namespace != "" && controller != "" {
		filter = bson.M{
			"namespace": namespace,
			"type":      controller,
		}
	} else if namespace != "" {
		filter = bson.M{
			"namespace": namespace,
		}
	} else if controller != "" {
		filter = bson.M{
			"type": controller,
		}
	}
	err := collection.Find(filter).Skip(skip).Limit(limit).All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) GetControllersByType(controller string, page int, perPage int) ([]cm.Controller, error) {
	var result []cm.Controller
	collection := kh.session.DB("kargos").C("controller")

	skip := (page - 1) * perPage
	limit := perPage

	filter := bson.M{"type": controller}

	err := collection.Find(filter).Skip(skip).Limit(limit).All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) NumberOfEvents(eventLevel string) (cm.Count, error) {
	var result cm.Count
	filter := bson.M{}
	collection := kh.session.DB("kargos").C("event")

	if eventLevel != "" {
		filter = bson.M{"eventlevel": strings.Title(eventLevel)}
	}
	count, err := collection.Find(filter).Count()
	if err != nil {
		log.Println(err)
		return result, err
	}
	result.Count = count
	return result, nil
}

func (kh K8sHandler) NumberOfControllers(namespace string, controllerType string) (cm.Count, error) {
	var result cm.Count
	filter := bson.M{}
	collection := kh.session.DB("kargos").C("controller")

	if namespace != "" && controllerType == "" {
		filter = bson.M{"namespace": namespace}
	} else if namespace != "" && controllerType != "" {
		filter = bson.M{"namespace": namespace, "type": controllerType}
	}
	count, err := collection.Find(filter).Count()
	if err != nil {
		log.Println(err)
		return result, err
	}
	result.Count = count
	return result, nil
}

func (kh K8sHandler) GetEventsByController(controllerName string) ([]cm.Event, error) {
	var result []cm.Event
	collection := kh.session.DB("kargos").C("event")

	limit := 10
	filter := bson.M{"name": controllerName}

	err := collection.Find(filter).Limit(limit).Sort("-created").All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) GetContainersOfPod(podName string) (cm.Containers, error) {
	var result cm.Containers

	collection := kh.session.DB("kargos").C("podusage")
	filter := bson.M{"name": podName}

	// quickly find the most recent data without having to sort the entire collection.
	sort := []string{"-timestamp"}
	limit := 1
	err := collection.Find(filter).Sort(sort...).Limit(limit).One(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}
