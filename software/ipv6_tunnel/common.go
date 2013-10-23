package main

type ClientIdentifier struct {
	Id string
}

type ClientPrefix struct {
	Prefix string
}

type ConfigIni struct {
	Server struct {
		Localhost string
		Listenport string
	}
	Client struct {
		Remotehost string
		Remoteport string
	}
}
