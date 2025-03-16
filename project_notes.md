# Project Overview

## ECS

### Entities

Entities are used purely to point to components. They are a way to group components together to define what an entity is.

### Components

Components are used to define the _data_ portion, whether it be for an interaction in systems, or purely just data (ie. item info).
Components can be separated into two categories

1. Pure Components
1. Intents

Pure components just hold data and attach to entities to help define what that entity is / can do.

Intents give action to an entity and are processed in systems.

### Systems

Systems are used to process intents and update the state of the game. They define the logic of the game.
