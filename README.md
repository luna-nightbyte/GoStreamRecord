# GoStreamRecord
 
## Core Features
- Pre-built docker images. Base and minimal image. See the [docker section](https://github.com/luna-nightbyte/GoStreamRecord?tab=readme-ov-file#docker) for example startup command.
### Recorder:
- Streamer status checks bypasses rate limits; recording rate limits still under testing.
- Start, stop, and restart all recordings.
- View active recorders in real-time
- Start, stop restart  individual recordings as needed.
- Add or delete streamers dynamically (no need to restart the recorder).
- Import/export streamer lists for easier management


### WebUI:
- Secure login using cookies for each client (prevents unauthorized access).
- ~~Manage api secret keys.~~
- ~~Manage multiple user accounts.~~
- Directly view ~~logs and~~ recorded videos through the WebUI.
- Watch live streams directly from the WebUI.
- Check streamer online status.

## Usage
**See the [release page](https://github.com/luna-nightbyte/GoStreamRecord/releases) to find pre-built binaries**
|Username|Password|
|-|-|
|`admin`|`password`|

### Prerequisites (app only)
This app depends on the following to be installed:
- `ffmpeg` - `sudo apt install ffmpeg`
- `ffprobe` - `sudo apt install ffprobe`
- `yt-dlp` - `sudo apt install yt-dlp`
- `curl` - `sudo apt install curl`

### Prerequisites (building the source)
This app depends on the following to be installed:
- `node` - `sudo apt install nodejs`
- `npm` - `sudo apt install npm`
- `golang` - See the [official](https://go.dev/doc/install) installation instructions
- `git` - `sudo apt install git`
- `curl` - `sudo apt install curl`

_Windows / Mac users will have to google or ask ChatGPT how to install these on their system._

__important__: You will still need to have the `settings` folder and it's content in the same folder structure when running this app. That means that you'll have to copy that along with any binary you build.

Download this repo and open a terminal in this folder:
```bash
user@hostname:~$ git clone https://github.com/luna-nightbyte/GoStreamRecord
user@hostname:~$ cd GoStreamRecord
user@hostname:~/GoStreamRecord$ make app
user@hostname:~/GoStreamRecord$ ./GoStreamRecord
```

#### Settings
The main settings can be found in [`settings.json`](https://github.com/luna-nightbyte/GoStreamRecord/blob/main/settings/settings.json):

__Notes__: 
- Cookie value should be re-created for a production app. But locally it works just fine.
- Google drive and Telegram haven't been fully tested.

```json
{
  "app": {
    "port": 8050,
    "loop_interval_in_minutes": 0,
    "output_folder": "videos",
    "rate_limit": {
      "enable": false,
      "time": 0
    },
    "cookie": " eqy\u0003!\ufffd\ufffd\ufffdW\ufffd{\u0014z\ufffdf\ufffdG\ufffd\u0012\ufffd\ufffd\ufffdb\u0011yDg.\ufffd\ufffd"
  },
  "google_drive": {
    "enabled": false,
    "path": ""
  },
  "telegram": {
    "chatID": "",
    "token": "",
    "enabled": false
  }
}
```
### Docker

#### Output files
**Note:**
The `output` folder path defined in the configuration file applies **only inside the Docker container**.
To ensure recorded files are saved correctly, the Docker volume path **must match** the folder specified in the configuration file.
  Example: 
  - `settings.json`:
    ```json
    "output_folder": "MyCustomFolder"
    ```
  - `docker-compose-yml`:
    ```docker-compose.yml
    volumes:
      - ./output:/app/MyCustomFolder`
    ```

The recorded files will in this example be available outside of docker in a folder called `output` (Or whatever else you call your output folder in the docker-compose file)
#### Docker image
There are two docker images available:
- [base](https://github.com/luna-nightbyte/GoRecord-WebUI/blob/main/docker/Dockerfile.base) (Full source code Ubuntu based image. Image size < 1GB )
- [run](https://github.com/luna-nightbyte/GoRecord-WebUI/blob/main/docker/Dockerfile.run) (Minimalistic image. Image size < 100MB )

#### Building images
Use the `docker-compose.yml` file to build images.
```bash
docker compose build base
docker compose build GoRecord
```
#### Startup

```bash
docker compose up GoRecord -d
# or
docker compose up dev -d
```

#### APP (Minimalistic image)
- [`settings`](https://github.com/luna-nightbyte/GoRecord-WebUI/tree/main/settings) save login, api and streamer lists.
- `output` folder for saving output videos.

App uses port __8050__ by default internally.
```bash
user@hostname:~$ docker run \
  -v ./settings:/app/settings \
  -v ./output:/app/videos \
  -v ./app.log:/app/remoteCtrl.log \ # Use this to access the logfile outside of docker
  -p 8080:8050 \
  docker.io/lunanightbyte/gorecord:latest
```


### Source
#### Build and run
Building the code wil create a binary for your os system. Although golang is [cross-compatible](https://go.dev/wiki/GccgoCrossCompilation) for windows, linux and mac, this app might not be fully compatible because of the difference in system patch and command execution. I haven't really tested it on anything else than Linux (armv6, armv7, x86_64)

##### Build binary: 
```bash
make vue # Only needed if frontend has been modified
make app

# Run the newly compiled binary:
./GoStreamRecord
```

##### Build & start :
```bash
make run
``` 

### Logs
Check the `app.log` to read any logs.

## Additional startup arguments
- `./GoStreamRecord reset-pwd admin MySecretPassword` - Resets password for the __admin__ user.
- `/GoStreamRecord add-user newUser newPassword` - Creates a new user and saves it to [users.json](https://github.com/luna-nightbyte/GoStreamRecord/blob/main/settings/users.json) 

## Other

### Todo

- ~~Select and delete videos~~
- Option for max video length (and size?)
- ~~headless mode without webui~~ (Abandoned because i will not create all the logic for handling the various arguments myself. Others can create a PR if they want to.)
- ~~Move frontend to Vue~~
  - ~~Better for organizing components being re-used~~
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
