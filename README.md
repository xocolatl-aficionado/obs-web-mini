# OBS-web-mini

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-17-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

#### The minified way to control [OBS](https://obsproject.com/) remotely on a private network

---

#### Screenshot:
- ![home-page-app](https://github.com/user-attachments/assets/675e0423-15d9-4e3a-8012-da2ba26e6673)
- ![working-obs-app](https://github.com/user-attachments/assets/25864c5b-7b53-4690-88ce-dfe6e2636949)


#### Features:

- No installation or extra software is needed (except npm), works in any modern browser (desktop + mobile)
- Start/Stop/Pause/Resume recordings
- Security to be revisited. Use at your own risk!
---

#### Requirements:

- [OBS](https://obsproject.com/) v28 or higher - this includes the latest version of the OBS-websocket plugin
- Enabling the OBS-websocket server in OBS under `Tools -> obs-websocket Settings -> Enable WebSocket Server`
- Download the `my-app-code.tar.gz` file onto a machine and cd into it, then run `npm install` && `npm run start:all`. That should make the web app available at port 5000. 
---

#### Build instructions:

```bash
npm ci
```

## Contributors✨

Based on work done by [Niek](https://github.com/Niek/obs-web?tab=readme-ov-file). This repo just makes a minified version, that aims to be more user-friendly. 

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!
