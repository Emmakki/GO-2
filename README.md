Projet ELP dijkstra_go:

Objectifs du projet :
Créer en go un client serveur (TCP) pour retourner la liste des chemins les plus courts de tous les noeuds vers tous les noeuds pour un graph donné par l'algorithme de Dijkstra.

Les grandes étapes de l'exécution :

-Génération et écriture du graph à traiter (par Graph_generator.go)
-Extraction des données du graph puis envoi au serveur (par client.go)
-Execution de l'algorithme de Dijkstra (par server.go)
-récupération et traitement des données envoyés par le client
-décomposition en go routine de l'éxecution de 1 vers tous les noeuds pour chaque noeud de notre graph.
-renvois des résultats au client
-Reception des résultats de l'algo par le client et écriture du fichier de sortie  (par client.go)
-le choix de l'implémentation client server TCP  plutot que UDP est justifié par le fait que nous ne pouvons pas nous permettre de perdre des données , comme chacune d'elle est importante.

Implémentation d'algorithme du Dijkstra détaillé (éxecuté par le serveur):

/**Structure utilisé:**/

Pour éxecuter cette algorithme nous allons utiliser deux structures facilitant le calcul. 

****************************************************************************
    type elementGraph struct { 
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
*****************************************************************************
La structure elementGraph nous permet de constituer notre graphe à analyser. Elle contitue un lien avec le sommet d'où il part (from), le sommet d'arrivé (to) et le poids/distance entre ces deux sommets (weight).
La structure chemin facilite l'écriture du résultat de l'algorithme étant constitué le sommet d'où l'on vient par le plus court chemin (from) et la distance minimale jusqu'à ce point depuis le sommet de départ depuis lequel on applique l'algorithme (weight).
La structure data nous permet de stocker notre résultat pour un seul noeud (départ) c'est le noeud depuis lequel  on veut calculer tous les plus courts chemin avec les autres noeuds du graphe. routes c'est pour stocker les noeuds depuis lequels on passe , et distances , c'est la distance la plus courte depuis le noeuds départ a tous les noueds de graphes.

********Implémentation:
Le graphe à traité est représenté par un slice de elementGraph, écrit lors de la lecture envoyé par le client.

Pour la bonne réalisation de l'algorithme, nous implémentons différents élements:

*****************************************************************************
    sommetParcouru := make(map[int]int)
    tabD := make(map[int]chemin)
    voisin := toutLesVoisins(g, noeuds)
*****************************************************************************
sommentParcouru nous permet de stocker les sommets par lequel nous sommes passés lors que la réalisation de l'algorithme.
La structure final est tabD stockant pour chaque sommet du graphe une structure chemin.
voisin récupère un map[int][]elementGrap retrouvant pour tout les sommets tout les somment qui lui sont lié.

*****************************************************************************
    for i := range noeuds {
		tabD[noeuds[i]] = chemin{-1, 0}
	}
*****************************************************************************
On initialise tabD pour tout les sommets de notre graphe son précédent sommet par -1 et 0 pour la distance, indiquant que le sommet n'a pas été parcouru.

/**Exécution:**/

Voici dans les grandes lignes ce que réalise l'algorithme:
    Il réalise le schéma suivant autant de fois qu'il y a de noeuds pour tous les parcourir
*Pour le sommet traité, il implémente tabD de ces voisins non encore parcouru et dont la distance est annoncé plus grande si le sommet traité n'est pas son précédent.
*Il recherche le sommet suivant à traité parmis les voisins non encore parcouru, il choisira le sommet dont la distance du lien est la plus faible. Si le sommet traité n'a plus de voisin, il parcourt le graphe en vérifiant s'il ne manque pas des sommets non traité
    Un fois tabD implémenté pour tout les noeuds, il remonte dans le tableau pour chaque sommet vers le sommet de départ en paramètre de la fonction.

/**Données renvoyées:**/

Les chemins à empreunter pour chaque sommet depuis le sommet de départ sont stocké dans toutChemins qui est un map[int][]int. Tandis que les distances minimales pour ses différents sommet sont stockées dans distanceMin qui est un map[int]int.
Ces structures sont ensuite écrite dans un channel, permettant la synchronisation des goroutines  et de la centralisation de  tout les résultants des exécutions de cet algorithme pour tout les sommets comme sommet initial.

******************************************************************
    distances := distanceMin
	routes := toutChemins
	dataCh <- data{depart, routes, distances}
******************************************************************
Après on ajoute deux structures dijk := make(map[int]map[int][]int et distance := make(map[int]map[int]int) où on stocke notre résultat final pour chaque noeud de notre graphe(clé de nos maps).