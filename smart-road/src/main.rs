use macroquad::miniquad::window::set_window_size;
use macroquad::prelude::*;
mod statistics;
mod vehicle;
mod vehiclemanager;
use vehicle::SpawnPoint;
use vehiclemanager::VehicleManager;

#[macroquad::main("InterSEXtion heheheahaha:-DD")]
async fn main() {
    let mut vehicle_manager = VehicleManager::new();
    set_window_size(1000, 1030);
    let mut pause = false;
    let mut random = false;
    let mut display_statistics = false;
    let mut lutka_tutka = false;

    let background_texture = load_texture("./pictures/intersection.png").await.expect(
        "Failed to load image"
    );

    let spawn_delay = 0.8; // 0.8 seconds delay for spawning vehicles
    let mut last_spawn_time = get_time() - spawn_delay; // initialize to allow immediate spawn at the beginning

    let car_texture = load_texture("./pictures/greycar.png").await.expect(
        "Failed to load car image"
    );

    loop {
        clear_background(BLACK);

        let current_time = get_time();

        draw_texture(&background_texture, 0.0, 0.0, WHITE);

        //LOGIC FOR SPAWNING VEHICLES WITH VEHICLE MANAGER
        // Handle arrow key inputs to spawn vehicles
        if
            is_key_pressed(KeyCode::Right) &&
            current_time - last_spawn_time >= spawn_delay &&
            !pause
        {
            vehicle_manager.spawn_vehicle(SpawnPoint::West, car_texture.clone());
            last_spawn_time = current_time;
        }
        if is_key_pressed(KeyCode::Left) && current_time - last_spawn_time >= spawn_delay && !pause {
            vehicle_manager.spawn_vehicle(SpawnPoint::East, car_texture.clone());
            last_spawn_time = current_time;
        }
        if is_key_pressed(KeyCode::Up) && current_time - last_spawn_time >= spawn_delay && !pause {
            vehicle_manager.spawn_vehicle(SpawnPoint::South, car_texture.clone());
            last_spawn_time = current_time;
        }
        if is_key_pressed(KeyCode::Down) && current_time - last_spawn_time >= spawn_delay && !pause {
            vehicle_manager.spawn_vehicle(SpawnPoint::North, car_texture.clone());
            last_spawn_time = current_time;
        }
        //apply the same logic for R which spawns a random vehicle
        if is_key_pressed(KeyCode::R) {
            if random {
                random = false;
            } else {
                random = true;
            }
        }
        if is_key_pressed(KeyCode::L) {
            if lutka_tutka {
                lutka_tutka = false;
            } else {
                lutka_tutka = true;
            }
        }
        if random && current_time - last_spawn_time >= spawn_delay && !pause {
            vehicle_manager.spawn_random_vehicle(car_texture.clone());
            last_spawn_time = current_time;
        }

        // if there are no cars dont run the update with collisions
        let number_of_vehicles = vehicle_manager.vehicles.len();

        if !pause && number_of_vehicles > 0 {
            //update vehicle manager
            vehicle_manager.update_with_collision_check(lutka_tutka);
            //draw the updated vehicles from vehicle managers vehicles vector
        }
        if is_key_pressed(KeyCode::Escape) {
            if display_statistics {
                // Close the program when "esc" is pressed again
                break;
            } else {
                pause = true;
                display_statistics = true;
            }
        }

        for vehicle in &mut vehicle_manager.vehicles {
            vehicle.draw(); // Draw the vehicle
            if lutka_tutka {
                draw_text(&vehicle.id.to_string(), vehicle.x, vehicle.y, 20.0, RED); // Draw the vehicle id
                draw_rectangle_lines(vehicle.x, vehicle.y, 40.0, 40.0, 2.0, RED); // Draw a box around the vehicle
            }
        }

        if display_statistics {
            // Draw a semi-transparent background for the statistics screen
            draw_rectangle(
                0.0,
                0.0,
                screen_width(),
                screen_height(),
                Color::new(0.0, 0.0, 0.0, 0.5)
            );

            // let mut close_calls = 0;

            // //loop through all the vehicles
            // for vehicle in &mut vehicle_manager.vehicles {
            //     if vehicle.long_radar {
            //         close_calls += 1;
            //     }
            // }

            // Calculate statistics
            let max_vehicles = vehicle_manager.statistics.max_vehicles;
            let max_velocity = vehicle_manager.statistics.max_velocity;
            let mut min_velocity = vehicle_manager.statistics.min_velocity;
            let max_time = vehicle_manager.statistics.max_time;
            let min_time = vehicle_manager.statistics.min_time;
            let collisions = vehicle_manager.statistics.collisions;
            let close_calls = vehicle_manager.statistics.close_calls;
            if max_velocity == 0.0 {
                min_velocity = 0.0;
            }

            // Draw statistics as a vertical list
            let stats_text = format!(
                "Max Vehicles: {}\nMax Velocity: {:.2}\nMin Velocity: {:.2}\nMax Time: {:.2}\nMin Time: {:.2}\nClose Calls: {}\nCollisions: {}",
                max_vehicles,
                max_velocity,
                min_velocity,
                max_time,
                min_time,
                close_calls,
                collisions
            );

            let mut y_position = screen_height() * 0.42;
            let line_height = 30.0;

            for line in stats_text.lines() {
                draw_text(line, screen_width() * 0.4, y_position, 30.0, WHITE);
                y_position += line_height;
            }
        }

        // draw_rectangle_lines(445.0, 445.0, 110.0, 110.0, 2.0, LIGHTGRAY); // Draw a box in the center
        next_frame().await;
    }
}
