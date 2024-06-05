<p align="center">
<a href="https://github.com/its-a-feature/tiger/pulse">
        <img src="https://img.shields.io/github/commit-activity/m/its-a-feature/tiger/v3.0.0" 
          alt="Activity"/></a>
<img src="https://img.shields.io/badge/version-3.0.0-blue" alt="version 3.0.0"/>
<img src="https://img.shields.io/github/commits-since/its-a-feature/tiger/latest?include_prereleases&color=orange" 
  alt="commits since last release"/>
<a href="https://twitter.com/its_a_feature_">
    <img src="https://img.shields.io/twitter/follow/its_a_feature_?style=social" 
      alt="@its_a_feature_ on Twitter"/></a>
<a href="https://ghst.ly/BHSlack">
    <img src="https://img.shields.io/badge/BloodHound Slack-4A154B?logo=slack&logoColor=white"
        alt="chat on Bloodhound Slack"></a>
<a href="https://github.com/specterops#tiger">
    <img src="https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fspecterops%2F.github%2Fmain%2Fconfig%2Fshield.json"
      alt="Sponsored by SpecterOps"/>
</a>
</p>

# tiger
A cross-platform, post-exploit, red teaming framework built with GoLang, docker, docker-compose, and a web browser UI. It's designed to provide a collaborative and user friendly interface for operators, managers, and reporting throughout red teaming. 

## Starting tiger

tiger is controlled via the `tiger-cli` binary. To generate the binary, run `sudo make` from the main tiger directory. 
From there, you can run `sudo ./tiger-cli start` to bring up all default tiger containers.

