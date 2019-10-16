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


## Implementation

# Network

Le package network contient toutes les fonctions relatives a la communication udp entre le maitre et les esclaves.
La fonction ClientWriter permet d'envoyer un datagramme vers l'adresse fournie en parametre.
Les fonctions ClientReader permettent de recevoir les datagrammes envoyes par les differentes applications
celles-ci ouvrent une connection et ecoutent soit les messages transmit en multicast ou en unicast.
Ces fonctions utilisent la fonction decrypt qui recoit le PacketConn ainsi qu'un channel et decrypte ensuite le datagramme encode grace au package gop.
Une fois le datagramme decode celui ci est stocke dans une structure message et est transmis en dehors de la routinne grace a un canal.

# Master

Le package master contient l'application maitre qui envoie periodiquement les messages de synchronisation en multicast.
Deux GoRoutines sont lancées dans le main,la premiere ecoute sur le port UDP et traite les messages de delay envoyés par les differents esclaves et la deuxieme s'occupe d'envoyer periodiquement les messages de synchronisation et de FOLLOWUP aux differents esclaves.

# slave

Le package slave contient l'application esclaves qui se synchronise sur l'horloge du maitre. La fonction main contient trois GoRoutines, une pour ecouter les messages diffusés en Multicast, une pour les communications point a point et une pour gerer l'envoie des delay request au maitre. L'utilisation des goroutines permet de gerer simultanement la reception et l'envoie de messages de maniere independante.

## Protocole de communications

Le Protocole de communication passe par l'envoie de datagrammes en UDP.Les messages de synchronisation sont envoyés depuis le maitre en multicast et sont stockes dans des structures Messages contenant un identifiant, le temps de l'envoie et le type du message stocké en binaires.Les communications point a point contiennent egalement l'adresse de l'expediteur pour savoir a quel esclave repondre.
L'utilisation de canal et de select permet de rendre la reception de message bloquant de ce fait si le maitre arrete d'envoyer des messages les esclaves vont juste attendre sur le select jusqu'a ce que le maitre soit de nouveau actif. De meme le maitre transmettant en multicast, il n'est pas affecté par la disparition ou l'arrivé d'un nouvel esclave. 
