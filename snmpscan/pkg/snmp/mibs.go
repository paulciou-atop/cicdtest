package snmp

const TString = 1
const TNumber = 2
const TMac = 3

const PrivateEnterprises = ".1.3.6.1.4.1"
const SystemObjectID = ".1.3.6.1.2.1.1.2.0"

var BasicOIDs = []string{
	".1.3.6.1.2.1.1.1.0", //sysDescr
	".1.3.6.1.2.1.1.2.0", //sysObjectID
	".1.3.6.1.2.1.1.3.0", //sysUpTime
	".1.3.6.1.2.1.1.4.0", //sysContact
	".1.3.6.1.2.1.1.5.0", //sysName

}

type mibDefined struct {
	Oid         string
	Name        string
	Description string
}
