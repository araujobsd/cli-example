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

package main

import (
	"fmt"
	"os"

	"github.com/araujobsd/cli-example/utils"
)

func help() {
	fmt.Println("[get_category] - Get a category of products")
}

func do(cmd []string) {
	var err error
	if len(cmd) > 0 && cmd[0] == "-h" {
		help()
	}

	if len(cmd) < 2 {
		fmt.Println("Error - You need to specify username and category")
		help()
	} else if len(cmd) >= 4 {
		if cmd[2] == "sort_price" || cmd[2] == "sort_time" &&
			cmd[3] == "dsc" || cmd[3] == "asc" {
			err = utils.GetCSVCategory(cmd[0], cmd[1], cmd[2], cmd[3])
		}
	} else {
		err = utils.GetCSVCategory(cmd[0], cmd[1])
	}

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	args := os.Args[1:]
	do(args)
}
