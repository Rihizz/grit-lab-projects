use macroquad::texture::Texture2D;

use crate::statistics::Statistics;
use crate::vehicle::{ CurrentDirection, Lane, SpawnPoint, Vehicle };
use macroquad::prelude::*;

//This manager is used to keep track of all vehicles in the simulation

pub struct VehicleManager {
    pub next_id: usize,
    pub vehicles: Vec<Vehicle>,
    pub statistics: Statistics,
}

pub struct Rectangle {
    x: f32,
    y: f32,
    width: f32,
    height: f32,
}

impl VehicleManager {
    pub fn new() -> Self {
        Self {
            next_id: 0,
            vehicles: Vec::new(),
            statistics: Statistics::new(),
        }
    }

    pub fn spawn_vehicle(&mut self, spawn_point: SpawnPoint, car_texture: Texture2D) {
        let prospective_vehicle = Vehicle::new(self.next_id, spawn_point, car_texture);
        let prospective_box = Rectangle {
            x: prospective_vehicle.x,
            y: prospective_vehicle.y,
            width: 40.0,
            height: 40.0,
        };
        if !self.would_collide_with_existing(&prospective_box) {
            self.vehicles.push(prospective_vehicle);
            self.next_id += 1;
            self.statistics.max_vehicles += 1;
        }
    }

    pub fn spawn_random_vehicle(&mut self, car_texture: Texture2D) {
        let prospective_vehicle = Vehicle::spawn_at_random(self.next_id, car_texture);
        let prospective_box = Rectangle {
            x: prospective_vehicle.x,
            y: prospective_vehicle.y,
            width: 40.0,
            height: 40.0,
        };
        if !self.would_collide_with_existing(&prospective_box) {
            self.vehicles.push(prospective_vehicle);
            self.next_id += 1;
            self.statistics.max_vehicles += 1;
        }
    }

    pub fn would_collide_with_existing(&self, new_vehicle_box: &Rectangle) -> bool {
        for vehicle in &self.vehicles {
            let existing_vehicle_box = Rectangle {
                x: vehicle.x,
                y: vehicle.y,
                width: 40.0,
                height: 40.0,
            };
            if Self::intersects(&new_vehicle_box, &existing_vehicle_box) {
                return true;
            }
        }
        false
    }

