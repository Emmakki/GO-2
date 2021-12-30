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
	fmt.Println(dijkstra_routine(slice,noeuds))
	fin := time.Now()
	fmt.Println(fin.Sub(debut))

}

//Cette fonction permet de récupérer tous les voisins de tous les noeuds
// La fonction retourne un map on peut donc appeler la liste des noeuds visins facilement

var wg sync.WaitGroup

func Dijkstra(depart int, g []elementGraph, wg *sync.WaitGroup, dataCh chan<- data, noeuds []int) { {
	//exemple de graphe:
	defer wg.Done()
	sommetParcouru := make(map[int]int) //stock les int des sommets par lesquels on ne peut plus passer
	sommetParcouru[0] = depart

	tabD := make(map[int]chemin) //création tu tableau final
	//noeuds := toutSommets(g)
	var somActuel int = depart
	//var prec int

	voisin := toutLesVoisins(g, noeuds)

	//Implémentation du tableau final
	for i := range noeuds {
		tabD[noeuds[i]] = chemin{-1, 0}
	}
	//tabD[depart] = chemin{0, 0}
	//fmt.Print("tabD", tabD, "\n")

	for i := 1; i < len(noeuds); i++ {

		//somVoisin := voisins(g, somActuel, sommetParcouru)
		//fmt.Print(somVoisin, "\n")

		for _, v := range voisin[somActuel] {
			fmt.Print("sommet actuel", somActuel, "\n")
			if !exist(sommetParcouru, v.to) && (tabD[v.to].weight == 0 || tabD[v.to].weight > tabD[somActuel].weight+v.weight /*poids(g, somActuel, v.to)*/) {
				tabD[v.to] = chemin{somActuel, tabD[somActuel].weight + v.weight}
				//fmt.Print("tabD[v.to]  ", tabD[v.to], "\n") /*poids(g, somActuel, v.to)*/

				//fmt.Print("tabD", tabD, "\n")
			}
		}

		if left_neighbors(somActuel,voisin,sommetParcouru) != 0 {
		
			somActuel = minimum(somActuel, sommetParcouru, voisin)
		}else {
		
			for i := range noeuds {
				if exist(sommetParcouru, noeuds[i]) && left_neighbors(noeuds[i],voisin,sommetParcouru) != 0 {
					//fmt.Print(voisins(g, noeuds[i], sommetParcouru))
					somActuel = minimum( noeuds[i], sommetParcouru, voisin)
					break
				}
			}
		
		}

	

		sommetParcouru[i] = somActuel
		//fmt.Print("sommetparcouru", sommetParcouru, "\n")

		if toutSommetsParcouru(noeuds, sommetParcouru) {
			break
		}
	}

	fmt.Print("Les sommets:", noeuds, "\n")
	fmt.Print("tabD", tabD, "\n")

	//Retrouver le chemin pour tout les noeuds jusqu'au point de départ
	//Structure à renvoyer: map[{to from}]{[]int(chemin)} + map [to]int(distancemin)
	toutChemins := make(map[int][]int)
	distanceMin := make(map[int]int)
	for i := 0; i < len(noeuds); i++ {
		noeud := noeuds[i]
		fmt.Print("noeud", noeud, ": \n")
		distanceMin[noeud] = tabD[noeud].weight
		fmt.Print("tabD [[noeud].Weight]", distanceMin[noeud], "\n")
		toutChemins[noeud] = append(toutChemins[noeud], noeud)
		noeudPrec := noeud

		for tabD[noeudPrec].from != -1 {
			toutChemins[noeud] = append(toutChemins[noeud], tabD[noeudPrec].from)

			//fmt.Print(tabD[noeudPrec].from, "\n")
			noeudPrec = tabD[noeudPrec].from
		}

	}
	toutChemins = renverse(toutChemins)
	fmt.Print("Pour chaque noeud chemin à prendre:", toutChemins, "\n")
	fmt.Print("distance du plus court chemin pour chaque noeud:", distanceMin, "\n")

	distances := distanceMin
	routes := toutChemins
	//return distanceMin,toutChemins
	dataCh <- data{depart, routes, distances}
}

}

func dijkstra_routine(graph []elementGraph, noeuds []int) (map[int]map[int][]int, map[int]map[int]int) {
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
	i := 0
	for _, node := range noeuds {
		fmt.Print(i, "\n")
		go Dijkstra(node, graph, &wg, datach, noeuds)
		i++
	} //

	wg.Wait()
	close(datach) //this closes the dataCh channel, which will make the for-range loop exit once all the data has been read
	<-done        //we wait for all of the data to get read and put into maps

	return dijk, distance

}


func left_neighbors(somActuel int, voisin map[int][]elementGraph, sommetParcouru map[int]int) int {
	nb_voisin:=0
	for _, s := range voisin[somActuel] {
		if !exist(sommetParcouru,s.to){
			nb_voisin+=1
		}

	}
	return nb_voisin
}

func Voisins(sommet int, g []elementGraph) []elementGraph { //[{1 2 6} {1 3 4} {1 8 12} … ]
	var voisin []elementGraph
	for _, s := range g {
		if s.from == sommet {
			voisin = append(voisin, s)
		}

		/*//Si graphe non orienté décommenter ça
		if s.to == sommet {
			var sr elementGraph = elementGraph{s.to, s.from, s.weight}
			voisin = append(voisin, sr)
		}
		*/
	}
	return voisin
}

//Cette fonction permet de récupérer tous les voisins de tous les noeuds
// La fonction retourne un map on peut donc appeler la liste des noeuds visins facilement
func toutLesVoisins(g []elementGraph, noeuds []int) map[int][]elementGraph { // [1] : [{1 2 6} {1 3 4} {1 8 12} … ], [2] : [ { … } …]
	toutVoisins := make(map[int][]elementGraph) //instantiation
	for _, noeud := range noeuds {              // parcours la liste des noeuds qui existe
		toutVoisins[noeud] = Voisins(noeud, g) // Ajout de la liste des voisins au map
	}
	return toutVoisins
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

func minimum(sommet int, sommetParcouru map[int]int, voisin map[int][]elementGraph) int {
	min := voisin[sommet][0].weight
	suivant := voisin[sommet][0].to
	for _, v := range voisin[sommet] {
        if !exist(sommetParcouru, v.to) {
            min = v.weight
            suivant = v.to
            break
        }
    }

    for _,s := range voisin[sommet] {
        if !exist(sommetParcouru, s.to) {
            if s.weight < min {
                min = s.weight
                suivant = s.to
            }
        }

    }
    return suivant
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
