package main

// Process the build arguments and execute build
func gadgetDelete(args []string, g *GadgetContext) {
	loadConfig(g)
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	for _, onboot := range g.Config.Onboot {
		runRemoteCommand(client, "docker", "rmi", onboot.ImageAlias)
	}

	for _, service := range g.Config.Services {
		runRemoteCommand(client, "docker", "rmi", service.ImageAlias)
	}
}
