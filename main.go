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
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/araujobsd/cli-example/utils"
)

const (
	thecarousell = "thecarousell# "
	banner       = "\t\t Find the best deals\n \t\t wherever you go :)\n"
)

var (
	prompt = utils.SetPrompt()
)

func runCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(strings.ToLower(commandStr), "\n")
	commandStr = strings.TrimSuffix(commandStr, "\n")
	rgx := regexp.MustCompile("'.*?'|\".*?\"|\\S+")
	argCommandStr := rgx.FindAllString(commandStr, -1)

	if len(argCommandStr) <= 0 {
		return nil
	}

	switch argCommandStr[0] {
	case "help":
	case "exit":
		if prompt != thecarousell {
			prompt = thecarousell
			return nil
		}
		os.Exit(0)
	case "su":
		if len(argCommandStr) < 2 {
			return errors.New("You need to specify the username")
		}

		if prompt != thecarousell {
			return errors.New("You need to be super user")
		}
		prompt = utils.SetPrompt(argCommandStr[1])
		return nil
	default:
		if len(argCommandStr) > 0 {
			cmd := []string{argCommandStr[0]}
			fcmd, err := utils.FindCmd(cmd)
			if err != nil {
				return err
			}

			runCmd := exec.Command(fcmd, argCommandStr[1:]...)
			runCmd.Stderr = os.Stderr
			runCmd.Stdout = os.Stdout

			err = runCmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// print banner
	banner_l := "\t\t " + strings.Repeat("-", 20)
	fmt.Fprintln(os.Stdout, banner_l)
	fmt.Fprintln(os.Stdout, banner)

	for {
		fmt.Print(prompt)

		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		err = runCommand(cmdString)
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
		}
	}
}
