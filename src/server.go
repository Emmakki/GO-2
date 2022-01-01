package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type elementGraph struct { //element contenant le départ, l'arrivée et le poids de notre chemin
	from   int
	to     int
	weight int
}
type chemin struct {
	from   int
	weight int
}
type data struct {
	depart    int
	routes    map[int][]int
	distances map[int]int
}

// On va pouvoir recevoir des données du programme client et les traiter avec l'algorithme Dijkstra et les renvoyer au client

var wg sync.WaitGroup

func Dijkstra(depart int, g []elementGraph, wg *sync.WaitGroup, dataCh chan<- data, noeuds []int) {

	defer wg.Done()
	sommetParcouru := make(map[int]int) //stock les int des sommets par lesquels on ne peut plus passer
	sommetParcouru[0] = depart

	tabD := make(map[int]chemin) //création tu tableau final
	var somActuel int = depart

	voisin := toutLesVoisins(g, noeuds)

	//Implémentation du tableau final
	for i := range noeuds {
		tabD[noeuds[i]] = chemin{-1, 0}
	}

	//Passe par tout les noeuds et détermine la distance minimale pour tout ses noeuds ainsi que le noeud précédent
	for i := 0; i < len(noeuds); i++ {

		//Implémente tabD pour les voisins de somActuel nous encore visté
		//et donc la distance minimale supposé dans tabD est supérieur à celle trouvée en passant par sommet actuel
		for _, v := range voisin[somActuel] {
			if !exist(sommetParcouru, v.to) && (tabD[v.to].weight == 0 || tabD[v.to].weight > tabD[somActuel].weight+v.weight) {
				tabD[v.to] = chemin{somActuel, tabD[somActuel].weight + v.weight}
			}
		}

		//Cherche le prochain noeud à parcourrir
		//Si somActuel n'a plus de voisin non parcouru, il cherche dans tout le graphe
		restant, v_possible := left_neighbors(somActuel, voisin, sommetParcouru)
		if restant != 0 {
			somActuel = minimum(sommetParcouru, v_possible)
		} else {
			for j := range noeuds {
				restant, v_possible = left_neighbors(noeuds[j], voisin, sommetParcouru)
				if exist(sommetParcouru, noeuds[j]) && restant != 0 {
					somActuel = minimum(sommetParcouru, v_possible)
					break
				}
			}

		}

		//Annonce ce prochain noeud comme parcouru
		sommetParcouru[i+1] = somActuel

	}

	//Retrouver le chemin pour tout les noeuds jusqu'au point de départ
	//Structure à renvoyer: map[{to from}]{[]int(chemin)} + map [to]int(distancemin)
	toutChemins := make(map[int][]int)
	distanceMin := make(map[int]int)
	for i := 0; i < len(noeuds); i++ {

		noeud := noeuds[i]
		distanceMin[noeud] = tabD[noeud].weight
		toutChemins[noeud] = append(toutChemins[noeud], noeud)
		noeudPrec := noeud

		for tabD[noeudPrec].from != -1 {
			toutChemins[noeud] = append(toutChemins[noeud], tabD[noeudPrec].from)
			noeudPrec = tabD[noeudPrec].from
		}

	}
	toutChemins = renverse(toutChemins) //permet de mettre la liste des sommet à parcourir dans l'ordre (le premier de la liste est le point de départ)

	distances := distanceMin
	routes := toutChemins
	//return distanceMin,toutChemins
	dataCh <- data{depart, routes, distances}

}

func dijkstra_routine(graph []elementGraph, noeuds []int) (map[int]map[int][]int, map[int]map[int]int) {
	datach := make(chan data)
	done := make(chan bool)
	dijk := make(map[int]map[int][]int)
	distance := make(map[int]map[int]int)

	wg.Add(len(noeuds))
	//i := 0
	for _, node := range noeuds {
		//fmt.Print(i, "\n")
		go Dijkstra(node, graph, &wg, datach, noeuds)
		//i++
	} //

	go func() {
		for d := range datach {
			dijk[d.depart] = d.routes
			distance[d.depart] = d.distances
		}
		done <- true //this is used to wait until all data has been read from the channel
	}()

	wg.Wait()
	close(datach) //this closes the dataCh channel, which will make the for-range loop exit once all the data has been read
	<-done        //we wait for all of the data to get read and put into maps

	return dijk, distance

}

//Cette fonction donne pour un sommet donnée le nombre et quels sont ses voisins non parcouru
func left_neighbors(somActuel int, voisin map[int][]elementGraph, sommetParcouru map[int]int) (int, map[int]elementGraph) {
	nb_voisin := 0
	non_parcouru := make(map[int]elementGraph)
	for _, s := range voisin[somActuel] {
		if !exist(sommetParcouru, s.to) {
			non_parcouru[nb_voisin] = s
			nb_voisin += 1
		}

	}
	return nb_voisin, non_parcouru
}

//Récupère tout les voisins pour un sommet donné
func Voisins(sommet int, g []elementGraph) []elementGraph { //[{1 2 6} {1 3 4} {1 8 12} … ]
	var voisin []elementGraph
	for _, s := range g {
		if s.from == sommet {
			voisin = append(voisin, s)
		}
	}
	return voisin
}

