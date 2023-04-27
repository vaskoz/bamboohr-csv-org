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

	idToPerson := make(map[string]*Person)
	lowerNameToPerson := make(map[string]*Person)
	managerToDirects := make(map[string][]*Person)
	var topOfOrg *Person

	for _, r := range records {
		p := Person{r[0], r[1], r[2]}
		if p.ManagerID == "" {
			topOfOrg = &p
		}
		key := strings.ToLower(r[1])
		lowerNameToPerson[key] = &p
		managerToDirects[p.ManagerID] = append(managerToDirects[p.ManagerID], &p)
		idToPerson[p.ID] = &p
	}

	var searchName string
	var start *Person

	if len(os.Args) > 1 {
		searchName = strings.ToLower(os.Args[1])
	}

	for name, p := range lowerNameToPerson {
		if strings.Contains(name, searchName) {
			start = p

			break
		}
	}

	if start == nil && searchName != "" {
		fmt.Printf("person '%s' not found\n", searchName)
		os.Exit(1)
	} else if searchName == "" {
		start = topOfOrg
	}

	var managers, ics []*Person

	for _, person := range managerToDirects[start.ID] {
		if _, manager := managerToDirects[person.ID]; manager {
			managers = append(managers, person)
		} else {
			ics = append(ics, person)
		}
	}

	allManagers := make([]*Person, 0)
	queue := append([]*Person{}, managers...)
	allIcs := append([]*Person{}, ics...)

	for len(queue) != 0 {
		var newQueue []*Person

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

	sort.Slice(managers, func(i, j int) bool { return managers[i].Name < managers[j].Name })
	sort.Slice(ics, func(i, j int) bool { return ics[i].Name < ics[j].Name })
	sort.Slice(allManagers, func(i, j int) bool { return allManagers[i].Name < allManagers[j].Name })
	sort.Slice(allIcs, func(i, j int) bool { return allIcs[i].Name < allIcs[j].Name })

	fmt.Println("Direct reports for:", start.Name)
	supervisor := idToPerson[start.ManagerID]
	fmt.Println("Supervisor:", supervisor.Name)
	fmt.Println("Number of directs:", len(managers)+len(ics))

	fmt.Println("Managers: ", PrintNames(managers))
	fmt.Println("Individual Contributors:", PrintNames(ics))
	fmt.Println("Entire Org Size:", len(allManagers)+len(allIcs))
	fmt.Println("Entire Org Manager Count:", len(allManagers))
	fmt.Println("Entire Org Individual Contributor Count:", len(allIcs))

	fmt.Println("Org-wide Managers: ", PrintNames(allManagers))
	fmt.Println("Org-wide Individual Contributors:", PrintNames(allIcs))
}

func PrintNames(people []*Person) string {
	if len(people) == 0 {
		return "N/A"
	}

	names := make([]string, 0, len(people))

	for _, person := range people {
		names = append(names, person.Name)
	}

	return strings.Join(names, ", ")
}
