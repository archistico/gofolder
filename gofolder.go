package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Cartella struct {
	Indirizzo	string
}

func (c *Cartella) GetIndirizzo() string {
	return c.Indirizzo
}

func (c *Cartella) Write() error {
	// err se non esiste il file/cartella
	_, err := os.Stat(c.Indirizzo)
	// Se non esiste creala
	if os.IsNotExist(err) {
		errMake := os.Mkdir(c.Indirizzo, os.ModeDir)

		// se errore nella creazione della cartella
		if errMake!=nil {
			return nil
		} else {
			return errMake
		}
	}
	return nil
}

type Cartelle struct {
	lista []Cartella
	testo []string
}

func (c *Cartelle) Mostra() {
	if len(c.lista) > 0 {
		for i, element := range c.lista {
			fmt.Printf("(%d) %s\n", i, element.Indirizzo)
		}
	} else {
		fmt.Println("Nessuna directory")
	}
}

func (c *Cartelle) Crea() {
	if len(c.lista) > 0 {
		for _, d := range c.lista {
			err := d.Write()
			if err != nil {
				fmt.Println("Errore:", d.GetIndirizzo())
			} else {
				fmt.Println("Creata:", d.GetIndirizzo())
			}
		}
	} else {
		fmt.Println("Nessuna directory")
	}

}

func (c *Cartelle) Add(n *Cartella) {
	c.lista = append(c.lista, *n)
}

func (c *Cartelle) Analizza() {
	var dirTesto string
	var livelloGenitore []string
	var numTab, numTabLast int = 0, 0
	var indirizzo string

	for _, testo := range c.testo {
		numTab = 0
		dirTesto = ""
		indirizzo = ""

		spazi := strings.NewReplacer("    ", "\t")
		testo = spazi.Replace(testo)
		numTab = Tab(testo)
		r := strings.NewReplacer("\t", "")
		dirTesto = r.Replace(testo)

		// Directory nella cartella principale
		if numTab == 0 {
			livelloGenitore = nil
		}

		// Riduce il livello
		if numTab < numTabLast {
			var diff = int(math.Abs(float64(numTab - numTabLast)))
			for x := 0; x <= diff; x++ {
				if len(livelloGenitore) > 0 {
					livelloGenitore = livelloGenitore[:len(livelloGenitore)-1]
				}
			}
			livelloGenitore = append(livelloGenitore, dirTesto)
		}

		// Cartelle ancora più nidificate
		if numTab > numTabLast {
			livelloGenitore = append(livelloGenitore, dirTesto)
		}

		// Cartelle stesso livello del genitore
		if numTab == numTabLast {
			if len(livelloGenitore) > 0 {
				livelloGenitore = livelloGenitore[:len(livelloGenitore)-1]
			}
			livelloGenitore = append(livelloGenitore, dirTesto)
		}

		// Concatena genitori per aggiungere alla lista cartelle
		if len(livelloGenitore)>0 {
			indirizzo = ""
			for _, d := range livelloGenitore {
				indirizzo += d + string(filepath.Separator)
			}
			temp := Cartella{ Indirizzo: indirizzo}
			c.Add(&temp)
		}
		numTabLast = numTab
	}
}


func Tab(text string) int {
	return strings.Count(text, "\t")
}

func main() {

	// LEGGERE
	file, err := os.Open("lista.txt")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var directories = []string{}
	for scanner.Scan() {
		directories = append(directories, scanner.Text())
	}

	// Principale
	cartelle:=Cartelle{}
	cartelle.testo = directories
	cartelle.Analizza()
	cartelle.Crea()
}

