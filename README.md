# broom

A Go program for scanning JAR files to uncover [the 29-09-2022 Minecraft malware](https://forums.papermc.io/threads/malware-announcement.529/) infections.

## Disclaimer

This software is by no means a sophisticated antimalware, it only catches one or possibly more variants of the 'Updater' malware; run [OpticFusion1's antimalware](https://github.com/OpticFusion1/MCAntiMalware) for a thorough check.

## Usage

Download a binary for your architecture and operating system in the [Releases](https://github.com/berrybyte-net/broom/releases) tab. (amd64 is x64 - 64-bit, 386 is x86 - 32-bit, Darwin is MacOS, arm64 is ARM - Apple Silicon for Apple users).

### Windows & MacOS

Drag and drop the binary into your server's directory (next to the plugins folder) and double click on it.

### Linux

Invoke the binary through your shell (most likely bash or a derivative).

Example:
````bash
curl -LO https://github.com/berrybyte-net/broom/releases/latest/download/broom_linux_amd64
chmod +x broom_linux_amd64
./broom_linux_amd64
````

## Support

This was hastily made and tested minimally, so please report any issues in the [issues](https://github.com/berrybyte-net/broom/issues) tab and/or reach out to us via Discord.
