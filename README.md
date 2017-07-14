# Gadgetcli Alpha

### Foreword:
This utility is currently in an alpha release state. It's completely capable of being used as a tool to build and orchestrate Docker containers on Nextthing's C.H.I.P. Pro. However, as it's only been used internally, there's a likelihood of bugs and pain points in the developer workflow. That's where you come in, we'd love to get your feedback!

So what works, and what doesn't? Follow the steps below to get a quick few examples on what works and how to do it.

As far as what doesn't work, the Gadget command-line utility currently doesn't support, but has milestones for the following:
 - Pulling images from private Docker registries [images can only be pulled from Dockerhub, or built from a local Dockerfile]
 - Porcelain commands for editing all fields in the gadget.yml [e.g. `gadget add-capability SYS_RAWIO mycontainer2`]
 - Local over the air[OTA] update testing
 - Preventing the user from filling the NAND
 - WiFi provisioning
 - Dependent commands [you must explicitly build a container before you can deploy it]
 - Local image creation/backups
 - Support for C.H.I.P.
 - Multistage Docker image builds

Setup:
 - Download the .CHP and the gadget.zip.
 - Flash your C.H.I.P. Pro, and either put the `gadget` binary in your
   path or precede all the commands with `./` eg. `./gadget`.
   [in public releases, there will be an installer for that]

Setup Notes:
 - The .CHP has a bug where the C.H.I.P. Pro won't immediately boot,
   you'll have to power cycle.
 - Targeted/tested platforms: Windows 10 Pro 64-bit, Mac OSX [undefined versions], Linux 64-bit.
 - Tested with Docker version 1.12 onward, but will likely put a hard limit on 17.05+ once multistage builds reach Docker's `stable` channel. 
 - For further documentation, check out our [Gadget docs page](https://docs.getchip.com/gadget.html).

## Gadget Command-Line Arguments:
 `gadget init`
   - This will create a template `gadget.yml` file that is ready to
     build/deploy/use
 
 `gadget build`
   - This reads the `gadget.yml`, and downloads/builds a new container, 
     from the specified image, the container is unique to your project.
 
 `gadget deploy`
   - This pushes over the `hello-world` container
 
 `gadget start`
   - This starts the `hello-world` container
 
 `gadget logs`
   - Prints the output of the `hello-world` container
 
 `gadget status`
   - Shows the status of running containers
 
 `gadget stop`
   - Stops running containers
 
 `gadget delete`
   - Deletes running containers


fin.

