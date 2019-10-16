# Guide d'utilisation
## Demo
le script DemoStarter.bat lance un maitre et autant d'esclave souhaité avec un NetworkDelay configurable ainsi qu'un ClockOffset aléatoire entre les offset minimum et maximum entré par l'utilisateur

## Master
le maitre se lance avec:<br/>
go run master.go P2PPort MulticastPort<br/>
ou<br/>
go run master.go P2PPort MulticastPort NetworkDelay Verbose<br/>
<br/>
où:<br/>
-P2Pport: le port qui sera utilisé pour les communication en P2P comme le REQUEST_DELAY et DELAY_RESPONSE<br/>
-MulticastPort: le port qui sera utilisé pour les multicast de SYNC et FOLLOWUP<br/>
-NetworkDelay: un délai réseau artificiel pour le test en local<br/>
-Verbose: (true/false) si le programme output toutes les communications entrantes ou sortantes<br/>

## Slave
le slave se lance avec:<br/>
go run slave.go Port MasterP2PIP MasterMultCastIp <br/>
ou <br/>
go run slave.go Port MasterP2PIP MasterMultCastIp ClockOffset NetworkDelay verbose<br/>
<br/>
où:<br/>
-Port: le port utilisé par le slave pour la communication<br/>
-MasterP2PIP: L'addresse du maitre pour la communication P2P<br/>
-MasterMultCastIp: L'addresse du multicast <br/>
-ClockOffset: un décallage artificiel de l'horloge interne de l'esclave<br/>
-NetworkDelay: un délai réseau artificiel pour le test en local<br/>
-Verbose: (true/false) si le programme output toutes les communications entrantes ou sortantes<br/>
