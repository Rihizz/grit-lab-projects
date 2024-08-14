# smart-road

Smart Intersection Management for Autonomous Vehicles

## Objectives

In this project, we aim to develop a smart intersection management strategy for Autonomous Vehicles (AVs).
Unlike traditional methods, which are designed for human drivers, this strategy aims to facilitate a collision-free and low-congestion environment for AVs at a cross intersection.

## Technologies

-Rust<br>
-macroquad Rust library

## Map of the simulation

```plaintext

               |   |   |   |   |   |   |
               |   |   |   |   |   |   |
               |   |   |   |   |   |   |
               |r  | s | l |   |   |   |
_______________| ← | ↓ | → |   |   |   |________________
                           |            ↑ r
_______________            |            ________________
                           |            ← s
_______________            |            ________________
                           |            ↓ l
___________________________|____________________________
           l ↑             |
_______________            |            ________________
           s →             |
_______________            |            ________________
           r ↓             |
_______________            |            ________________
               |   |   |   | ← | ↑ | → |
               |   |   |   | l | s | r |
               |   |   |   |   |   |   |
               |   |   |   |   |   |   |
               |   |   |   |   |   |   |
               |   |   |   |   |   |   |

```

Turning right (r)<br>
Going straight (s)<br>
Turning left (l)<br>

## Physics Rules:

AVs must strictly follow their designated lanes and routes.
Must have at least 3 different velocities.
Must maintain a safety distance from other AVs.
Animation
Animation is an essential part of this project. You need to find or create assets for vehicles and roads. The animation should not only render the image but also simulate the motion and direction of AVs.

## Commands

Arrow Up: Generate vehicles from south to north.<br>
Arrow Down: Generate vehicles from north to south.<br>
Arrow Right: Generate vehicles from west to east.<br>
Arrow Left: Generate vehicles from east to west.<br>
R: Continually generate random vehicles.<br>
Esc: End the simulation and display statistics.<br>
L: show the car boundaries and id. also the "radar"<br>

## Statistics

The program should be able to generate the following statistics:

Max/Min number of vehicles that passed the intersection.<br>
Max/Min velocity of all vehicles.<br>
Max/Min time taken to pass the intersection.<br>
Close calls, or violations of the safety distance.<br><br>
Bonus Features (Optional)<br>
Create your own assets for animation.<br>
Add more statistics.<br>
Consider acceleration and deceleration in the physics.<br>

## Learning Outcomes

Rust programming language<br>
sdl2 library<br>
Animation techniques<br>
Algorithm development<br>
Mathematics<br>
Event handling<br>

## Made by

Oskar<br>
Santeri<br>
Ville<br>

## auto koordinatit

VASEN tie vasemmalle menevä auto:

car_x = Some(0.0);
car_y = Some(508.0);
x_target = 508.0;
y_target = 0.0;

VASEN tie oikealle menevä auto:

car_x = Some(0.0);
car_y = Some(614.0);
x_target = 344.0;
y_target = 1000.0;

VASEN tie suoraan menevä auto:

car_x = Some(0.0);
car_y = Some(564.0);
x_target = 1000.0;
y_target = 564.0;

OIKEA tie suoraan menevä auto:

car_x = Some(1000.0);
car_y = Some(394.0);
x_target = 0.0;
y_target = 394.0;

OIKEA tie vasemmalle menevä auto:

car_x = Some(1000.0);
car_y = Some(450.0);
x_target = 450.0;
y_target = 1000.0;

OIKEA tie oikealle menevä auto:

car_x = Some(1000.0);
car_y = Some(344.0);
x_target = 614.0;
y_target = 0.0;

ALA tie suoraan menevä auto:

car_x = Some(564.0);
car_y = Some(1000.0);
x_target = 564.0;
y_target = 0.0;

ALA tie vasemmalle menevä auto:

car_x = Some(508.0);
car_y = Some(1000.0);
x_target = 0.0;
y_target = 450.0;

ALA tie oikealle menevä auto:

car_x = Some(614.0);
car_y = Some(1000.0);
x_target = 1000.0;
y_target = 613.0;

YLÄ tie suoraan menevä auto:

car_x = Some(395.0);
car_y = Some(0.0);
x_target = 395.0;
y_target = 1000.0;

YLÄ tie oikealle menevä auto:

car_x = Some(345.0);
car_y = Some(0.0);
x_target = 0.0;
y_target = 342.0;

YLÄ tie vasemmalle menevä auto:

car_x = Some(450.0);
car_y = Some(0.0);
x_target = 1000.0;
y_target = 510.0;
