// +build !windows

// SPDX-License-Identifier: BSD-2-Clause
/*-
 * Copyright 2019 by Marcelo Araujo <araujo@FreeBSD.org>
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted providing that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 */

package utils

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	csvUserPath  = "/tmp/users.csv"
	csvItemsPath = "/tmp/items.csv"
	timeFormat   = "02-01-2006-15:04PM"
)

var (
	ErrPAE  = errors.New("Error - Product already exist")
	ErrUNKU = errors.New("Error - Unknow user")
	ErrPLE  = errors.New("Error - Product list is empty")
)

// ProductListing - Structure used to organize the item
type ProductListing struct {
	Id          int
	Username    string
	Title       string
	Description string
	Price       int
	Category    string
	CreatedAt   string
}

func trimQuotes(word string) string {
	word = strings.Replace(word, "'", "", -1)
	word = strings.Replace(word, `"`, "", -1)

	return word
}

func sortMap(data map[int]string, args ...string) {
	product := []ProductListing{}

	for _, v := range data {
		splEntry := strings.Split(v, "|")
		_price, _ := strconv.Atoi(splEntry[2])
		n := ProductListing{Title: splEntry[0], Description: splEntry[1],
			Price: _price, CreatedAt: splEntry[3]}
		product = append(product, n)
	}

	if len(args) >= 2 {
		if args[0] == "sort_price" && args[1] == "dsc" {
			sort.SliceStable(product, func(i, j int) bool {
				return product[i].Price < product[j].Price
			})
		}

		if args[0] == "sort_price" && args[1] == "asc" {
			sort.SliceStable(product, func(i, j int) bool {
				return product[i].Price > product[j].Price
			})
		}

		if args[0] == "sort_time" && args[1] == "dsc" {
			sort.SliceStable(product, func(i, j int) bool {
				ti, _ := time.Parse(timeFormat, product[i].CreatedAt)
				tj, _ := time.Parse(timeFormat, product[j].CreatedAt)
				return ti.Before(tj)
			})
		}
		if args[0] == "sort_time" && args[1] == "asc" {
			sort.SliceStable(product, func(i, j int) bool {
				ti, _ := time.Parse(timeFormat, product[i].CreatedAt)
				tj, _ := time.Parse(timeFormat, product[j].CreatedAt)
				return ti.After(tj)
			})
		}
	}

	for i := 0; i < len(product); i++ {
		fmt.Println(string(fmt.Sprintf("%s|%s|%d|%s",
			product[i].Title, product[i].Description,
			product[i].Price, product[i].CreatedAt)))
	}
}

// DoesProductExist - Verify if a product exist
func DoesProductExist(product ProductListing) bool {
	entries := ReadCSVProduct()
	if len(entries) == 0 {
		return false
	}

	for _, entry := range entries {
		splEntry := strings.Split(entry[0], "|")
		title := splEntry[2]
		desc := splEntry[3]
		cat := splEntry[5]
		if trimQuotes(product.Title) == trimQuotes(title) &&
			trimQuotes(product.Description) == trimQuotes(desc) &&
			trimQuotes(product.Category) == trimQuotes(cat) {

			return true
		}
	}

	return false
}

// LastProductId - Gets the latest product ID
func LastProductId() int {
	var lastID int = 1
	file, err := os.Open(csvItemsPath)
	if err != nil {
		return lastID
	}
	defer file.Close()

	csv := csv.NewReader(file)
	csv.LazyQuotes = false
	for {
		record, err := csv.Read()
		if err == io.EOF {
			break
		}
		lastID, _ = strconv.Atoi(strings.Split(record[0], "|")[0])
	}

	return lastID + 1
}

// IsUsernameExist - Check if user exist
func IsUsernameExist(username string) bool {
	file, err := os.Open(csvUserPath)
	if err != nil {
		fmt.Println("Error - csv users file does not exist")
		return false
	}

	csv := bufio.NewScanner(file)
	for csv.Scan() {
		if username == csv.Text() {
			return true
		}
	}

	return false
}

