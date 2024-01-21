package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func ReadPBM(filename string) (*PBM, error) {

	var pbmIn = &PBM{}

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
		if strings.Contains(line, "#") { //On ignore les commentaires
			continue
		}
		if i == 0 { //On lit actuellement le magic number
			pbmIn.magicNumber = line
		}
		if i == 1 {
			size := strings.Fields(scanner.Text())
			if len(size) != 2 {
				return nil, fmt.Errorf("Taille du format invalide") //On créé une erreur
			}
			pbmIn.width, err = strconv.Atoi(size[0])
			if err != nil {
				return nil, fmt.Errorf("longueur invalide") //On créé une erreur
			}
			pbmIn.height, err = strconv.Atoi(size[1])
			if err != nil {
				return nil, fmt.Errorf("hauteur invalide") //On créé une erreur
			}

			if pbmIn.magicNumber == "P1" {
				pbmIn.data = make([][]bool, pbmIn.height)
				for j := range pbmIn.data {
					pbmIn.data[j] = make([]bool, pbmIn.width)
				}
			}

		}
		if i > 1 {
			// Lire le "body"
			if pbmIn.magicNumber == "P1" {
				// Initialiser la matrice de données

				lineData := strings.Fields(line)
				if len(lineData) != pbmIn.width {
					return nil, fmt.Errorf("Largeur de la ligne du body invalide") // On créé une erreur
				}
				for j, pixel := range lineData {
					val, err := strconv.Atoi(pixel)
					if err != nil {
						return nil, fmt.Errorf("Valeur de pixel invalide") // On créé une erreur
					}
					pbmIn.data[i-2][j] = val == 1
					println(i-2, j, val, val == 1, pbmIn.data[i-2][j])
				}
			} else if pbmIn.magicNumber == "P4" && pbmIn.data == nil {
				// Initialiser la matrice de données
				pbmIn.data = make([][]bool, pbmIn.height)
				for j := range pbmIn.data {
					pbmIn.data[j] = make([]bool, pbmIn.width)
				}
				p4Buff := 0
				p4LineIn := 0
				for _, asciiCode := range []byte(line) {
					binaryCode := DecimalToBinary(int(asciiCode), 8)
					for b := 0; b < len(binaryCode); b++ {
						if p4Buff >= pbmIn.width {
							p4Buff = 0
							p4LineIn++
							continue
						}
						if p4LineIn >= pbmIn.height {
							break
						}

						pbmIn.data[p4LineIn][p4Buff] = binaryCode[b] == 1
						p4Buff++
					}
				}
			} else {
				return nil, fmt.Errorf("Magic Number invalide") // On créé une erreur
			}
		}
		i++
	}
	file.Close()
	return pbmIn, nil
}

func DecimalToBinary(decimal int, fixedLength int) []int {
	binaryArray := []int{}

	for decimal > 0 {
		remainder := decimal % 2
		binaryArray = append([]int{remainder}, binaryArray...)
		decimal = decimal / 2
	}

	// Remplir avec des zéros à gauche pour atteindre la longueur fixe
	for len(binaryArray) < fixedLength {
		binaryArray = append([]int{0}, binaryArray...)
	}

	return binaryArray
}

func BinaryToDecimal(binaryArray []int) int {
	decimal := 0
	power := len(binaryArray) - 1

	for _, bit := range binaryArray {
		decimal += bit * (1 << power)
		power--
	}

	return decimal
}

func BinaryToWindows1252(binaryArray []int) (string, error) {
	// Vérifier que la longueur du tableau binaire est un multiple de 8
	if len(binaryArray)%8 != 0 {
		return "", fmt.Errorf("la longueur du tableau binaire doit être un multiple de 8")
	}

	// Convertir la séquence binaire en un tableau d'octets
	var byteArray []byte
	for i := 0; i < len(binaryArray); i += 8 {
		byteValue := BinaryToDecimal(binaryArray[i : i+8])
		byteArray = append(byteArray, byte(byteValue))
	}

	return string(byteArray), nil
}

func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

func (pbm *PBM) At(x, y int) bool {
	if x < 0 || y < 0 || x >= pbm.width || y >= pbm.height {
		return false
	}
	return pbm.data[x][y]
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[x][y] = value
}

func (pbm *PBM) Invert() {
	for i := 0; i < len(pbm.data); i++ {
		for j := 0; j < len(pbm.data[i]); j++ {
			if (pbm.data[i][j]) == false {
				pbm.data[i][j] = true
			} else if (pbm.data[i][j]) == true {
				pbm.data[i][j] = false
			}
		}
	}
}

func (pbm *PBM) Flip() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width/2; x++ {
			// Échanger les pixels de manière symétrique par rapport à l'axe vertical central
			oppositeX := pbm.width - 1 - x
			pbm.data[y][x], pbm.data[y][oppositeX] = pbm.data[y][oppositeX], pbm.data[y][x]
		}
	}
}

func (pbm *PBM) Flop() {
	for x := 0; x < pbm.width; x++ {
		for y := 0; y < pbm.height/2; y++ {
			// Échanger les pixels de manière symétrique par rapport à l'axe horizontal central
			oppositeY := pbm.height - 1 - y
			pbm.data[y][x], pbm.data[oppositeY][x] = pbm.data[oppositeY][x], pbm.data[y][x]
		}
	}
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(pbm.magicNumber + "\n")
	if err != nil {
		return err
	}
	_, err = writer.WriteString(fmt.Sprintf("%d %d\n", pbm.width, pbm.height))
	if err != nil {
		return err
	}

	for _, row := range pbm.data {
		for _, pixel := range row {
			var pixelValue string
			if pixel {
				pixelValue = "1"
			} else {
				pixelValue = "0"
			}
			_, err = writer.WriteString(pixelValue)
			if err != nil {
				return err
			}
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
