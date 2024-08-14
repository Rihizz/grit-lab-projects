// Public Struct to keep track of statistics during the simulation
//implement debug
#[derive(Debug, Clone)]
pub struct Statistics {
    pub max_vehicles: usize, // Max number of vehicles that passed the intersection
    pub max_velocity: f32, // Max velocity among all vehicles
    pub min_velocity: f32, // Min velocity among all vehicles
    pub max_time: f32, // Max time taken to pass the intersection
    pub min_time: f32, // Min time taken to pass the intersection
    pub close_calls: usize, // Number of close calls (violations of safe distance)
    pub collisions: usize, // Number of collisions
}

impl Statistics {
    pub fn new() -> Self {
        Self {
            max_vehicles: 0,
            max_velocity: 0.0,
            min_velocity: 1.5,
            max_time: 0.0,
            min_time: 0.0,
            close_calls: 0,
            collisions: 0,
        }
    }
}