    pub fn check_for_others(&mut self, last_id: usize, lutka_tutka: bool) {
        let others = self.vehicles.clone();
        for vehicle in &mut self.vehicles {
            if vehicle.lane == Lane::Right {
                continue;
            }

            // Calculate the safe distance based on current speed.
            let safe_distance = vehicle.current_velocity * 40.0;

            let mut front_box = Rectangle {
                x: vehicle.x,
                y: vehicle.y,
                width: 40.0,
                height: 40.0 + safe_distance,
            };

            if vehicle.current_direction == CurrentDirection::North {
                front_box.y -= 0.0 + safe_distance;
            } else if vehicle.current_direction == CurrentDirection::East {
                front_box.height = 40.0;
                front_box.width = 40.0 + safe_distance;
            } else if vehicle.current_direction == CurrentDirection::West {
                front_box.height = 40.0;
                front_box.width = 40.0 + safe_distance;
                front_box.x -= 0.0 + safe_distance;
            }

            if lutka_tutka {
                macroquad::prelude::draw_rectangle(
                    front_box.x,
                    front_box.y,
                    front_box.width,
                    front_box.height,
                    macroquad::prelude::Color::new(1.0, 0.0, 0.0, 0.5)
                );
            }

            let mut x_for_box = vehicle.x;
            let mut y_for_box = vehicle.y;

            if vehicle.current_direction == CurrentDirection::North {
                y_for_box -= 60.0;
            } else if vehicle.current_direction == CurrentDirection::East {
                x_for_box += 60.0;
            } else if vehicle.current_direction == CurrentDirection::West {
                x_for_box -= 60.0;
            } else if vehicle.current_direction == CurrentDirection::South {
                y_for_box += 60.0;
            }

            if
                Self::is_car_inside_boundary(x_for_box, y_for_box) &&
                vehicle.lane == Lane::Left &&
                vehicle.id > last_id &&
                last_id != 0
            {
                vehicle.stop = true;
                vehicle.stopped_because_box = true;
                continue;
            }
            if vehicle.stopped_because_box && vehicle.id <= last_id {
                // If a car was stopped due to the boundary rule but there's space now.
                vehicle.stop = false;
                vehicle.stopped_because_box = false;
            }

            // let closest_vehicle_distance: Option<f32> = None;
            let mut closest_vehicle_info: Option<(f32, usize)> = None;
            // Define this variable before the loop:
            let mut closest_vehicle: Option<Vehicle> = None;
            for other in &others {
                if other.id == vehicle.id {
                    continue;
                }
                let car_b = Rectangle {
                    x: other.x,
                    y: other.y,
                    width: 40.0,
                    height: 40.0,
                };
                if Self::intersects(&front_box, &car_b) {
                    let distance = match vehicle.current_direction {
                        CurrentDirection::North => vehicle.y - (other.y + 40.0),
                        CurrentDirection::South => other.y - (vehicle.y + 40.0),
                        CurrentDirection::East => other.x - (vehicle.x + 40.0),
                        CurrentDirection::West => vehicle.x - (other.x + 40.0),
                    };
                    if closest_vehicle_info.is_none() || distance < closest_vehicle_info.unwrap().0 {
                        closest_vehicle_info = Some((distance, other.id)); // Store distance and other vehicle's ID
                        closest_vehicle = Some(other.clone()); // assuming VehicleType implements Clone
                    }
                } else {
                    vehicle.stop = false;
                }
            }

            if let Some((distance, closest_vehicle_id)) = closest_vehicle_info {
                let distancexd = distance.abs();
                if distancexd < 2.0 || vehicle.can_go {
                    vehicle.stop = true;
                    vehicle.long_radar = true;
                    vehicle.stopped_because = closest_vehicle_id;
                    if let Some(closest) = &closest_vehicle {
                        if
                            closest.stop &&
                            vehicle.stop &&
                            vehicle.stopped_because == closest.id &&
                            closest.stopped_because == vehicle.id &&
                            vehicle.id < closest.id
                        {
                            vehicle.stop = false;
                            vehicle.can_go = true;
                            vehicle.current_velocity = vehicle.max_velocity;
                        } else if vehicle.can_go {
                            vehicle.can_go = false;
                        }
                    }
                    continue;
                } else {
                    vehicle.stop = false;
                }
                if vehicle.stop && distancexd >= 2.0 && !vehicle.circle_stop {
                    vehicle.stop = false;
                }
                if distance < safe_distance {
                    // Decelerate if the closest vehicle is within the safe distance.
                    vehicle.current_velocity *= 0.5; // Decelerate by reducing velocity to half.
                    if vehicle.current_velocity < 0.05 {
                        vehicle.current_velocity = 0.0; // Stop the vehicle if its speed is very low.
                    }
                } else if distance > safe_distance && distance < safe_distance * 1.5 {
                    vehicle.stop = false;
                } else if vehicle.current_velocity < vehicle.max_velocity {
                    vehicle.stop = false;
                    // If there's enough room ahead, then accelerate more aggressively.
                    if vehicle.current_velocity < 0.1 {
                        vehicle.current_velocity = 0.1;
                    }
                    vehicle.current_velocity *= 1.5;
                    if vehicle.current_velocity > vehicle.max_velocity {
                        vehicle.current_velocity = vehicle.max_velocity;
                    }
                }
            } else if vehicle.current_velocity < vehicle.max_velocity {
                vehicle.stop = false;
                // If there's no vehicle ahead, accelerate.
                if vehicle.current_velocity < 0.1 {
                    vehicle.current_velocity = 0.1;
                }
                vehicle.current_velocity *= 1.5;
                if vehicle.current_velocity > vehicle.max_velocity {
                    vehicle.current_velocity = vehicle.max_velocity;
                }
            }
        }
    }

    fn is_car_inside_boundary(x: f32, y: f32) -> bool {
        let boundary = Rectangle {
            x: 445.0,
            y: 445.0,
            width: 110.0,
            height: 110.0,
        };

        let car = Rectangle {
            x: x,
            y: y,
            width: 40.0,
            height: 40.0,
        };

        // Check for overlaps on each axis
        let x_overlap = car.x < boundary.x + boundary.width && car.x + car.width > boundary.x;
        let y_overlap = car.y < boundary.y + boundary.height && car.y + car.height > boundary.y;

        x_overlap && y_overlap
    }

    pub fn update_with_collision_check(&mut self, lutka_tutka: bool) {
        // Reset vehicle counts
        let mut left_queue = 0;
        let mut last_id = 0;
        for vehicle in &mut self.vehicles {
            if Self::is_car_inside_boundary(vehicle.x, vehicle.y) {
                vehicle.is_in_center = true;
            }
            if vehicle.is_in_center && !Self::is_car_inside_boundary(vehicle.x, vehicle.y) {
                vehicle.has_exited_center = true;
            }
            // Counting vehicles in the left lane that haven't exited the center
            if vehicle.lane == Lane::Left && !vehicle.has_exited_center {
                left_queue += 1;
                if left_queue == 2 {
                    last_id = vehicle.id;
                }
            }
        }

        self.check_for_others(last_id, lutka_tutka);

        // Then run the regular update logic
        self.update();
    }

