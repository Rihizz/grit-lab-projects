use std::{ fs::File, io::Read, collections::HashMap };
use serde::{ Serialize, Deserialize };
use std::io::{ self };
use serde_json::{ Result, Value };

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct ServerInfo {
    pub ports: Vec<PortOptions>,
    pub page_path: String,
    pub error_page_path: String,
    pub timeout_duration: usize,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PortOptions {
    pub port: String,
    pub paths: HashMap<String, PathOptions>,
    pub error_page_path: String,
    pub upload_limit: usize,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PathOptions {
    pub allowed_methods: Vec<String>,
    pub route: String,
    pub default_file: String,
    pub directory_listing_mode: bool,
    pub cgi_path: String,
}

impl PathOptions {
    pub fn new() -> Self {
        return PathOptions {
            allowed_methods: vec![],
            route: String::new(),
            default_file: String::new(),
            directory_listing_mode: false,
            cgi_path: String::new(),
        };
    }
}

impl PortOptions {
    pub fn new() -> Self {
        return PortOptions {
            paths: HashMap::new(),
            port: String::new(),
            error_page_path: String::new(),
            upload_limit: 1048576,
        };
    }
}

impl ServerInfo {
    pub fn new() -> Self {
        return ServerInfo {
            ports: vec![],
            page_path: "src/content".to_owned(),
            error_page_path: "src/error-pages".to_owned(),
            timeout_duration: 25,
        };
    }
    pub fn deserialize(&mut self) -> io::Result<()> {
        let mut buf = String::new();
        let mut f = File::open("server_config.json")?;
        f.read_to_string(&mut buf)?;
        let conf: ServerInfo = serde_json::from_str(
            &buf
        )?; /* serde_json::from_str(&buf).unwrap(); */
        *self = conf;
        Ok(())
    }
    pub fn default(&mut self) {
        self.ports = vec![PortOptions::new()]; //vec!["0.0.0.0:7878".to_owned()];
        let mut h_map: HashMap<String, PathOptions> = HashMap::new();
        h_map.insert("/index.html".to_owned(), PathOptions {
            allowed_methods: vec!["GET".to_owned(), "POST".to_owned(), "DELETE".to_owned()],
            route: "src/content/index.html".to_owned(),
            default_file: String::new(),
            directory_listing_mode: false,
            cgi_path: String::new(),
        });
        //let h_map2 = h_map.clone();
        self.ports[0].paths = h_map;
        self.page_path = "src/content".to_owned();
        self.error_page_path = "src/error-pages".to_owned();
        self.timeout_duration = 25;
        //let js = serde_json::to_string(&self);
        //println!("serialized {:#?}", js);
    }
}
