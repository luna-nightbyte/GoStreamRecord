# GoStreamRecord
<p align="center">
 <img src="https://github.com/user-attachments/assets/1eb7fd04-e421-47f6-bff0-fc75fdfd6f21" alt="Login page"/>
</p>

__API NOTE__: The API is still in early development. I've added checks for login/API, but i would not recommend exposing the port for this app on your router as of now. 
## Core Features
- Pre-built docker images. Base and minimal image. See the [docker section](https://github.com/luna-nightbyte/GoStreamRecord?tab=readme-ov-file#docker) for example startup command.
### Recorder:
- Streamer status checks bypasses rate limits; recording rate limits still under testing.
- Start, stop, and restart all recordings.
- View active recorders in real-time (More data will be added)
- Start, stop restart  individual recordings as needed.
- Add or delete streamers dynamically (no need to restart the recorder).
- Import/export streamer lists for easier management


### WebUI:
- Secure login using cookies for each client (prevents unauthorized access).
- Manage api secret keys.
- Manage multiple user accounts.
- Directly view logs and recorded videos through the WebUI.
- Watch live streams directly from the WebUI.
- Check streamer online status.
### Setup & Deployment:
- Docker configuration and service examples for straightforward deployment.
- Install [Golang](https://go.dev/doc/install) to run the source or build binary.
## Usage
|Username|Password|
|-|-|
|`admin`|see [this](https://github.com/luna-nightbyte/Recordurbate-WebUI/tree/main?tab=readme-ov-file#reset-password)|

### Setup
__important__: You will still need to have the `internal/settings` folder and it's content in the same folder structure when running this app. That means that you'll have to copy that along with any binary you build.

- Download this repo and open a terminal in this folder. Ask ChatGPT how to find the folder path and how to move into it via cli if you dont know.

#### Optional config settings
The main settings can be found in [`settings.json`](https://github.com/luna-nightbyte/Recordurbate-WebUI/blob/main/internal/app/db/settings/settings.json):
```json
{
  "app": {
    "port": 8055,
    "loop_interval_in_minutes": 2,
    "video_output_folder": "output/videos",
    "rate_limit": {
      "enable": true,
      "time": 5
    },
    "default_export_location": "./output/list.txt"

  },
  "youtube-dl": {
    "binary": "youtube-dl"
  },
  "auto_reload_config": true
}
```
#### Reset password
To change forgotten password, start the program with the `reset-pwd` argument. I.e:
```
./GoRecordurbate reset-pwd admin newpassword 
```
New login for the user `admin` would then be `newpassword`

### Docker

There is two docker images available:
- [base](https://github.com/luna-nightbyte/GoRecord-WebUI/blob/main/docker/Dockerfile.base) (Full source code Ubuntu based image. Image size > 1.5GB )
- [run](https://github.com/luna-nightbyte/GoRecord-WebUI/blob/main/docker/Dockerfile.run) (Minimalistic image. Image size < 500MB )

#### Building images
Use the Makefile to build images to ensure proper tagging.
```bash
make build # Builds all

# or 

make base # Only base

# or 

make app # Only app
```

#### docker-compose.yml

##### Usage

```bash
docker compose up GoRecord -d

# or

docker compose up dev -d

```

##### Logs

Docker logs can be found using `docker logs --tail 200 -f CONTAINER_NAME`. 

#### APP (Minimalistic image)
Files / folders needed to save app settings is (only need env file to just test the container):
- [`settings`](https://github.com/luna-nightbyte/GoRecord-WebUI/tree/main/internal/app/settings) save login, api and streamer lists.
- `output` folder for saving output videos.

App uses port __80__ by default internally.
```bash
user@hostname:~$ docker run \
  -v ./internal/db:/app/internal/db \
  -p 8050:80 \
  docker.io/lunanightbyte/gorecord:latest
```

#### Ubuntu based image


```bash
user@hostname:~$ docker run \
  -v ./:/app  \
  -p 8050:80 \
  docker.io/lunanightbyte/gorecord-base:latest
```


### Source
#### Build
Building the code wil create a binary for your os system. Golang is [cross-compatible](https://go.dev/wiki/GccgoCrossCompilation) for windows, linux and mac.
```bash
go mod init GoStreamRecord # Only run this line once
go mod tidy
go build
./GoStreamRecord #windows will have 'GoStreamRecord.exe'
```
#### Source
```bash
go mod init GoStreamRecord # Only run this line once
go mod tidy
go run main.go
```

## WebUI (v0.1.x)


<p align="center">
  <img src="https://github.com/user-attachments/assets/edf30517-de6a-4f91-9ab4-89f9c91d7779" alt="Login page"/>
  <img src="https://github.com/user-attachments/assets/5d939bc0-778b-42c8-a453-eb30c13e95e2" alt="Video tab"/>
  <img src="https://github.com/user-attachments/assets/0ce5b2c1-e7f3-47bb-96e9-1532915dd5e4" alt="individual tab"/>
  <img src="https://github.com/user-attachments/assets/7736fac5-5ce8-4634-8179-6ea2cf03969b" alt="User settings tab"/>
  
  <img src="https://github.com/user-attachments/assets/ced11119-8e74-4c15-8aff-6c31242f8fe5" alt="Streamers tab"/>
  <img src="https://github.com/user-attachments/assets/edc136e5-0238-463e-b8f3-d4b1b7e74687" alt="Livestream tab"/>
</p>

_Online status with a small bug at the time of uploading this.._

## Make command list

__Note:__ This is intended to be used together with the source files.

- `make reset-pwd USERNAME=admin PASSWORD=MySecretPassword` - Resets password for the __admin__ user.
- `make app` - Builds and starts the app within the output folder
- `make build-and-push` - Builds all docker images and pushes them
- `make build-base` - Builds [base](https://github.com/luna-nightbyte/GoRecord-WebUI/blob/main/docker/Dockerfile.base) image
- `make build-app` - Builds [app](https://github.com/luna-nightbyte/GoRecord-WebUI/blob/main/docker/Dockerfile.run)  image
- `make push-base` - Pushes [base](https://github.com/luna-nightbyte/GoRecord-WebUI/blob/main/docker/Dockerfile.base) image
- `make push-app` - Pushes [app](https://github.com/luna-nightbyte/GoRecord-WebUI/blob/main/docker/Dockerfile.run)  image

## Other

### Todo

- ~~Select and delete videos~~
- Option for max video length (and size?)
- ~~headless mode without webui~~ (Abandoned because i will not create all the logic for handling the various arguments myself. Others can create a PR if they want to.)
- Move frontend to Vue
  - Btter for organizing components being re-used
- ~~Build a default docker image~~
- Individual recorders in UI
  - ~~Start/Restart individual recorders (in progress)~~
  - set max lenght/size (could be optional to use one of either)
  - view current recording length
- ~~Better video view~~
- Add support for multiple websites
### Ideas (not planned
- Log online-time of streamers and save to csv for graph plotting. Can help understand the work-hours of different streamers.
- Option to login to the streaming site and use follower list instead of config? (Unsure)
- View current videos in progress.
- Add option to try and use a custom url.
### Disclaimer 
Unauthorized resale, redistribution, or sharing of recorded content that you do not own or have explicit permission to distribute is strictly prohibited. Users are solely responsible for ensuring compliance with all applicable copyright and privacy laws. The creator of this recorder assumes no liability for any misuse or legal consequences arising from user actions.

## Thanks

Special thanks to [oliverjrose99](https://github.com/oliverjrose99) for the initial inspiration and their work on [Recordurbate](https://github.com/oliverjrose99/Recordurbate). Initial code of this project was directly inspired by their project.
