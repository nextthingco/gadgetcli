/*
This file is part of the Gadget command-line tools.
Copyright (C) 2017 Next Thing Co.

Gadget is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

Gadget is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Gadget.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"github.com/nextthingco/libgadget"
	log "github.com/sirupsen/logrus"
	"strings"
)

func GadgetRun(args []string, g *libgadget.GadgetContext) error {

	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)
	if err != nil {
		return err
	}

	stdout, stderr, err := libgadget.RunRemoteCommand(client, strings.Join(args, " "))

	if err != nil {
		log.Errorf("\n%s", stdout)
		log.Errorf("\n%s", stderr)
		return err
	}

	log.Infof("\n%s", stdout)
	log.Debugf("\n%s", stderr)

	return err
}
