use ::rand::Rng;
use macroquad::texture::Texture2D;
use std::time::Instant;

use macroquad::prelude::*;

// Public Enums for Vehicle Directions and Lanes
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum Lane {
    Right,
    Left,
    Straight,
}
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum CurrentDirection {
    North,
    East,
    South,
    West,
}
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum SpawnPoint {
    North,
    East,
    South,
    West,
}

// Public Struct to represent each vehicle in the simulation
//implement debug for vehicle
#[derive(Debug, Clone, PartialEq)]
pub struct Vehicle {
    pub id: usize, // Unique identifier for each vehicle
    pub x: f32,
    pub y: f32,
    pub x_target: f32,
    pub y_target: f32,
    pub max_velocity: f32,
    pub min_velocity: f32,
    pub current_velocity: f32,
    pub acceleration: f32,
    pub deceleration: f32,
    pub current_direction: CurrentDirection,
    pub lane: Lane,
    pub spawn_point: SpawnPoint,
    pub time_on_map: Instant, // Time the vehicle has been on the map
    pub initial_axis: &'static str,
    pub switch_to_next_axis: bool,
    pub texture: Texture2D,
    pub rotation: f32,
    pub prevent_acceleration: bool,
    pub long_radar: bool,
    pub is_in_center: bool,
    pub has_exited_center: bool,
    pub stop: bool,
    pub stopped_because_box: bool,
    pub circle_stop: bool,
    pub stopped_because: usize,
    pub can_go: bool,
    pub counted: bool,
}

impl Vehicle {
    pub fn new(id: usize, spawn_point: SpawnPoint, car_texture: Texture2D) -> Self {
        // Generate a random lane
        let mut rng = ::rand::thread_rng(); // Note the :: before rand
        let lane = match rng.gen_range(0..3) {
            0 => Lane::Right,
            1 => Lane::Left,
            _ => Lane::Straight,
        };
        //let lane left
        // let lane = Lane::Left;

        let (x, y, x_target, y_target) = match (&spawn_point, &lane) {
            (SpawnPoint::South, Lane::Straight) => (564.0, 1000.0, 564.0, 0.0),
            (SpawnPoint::South, Lane::Left) => (508.0, 1000.0, 0.0, 450.0),
            (SpawnPoint::South, Lane::Right) => (614.0, 1000.0, 1000.0, 613.0),
            (SpawnPoint::North, Lane::Straight) => (395.0, 0.0, 395.0, 1000.0),
            (SpawnPoint::North, Lane::Left) => (450.0, 0.0, 1000.0, 510.0),
            (SpawnPoint::North, Lane::Right) => (345.0, 0.0, 0.0, 342.0),
            (SpawnPoint::East, Lane::Straight) => (1000.0, 394.0, 0.0, 394.0),
            (SpawnPoint::East, Lane::Left) => (1000.0, 450.0, 450.0, 1000.0),
            (SpawnPoint::East, Lane::Right) => (1000.0, 344.0, 614.0, 0.0),
            (SpawnPoint::West, Lane::Straight) => (0.0, 564.0, 1000.0, 564.0),
            (SpawnPoint::West, Lane::Left) => (0.0, 508.0, 508.0, 0.0),
            (SpawnPoint::West, Lane::Right) => (0.0, 614.0, 344.0, 1000.0),
        };
        let current_velocity = 1.5;
        let max_velocity = current_velocity;
        let min_velocity = current_velocity;

        //initialize initial axis and switch to next axis as false
        //initial axis will be x if the spawn point is east or west, and y if the spawn point is north or south
        let initial_axis = if x == 0.0 || x == 1000.0 { "x" } else { "y" };
        let switch_to_next_axis = false;

        //initialize texture based on the variable car texture
        let texture = car_texture;

        //initialize acceleration and deceleration as 0 for now
        let acceleration = 0.0;
        let deceleration = 0.0;

        // Current direction will be the opposite of the spawn point
        let current_direction = match spawn_point {
            SpawnPoint::North => CurrentDirection::South,
            SpawnPoint::East => CurrentDirection::West,
            SpawnPoint::South => CurrentDirection::North,
            SpawnPoint::West => CurrentDirection::East,
        };

        let rotation = match current_direction {
            CurrentDirection::North => 0.0,
            CurrentDirection::East => 90.0,
            CurrentDirection::South => 180.0,
            CurrentDirection::West => 270.0,
        };

        //initialize time on map to be current time using instant
        let time_on_map = Instant::now();

        Self {
            id,
            x,
            y,
            x_target,
            y_target,
            max_velocity,
            min_velocity,
            current_velocity,
            acceleration,
            deceleration,
            current_direction,
            lane, // Set the random lane
            spawn_point,
            time_on_map, // Initialize to 0
            initial_axis,
            switch_to_next_axis,
            texture,
            rotation,
            prevent_acceleration: false,
            long_radar: false,
            is_in_center: false,
            has_exited_center: false,
            stop: false,
            stopped_because_box: false,
            circle_stop: false,
            stopped_because: 696969,
            can_go: false,
            counted: false,
        }
    }

    pub fn spawn_at_random(id: usize, car_texture: Texture2D) -> Self {
        let mut rng = ::rand::thread_rng();

        // Randomly choose a spawn point
        let spawn_point = match rng.gen_range(0..4) {
            0 => SpawnPoint::North,
            1 => SpawnPoint::East,
            2 => SpawnPoint::South,
            _ => SpawnPoint::West,
        };

        // Use the existing new method to construct the Vehicle object with the random spawn point
        Self::new(id, spawn_point, car_texture)
    }
    //implementing draw function for vehicle using macroquad
    pub fn draw(&self) {
        let rotation_degrees = match self.current_direction {
            CurrentDirection::North => 0.0,
            CurrentDirection::East => 90.0,
            CurrentDirection::South => 180.0,
            CurrentDirection::West => 270.0,
        };

        let rotation_radians = (rotation_degrees * std::f32::consts::PI) / 180.0;

        draw_texture_ex(&self.texture, self.x, self.y, WHITE, DrawTextureParams {
            rotation: rotation_radians,
            ..Default::default()
        });
    }
}
