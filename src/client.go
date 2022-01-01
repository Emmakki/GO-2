package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*
Le programme client va nous permettre de se connecter au serveur, de recupérer les données de notre graph généré précedemment,
d'envoyer ses données au serveur afin qu'ils soient traités par l'algorithme de Dijkstra.
Les données finales vont etre réenvoyer au client qui va les écrire dans un fichier texte.
*/
func getargs() (int, string) {
	// Vérifie si le client a donné le nom du graphe, le numero de port qu'il veut utilisé
	if len(os.Args) != 3 {
		fmt.Println("Erreur : l'usage de Client.go nécessite l'appel suivant : go run Client.go <graph.txt> <portNumber>. ")
		os.Exit(1)

	} else {
		//récupère le port et vérifie si une erreur est intervenue lors de la conversion
		fmt.Printf("Port numéro : %s\n", os.Args[2])
		portNumber, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Le port donné n'est pas valide, vérifiez qu'il soit bien un nombre.\n")
			os.Exit(1)
		} else { //Si pas d'erreurs
			//récupère le nom du fichier et vérifie que le fichier existe bien
			filename := os.Args[1]
			_, err := os.Stat(filename) //Nous permet d'obtenir des info (permissions, name, size)
			if os.IsNotExist(err) {     //s'execute si le fichier n'a pas été trouvé
				fmt.Printf("Erreur : le fichier %v n'existe pas, ou il fait référence à un dossier.\n", filename)
				os.Exit(1)
			} else {
				return portNumber, filename
			}
		}
	}
	return -1, ""
	//Ne devrait pas etre affiché
}

func main() {
	//On lance 2 timers
	start := time.Now()
	s := time.Now()
	port, filename := getargs()
	//Connection au serveur
	fmt.Printf("TCP sur le port %d\n", port)
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port)) //formatage selon 127.0.0.1:xxxx ex: 127.0.0.1:4000
	fmt.Printf("PORT STRING |%s|\n", portString)
	connection, err := net.Dial("tcp", portString) //net.Dial is the general-purpose connect command. First parameter is a string specifying the network. In this case we are using tcp. Second parameter is a string with the address of the endpoint in format of host:port

	if err != nil { //Une erreur a été détectée
		fmt.Printf("Connection echouée \n")
		os.Exit(1)
	} else {
		defer connection.Close() //on assure la fermeture de la connection une fois toutes les actions terminées

		serveur := bufio.NewReader(connection) // On ajoute un reader pour écouter le serveur en retour
		fmt.Printf("Vous etes connecté\n")
		fmt.Printf("Connection au serveur en : %s\n", time.Since(s)) //temps de connection
		s = time.Now()                                               //on relance un autre timer
		file, err := os.Open(filename)                               // On ouvre le fichier qui a été donné en argument
		if err != nil {
			fmt.Printf("Ouverture du fichier impossible \n")
			os.Exit(1)
		} else {
			defer file.Close()
			scanner := bufio.NewScanner(file) // On met un scanneur pour scanner notre fichier
			for scanner.Scan() {              //On va parcourir le fichier ligne par ligne jusqu'a ce que le fichier se termine
				txt := scanner.Text() // On récupère la ligne
				//Text reads each line
				io.WriteString(connection, txt+"\n") //Ici on envoie la ligne plus un retour à la ligne au serveur
			}
			err := scanner.Err()
			if err != nil {
				fmt.Printf("Erreur avec le scanneur \n")
				os.Exit(1)
			}
			fmt.Printf("Fichier parsé et envoyé en : %s\n", time.Since(s))
			s = time.Now() //encore un timer pour la réponse
			//Après avoir tout envoyé on récupère la réponse du serveur
			res := fmt.Sprintf("res_%v", filepath.Base(filename)) //filepath.Base nous donne le dernier element du path donné (dans le cas où un path a été donné)
			resultat := ""
			for {
				results, err := serveur.ReadString('#') //On attend la réponse du serveur par le reader instancié précédemment

				if err != nil {
					fmt.Printf("Fin de traitement du serveur \n")
					break //on sort de la boucle une fois qu'on est arrivé à la fin du fichier
				}

				resultat += strings.TrimSuffix(results, "#") //on ajoute a chaque boucle les resultat du serveur dans la variable

			}
			fmt.Printf("Réception et traitement des données in : %s\n", time.Since(s))
			s = time.Now() //encore un timer pour la réponse
			solution, err := os.OpenFile(res, os.O_CREATE|os.O_RDWR, 0755)
			defer solution.Close() //defer la fermeture du fichier de sortie
			if err != nil {
				log.Fatalf("Fin de traitement du serveur \n")
			}
			_, err = solution.WriteString(resultat) //S'il n'y a pas eu d'érreurs on stock les resultats obtenus dans le fichier texte qu'on a créer précédemment
			if err != nil {
				fmt.Printf("Échec de l'écriture des données \n")
				os.Exit(1)
			}
			//pour évaluer la vitesse d'éxécution
			fmt.Printf("L'analyse de dijkstra est contenu dans : %v \n", res)
			fmt.Printf("Écriture des données en : %s\n", time.Since(s))
			fmt.Printf("Éxécuté en : %s\n", time.Since(start))

		}

	}
}