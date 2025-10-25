# GoStreamRecord
 __NOTE__: v0.3.5 and up uses SQL database. The readme and usage will be updated on the next release (not pre-releases).
__Always use the latest release (nit pre-release) and not the master branch__

## Core Features
### Recorder:
- Streamer status checks bypasses rate limits. _Recording ratelimit can still occur if recordings are restarted too often in a short period of time._
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

### To be implemented
- Google drive
- Telegram bot
- User groups: Videos and streamers only visible to owners and/or members of common groups.
- Role based access to web functions like gallery, download tool, recording tool, settings and so on
## Usage


**Always download the latest release from the [release page](https://github.com/luna-nightbyte/GoStreamRecord/releases) as this is ensured to have the files needed. Pre releases are not included in this.**
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
- `node` - `sudo apt install nodejs`
- `npm` - `sudo apt install npm`
- `golang` - See the [official](https://go.dev/doc/install) installation instructions
- `git` - `sudo apt install git` 

_Windows / Mac users will have to google or ask ChatGPT how to install these on their system._

__important__: You will still need to have the `settings` folder and it's content in the same folder structure when running this app. That means that you'll have to copy that along with any binary you build.

Download this repo and open a terminal in this folder:
```bash
user@hostname:~$ git clone https://github.com/luna-nightbyte/GoStreamRecord
user@hostname:~$ cd GoStreamRecord
user@hostname:~/GoStreamRecord$ make app
user@hostname:~/GoStreamRecord$ ./GoStreamRecord
```

#### Settings (Only up to release v0.3.4)
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
  Example (Only up to release v0.3.4): 
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

__NOTE__: v0.3.5 and up uses SQL database. The readme will be updated on the next release.
```bash
user@hostname:~$ docker run \
  -v ./settings:/app/settings \
  -v ./output:/app/videos \
  -v ./app.log:/app/remoteCtrl.log \
  -p 8080:8050 \
  docker.io/lunanightbyte/gorecord:v0.3.3
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
### Screenshots
<img width="2559" height="1026" alt="gallery_2" src="https://github.com/user-attachments/assets/2b8abcf5-8b26-4112-83d7-2905f35e8b3d" />
<img width="1383" height="733" alt="Download_2" src="https://github.com/user-attachments/assets/2cdfc58d-0fd9-40d5-a5ff-2a75a4734796" />
<img width="1550" height="734" alt="Livestream_2" src="https://github.com/user-attachments/assets/0d30a846-d6b8-425d-9c03-ad7e67d678b5" />
<img width="2152" height="734" alt="Recorder_2" src="https://github.com/user-attachments/assets/0536de95-049c-4c5f-a178-b641d4b3f1fe" />

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
