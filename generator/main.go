package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
)

func main() {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	ldapEntries := []string{}
	ldapEntries = append(ldapEntries, `
debug = true
#################
# Server configuration.
[ldap]
	enabled = true
# run on a non privileged port
	listen = "0.0.0.0:3893"

[ldaps]
# to enable ldaps genrerate a certificate, eg. with:
	enabled = false
	listen = "0.0.0.0:3894"
	cert = "glauth.crt"
	key = "glauth.key"

#################
# The backend section controls the data store.
[backend]
	datastore = "config"
	baseDN = "dc=glauth,dc=com"
	nameformat = "cn"
	groupformat = "ou"

[behaviors]
	# Ignore all capabilities restrictions, for instance allowing every user to perform a search
	IgnoreCapabilities = false
	# Enable a "fail2ban" type backoff mechanism temporarily banning repeated failed login attempts
	LimitFailedBinds = true
	# How many failed login attempts are allowed before a ban is imposed
	NumberOfFailedBinds = 3
	# How long (in seconds) is the window for failed login attempts
	PeriodOfFailedBinds = 10
	# How long (in seconds) is the ban duration
	BlockFailedBindsFor = 60
	# Clean learnt IP addresses every N seconds
	PruneSourceTableEvery = 600
	# Clean learnt IP addresses not seen in N seconds
	PruneSourcesOlderThan = 600

[[users]]
	name = "serviceuser"
	mail = "serviceuser@example.com"
	uidnumber = 5003
	primarygroup = 5502
	passsha256 = "652c7dc687d98c9889304ed2e408c74b611e86a40caa51c4b43f1dd5913c5cd0" # mysecret
	[[users.capabilities]]
	action = "search"
	object = "*"


#################
# The groups section contains a hardcoded list of valid users.
[[groups]]
	name = "superheros"
	gidnumber = 5501

[[groups]]
	name = "svcaccts"
	gidnumber = 5502

[[groups]]
	name = "vpn"
	gidnumber = 5503
	includegroups = [ 5501 ]

#################
# Enable and configure the optional REST API here.
[api]
	enabled = true
	internals = true # debug application performance
	tls = false # enable TLS for production!!
	listen = "0.0.0.0:5555"
	cert = "cert.pem"
	key = "key.pem"
		`)

	for i := 1; i < 20000; i++ {
		userName := fmt.Sprintf("%s %v", nameGenerator.Generate(), i)
		ldapEntries = append(ldapEntries, fmt.Sprintf(ldapUser, userName, userName, userName, userName))
	}

	for i := 1; i < 20000; i++ {
		groupName := fmt.Sprintf("%s %v", nameGenerator.Generate(), i)
		ldapEntries = append(ldapEntries, fmt.Sprintf(ldapGroup, groupName))
	}

	d1 := []byte(strings.Join(ldapEntries, "\n"))
	err := os.WriteFile("/Users/ddevi/EA/github.com/dayadev/openldap-mock/sample.cfg", d1, 0644)
	if err != nil {
		println(err.Error())
	}
}

const (
	ldapUser = `[[users]]
	name = "%s"
	givenname="%s"
	sn="%s"
	mail = "%s@example.com"
	uidnumber = 5002
	primarygroup = 5501
	loginShell = "/bin/sh"
	homeDir = "/root"
	passsha256 = "6478579e37aff45f013e14eeb30b3cc56c72ccdc310123bcdf53e0333e3f416a" # dogood
	sshkeys = ["ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA3UKCEllO2IZXgqNygiVb+dDLJJwVw3AJwV34t2jzR+/tUNVeJ9XddKpYQektNHsFmY93lJw5QDSbeH/mAC4KPoUM47EriINKEelRbyG4hC/ko/e2JWqEclPS9LP7GtqGmscXXo4JFkqnKw4TIRD52XI9n1syYM9Y8rJ88fjC/Lpn+01AB0paLVIfppJU35t0Ho9doHAEfEvcQA6tcm7FLJUvklAxc8WUbdziczbRV40KzDroIkXAZRjX7vXXhh/p7XBYnA0GO8oTa2VY4dTQSeDAUJSUxbzevbL0ll9Gi1uYaTDQyE5gbn2NfJSqq0OYA+3eyGtIVjFYZgi+txSuhw== rsa-key-20160209"]
	passappsha256 = [
	  "c32255dbf6fd6b64883ec8801f793bccfa2a860f2b1ae1315cd95cdac1338efa", # TestAppPw1
	  "c9853d5f2599e90497e9f8cc671bd2022b0fb5d1bd7cfff92f079e8f8f02b8d3", # TestAppPw2
	  "4939efa7c87095dacb5e7e8b8cfb3a660fa1f5edcc9108f6d7ec20ea4d6b3a88", # TestAppPw3
	]`

	ldapGroup = `[[groups]]
	name = "%s"`
)
