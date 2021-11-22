package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	file, err := os.Open("C:\\Users\\safae\\Desktop\\ELP\\graphe.txt") // permet d'ouvir le fichier txt en appelant file
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // ferme le fichier quand le main est fini

	compteurLiens := 0

	scanner := bufio.NewScanner(file)

	scanner.Scan() //On lit la première ligne du doc pour avoir le nombre de liens avec .scan()

	n, err := strconv.Atoi(scanner.Text()) //On le convertie en entier pour pouvoir creer nos slice de slice : "tableau de liens"
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	liens := make([][]string, n) // notre slice de slice

	//On remplit un slice avec les différents liens
	for scanner.Scan() { // on parcours tous le fichier ligne par ligne

		lien := strings.Split(scanner.Text(), " ") // on lit chaque ligne avec la fonction scanner.text() et on créer un slice de string avec la fonction strings.Split avec comme indication " " comme séparateur
		liens[compteurLiens] = lien                // on associe pour chaque item de la "slice de slice" "un lien"(sous forme d'un slice) du fichier texte dans l'ordre chronologique
		compteurLiens++

	}
	//fmt.Println(liens)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	m := make(map[string][][]string) // création d'une map
	//On remplit la map avec les élément de notre " slice de slice "
	for _, lien := range liens { // on parcours chaque "lien" de "liens" (slice de slice)
		if m[lien[0]] == nil {
			m[lien[0]] = make([][]string, 0) // on associe une clé à un tableau
		}
		if m[lien[1]] == nil {
			m[lien[1]] = make([][]string, 0)
		}
		//Enfin on rempli chaque tableau de la clé avec le lien associé(càd l'iteration en cours)
		m[lien[0]] = append(m[lien[0]], lien) // append permet de pouvoir aggrandir la taille de ce tableau en lui mettant un lien
		m[lien[1]] = append(m[lien[1]], lien)
		// exemple : le tableau de "B" va être creer à l'itération 1, ensuite il va être rempli avec le lien 1 ET à l'itération 2 et 3 on va l'agrendir en lui mettant le lien 2 et le lien 3
		// B :[[lien1][lien2][lien3]]
	

	fmt.Println("\nmap : <<sommet et liens>> a envoyer : \n", m)
	

	

	}
}