// WriteCSVProduct - Writes the item into the csv file
func WriteCSVProduct(product ProductListing) error {
	// Check if product already exist
	if DoesProductExist(product) {
		return ErrPAE
	}

	// Check if user exist
	if !IsUsernameExist(trimQuotes(product.Username)) {
		return ErrUNKU
	}

	file, err := os.OpenFile(csvItemsPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	wr := csv.NewWriter(file)
	defer wr.Flush()

	item := []string{fmt.Sprintf("%d|%s|%s|%s|%d|%s|%s",
		product.Id,
		trimQuotes(product.Username),
		trimQuotes(product.Title),
		trimQuotes(product.Description),
		product.Price,
		trimQuotes(product.Category),
		trimQuotes(product.CreatedAt))}

	res := wr.Write(item)

	return res
}

// ReGenerateCSVProduct - Regenerates the csv item file
//                        for delete/update operations
func ReGenerateCSVProduct(lines [][]string) error {
	file, err := os.Create(csvItemsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	for _, line := range lines {
		if len(line[0]) > 0 {
			err = w.Write(line)
			if err != nil {
				return err
			}
		}
	}
	w.Flush()
	return nil
}

// ReadCSVProduct - Read all items from a csv item file
func ReadCSVProduct() [][]string {
	file, err := os.OpenFile(csvItemsPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil
	}
	r := csv.NewReader(file)
	r.LazyQuotes = true
	lines, err := r.ReadAll()
	if err != nil {
		return nil
	}

	if len(lines) == 0 {
		return nil
	}

	return lines
}

// DeleteCSVItem - Remove an item from csv item file
func DeleteCSVItem(username string, id int) (err error) {
	entries := ReadCSVProduct()
	if len(entries) == 0 {
		return ErrPLE
	}

	for index, entry := range entries {
		splEntry := strings.Split(entry[0], "|")
		_id, _ := strconv.Atoi(splEntry[0])
		_username := splEntry[1]
		if id == _id {
			if trimQuotes(username) != trimQuotes(_username) {
				err = errors.New("Error - listing owner mismatch")
				break
			} else {
				entries[index][0] = ""
				err = nil
				break
			}
		} else {
			err = errors.New("Error - listing does not exist")
		}
	}

	if err != nil {
		return err
	}

	err = ReGenerateCSVProduct(entries)
	if err != nil {
		return err
	}

	return errors.New("Success")
}

// GetCSVItem - Find and return an item from csv item file
func GetCSVItem(username string, id int) (err error) {
	entries := ReadCSVProduct()
	if len(entries) == 0 {
		return ErrPLE
	}

	for index, entry := range entries {
		splEntry := strings.Split(entry[0], "|")
		_id, _ := strconv.Atoi(splEntry[0])
		_username := splEntry[1]

		if id == _id {
			if trimQuotes(username) != trimQuotes(_username) {
				err = ErrUNKU
				break
			} else {
				a := strings.Split(entries[index][0], "|")
				item := string(fmt.Sprintf("%s|%s|%s|%s|%s|%s",
					a[2], a[3], a[4], a[6], a[5], a[1]))

				return errors.New(item)
			}
		} else {
			err = errors.New("Error - not found")
		}
	}

	if err != nil {
		return err
	}

	return errors.New("Success")
}

// GetCSVTopCategory - Show the top category with most items
func GetCSVTopCategory(username string) (err error) {
	var topCategory string

	top := make(map[string]int)
	entries := ReadCSVProduct()
	if len(entries) == 0 {
		return ErrPLE
	}

	for _, entry := range entries {
		splEntry := strings.Split(entry[0], "|")
		_username := strings.ToLower(splEntry[1])
		if trimQuotes(strings.ToLower(username)) == trimQuotes(_username) {
			category := strings.ToLower(splEntry[5])
			top[category] = top[category] + 1
		} else {
			err = errors.New("Error - unknown user")
		}
	}

	if err != nil && len(top) == 0 {
		return err
	}
	max := 0
	for k, v := range top {
		if v >= max {
			max = v
			topCategory = k
		}
	}
	fmt.Println(topCategory)

	return nil
}

// GetCSVCategory - Show items from a follow category
func GetCSVCategory(username string, category string, args ...string) (err error) {
	allitems := make(map[int]string)
	entries := ReadCSVProduct()
	if len(entries) == 0 {
		return ErrPLE
	}

	for index, entry := range entries {
		splEntry := strings.Split(entry[0], "|")
		_username := strings.ToLower(splEntry[1])
		_category := strings.ToLower(splEntry[5])

		if trimQuotes(strings.ToLower(username)) == trimQuotes(_username) {
			if trimQuotes(strings.ToLower(category)) == trimQuotes(_category) {
				a := strings.Split(entries[index][0], "|")
				item := string(fmt.Sprintf("%s|%s|%s|%s",
					a[2], a[3], a[4], a[6]))
				allitems[index] = item
			} else {
				err = errors.New("Error - category not found")
			}
		} else {
			if err == nil {
				err = ErrUNKU
			}
		}
	}

	if len(args) >= 2 {
		if args[0] == "sort_price" && args[1] == "dsc" {
			sortMap(allitems, "sort_price", "dsc")
		}
		if args[0] == "sort_price" && args[1] == "asc" {
			sortMap(allitems, "sort_price", "asc")
		}
		if args[0] == "sort_time" && args[1] == "dsc" {
			sortMap(allitems, "sort_time", "dsc")
		}
		if args[0] == "sort_time" && args[1] == "asc" {
			sortMap(allitems, "sort_time", "asc")
		}
	} else {
		for _, v := range allitems {
			fmt.Println(v)
		}
	}

	if len(allitems) > 0 {
		return nil
	}

	return err
}

// WriteCSVUser - Write username into csv user file
func WriteCSVUser(username string) error {
	file, err := os.OpenFile(csvUserPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error - csv users file does not exist")
		return err
	}
	defer file.Close()

	if IsUsernameExist(username) {
		return errors.New("EUSERS")
	}

	wr := csv.NewWriter(file)
	defer wr.Flush()

	_username := strings.Fields(username)
	res := wr.Write(_username)

	return res
}

// UpdateCSVItem - Find and update an item from csv item file
func UpdateCSVItem(username string, id int, args []string) (err error) {
	var newproduct ProductListing
	entries := ReadCSVProduct()
	if len(entries) == 0 {
		return errors.New("Warning - Product list is empty")
	}

	for index, entry := range entries {
		splEntry := strings.Split(entry[0], "|")
		_id, _ := strconv.Atoi(splEntry[0])
		_username := splEntry[1]

		if id == _id {
			if trimQuotes(username) != trimQuotes(_username) {
				err = errors.New("Error - unknow user")
				break
			} else {
				a := strings.Split(entries[index][0], "|")
				_price, _ := strconv.Atoi(a[4])
				switch {
				case len(args) == 3:
					newproduct = ProductListing{Id: id, Username: username,
						Title: args[2], Description: a[3],
						Price: _price, Category: a[5],
						CreatedAt: a[6]}
				case len(args) == 4:
					newproduct = ProductListing{Id: id, Username: username,
						Title: args[2], Description: args[3],
						Price: _price, Category: a[5],
						CreatedAt: a[6]}
				case len(args) == 5:
					_price, _ := strconv.Atoi(args[4])
					newproduct = ProductListing{Id: id, Username: username,
						Title: args[2], Description: args[3],
						Price: _price, Category: a[5],
						CreatedAt: a[6]}
				case len(args) == 6:
					_price, _ := strconv.Atoi(args[4])
					newproduct = ProductListing{Id: id, Username: username,
						Title: args[2], Description: args[3],
						Price: _price, Category: args[5],
						CreatedAt: a[6]}
				}
			}
		} else {
			err = errors.New("Error - item not found for update")
		}
	}

	if err != nil {
		return err
	} else {
		_ = DeleteCSVItem(username, id)
		_ = WriteCSVProduct(newproduct)
		err = errors.New("Item updated")
	}

	return err
}
