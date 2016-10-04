package core

import (
	"cto-github.cisco.com/appstate/logging"
	"cto-github.cisco.com/appstate/models/composition"
	"cto-github.cisco.com/appstate/models/credentials"
	"cto-github.cisco.com/appstate/openstack-client/constants"
	"cto-github.cisco.com/appstate/openstack-client/global/globalvar"
	"cto-github.cisco.com/appstate/openstack-client/osclient"
	"cto-github.cisco.com/appstate/util"
	"encoding/json"
	"io/ioutil"
)

func LoadComposition(filename string) *composition.Composition {
	logging.Info.Println("Enter LoadComposition():", filename)
	cmp := composition.NewComposition()

	if filename != "" {
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			logging.Error.Fatalln("ERROR: docker-client composition.json can not be read !", err)
		}
		//logging.Info.Println("ReadFile orignal data:", string(file))
		err = json.Unmarshal(file, &cmp)
		if err != nil {
			logging.Error.Fatalln("ERROR:unmarshal composition json", err)

		}

		logging.Info.Println("Retrieve composition from file.")
	} else {
		logging.Error.Fatalln("ERROR: docker-client composition.json can not be found !")
	}
	str, _ := json.Marshal(cmp)
	logging.Info.Println("Convert back from structure: ", string(str))
	if cmp.Components == nil {
		logging.Error.Fatalln("ERROR:read Components are nil ")
	} else {
		logging.Info.Println("Read construct config successfully.")
	}
	return cmp
}

func RealMain(argv map[string]interface{}) int {
	logging.Info.Println("AppState OpenStack-Client: RealMain: start*********")
	globalvar.InitConfig(argv)

	//call different function switch
	optype := argv[constants.OPERATION].(string)
	para := argv[constants.PARAMETER].(string)
	os_client := osclient.NewOSClient(globalvar.CONFIG.OSCfg)
	tenantid := ""
	cps := LoadComposition(argv[constants.COMPOSITIONFILE].(string))
	creds := credentials.NewCredentials()

	switch optype {
	case "deploynova":
		logging.Info.Println(" RealMain: call deploy()")
		//os_client.Deploy(cmp)
	case "crnet":
		logging.Info.Println(" RealMain: call CreateNetwork()")
		os_client.CreateNetwork(para, creds, tenantid)
	case "crport":
		logging.Info.Println(" RealMain: call CreatePort()")
		//os_client.CreatePort(para)
	case "crsub":
		logging.Info.Println(" RealMain: call CreateSubnet()")
		//os_client.CreateSubnet(para)
	case "crsvr":
		logging.Info.Println(" RealMain: call CreateServer()")
		//os_client.CreateServer(para)
	case "lsnet":
		logging.Info.Println(" RealMain: call ListNetworks()")
		os_client.GetNetworkID(para, creds, tenantid)
	case "lsport":
		logging.Info.Println(" RealMain: call ListPorts()")
		os_client.GetPortID(para, creds, tenantid)
	case "lssub":
		logging.Info.Println(" RealMain: call ListSubnets()")
		os_client.GetSubnetID(para, creds, tenantid)
	case "lssvr":
		logging.Info.Println(" RealMain: call ListServers()")
		os_client.GetServers(creds, tenantid)
	case "delnet":
		logging.Info.Println(" RealMain: call DeleteNetwork()")
		os_client.DeleteNetwork("", creds, tenantid)
	case "delport":
		logging.Info.Println(" RealMain: call DeletePort()")
		os_client.DeletePort("", creds, tenantid)
	case "delsub":
		logging.Info.Println(" RealMain: call DeleteSubnet()")
		os_client.DeleteSubnet("", creds, tenantid)
	case "delserver":
		logging.Info.Println(" RealMain: call DeleteServer()")
		os_client.DeleteServer("", creds, tenantid)
	case "deployceilometer":
		logging.Info.Println(" RealMain: call Deployceilometer()")
		os_client.DeployCeilometer(nil)
	case "deploy":
		logging.Info.Println(" RealMain:Default call deploy()")
		//clear the all staff
		for _, each := range cps.Components {
			if each.Type == "OPENSTACK" {
				logging.Info.Println("in undeploy oscfg.osaccount:%v\n", each.OpenStackCfg.OSAccount)
				os_client.Undeploy(each)
			}
		}
		for _, each := range cps.Components {
			if each.Type == "OPENSTACK" {
				logging.Info.Println(util.PointerToString(each))
				if each.OpenStackCfg.OSAccount.UserName == "" || each.OpenStackCfg.OSAccount.Password == "" {
					logging.Error.Fatalln("ERROR: username and/or password was not supplied!")
				}
				if each.OpenStackCfg.TenantID == "" {
					logging.Error.Fatalln("ERROR: tenanid was not supplied!")
				}
				logging.Info.Println("Operation type:**********", optype)
				os_client.Deploy(each)
			} else {
				logging.Error.Fatalln("UNSUPPORTED COMPONENT")
			}
		}
	}
	logging.Info.Println("AppState OpenStack-Client: RealMain: end")
	return 0
}