//Cette fonction permet de récupérer tous les voisins de tous les noeuds
// La fonction retourne un map on peut donc appeler la liste des noeuds visins facilement
func toutLesVoisins(g []elementGraph, noeuds []int) map[int][]elementGraph { // [1] : [{1 2 6} {1 3 4} {1 8 12} … ], [2] : [ { … } …]
	toutVoisins := make(map[int][]elementGraph)
	for _, noeud := range noeuds {
		toutVoisins[noeud] = Voisins(noeud, g)
	}
	return toutVoisins
}

//Annonce si une valeur existe dans un map[int]int
func exist(tab map[int]int, valeur int) bool {
	var existe bool = false
	for i := 0; i < len(tab); i++ {
		if valeur == tab[i] {
			existe = true
		}
	}
	return existe
}

//Renvoie le int du sommet dont le poids est le plus petit parmis tout les voisins non parcourue d'un sommet
func minimum(sommetParcouru map[int]int, voisin_possible map[int]elementGraph) int {
	min := voisin_possible[0].weight
	suivant := voisin_possible[0].to

	for _, s := range voisin_possible {
		if !exist(sommetParcouru, s.to) {
			if s.weight < min {
				min = s.weight
				suivant = s.to
			}
		}

	}
	return suivant
}

//Retourne les valeur des slices dans un map[int][]int
func renverse(tab map[int][]int) map[int][]int {
	tabR := make(map[int][]int)
	for s := range tab {
		taille := len(tab[s])
		for i := 0; i < taille; i++ {
			tabR[s] = append(tabR[s], tab[s][taille-(i+1)])
		}
	}
	return tabR
}

func removedup(lst []int) []int {

	keys := make(map[int]bool) // on fait une map qui associe un bool à chaque entier
	list := []int{}            // slice d'int sans taille
	for _, entry := range lst {
		if _, value := keys[entry]; !value { //on rentre dans le if que quand le noeud apparait pour la première fois
			keys[entry] = true         //on passe à true pour indiquer qu'on à déjà vérifier ce noeud
			list = append(list, entry) // On rajoute à notre liste de sortie le noeud.
		}
	}
	return list
}
func getPort() int {
	if len(os.Args) != 2 {
		fmt.Printf("Veuillez mettre uniquement le numéro de port. Exemple: go run Server.go <port>")
		os.Exit(1)
	} else {
		port, err := strconv.Atoi(os.Args[1]) // String -> Int
		if err != nil {
			fmt.Printf("Le port n'est pas valide")
			os.Exit(1)
		} else {
			return port
		}
	}
	return -1
	//Ne doit pas etre retourné
}

func handleConnection(connection net.Conn) {
	defer connection.Close() //Assure la fermeture de la connection
	fmt.Println("Connection réussie")
	connReader := bufio.NewReader(connection) // On place un reader
	var listelm []elementGraph
	var nodes []int
	var out string
	start := time.Now()
	for {
		inputLine, err := connReader.ReadString('\n') // On lit les données envoyés par le client, on lit jusqu'a un retour a la ligne (\n)
		if err != nil {
			fmt.Printf("Error :|%s|\n", err.Error())
			break
		}
		inputLine = strings.TrimSuffix(inputLine, "\n") //Nous donne la ligne sans le retour à la ligne (\n)
		lst := strings.Split(inputLine, " ")            //Nous renvoie une liste obtenu grace au separateur " " (from to weight)
		if lst[0] != "." {                              //C'est uniquement le cas quand c'est la fin du fichier txt
			from, _ := strconv.Atoi(lst[0]) //String to int
			to, _ := strconv.Atoi(lst[1])
			nodes = append(nodes, from, to)
			weight, _ := strconv.Atoi(lst[2])
			listelm = append(listelm, elementGraph{from, to, weight})
		} else {
			break // Sortie de la boucle a la fin du fichier txt
		}
	}
	fmt.Printf("Donnée traitée et répartie in : %s\n", time.Since(start))
	nodes = removedup(nodes) //Nous permet d'avoir une liste sans les sommets dupliqués
	sort.Ints(nodes)         //On trie la liste par ordre croissant
	ways, distances := dijkstra_routine(listelm, nodes)
	fmt.Printf("Dijkstra done in : %s\n", time.Since(start))
	start = time.Now()
	out = ""
	for letter, graph := range ways {
		for l, way := range graph {
			if letter != l {
				out += fmt.Sprintf("%v %v %v %v \n", letter, l, way, distances[letter][l]) // départ | arrivé | chemin | poids
			}
		}
		out += "#"
	}

	io.WriteString(connection, fmt.Sprintf("%s ", out))
	fmt.Printf("Envoie des données en : %s\n", time.Since(start))

}

func main() {

	port := getPort() //On récupère le port
	fmt.Printf("Serveur TCP sur le port: %d\n", port)
	portString := fmt.Sprintf(":%s", strconv.Itoa(port)) // Ca nous permettra d'écouter le client
	ecoute, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Printf("Le serveur n'a pas pu écouter le client\n")
		os.Exit(1)
	}
	cptcon := 1 //nombre de connection
	for {
		fmt.Printf("Accepte la connection\n")
		con, errcon := ecoute.Accept() //accepte la connection
		if errcon != nil {
			fmt.Printf("Echec de l'acceptation de la connection\n")
			panic(errcon)
		}
		go handleConnection(con)
		cptcon += 1 //incremente le compteur
	}
}