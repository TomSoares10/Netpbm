package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

func ReadPGM(filename string) (*PGM, error) {
	var pgmIn = &PGM{}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur à l'ouverture du fichier:", err)
		return nil, err
	}

	// Créer un scanner pour lire le fichier ligne par ligne
	scanner := bufio.NewScanner(file)

	// Lire le fichier ligne par ligne
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "#") { // On ignore les commentaires
			continue
		}
		if i == 0 { // On lit actuellement le magic number
			pgmIn.magicNumber = line
			if pgmIn.magicNumber != "P2" && pgmIn.magicNumber != "P5" {
				return nil, fmt.Errorf("Format non pris en charge")
			}
		}
		if i == 1 {
			size := strings.Fields(scanner.Text())
			if len(size) != 2 {
				return nil, fmt.Errorf("Taille du format invalide") // On créé une erreur
			}
			pgmIn.width, err = strconv.Atoi(size[0])
			if err != nil {
				return nil, fmt.Errorf("largeur invalide") // On créé une erreur
			}
			pgmIn.height, err = strconv.Atoi(size[1])
			if err != nil {
				return nil, fmt.Errorf("hauteur invalide") // On créé une erreur
			}

			// Initialiser la matrice de données
			pgmIn.data = make([][]uint8, pgmIn.height)
			for j := range pgmIn.data {
				pgmIn.data[j] = make([]uint8, pgmIn.width)
			}
		}
		if i == 2 {
			// Lire la valeur maximale autorisée.
			maxValue, err := strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("valeur maximale invalide") // On créé une erreur
			}
			pgmIn.max = maxValue
		}
		if i > 2 {
			// Lire le "body"
			if pgmIn.magicNumber == "P2" {
				lineData := strings.Fields(line)
				if len(lineData) != pgmIn.width {
					return nil, fmt.Errorf("Largeur de la ligne du body invalide") // On créé une erreur
				}
				for j, pixel := range lineData {
					val, err := strconv.Atoi(pixel)
					if err != nil {
						return nil, fmt.Errorf("Valeur de pixel invalide") // On créé une erreur
					}
					pgmIn.data[i-3][j] = uint8(val)
				}
			} else if pgmIn.magicNumber == "P5" {
				x, y := 0, 0

				for _, asciiCode := range line {

					if x == pgmIn.width {
						x = 0
						y++
					}
					pgmIn.data[y][x] = uint8(asciiCode)
					x++
				}
			} else {
				return nil, fmt.Errorf("Magic Number invalide") // On créé une erreur
			}

		}
		i++
	}
	file.Close()
	return pgmIn, nil
}

func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[x][y]
}
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[x][y] = value
}
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the header
	_, err = fmt.Fprintf(writer, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)
	if err != nil {
		return err
	}

	// Write the pixel data
	if pgm.magicNumber == "P2" {
		// Write ASCII data for P2 format
		for _, row := range pgm.data {
			for j, pixel := range row {
				if j > 0 {
					_, err = writer.WriteString(" ")
					if err != nil {
						return err
					}
				}
				_, err = writer.WriteString(strconv.Itoa(int(pixel)))
				if err != nil {
					return err
				}
			}
			_, err = writer.WriteString("\n")
			if err != nil {
				return err
			}
		}
	} else if pgm.magicNumber == "P5" {
		// Write binary data for P5 format
		for _, row := range pgm.data {
			for _, pixel := range row {
				err = writer.WriteByte(pixel)
				if err != nil {
					return err
				}
			}
		}
	} else {
		return fmt.Errorf("unsupported PGM format: %s", pgm.magicNumber)
	}

	return writer.Flush()
}

func (pgm *PGM) Invert() {
	maxVal := uint8(pgm.max)
	for y := range pgm.data {
		for x := range pgm.data[y] {
			pgm.data[y][x] = maxVal - pgm.data[y][x]
		}
	}
}

func (pgm *PGM) Flip() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width/2; x++ {
			// Échanger les pixels de manière symétrique par rapport à l'axe vertical central
			oppositeX := pgm.width - 1 - x
			pgm.data[y][x], pgm.data[y][oppositeX] = pgm.data[y][oppositeX], pgm.data[y][x]
		}
	}
}

func (pgm *PGM) Flop() {
	for x := 0; x < pgm.width; x++ {
		for y := 0; y < pgm.height/2; y++ {
			// Échanger les pixels de manière symétrique par rapport à l'axe horizontal central
			oppositeY := pgm.height - 1 - y
			pgm.data[y][x], pgm.data[oppositeY][x] = pgm.data[oppositeY][x], pgm.data[y][x]
		}
	}
}

func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	// Mettre à jour la valeur maximale
	pgm.max = int(maxValue)

	// Ajuster les valeurs des pixels pour s'assurer qu'elles ne dépassent pas la nouvelle valeur maximale
	for i := range pgm.data {
		for j := range pgm.data[i] {
			if pgm.data[i][j] > maxValue {
				pgm.data[i][j] = maxValue
			}
		}
	}
}

func (pgm *PGM) Rotate90CW() {
	// Créer une nouvelle matrice de la taille de l'image pivotée
	newData := make([][]uint8, pgm.width)
	for i := range newData {
		newData[i] = make([]uint8, pgm.height)
	}

	// Effectuer la rotation
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			newData[j][pgm.height-1-i] = pgm.data[i][j]
		}
	}

	// Mettre à jour la structure PGM
	pgm.data = newData
	pgm.width, pgm.height = pgm.height, pgm.width
}

func (pgm *PGM) ToPBM() *PBM {
	threshold := pgm.max / 2
	pbm := &PBM{
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1", // Assuming ASCII format for PBM
	}

	pbm.data = make([][]bool, pbm.height)
	for y := range pbm.data {
		pbm.data[y] = make([]bool, pbm.width)
		for x := range pbm.data[y] {
			pbm.data[y][x] = pgm.data[y][x] > uint8(threshold)
		}
	}

	return pbm
}
