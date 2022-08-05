package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Person struct {
	ID        string
	Name      string
	ManagerID string
}

func main() {
	csvFile := os.Getenv("BMB_CSV_FILE")

	if csvFile == "" {
		fmt.Println("Usage: Environment variable BMB_CSV_FILE is not set")
		os.Exit(1)
	}

	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	lowerNameToPerson := make(map[string]Person)
	managerToDirects := make(map[string][]Person)

	for _, r := range records {
		p := Person{r[0], r[1], r[2]}
		key := strings.ToLower(r[1])
		lowerNameToPerson[key] = p
		managerToDirects[p.ManagerID] = append(managerToDirects[p.ManagerID], p)
	}

	searchName := os.Args[1]
	var start Person

	for name, p := range lowerNameToPerson {
		if strings.Contains(name, searchName) {
			start = p
			break
		}
	}

	var managers, ics []Person

	for _, person := range managerToDirects[start.ID] {
		if _, manager := managerToDirects[person.ID]; manager {
			managers = append(managers, person)
		} else {
			ics = append(ics, person)
		}
	}

	allManagers := append([]Person{}, managers...)
	queue := append([]Person{}, managers...)
	allIcs := append([]Person{}, ics...)

	for len(queue) != 0 {
		var newQueue []Person

		for _, person := range queue {
			if directs, isManager := managerToDirects[person.ID]; isManager {
				allManagers = append(allManagers, person)
				newQueue = append(newQueue, directs...)
			} else {
				allIcs = append(allIcs, person)
			}
		}

		queue = newQueue
	}

	sort.Slice(managers, func(i, j int) bool {
		return managers[i].Name < managers[j].Name
	})

	sort.Slice(ics, func(i, j int) bool {
		return ics[i].Name < ics[j].Name
	})

	fmt.Println("Direct reports for:", start.Name)
	fmt.Println("Number of directs:", len(managers)+len(ics))

	if len(managers) == 0 {
		fmt.Println("Managers: N/A")
	} else {
		fmt.Println("Managers: ", PrintNames(managers))
	}

	if len(ics) == 0 {
		fmt.Println("Individual Contributors: N/A")
	} else {
		fmt.Println("Individual Contributors:", PrintNames(ics))
	}

	fmt.Println("Entire Org Size:", len(allManagers)+len(allIcs))
}

func PrintNames(people []Person) string {
	names := make([]string, 0, len(people))

	for _, person := range people {
		names = append(names, person.Name)
	}

	return strings.Join(names, ", ")
}
