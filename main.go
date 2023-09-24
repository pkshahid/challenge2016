package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// Structure to hold permission data
type Permission struct {
	Include []string
	Exclude []string
}

// Load permissions from permissions CSV
func loadPermissions(filePath string) (map[string]Permission, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read() // skip header row

	permissions := make(map[string]Permission)
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}
		distributor := line[0]
		action := line[1]
		region := line[2]

		p := permissions[distributor]

		if action == "INCLUDE" {
			p.Include = append(p.Include, region)
		} else if action == "EXCLUDE" {
			p.Exclude = append(p.Exclude, region)
		}

		permissions[distributor] = p

	}
	return permissions, nil
}

// Check if a distributor has permission for a given region
func hasPermission(permissions map[string]Permission, distributor, region string) bool {
	p, exists := permissions[distributor]
	if !exists {
		return false
	}

	for _, excluded := range p.Exclude {
		if strings.HasSuffix(region, excluded) {
			return false
		}
	}

	for _, included := range p.Include {
		if strings.HasSuffix(region, included) {
			return true
		}
	}
	fmt.Printf("No any permissions found for the distributor %s for the region %s.\n", distributor, region)
	return false
}

// Read Initial Options
func readOption() int {
	var option int
	fmt.Println("Available Options :")
	fmt.Println("1. Check Distribution Permission.")
	fmt.Println("2. Add Permission.")
	fmt.Println("3. Add Sub-distributor")
	fmt.Print("Choose any option: ")
	fmt.Scanln(&option)
	return option
}

func main() {

	// Load Cities,States and Countries
	cityDataArray := make([]string, 0)

	// Read Data CSV file
	csvfile, err := os.Open("cities.csv")
	if err != nil {
		panic(err)
	}

	// Close CSV file at the end
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.Read() // skip header row

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		countryCode := strings.ToUpper(row[5])
		stateCode := strings.ToUpper(row[4])
		city := strings.ToUpper(row[3])

		cityDataArray = append(cityDataArray, fmt.Sprintf("%s-%s-%s", city, stateCode, countryCode))

	}

	//Select Option
	for {
		permissions, err := loadPermissions("permissions.csv")
		if err != nil {
			fmt.Println("Error loading permissions:", err)
			return
		}
		option := readOption()
		switch option {
		case 1:
			// Check Permission
			var regionToCheck, distributor string
			fmt.Print("Enter a distributor name : ")
			fmt.Scanln(&distributor)
			fmt.Print("Enter region : ")
			fmt.Scanln(&regionToCheck)

			regionToCheck = strings.ToUpper(regionToCheck)
			distributor = strings.ToUpper(distributor)
			selectedCity := ""
			for _, city_string := range cityDataArray {
				if strings.HasPrefix(city_string, regionToCheck) {
					selectedCity = city_string
				}
			}

			if len(selectedCity) == 0 {
				fmt.Println("No Matching City found in dataset.")
				continue
			}

			if hasPermission(permissions, distributor, selectedCity) {
				fmt.Println("YES")
			} else {
				fmt.Println("NO")
			}
		case 2:
			// Add Permission
			var regionToAdd, distributorToAdd, permissionToAdd string
			fmt.Print("Enter distributor name : ")
			fmt.Scanln(&distributorToAdd)
			fmt.Print("Enter permission (INCLUDE/EXCLUDE): ")
			fmt.Scanln(&permissionToAdd)
			fmt.Print("Enter region : ")
			fmt.Scanln(&regionToAdd)

			distributorToAdd = strings.ToUpper(distributorToAdd)
			permissionToAdd = strings.ToUpper(permissionToAdd)
			regionToAdd = strings.ToUpper(regionToAdd)

			// Permission data to be appended to the CSV file
			newPermission := fmt.Sprintf("%s,%s,%s\n", distributorToAdd, permissionToAdd, regionToAdd)

			// Open or create the file for appending
			file, err := os.OpenFile("permissions.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			// Append data to the file
			_, err = file.WriteString(newPermission)
			if err != nil {
				panic(err)
			}

			// Output success message
			fmt.Println("New Permission Added.")

		case 3:
			// Add Sub-distributor
			var subDistributor, parentDistributor string
			fmt.Print("Enter Sub-Distributor Name : ")
			fmt.Scanln(&subDistributor)
			fmt.Print("Enter Upline Distributor Name: ")
			fmt.Scanln(&parentDistributor)

			subDistributor = strings.ToUpper(subDistributor)
			parentDistributor = strings.ToUpper(parentDistributor)

			// Permission data to be appended to the CSV file
			distributorData := fmt.Sprintf("%s,%s\n", subDistributor, parentDistributor)

			// Open or create the file for appending
			file, err := os.OpenFile("distributors.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			// Append data to the file
			_, err = file.WriteString(distributorData)
			if err != nil {
				panic(err)
			}

			// Add Permissions of Upline distributors
			p, exists := permissions[parentDistributor]
			if exists {
				appendLine := ""
				for _, excluded := range p.Exclude {
					appendLine += fmt.Sprintf("%s,EXCLUDE,%s\n", subDistributor, excluded)
				}

				for _, included := range p.Include {
					appendLine += fmt.Sprintf("%s,INCLUDE,%s\n", subDistributor, included)
				}

				if len(appendLine) > 0 {
					permissionFile, err := os.OpenFile("permissions.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
					if err != nil {
						panic(err)
					}
					defer permissionFile.Close()

					_, err = permissionFile.WriteString(appendLine)
					if err != nil {
						panic(err)
					}
				}
			}

			fmt.Println("New Sub-Distributor Added")

		default:
			fmt.Println("Invalid Option.")
		}
	}
}
