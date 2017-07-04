# Gadgetcli internal demo

### Foreword:
This utility is currently in an alpha release state. It's completely capable of being used as a tool to build and orchestrate docker containers on Nextthing's Chip Pro. However, as it's only been used internally, there's a likelyhood of bugs and pain points in the developer work flow. That's where you come in, we'd love to get your feedback!

So what works, and what doesn't? Follow the steps below to get a quick few examples on what works and how to do it.

As far as what doesn't work, the gadget command line utility currently doesn't support, but has milestones for the following:
 - Pulling images from private Docker registries [images can only be pulled from Dockerhub, or built from a local Dockerfile]
 - Porcelain commands for editing all fields in the gadget.yml [e.g. `gadget add-capability SYS_RAWIO mycontainer2`]
 - Local over the air[OTA] update testing
 - Preventing the user from filling the NAND
 - WiFi provisioning
 - Dependent commands [you must explicitly build a container before you can deploy it]
 - Local image creation/backups
 - Support for CHIP
 - Multistage Docker image builds

Setup:
 - download the .chp and the gadget.zip
 - flash your chippro, and either put the `gadget` binary in your
   path or precede all the commands with `./` eg. `./gadget`
   [in public releases, there will be an installer for that]

Setup Notes:
 - The .chp has a bug where the chippro won't immediately boot,
   you'll have to power cycle
 - Targeted/tested platforms: Windows 10 Pro 64-bit, Mac OSX [undefined versions], Linux 64-bit
 - Tested with Docker version 1.12 onward, but will likely put a hard limit on 17.05+ once multistage builds reach Docker's `stable` channel. 

## Demo_0 [first use]:
 `./gadget init`
   - This will create a template `gadget.yml` file that is ready to
     build/deploy/use
 
 `./gadget build`
   - This reads the `gadget.yml`, and downloads/builds a new container, 
     from the specified image, the container is unique to your project.
 
 `./gadget deploy`
   - This pushes over the `hello-world` container
 
 `./gadget start`
   - This starts the `hello-world` container
 
 `./gadget logs`
   - Prints the output of the `hello-world` container
 
 `./gadget status`
   - Shows the status of running containers
 
 `./gadget stop`
   - Stops running containers
 
 `./gadget delete`
   - Deletes running containers

## Demo_1 [seeking an example]:
 `./gadget add service blink`
   - creates a new entry under `services:`
   - gives it the name "blink"
   - generates new uuid
   - generates image name "$parent_directory/$name"
   
     NOTES:
       - containers start from top to bottom listing in `gadget.yml`
       - `onboot` containers run once, and when they exit, stay exited
         until reboot.
       - `service` containers get restarted if they ever exit with a 
         non-zero code.
   
 `nano gadget.yml`
   - change `image` to `pushreset/pythonio:v2`
   - change `command` from `[]` to `['python', 'blink.py']`
   - change `binds` from `[]` to `['/sys:/sys']`
   - change `capabilities` from `[]` to `['SYS_RAWIO']`
   - change `devices` from `[]` to `['/dev/mem']`
     
     NOTES:
       - image is a container that will be pulled from Dockerhub
       - command is what gets run automatically upon `./gadget start`
         [or upon reboot]
       - binds are directories that get mounted from the host, into the
         running container. The first argument is what to mount from 
         the host[gadget], and the second argument is where to mount it 
         in the running container.
       - capabilites are special permissions you can give to the
         container. In our case here, we need SYS_RAWIO to access the
         /sys files. Capabilites are for saving yourself from having to
         run --privileged containers. --privileged is a shortcut for 
         "give me every single capability". This is insecure and bad
         practice. A list of all capabilities can be found here:
         http://man7.org/linux/man-pages/man7/capabilities.7.html
       - devices are raw devices that can pass through to the container.
         These are different from binds, because Linux devices have 
         several different modes of access.
 
 `./gadget build`
   - the new entry is spotted, the image is pulled from Dockerhub, 
     tagged with the uuid.
 
 `./gadget deploy`
   - both containers are pushed over
 
 `./gadget start`
   - hello-world container starts
   - blink container starts
 
 `./gadget status`
   - shows status of both containers
 
 `./gadget logs`
   - shows output of hello-world
   - shows empty output of blink container
 
 NOTES:
   - At this point you should be able to power cycle your chippro,
     and the hello-world container will start on the next boot, followed
     by the blink container.
   - If for some reason, the hello-world container exited 1, it would
     stay dead.
   - If the blink container exited 1, it would be restarted.
 

## Demo_2 [modifying an example]:
 I don't have the actual resources to do this, but here's the
 scenario I imagine:
 
 - user performs `git clone https://github.com/pushreset/pythonio`
 - user changes `gadget.yml` blink container section :
   - change `image` from `pushreset/pythonio:v2` to `""` [it's empty]
   - change `directory` from `""` to `pythonio`
     NOTES:
       - now when you `./gadget build` it will build `blink` from the
         Dockerfile in the `pythonio` directory, instead of pulling
         from dockerhub.
   - user edits the python script [to make it blink more rapidly],
     which gets added to the container in `blink`'s `./gadget build` 
     step.
 - user runs `./gadget build blink`
 - user runs `./gadget deploy blink`
 - user runs `./gadget start blink`
 - user sees the LED blink more rapidly, bitchin'
 
 NOTES:
   - specifying a container name after `./gadget <command>` will only
     perform that action for the specified container. In our case,
     it means we don't have to rebuild/redeploy the hello-world
     container.
     

fin.

