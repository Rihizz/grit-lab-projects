FROM rust:1.63-buster

COPY ./maps			        /filler/maps
COPY ./m1_robots		    /filler/m1_robots
COPY ./m1_game_engine	    /filler/m1_game_engine
COPY ./solution             /filler/solution

WORKDIR /filler/solution/xdd

RUN cargo build --release

WORKDIR /filler/solution

ENTRYPOINT /bin/bash
