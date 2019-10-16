# Guide d'utilisation

## Master
le maitre se lance avec:<br/>
go run master.go P2PPort MulticastPort<br/>
ou<br/>
go run master.go P2PPort MulticastPort NetworkDelay Verbose
<br/>
<br/>
où:<br/>
P2Pport: le port qui sera utilisé pour les communication en P2P comme le REQUEST_DELAY et DELAY_RESPONSE<br/>
MulticastPort: le port qui sera utilisé pour les multicast de SYNC et FOLLOWUP<br/>
NetworkDelay: un délai réseau artificiel pour le test en local<br/>
Verbose: (true/false) si le programme output toutes les communications entrantes ou sortantes<br/>

## Slave
le slave se lance avec:<br/>
go run slave.go Port MasterP2PIP masterMultCastIp <br/>
ou go run slave.go Port MasterP2PIP MasterMultCastIp ArtificialClockOffset ArtificialNetworkDelay verbose
