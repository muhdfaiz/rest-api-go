## Shoppermate API

### Development Enviroment

### Prerequisite
- Go Languange 1.7 above, MariaDB, Glide (Package Management For Go), GIT, supervisor

#### Install MariaDB 10.x
```
Refer here - https://www.linuxbabe.com/mariadb/install-mariadb-10-1-ubuntu14-04-15-10
```

#### Install Git
```
sudo apt-get update
sudo apt-get install git
```

#### Install Glide - Package Management for Go (https://github.com/Masterminds/glide)
- Install Glide

For Ubuntu
```
sudo add-apt-repository ppa:masterminds/glide && sudo apt-get update
sudo apt-get install glide
```

For Mac Os X
```
brew install glide
```

#### Install Go Language 1.7.x
- Download Go Language 1.7.x
```
sudo apt-get update
sudo wget https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz
```

- Extract Go Language 1.7.x
```
sudo tar -xvf go1.7.linux-amd64.tar.
sudo mv go /usr/local
```

#### Setup Go Environment.

Edit file `~/.profile` and include 3 environment variables below.

- Set GOROOT (location when Go package is installed on your system)
```
export GOROOT=/usr/local/go
```

- Set GOPATH. Location of your project path. For example 
```
export GOPATH=$HOME/golang
```

- Set PATH variable to access go binary system wide.
```
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

#### Verify Installation
- Check Go Version
```
go version
```

- Verify all environment variable. Make sure GOROOT and GOPATH not empty and set to the correct folder.
```
go env
```

#### Setting Up Shoppermate API
- Allow go get to retrieve shoppermate API from private bitbucket repositories. Enter code below in command line.

```
git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"
```

Here, we say use git@bitbucket.org any time you’d use https://bitbucket.org. This works for everything, not just go get. It just has the nice side effect of using your SSH key any time you run go get too.

- Verify .gitconfig file contain information like below
```
[url "git@bitbucket.org:"]
        insteadOf = https://bitbucket.org/
```

- Generate Your SSH Key. Paste code below in your command line and press enter until finish.
```
ssh-keygen
```

- Copy the ssh key out from the command below and add the SSH key in Bitbucket Repository.
```
cat ~/.ssh/id_rsa.pub
```

- Go to shoppermate-api project path. For example `~/golang/src/bitbucket.org/cliqers/shoppermate-api` and then install package dependencies using Glide.
Glide will install another package dependencies into `~/golang/src/bitbucket.org/cliqers/shoppermate-api/vendor` folder
```
glide install
```

- Create new .env file and copy the content from .env.example file in root directory. Update all setting in .env file.

- Go to project root directory and run Shoppermate API.
```
go run api.go
```
