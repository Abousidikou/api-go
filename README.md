# api-go
  api-go est une API de visualisation de BjNet Monitor
 
# Mise en place
  
  ## Installer Golang
    Visit [here](https://tecadmin.net/install-go-on-ubuntu/) to install golang
  
  ## Clone repository
  ```bash
  git clone https://github.com/Abousidikou/api-go.git && cd api-go
  ```
  
  ## Initialize go module
  ```bash
  go mod init api-go
  ```
  
  ## Download missing package
  ```bash
  go mod tidy
  ```
  
   ## Make symbolic link to certificate and privkey
   ```bash
   ln -s /etc/letsencrypt/live/emes.bj/fullchain.pem  fullchain.pem
   ln -s /etc/letsencrypt/live/emes.bj/privkey.pem  privkey.pem
   ```
   NB: Replace  the path with your own path
   
  ## Create service for restarting service if failed
  Create file into systemd
  ```bash
  sudo nano /etc/systemd/system/api-go.service
  ```
  Fill this file with this:
  ```bash
  [Unit]
  Description = API-GO 
  After = network.target

  [Service]
  Environment="GOPATH=/home/emes/go"
  Environment="GOROOT=/usr/local/go"
  Environment="GOCACHE=/home/emes/.cache/go-build"
  Environment="PATH=/home/emes/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin"
  WorkingDirectory = /home/emes/ndt/api-go/
  ExecStart = /usr/local/go/bin/go run api.go
  Restart=always
  RestartSec=3

  [Install]
  WantedBy = multi-user.target
  ```
  
  ## Run api
  ```bash
  sudo systemctl daemon-reload
  sudo systemctl start api-go.service
  ```

  ## Verify that api is running
  ```bash
  sudo systemctl status api-go.service
  ```