More specific setup instructions, configurations, examples, screenshots, and more can be found on the [tiger Documentation](https://docs.tiger-c2.net) website.

## Installing Agents and C2 Profiles

The tiger repository itself does not host any Payload Types or any C2 Profiles. Instead, tiger provides a command, `./tiger-cli install github <url> [branch name] [-f]`, that can be used to install agents into a current tiger instance.

Payload Types and C2 Profiles can be found on the [overview](https://tigermeta.github.io/overview) page.

To install an agent, simply run the script and provide an argument of the path to the agent on GitHub:
```bash
sudo ./tiger-cli install github https://github.com/tigerAgents/apfell
```

The same is true for installing C2 Profiles:
```bash
sudo ./tiger-cli install github https://github.com/tigerC2Profiles/http
```

This allows the agents and c2 profiles to be updated at a much more regular pace and separates out the tiger Core components from the rest of tiger.

## Updating

Use the `./tiger-cli update` command to check for available updates across `tiger-cli`, `tiger_server`, and `tiger_react`'s UI. 
This will _NOT_ do the update for you, but let you know if an update exists. To check for updates against a specific branch, use `./tiger-cli update -b [branch name]`.


## tiger Docker Containers
<p align="left">
  <img src="https://img.shields.io/docker/v/itsafeaturetiger/tiger_go_base?color=green&label=Latest Release&sort=semver" alt="latest docker versions"/> 
  <img src="https://img.shields.io/github/v/release/tigerMeta/tiger_Docker_Templates?include_prereleases&label=Latest%20Pre-Release"/>
</p>

tiger uses Docker and Docker-compose for all of its components, which allows tiger to provide a wide range of components and features without having requirements exist on the host. However, it can be helpful to have insight into how the containers are configured. All of tiger's docker containers are hosted on DockerHub under [itsafeaturetiger](https://hub.docker.com/search?q=itsafeaturetiger&type=image).

- [tiger_go_base](https://hub.docker.com/repository/docker/itsafeaturetiger/tiger_go_base/general) - [Dockerfile](https://github.com/tigerMeta/tiger_Docker_Templates/tree/master/tiger_go_base)
  - <img src="https://img.shields.io/docker/image-size/itsafeaturetiger/tiger_go_base/latest" alt="image size"/>
  - <img src="https://img.shields.io/docker/pulls/itsafeaturetiger/tiger_go_base" alt="docker pull count" />
- [tiger_go_dotnet](https://hub.docker.com/repository/docker/itsafeaturetiger/tiger_go_dotnet/general) - [Dockerfile](https://github.com/tigerMeta/tiger_Docker_Templates/tree/master/tiger_go_dotnet)
  - <img src="https://img.shields.io/docker/image-size/itsafeaturetiger/tiger_go_dotnet/latest" alt="image size"/>
  - <img src="https://img.shields.io/docker/pulls/itsafeaturetiger/tiger_go_dotnet" alt="docker pull count"/>
- [tiger_go_macos](https://hub.docker.com/repository/docker/itsafeaturetiger/tiger_go_macos/general) - [Dockerfile](https://github.com/tigerMeta/tiger_Docker_Templates/tree/master/tiger_go_macos)
  - <img src="https://img.shields.io/docker/image-size/itsafeaturetiger/tiger_go_macos/latest" alt="image size"/>
  - <img src="https://img.shields.io/docker/pulls/itsafeaturetiger/tiger_go_macos" alt="docker pull count"/>
- [tiger_python_base](https://hub.docker.com/repository/docker/itsafeaturetiger/tiger_python_base/general) - [Dockerfile](https://github.com/tigerMeta/tiger_Docker_Templates/tree/master/tiger_python_base)
  - <img src="https://img.shields.io/docker/image-size/itsafeaturetiger/tiger_python_base/latest" alt="image size"/>
  - <img src="https://img.shields.io/docker/pulls/itsafeaturetiger/tiger_python_base" alt="docker pull count"/>
- [tiger_python_dotnet](https://hub.docker.com/repository/docker/itsafeaturetiger/tiger_python_dotnet/general) - [Dockerfile](https://github.com/tigerMeta/tiger_Docker_Templates/tree/master/tiger_python_dotnet)
  - <img src="https://img.shields.io/docker/image-size/itsafeaturetiger/tiger_python_dotnet/latest" alt="image size"/>
  - <img src="https://img.shields.io/docker/pulls/itsafeaturetiger/tiger_python_dotnet" alt="docker pull count"/>
- [tiger_python_macos](https://hub.docker.com/repository/docker/itsafeaturetiger/tiger_python_macos/general) - [Dockerfile](https://github.com/tigerMeta/tiger_Docker_Templates/tree/master/tiger_python_macos)
  - <img src="https://img.shields.io/docker/image-size/itsafeaturetiger/tiger_python_macos/latest" alt="image size"/>
  - <img src="https://img.shields.io/docker/pulls/itsafeaturetiger/tiger_python_macos" alt="docker pull count"/>
- [tiger_python_go](https://hub.docker.com/repository/docker/itsafeaturetiger/tiger_python_go/general) - [Dockerfile](https://github.com/tigerMeta/tiger_Docker_Templates/tree/master/tiger_python_go)
  - <img src="https://img.shields.io/docker/image-size/itsafeaturetiger/tiger_python_go/latest" alt="image size"/>
  - <img src="https://img.shields.io/docker/pulls/itsafeaturetiger/tiger_python_go" alt="docker pull count"/>

Additionally, tiger uses a custom PyPi package (tiger_container) and a custom Golang package (https://github.com/tigerMeta/tigerContainer) to help control and sync information between all the containers as well as providing an easy way to script access to the server.

Dockerfiles for each of these Docker images can be found on [tigerMeta](https://github.com/tigerMeta/tiger_Docker_Templates).

### tiger-container PyPi
<p align="left">
  <img src="https://img.shields.io/pypi/dm/tiger-container" alt="tiger-container downloads" />
  <img src="https://img.shields.io/pypi/pyversions/tiger-container" alt="tiger-container python version" />
  <img src="https://img.shields.io/pypi/v/tiger-container?color=green&label=Latest%20stable%20PyPi" alt="tiger-container version" />
  <img src="https://img.shields.io/github/v/release/tigerMeta/tigerContainerPypi?include_prereleases&label=Latest Pre-Release&color=orange" alt="latest release" />
</p>

The `tiger-container` PyPi package source code is available on [tigerMeta](https://github.com/tigerMeta/tigerContainerPyPi) and is automatically installed on all of the `tiger_python_*` Docker images.

This PyPi package is responsible for connecting to RabbitMQ, syncing your data to tiger, and responding to things like Tasking, Webhooks, and configuration updates.

### github.com/tigerMeta/tigerContainer
<p align="left">
  <img src="https://img.shields.io/github/go-mod/go-version/tigerMeta/tigerContainer" alt="tigerContainer go version"/>
  <img src="https://img.shields.io/github/v/release/tigerMeta/tigerContainer?label=Latest%20Stable&color=green" alt="tigerContainer latest stable version" />
  <img src="https://img.shields.io/github/v/release/tigerMeta/tigerContainer?include_prereleases&label=Latest Pre-Release&color=orange" alt="tigerContainer latest version" />
</p>

The `github.com/tigerMeta/tigerContainer` Golang package source code is available on [tigerMeta](https://github.com/tigerMeta/tigerContainer).

This Golang package is responsible for connecting to RabbitMQ, syncing your data to tiger, and responding to things like Tasking, Webhooks, and configuration updates.

## tiger Scripting
<p align="left">
  <img src="https://img.shields.io/pypi/dm/tiger" alt="tiger scripting downloads" />
  <img src="https://img.shields.io/pypi/pyversions/tiger" alt="tiger scripting python version" />
  <img src="https://img.shields.io/pypi/v/tiger?color=green&label=Latest%20Stable%20PyPi" alt="tiger scripting latest pypi version" />
<img src="https://img.shields.io/github/v/release/tigerMeta/tiger_Scripting?include_prereleases&label=Latest Pre-Release&color=orange" alt="latest release" />
</p>


* Scripting source code (https://github.com/tigerMeta/tiger_Scripting)

## Documentation

All documentation for the tiger project is being maintained on the [docs.tiger-c2.net](https://docs.tiger-c2.net) website.


## Contributions

A bunch of people have suffered through bug reports, changes, and fixes to help make this project better. Thank you!

The following people have contributed a lot to the project. As you see their handles throughout the project on Payload Types and C2 Profiles, be sure to reach out to them for help and contributions:
- [@djhohnstein](https://twitter.com/djhohnstein)
- [@xorrior](https://twitter.com/xorrior)
- [@Airzero24](https://twitter.com/airzero24)
- [@SpecterOps](https://twitter.com/specterops)

## Liability

This is an open source project meant to be used with authorization to assess the security posture and for research purposes.

## Historic References

* Check out a [series of YouTube videos](https://www.youtube.com/playlist?list=PLHVFedjbv6sNLB1QqnGJxRBMukPRGYa-H) showing how tiger looks/works and highlighting a few key features
* Check out the [blog post](https://posts.specterops.io/a-change-of-tiger-proportions-21debeb03617) on the rebranding.
* BSides Seattle 2019 Slides: [Ready Player 2: Multiplayer Red Teaming against macOS](https://www.slideshare.net/CodyThomas6/ready-player-2-multiplayer-red-teaming-against-macos)
* BSides Seattle 2019 Demo Videos: [Available on my Youtube](https://www.youtube.com/playlist?list=PLHVFedjbv6sOz8OGuLdomdkr6-7VdMRQ9)
* Objective By the Sea 2019 talk on JXA: https://objectivebythesea.com/v2/talks/OBTS_v2_Thomas.pdf
* Objective By the sea 2019 Video: https://www.youtube.com/watch?v=E-QEsGsq3uI&list=PLliknDIoYszvTDaWyTh6SYiTccmwOsws8&index=17  

## File Icon Attribution

* [bin/txt file icons](https://www.flaticon.com/packs/file-types-31?word=file%20extension) - created by Icon home - Flaticon