# space-game_mk4

## Overview 

space-game_mk4 is an space station manager, idler game built on top of an mqtt-broker for out of the box multiplayer features.

build with:
- [ebiten](https://ebitengine.org/)
- [donburi ecs](https://github.com/yottahmd/donburi)
- [mochi mqtt](https://github.com/mochi-mqtt/server) (server)
- [pebble db](https://github.com/cockroachdb/pebble)

## Ideas, Systems, Components

### User
  A user profile consists of mainly user information used in authentication and broker message configuration.
  The [user hook](https://github.com/caleb-noodahl/space-game_mk4/blob/main/mqtt-server/hooks/user.go) handles the mqtt-client's auth request and stubs in the user profile data.

### Station
  
  The station serves as the primary interactable focal point of the game. Through it players can manage environment resources, employees, manufacture materials to construct components and machinery, conduct research, trade, and interact with other players.

### Employees
  
  Workers can be bought and sold on the employee market (refreshing occurs in the server's [gamehook's sysinfo tick](https://github.com/caleb-noodahl/space-game_mk4/blob/main/mqtt-server/hooks/game.go#L117))

### Research Lab

  The research lab exposes a tree that generates research tasks for employees and unlocks different gameplay functionality

### Machine Shop

  Consumes materials and generates new materials and components depending on research level and available employees.

### Docks

  Docks allow the station to accommodate new visitors and load/ unload cargo.

### Quests
   
   Quests relating to Station activities and game help will occasionally pop up from the user profile window.

## Feature Overview Checklist

- [x] Gamestate persistence
- [x] POC User Profile state and simple auth
- [x] Employee Marketplace poc with consistency across all clients
- [x] Employee tasks and XP gain
- [x] Auto complete quests based on world state
- [x] Wallet transactions 
- [x] Debug UI 
- [x] Research system poc
- [x] Basic quest system poc 

## Currently tinkering with

- [ ] Machine shop / manufacturing
- [ ] Commodities markets
- [ ] Docks & Trade
