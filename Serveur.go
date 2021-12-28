package main

import (
	"bufio"
	"fmt"
	"log"
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

//Cette fonction premet de retirer les valeurs dupliquées dans un slice
func unique(slice []int) []int {
	keys := make(map[int]bool)    // on fait une map qui associe un bool à chaque entier
	list := []int{}               // slice d'int sans taille
	for _, entry := range slice { //foreach
		if _, value := keys[entry]; !value { //on vérifie si la clé booléenne de l'entier éxiste, sinon on la créer, et on passe dans le if si la valeur est false, c'est à dire si on est jamais passé par celle ci
			keys[entry] = true         //on passe à true pour indiquer qu'on à déjà vérifier ce noeud
			list = append(list, entry) // On rajoute à notre liste de sortie le noeud.
		}
	}
	return list //on retourne notre tableau avec les noeuds uniques
}


func main() {

	file, err := os.Open("graph_test") // permet d'ouvir le fichier txt en appelant file
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // ferme le fichier quand le main est fini

	//compteurLiens := 0

	scanner := bufio.NewScanner(file)

	//On lit la première ligne du doc pour avoir le nombre de liens avec .scan()

	var slice []elementGraph
	var noeuds []int
	//var out string
	//start := time.Now() // notre slice de slice

	//On remplit un slice avec les différents liens
	for scanner.Scan() { // on parcours tous le fichier ligne par ligne

		splitted := strings.Split(scanner.Text(), " ") // on lit chaque ligne avec la fonction scanner.text() et on créer un slice de stringraphavec la fonction strings.Split avec comme indication " " comme séparateur
		if splitted[0] != "." {                        // si on a un point on est en EOF donc on ne prend pas
			// Je convertis mes entiers pcq il était stocké comme un string
			from, _ := strconv.Atoi(splitted[0]) //point de départ converti en int
			to, _ := strconv.Atoi(splitted[1])   // point d'arrivé converti en int
			noeuds = append(noeuds, from, to)    // ajout de l'entièreté des noeuds
			weight, _ := strconv.Atoi(splitted[2])
			// J'ajoute à mon slice un elementGraph
			noeuds = append(noeuds, from, to)
			slice = append(slice, elementGraph{from, to, weight})
		} else {

			break // pour une sortie de fichier en EOF (. . .)
		}

	}

	noeuds = unique(noeuds) //pour avoir un tableau contenant un exemplaire de tous les noeuds de notre graph
	sort.Ints(noeuds)
	debut := time.Now()
	fmt.Println(djikstra_routine(slice,noeuds))
	fin := time.Now()
	fmt.Println(fin.Sub(debut))
}


//Cette fonction permet de récupérer tous les voisins de tous les noeuds
// La fonction retourne un map on peut donc appeler la liste des noeuds visins facilement


var wg sync.WaitGroup

func Dijkstra(depart int, g []elementGraph, wg *sync.WaitGroup, dataCh chan<- data, noeuds []int) {
	//exemple de graphe:
	defer wg.Done()
	sommetParcouru := make(map[int]int) //stock les int des sommets par lesquels on ne peut plus passer
	sommetParcouru[0] = depart

	tabD := make(map[int]chemin) //création tu tableau final

	var somActuel int = depart
	var prec int

	//Implémentation du tableau final
	for i := range noeuds {
		tabD[noeuds[i]] = chemin{-1, 0}//-1 pour dire que le noeud n'a pas de noeud parent
	}
	//fmt.Print(tabD, "\n")

	for i := 1; i < len(noeuds); i++ {

		for _, v := range voisins(g, somActuel, sommetParcouru) {
			if !exist(sommetParcouru, v) && (tabD[v].weight == 0 || tabD[v].weight > tabD[somActuel].weight+poids(g, somActuel, v)) {
				tabD[v] = chemin{somActuel, tabD[somActuel].weight + poids(g, somActuel, v)}
			}
		}

		if len(voisins(g, somActuel, sommetParcouru)) > 1 {
			prec = somActuel
			somActuel = minimum(g, somActuel, sommetParcouru)
			//fmt.Print(prec, somActuel, "\n")
		} else {
			somActuel = minimum(g, prec, sommetParcouru)
			//fmt.Print(prec, somActuel, "\n")
		}

		sommetParcouru[i] = somActuel
		//fmt.Print(sommetParcouru, "\n")

		if toutSommetsParcouru(noeuds, sommetParcouru) {
			//fmt.Print(toutSommetsParcouru(noeuds, sommetParcouru), "\n")
			break
		}
	}
	///fmt.Print("Les sommets:", noeuds, "\n")
	//fmt.Print(tabD, "\n")

	//Retrouver le chemin pour tout les noeuds jusqu'au point de départ
	//Structure à renvoyer: map[{to from}]{[]int(chemin)} + map [to]int(distancemin)
	toutChemins := make(map[int][]int)
	distanceMin := make(map[int]int)
	for i := 0; i < len(noeuds); i++ {
		noeud := noeuds[i]
		//fmt.Print(noeud, ": \n")
		distanceMin[noeud] = tabD[noeud].weight
		toutChemins[noeud] = append(toutChemins[noeud], noeud)
		noeudPrec := noeud
		for tabD[noeudPrec].from != -1 {
			toutChemins[noeud] = append(toutChemins[noeud], tabD[noeudPrec].from)
			//fmt.Print(tabD[noeudPrec].from, "\n")
			noeudPrec = tabD[noeudPrec].from

		}
	}
	toutChemins = renverse(toutChemins)
	distances := distanceMin
	routes := toutChemins
	dataCh <- data{depart, routes, distances}

}

func djikstra_routine(graph []elementGraph, noeuds []int) (map[int]map[int][]int, map[int]map[int]int) {
	datach := make(chan data)
	done := make(chan bool)
	dijk := make(map[int]map[int][]int)
	distance := make(map[int]map[int]int)

	go func() {
		for d := range datach {
			dijk[d.depart] = d.routes
			distance[d.depart] = d.distances
		}
		done <- true //this is used to wait until all data has been read from the channel
	}()
	wg.Add(len(noeuds))
	for _, node := range noeuds {
		go Dijkstra(node, graph, &wg, datach, noeuds)
	} //

	wg.Wait()
	close(datach) //this closes the dataCh channel, which will make the for-range loop exit once all the data has been read
	<-done        //we wait for all of the data to get read and put into maps

	return dijk, distance

}
func voisins(graph []elementGraph, sommet int, sommetParcouru map[int]int) map[int]int {
	voisin := make(map[int]int)
	i := 0
	for _, s := range graph {
		if sommet == s.from && !exist(sommetParcouru, s.to) {
			voisin[i] = s.to
			i++
		}
	}
	return voisin
}

func toutSommetsParcouru(sommet []int, sommetParcouru map[int]int) bool {
	var fini bool = true
	for i := 0; i < len(sommet); i++ {
		if !exist(sommetParcouru, sommet[i]) {
			fini = false
		}
	}
	return fini

}



func exist(tab map[int]int, valeur int) bool {
	var existe bool = false
	for i := 0; i < len(tab); i++ {
		if valeur == tab[i] {
			existe = true
		}
	}
	return existe
}

func minimum(g []elementGraph, sommet int, sommetParcouru map[int]int) int {
	min := poids(g, sommet, voisins(g, sommet, sommetParcouru)[0])
	suivant := voisins(g, sommet, sommetParcouru)[0]
	for _, s := range g {
		if s.from == sommet && !exist(sommetParcouru, s.to) {
			if s.weight < min {
				min = s.weight
				suivant = s.to
			}
		}

	}
	return suivant
}

func poids(graph []elementGraph, depart int, arrive int) int {
	var poid int
	for _, s := range graph {
		if s.from == depart && s.to == arrive {
			poid = s.weight
		}
	}
	return poid
}

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
