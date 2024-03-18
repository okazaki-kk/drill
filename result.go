package main

type ResultCode int

const (
	NoError ResultCode = iota
	Formerr
	Servfail
	NxDomain
	NotImp
	Refused
)