    pub fn intersects(a: &Rectangle, b: &Rectangle) -> bool {
        if a.x + a.width < b.x || b.x + b.width < a.x {
            return false;
        }
        if a.y + a.height < b.y || b.y + b.height < a.y {
            return false;
        }
        true
    }

    pub fn update(&mut self) {
        // Update logic for all vehicles
        for vehicle in self.vehicles.iter_mut() {
            //check for collisions
            if vehicle.stop {
                self.statistics.min_velocity = 0.0;
                continue;
            }

            let distance_to_x_target = vehicle.x_target - vehicle.x;
            let distance_to_y_target = vehicle.y_target - vehicle.y;

            // Use vehicle's current_velocity instead of global_car_speed
            let current_speed = vehicle.current_velocity;

            if !vehicle.switch_to_next_axis {
                match vehicle.initial_axis {
                    "x" => {
                        if distance_to_x_target.abs() > 1.0 {
                            vehicle.x += current_speed * distance_to_x_target.signum();
                        } else {
                            vehicle.switch_to_next_axis = true;
                            vehicle.rotation = std::f32::consts::PI / 2.0;
                        }
                    }
                    "y" => {
                        if distance_to_y_target.abs() > 1.0 {
                            vehicle.y += current_speed * distance_to_y_target.signum();
                        } else {
                            vehicle.switch_to_next_axis = true;
                            vehicle.rotation = 0.0;
                        }
                    }
                    _ => {}
                }
                vehicle.rotation = 0.0; // Reset rotation when switching axis
            } else {
                if vehicle.initial_axis == "y" {
                    if vehicle.x < vehicle.x_target {
                        vehicle.current_direction = CurrentDirection::East;
                    } else if vehicle.x > vehicle.x_target {
                        vehicle.current_direction = CurrentDirection::West;
                    }
                } else {
                    if vehicle.y < vehicle.y_target {
                        vehicle.current_direction = CurrentDirection::South;
                    } else if vehicle.y > vehicle.y_target {
                        vehicle.current_direction = CurrentDirection::North;
                    }
                }

                // switch to the other axis
                match vehicle.initial_axis {
                    "x" => {
                        if distance_to_y_target.abs() > 1.0 {
                            vehicle.y += current_speed * distance_to_y_target.signum();
                        }
                    }
                    "y" => {
                        if distance_to_x_target.abs() > 1.0 {
                            vehicle.x += current_speed * distance_to_x_target.signum();
                        }
                    }
                    _ => {}
                }
            }
            //update statistics with max / min velocity
            if vehicle.current_velocity > self.statistics.max_velocity {
                self.statistics.max_velocity = vehicle.current_velocity;
            }
            if vehicle.current_velocity < self.statistics.min_velocity {
                self.statistics.min_velocity = vehicle.current_velocity;
            }
            if vehicle.long_radar && !vehicle.counted {
                self.statistics.close_calls += 1;
                vehicle.counted = true;
            }
        }
        // Initialize an empty Vec to store vehicles that you wish to retain
        let mut retained_vehicles: Vec<Vehicle> = Vec::new();

        // Iterate through the existing vehicles
        for vehicle in &self.vehicles {
            let distance_to_x_target = vehicle.x_target - vehicle.x;
            let distance_to_y_target = vehicle.y_target - vehicle.y;

            // If the vehicle is far from the target, keep it
            if distance_to_x_target.abs() > 1.0 || distance_to_y_target.abs() > 1.0 {
                retained_vehicles.push(vehicle.clone()); // Assumes your Vehicle struct derives Clone
            } else {
                // If the vehicle is close to the target, act upon it
                let elapsed_time = vehicle.time_on_map.elapsed();

                //check if the time is greater than max time
                if elapsed_time.as_secs_f32() > self.statistics.max_time {
                    self.statistics.max_time = elapsed_time.as_secs_f32();
                }
                //check if min time is 0 and update it
                if self.statistics.min_time == 0.0 {
                    self.statistics.min_time = elapsed_time.as_secs_f32();
                }
                //check if the time is less than min time
                if elapsed_time.as_secs_f32() < self.statistics.min_time {
                    self.statistics.min_time = elapsed_time.as_secs_f32();
                }
            }
        }

        // Replace the old Vec with the new one
        self.vehicles = retained_vehicles;
    }
}
