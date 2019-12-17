package edurouter

type ReqMsg struct{
	uid				string
	data  		[]byte
	checksum	[]byte
}

type ResMsg struct{
	status		string
	data			[]byte
	checksum  []byte
}
