package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/*
Ce fichier à pour but de générer un graph de taille et de path donnés en entrés, dans le but d'être utilisés en fichier d'entrée pour Client.go
	- DEBUG commentaires de debug
*/

// go run rand.go 300 /in/300.txt
//récupère les arguments fournis à l'éxecution du fichier
func getArgs() (int, string) {
	// Vérifie qu'il y ai bien un argument
	if len(os.Args) != 3 {
		fmt.Println("Erreur : l'usage de rand.go nécessite l'appel suivant : go run rand.go <size (nombre de lien)> <graph.txt>")
		os.Exit(1) //sinon exit
	} else {
		//récupère le nom du fichier et vérifie que le fichier existe bien
		size, err := strconv.Atoi(os.Args[1])
		size=size+1
		if err != nil {
			fmt.Printf("Vous devez utiliser le générateur ainsi : go run rand.go <size>\n")
			os.Exit(1) //sinon exit
		} else {
			filename := os.Args[2]
			// Tout est ok, je retourne le nom du fichier pour la suite du script
			return size, filename
		}
		// Ne devrait jamais retourner
	}
	return -1, ""
}

//Génère un poids aléatoire entre 1 et 16
func randWeight(min, max int) int {
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + rand.Intn(max-min)
}

//Supprime l'élement à l'index s du slice donné en entrée
func remove(slice []int, s int) []int {
	end := append(slice[:s], slice[s+1:]...)
	return end
}

//méthode appelant remove pour supprimer un élément connu (mais index inconnu) d'un slice à un index inconnu (parcours le slice)
func remove_element(slice []int, elt int) []int {
	i := 0
	for slice[i] != elt {
		i++
	}
	return remove(slice, i)
}

//Permet de générer le string représentant le graph pour une taille donnée
func generator(size int) string { //param : nb de lien voulu
	var alphabet []int          //la variable alphabet est un tableau d'entiers
	for i := 0; i < size; i++ { //On remplit le tableau d'entiers par les entiers consécutifs de 0 à la taille voulue
		alphabet = append(alphabet, i)
	}

	//fmt.Printf("alphabet %d \n", alphabet) DEBUG
	neighb := make(map[int][]int, size) //crée un map associant à un entier (le noeud) un tableau d'entiers contenant les noeuds avec lequels il est possible de matcher
	/*	A terme on aura quelque chose ainsi:
		[1] = [2],[3],[...],[n]
		[2] = [1],[3],[...],[n]
		[n] = [1], ... [n-1]
	*/

	//boucle d'initialisation de neighb (voisins disponibles)
	//fmt.Printf("Debug génération neighb\n") DEBUG
	for _, letter := range alphabet {
		neighb[letter] = make([]int, len(alphabet)) //Pour chaque clé du map on lui donne comme valeur un tableau d'int de la taille d'alphabet #? pourquoi pas alphabet-1 ?
		copy(neighb[letter], alphabet)              //On copie de telle sorte que le tableau d'entier dans neighb[valeur] soit identitique à alphabet on fait ça pour éviter les passages par référence
		//fmt.Printf("neighb[%d] :  %d \n",letter, neighb[letter]) DEBUG
		remove(neighb[letter], letter)                          //On enleve l'élément de la liste de ses voisins possibles
		neighb[letter] = neighb[letter][:len(neighb[letter])-1] //pour résoudre un problème : quand on retirait notre lettre, le tableau bourrait avec le dernière élément pour avoir notre taille de l'alphabet. On fait donc ici notre tableau étant de taille alphabet-1
	}
	//fmt.Printf("neighb %d \n", neighb) DEBUG
	//contient la liste des lettres de l'alphabet pour lesquels il reste des voisins à tirer
	var from, to int
	var toWrite string //noeud de départ -> noeud d'arrivé -> résultat de la fonction
	n := 0
	reste := alphabet
	for k := 1; k < len(reste); k++ {
		var taille int
		from = reste[n]
		 //[0 1 2 3 4]
		if len(neighb[from]) > 3 {
			taille = 3
		} else {
			taille = len(neighb[from])
		}
		for i := 0; i < taille; i++ {

			fmt.Print(neighb[from])
			to = neighb[from][0]
			//fmt.Print(to,"\n")
			if len(neighb[from]) > 1 {
				neighb[from] = remove_element(neighb[from], to)
			}
			if len(neighb[to]) > 1 {
				neighb[to] = remove_element(neighb[to], from)
				fmt.Print(neighb[to], "\n") // retrait de l'inverse selon les memes conditions
			}

			weight := randWeight(1, 5)
			if k%2 == 0 {
				weight = randWeight(1, 9)
			}
			if k%3 == 0 {
				weight = randWeight(9, 16)
			}

			toWrite += fmt.Sprintf("%d %d %d\n", from, to, weight)

			//toWrite += fmt.Sprintf("%d %d %d\n", to, from, weight) //j'ai ajouté cette ligne pour la fonction voisins dans serveur
		}
		n+=1
	}
	toWrite += ". . ."

	return toWrite
}
func randLetter(alphabet []int) int {
	rand.NewSource(time.Now().UnixNano())
	//liste des noeuds possibles
	return alphabet[rand.Intn(len(alphabet))] //je prends une "lettre" aléatoire dans l'objet alphabet
}

// fonction principale qui à pour rôle d'écrire le graph d'une taille donnée dans un fichier à un path donné
func writeGraph(size int, path string) {

	fmt.Printf("Création du fichier %v et génaration d'un graph de taille %d \n", path, size)
	//affichage et résumé de l'opération
	f, err := os.OpenFile(path, //ouvre le fichier donné en argument (méthode d'ouverture de fichier généralisée (plus précise que os.Open ou os.Create))
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644) /*ouvre le fichier avec les tag 	CREATE (crée le fichier si il n'existe pas, avec les permissions données en dernier argument)
	WRONLY (ouvre le fichier en écriture seulement)
	TRUNC (si possible, tronque le fichier à l'ouverture #? ) overwrite ?
	la permission 0644 représente l'équivalent octal du FileMod #?  permission d'écriture 6: rw 4 r'*/
	if err != nil { //Si l'erreur est non nulle l'afficher
		log.Println(err)
	}
	defer f.Close()                                           //L'utilisation de defer sur Close permet de s'assurer que le fichier se fermera quand toutes les actions seront effectuées
	if _, err := f.WriteString(generator(size)); err != nil { //On écrit dans le fichier le résultat de la fonction generateTie en fonction de la taille de graph voulue, seulement si cette écriture ne produit pas d'erreur
		log.Println(err) //Sinon afficher l'erreur
	}
}

// fonction main
func main() {
	s := time.Now()
	rand.NewSource(time.Now().UnixNano())             //lance un timer
	writeGraph(getArgs())                             //appelle la fonction writeGraph selon la taille et le chemin d'enregistrement du fichier donnés en arguments
	fmt.Printf("Éxécution en  : %s\n", time.Since(s)) //affiche le temps d'execution
}
