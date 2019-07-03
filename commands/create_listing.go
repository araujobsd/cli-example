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

// +build !windows
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/araujobsd/cli-example/utils"
)

const (
	timeFormat = "02-01-2006-15:04PM"
)

func help() {
	fmt.Println("[create_listing] - Listing a new product")
}

func do(cmd []string) {
	var product utils.ProductListing

	if len(cmd) > 0 && cmd[0] == "-h" {
		help()
	}

	if len(cmd) < 5 {
		fmt.Println("Error - Some items missing")
		help()
	} else {
		product.Id = utils.LastProductId()
		product.Username = cmd[0]
		product.Title = cmd[1]
		product.Description = cmd[2]
		product.Price, _ = strconv.Atoi(cmd[3])
		product.Category = cmd[4]

		t := time.Now()
		product.CreatedAt = t.Format(timeFormat)

		err := utils.WriteCSVProduct(product)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(product.Id)
		}
	}
}

func main() {
	args := os.Args[1:]
	do(args)
}
