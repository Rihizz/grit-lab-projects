use std::collections::HashMap;
use std::fs::File;
use std::io::{ self, Read, Write };
use std::net::SocketAddr;
use std::str;
use std::time::{ Duration, Instant };
use std::fs;
use base64;
use mio::event::{ Source, self };
use mio::net::{ TcpListener, TcpStream };
use mio::{ Events, Interest, Poll, Token };
use cookie::Cookie;
use std::process::Command;
use localhost::localhost::config::{ ServerInfo, PortOptions, PathOptions };

fn list_directory(path: &str) -> Result<String, io::Error> {
    let paths = fs::read_dir(path)?;

    let mut paths_str = String::new();
    for path in paths {
        match path {
            Ok(p) => {
                paths_str = format!(
                    "{paths_str}\n <a href=\"{}\">{}</a>",
                    p.path().display(),
                    p.path().display()
                );
            }
            Err(e) => {
                return Err(e);
            }
        }
    }
    let body = format!(
        "<!DOCTYPE html>
    <html lang=\"en\">
      <head>
        <meta charset=\"UTF-8\" />
        <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\" />
        <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />
        <link rel=\"shortcut icon\" href=\"data:image/x-icon;,\" type=\"image/x-icon\" />
        <title>directory listing</title>
      </head>
      <body>
        <h1>Directory listing</h1>
        {paths_str}
      </body>
    </html>"
    );
    let res = format!(
        "HTTP/1.1 200 OK\r\nContent-Length: {}\r\nContent-Type: text/html\r\nConnection: keep-alive\r\n\r\n {}",
        body.len(),
        body
    );
    return Ok(res);
}

fn redirect(path: String, socket: &mut TcpStream) -> Result<(), io::Error> {
    let status = format!("HTTP/1.1 301 Moved Permanently\r\nLocation: {path}\r\n\r\n");
    let res = socket.write_all(status.as_bytes());
    return res;
}

fn handle_cgi(socket: &mut TcpStream, path: &str) -> io::Result<()> {
    let output = Command::new("php-cgi").arg(path).output()?;
    //.expect("Failed to execute CGI script");

    let output_str = String::from_utf8_lossy(&output.stdout);
    let parts: Vec<&str> = output_str.splitn(2, "\r\n\r\n").collect();

    let first_part = parts[0];
    println!("header? {:?}", first_part);

    let content = if parts.len() > 1 {
        parts[1] // The part after the headers
    } else {
        "" // No content found
    };

    //case if cgi is small enough
    let content_bytes = content.as_bytes();

    if content_bytes.len() < 1024 {
        let status_line = "HTTP/1.1 200 OK";
        let response = format!(
            "{}\r\nContent-Type: text/html\r\nContent-Length: {}\r\nConnection: keep-alive\r\n\r\n{}",
            status_line,
            content.len(),
            content
        );

        socket.write_all(response.as_bytes())?;

        return Ok(());
    }
    let status_line = "HTTP/1.1 200 OK";
    let response =
        format!("{}\r\nTransfer-Encoding: chunked\r\nContent-Type: text/html\r\n\r\n", status_line);

    socket.write_all(response.as_bytes())?;

    for chunk in content.as_bytes().chunks(1024) {
        let chunk_size = format!("{:X}\r\n", chunk.len());
        socket.write_all(chunk_size.as_bytes())?;
        socket.write_all(chunk)?;
        socket.write_all(b"\r\n")?;
    }
    socket.write_all(b"0\r\n\r\n")?;

    return Ok(());
}

fn check_cookie(request: &str) -> Option<String> {
    let lines: Vec<&str> = request.split("\r\n").collect();
    for line in lines {
        if line.starts_with("Cookie: ") {
            let cookie_string = line.strip_prefix("Cookie: ").unwrap();
            let cookies: Vec<&str> = cookie_string.split("; ").collect();
            for cookie in cookies {
                let parts: Vec<&str> = cookie.split('=').collect();
                if parts[0] == "name" {
                    return Some(parts[1].to_string());
                }
            }
        }
    }
    None
}

fn get_error_page(
    error_path: &str,
    error_number: usize
) -> Result<(String, File), (String, String)> {
    match error_number {
        500 => {
            let f = File::open(format!("{error_path}500.html")); //.unwrap_or(create_error_page())
            match f {
                Ok(c) => Ok((get_response_message(500), c)),
                _ => Err((get_response_message(500), create_error_page())),
            }
        }
        _ => {
            let path = format!("{error_path}{error_number}.html");
            let f = File::open(path);
            match f {
                Ok(c) => Ok((get_response_message(error_number), c)),
                _ => get_error_page(error_path, 500),
            }
        }
    }
}

fn create_error_page() -> String {
    "<!DOCTYPE html>
    <html lang=\"en\">
      <head>
        <meta charset=\"UTF-8\" />
        <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\" />
        <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />
        <title>500</title>
      </head>
      <body>
        <h1>Oops!</h1>
        <p>Sorry, internal server error.</p>
      </body>
    </html>".to_owned()
}

fn get_response_message(response_code: usize) -> String {
    match response_code {
        200 => "200 OK".to_owned(),
        400 => "400 Bad Request".to_owned(),
        404 => "404 Not Found".to_owned(),
        403 => "403 Forbidden".to_owned(),
        405 => "405 Method Not Allowed".to_owned(),
        413 => "413 Payload Too Large".to_owned(),
        _ => "500 Internal Server Error".to_owned(),
    }
}

fn get_status_from_string(path: String) -> usize {
    if path.contains("400.html") {
        return 400;
    }
    if path.contains("403.html") {
        return 403;
    }
    if path.contains("404.html") {
        return 404;
    }
    if path.contains("405.html") {
        return 405;
    }
    if path.contains("413.html") {
        return 413;
    }
    if path.contains("500.html") {
        return 500;
    }
    return 404;
}

fn get_page(path: &str, error_page_path: &str) -> String {
    let new_path = if path == "/" { format!("src/content/index.html") } else { format!("{path}") };
    let mut response_status = String::from("200 OK");
    let file = File::open(new_path);
    let mut contents = String::new();

    match file {
        Ok(mut f) => {
            let res = f.read_to_string(&mut contents);
            if res.is_err() {
                contents = create_error_page();
                response_status = get_response_message(500);
            }
        }
        _ => {
            println!("file did not exist");
            let error_number = get_status_from_string(path.to_owned());

            match get_error_page(error_page_path, error_number) {
                Ok(mut f2) => {
                    let res = f2.1.read_to_string(&mut contents);
                    response_status = f2.0;
                    if res.is_err() {
                        contents = create_error_page();
                        response_status = get_response_message(500);
                    }
                }
                Err(s) => {
                    println!("404 page did not exist");
                    contents = s.1;
                    response_status = s.0;
                }
            }
        }
    }
    let cookie = Cookie::new("name", "value");
    let cookie_header = format!("Set-Cookie: {}", cookie.to_string());
    let response = format!(
        "HTTP/1.1 {response_status}\r\nContent-Type: text/html\r\n{}\r\nContent-Length: {}\r\nConnection: keep-alive\r\n\r\n{}",
        cookie_header,
        contents.len(),
        contents
    );
    response
}

fn create_listener(address: &str) -> Result<TcpListener, Box<dyn std::error::Error>> {
    let socket_addr: SocketAddr = address.parse()?;
    let listener = TcpListener::bind(socket_addr)?;
    Ok(listener)
}

fn is_double_crnl(window: &[u8]) -> bool {
    window.len() >= 4 &&
        window[0] == b'\r' &&
        window[1] == b'\n' &&
        window[2] == b'\r' &&
        window[3] == b'\n'
}

fn handle_get(socket: &mut TcpStream, path: &str, server_info: &PortOptions) {
    let page_info = server_info.paths.get(path);

    if path == "/view_image" {
        // Specify the path to your image here
        let image_path = "./src/content/uploaded_image.jpg"; // Update this to the actual image path

        if let Err(e) = send_image(socket, image_path) {
            eprintln!("Failed to send image: {}", e);
        }
        return;
    }

    match page_info {
        Some(p_info) => {
            if !p_info.allowed_methods.contains(&"GET".to_owned()) {
                handle_error(socket, 405, &server_info.error_page_path);
            }
            let mut page = String::new();
            if p_info.directory_listing_mode == true {
                println!("got to dir listing");
                match list_directory(&p_info.route) {
                    Ok(pg) => {
                        page = pg;
                    }
                    Err(e) => {
                        println!("directory listing error {e}");
                        handle_error(socket, 400, &server_info.error_page_path);
                    }
                }

                if let Err(e) = socket.write_all(page.as_bytes()) {
                    eprintln!("Failed to write to socket: {}", e);
                }
                return;
            }
            let route = if p_info.route.contains(".html") || p_info.route.contains(".php") {
                p_info.route.clone()
            } else {
                format!("{}{}", p_info.route.clone(), p_info.default_file.clone())
            };
            let page = get_page(&route, &server_info.error_page_path);

            if let Err(e) = socket.write_all(page.as_bytes()) {
                eprintln!("Failed to write to socket: {}", e);
            }
        }
        None => {
            let page = get_page(path, &server_info.error_page_path);
            if let Err(e) = socket.write_all(page.as_bytes()) {
                eprintln!("Failed to write to socket: {}", e);
            }
        }
    }
}

fn send_image(socket: &mut TcpStream, image_path: &str) -> io::Result<()> {
    // Open the image file
    let mut file = match File::open(image_path) {
        Ok(file) => file,
        Err(e) => {
            // Log the error and return a 404 Not Found response or other appropriate error response
            eprintln!("Failed to open image: {}", e);
            let _ = redirect("/400.html".to_owned(), socket);
            return Err(e);
        }
    };

    // Determine the content type based on file extension
    let content_type = if image_path.ends_with(".jpg") || image_path.ends_with(".jpeg") {
        "image/jpeg"
    } else if image_path.ends_with(".png") {
        "image/png"
    } else {
        // Default content type, or you could handle more types or return an error
        "application/octet-stream"
    };

    // get content length
    let metadata = fs::metadata(image_path)?;
    let content_length = metadata.len();

    // Prepare the HTTP response headers
    let response_headers = format!(
        "HTTP/1.1 200 OK\r\nContent-Type: {}\r\nContent-Length: {}\r\n\r\n",
        content_type,
        content_length
    );

    // Write the headers to the socket
    socket.write_all(response_headers.as_bytes())?;

    // Create a buffer and write the file to the socket
    let mut buffer = [0; 1024];
    while let Ok(count) = file.read(&mut buffer) {
        if count == 0 {
            break; // End of file reached
        }
        socket.write_all(&buffer[..count])?;
    }

    Ok(())
}

fn handle_post(
    socket: &mut TcpStream,
    path: &str,
    body: &[u8],
    server_info: &PortOptions,
    cookie_exists: bool
) {
    println!("Post request {:?}", path);

    if path == "/upload_image" {
        if !cookie_exists {
            let r = redirect("/403.html".to_owned(), socket);
            if r.is_err() {
                println!("redirect failed lol");
            }
            return;
        }
        // Find the start of the base64 data after the headers (indicated by "\r\n\r\n")
        let base64_start = body
            .windows(4)
            .position(|window| window == [13, 10, 13, 10])
            .map(|pos| pos + 4)
            .unwrap_or(0);
        let encoded_image = &body[base64_start..];

        // Decode the base64 encoded image data
        let image_data = match base64::decode(encoded_image) {
            Ok(data) => data,
            Err(e) => {
                eprintln!("Failed to decode image: {}", e);
                let response = "HTTP/1.1 400 Bad Request\r\n\r\nInvalid image data";
                let _ = socket.write_all(response.as_bytes());
                return;
            }
        };

        println!("image data length {}", image_data.len());

        if image_data.len() > server_info.upload_limit {
            println!("image too large");
            let _ = redirect("/413.html".to_owned(), socket);
            //handle_error(socket, 413, &server_info.error_page_path);
            return;
        }

        if image_data.len() > 3 * server_info.upload_limit {
            let _ = redirect("/413.html".to_owned(), socket);
            //handle_error(socket, 413, &server_info.error_page_path);
            return;
        }

        // Specify the path where you want to save the image
        let image_path = "src/content/uploaded_image.jpg";

        // Delete the image if it exists
        if std::path::Path::new(image_path).exists() {
            if let Err(e) = fs::remove_file(image_path) {
                eprintln!("Failed to delete existing image file: {}", e);
                handle_error(socket, 500, &server_info.error_page_path);
                return;
            }
        }

        // Create and write to the file
        let mut file = match File::create(image_path) {
            Ok(f) => f,
            Err(e) => {
                eprintln!("Failed to create image file: {}", e);
                handle_error(socket, 500, &server_info.error_page_path);
                return;
            }
        };

        if let Err(e) = file.write_all(&image_data) {
            eprintln!("Failed to write to image file: {}", e);
            handle_error(socket, 500, &server_info.error_page_path);
            return;
        }

        // Respond with a success message
        let response = "HTTP/1.1 200 OK\r\n\r\nImage uploaded successfully";
        let _ = socket.write_all(response.as_bytes());
    } else if path == "/benis" {
        // give bad request eror
        //handle_error(socket, 400);
        let r = redirect("/403.html".to_owned(), socket);
        if r.is_err() {
            println!("redirect failed lol");
        } else {
            println!("redirect success");
        }
    } else {
        let mut path_info = &PathOptions::new();
        _ = path_info;

        match server_info.paths.get(path) {
            Some(p) => {
                path_info = p;
            }
            None => {
                handle_error(socket, 404, &server_info.error_page_path);
                return;
            }
        }
        if !path_info.allowed_methods.contains(&"POST".to_owned()) {
            handle_error(socket, 405, &server_info.error_page_path);
            return;
        }

        if !cookie_exists {
            let r = redirect("/403.html".to_owned(), socket);
            if r.is_err() {
                println!("redirect failed lol");
            }
            return;
        }

        if !path_info.cgi_path.is_empty() {
            let res = handle_cgi(socket, &path_info.cgi_path);
            if res.is_err() {
                println!("cgi error {:?}", res.err());
                handle_error(socket, 400, &server_info.error_page_path);
                return;
            }
        }
    }
}

fn handle_delete(
    socket: &mut TcpStream,
    path: &str,
    server_info: &PortOptions,
    cookie_exists: bool
) {
    // Process DELETE request to delete the uploaded image

    if !cookie_exists {
        let r = redirect("/403.html".to_owned(), socket);
        if r.is_err() {
            println!("redirect failed lol");
        }
        return;
    }

    let mut path_info = &PathOptions::new();
    _ = path_info;

    match server_info.paths.get(path) {
        Some(p) => {
            path_info = p;
        }
        None => {
            handle_error(socket, 404, &server_info.error_page_path);
            return;
        }
    }
    if !path_info.allowed_methods.contains(&"DELETE".to_owned()) {
        handle_error(socket, 405, &server_info.error_page_path);
        return;
    }

    let image_path = "src/content/uploaded_image.jpg";
    if path == "/delete_image" {
        if std::fs::metadata(image_path).is_ok() {
            match std::fs::remove_file(image_path) {
                Ok(_) => {
                    // Respond with a success status (204 No Content) if the deletion is successful
                    let response = "HTTP/1.1 204 No Content\r\n\r\n";
                    if let Err(e) = socket.write_all(response.as_bytes()) {
                        eprintln!("Failed to write to socket: {}", e);
                    }
                }
                Err(e) => {
                    // If there's an error deleting the file, respond with a 500 Internal Server Error status
                    eprintln!("Failed to delete image: {}", e);
                    let response =
                        "HTTP/1.1 500 Internal Server Error\r\n\r\nFailed to delete image";
                    if let Err(e) = socket.write_all(response.as_bytes()) {
                        eprintln!("Failed to write to socket: {}", e);
                    }
                }
            }
        } else {
            // Respond with a message if the image does not exist
            let response = "HTTP/1.1 404 Not Found\r\n\r\nImage not found";
            if let Err(e) = socket.write_all(response.as_bytes()) {
                eprintln!("Failed to write to socket: {}", e);
            }
        }
    }
}

fn handle_error(socket: &mut TcpStream, code: usize, error_path: &str) {
    let mut response_message = String::new();
    let page = get_error_page(error_path, code);
    let mut error_page = String::new();
    match page {
        Ok(mut p) => {
            response_message = p.0;
            let e = p.1.read_to_string(&mut error_page);
            if e.is_err() {
                error_page = create_error_page();
                response_message = "500 Internal Server Error".to_owned();
            }
        }
        Err(e) => {
            error_page = e.1;
            response_message = e.0;
        }
    }
    let response = format!(
        "HTTP/1.1 {}\r\nContent-Type: text/html\r\nContent-Length: {}\r\nConnection: keep-alive\r\n\r\n{}",
        response_message,
        error_page.len(),
        error_page
    );
    if let Err(e) = socket.write_all(response.as_bytes()) {
        eprintln!("Failed to write to socket: {}", e);
    }
}

fn register_servers(
    addresses: &Vec<PortOptions>,
    poll: &mut Poll,
    servers: &mut HashMap<Token, TcpListener>
) -> usize {
    let mut counter = 0;
    for address in addresses.iter() {
        match create_listener(&address.port) {
            Ok(mut l) => {
                let result = poll.registry().register(&mut l, Token(counter), Interest::READABLE);
                if result.is_err() {
                    println!("failed to register to poll");
                } else {
                    servers.insert(Token(counter), l);
                    counter += 1;
                }
            }
            Err(e) => println!("failed to start server at address: {}, {:?}", address.port, e),
        };
    }
    return counter;
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut server_info = ServerInfo::new();
    match server_info.deserialize() {
        Ok(_) => {
            println!("deserialization success");
            println!("{:#?}", server_info);
        }
        Err(e) => {
            println!("serverinfo error {e}");
            server_info.default();
        }
    }

    let mut poll = Poll::new()?;
    let mut servers: HashMap<Token, TcpListener> = HashMap::new();
    let mut sockets: HashMap<Token, TcpStream> = HashMap::new();
    let mut requests: HashMap<Token, Vec<u8>> = HashMap::new();
    let mut request_times: HashMap<Token, Instant> = HashMap::new();
    let mut request_server: HashMap<Token, usize> = HashMap::new();
    let mut buffer = [0_u8; 1024];
    let timeout_duration = Duration::from_secs(server_info.timeout_duration as u64); // Set your desired timeout duration

    let mut events = Events::with_capacity(1024);

    let c = register_servers(&server_info.ports, &mut poll, &mut servers);
    let mut counter = c.clone();

    loop {
        if poll.poll(&mut events, None).is_err() {
            continue;
        }

        for event in &events {
            match event.token() {
                e if e <= Token(c) =>
                    loop {
                        match servers.get(&event.token()).unwrap().accept() {
                            Ok((mut socket, _)) => {
                                counter += 1;
                                let token = Token(counter);
                                if
                                    socket
                                        .register(poll.registry(), token, Interest::READABLE)
                                        .is_err()
                                {
                                    println!("sockets register error");
                                    break;
                                }
                                sockets.insert(token, socket);
                                requests.insert(token, Vec::with_capacity(192));
                                request_server.insert(token, event.token().into());
                                request_times.insert(token, Instant::now());
                            }
                            Err(ref e) if e.kind() == io::ErrorKind::WouldBlock => {
                                break;
                            }
                            // Unexpected error
                            e => {
                                println!("unexpected error {:?}", e);
                                break;
                            }
                        }
                    }
                token if event.is_readable() => {
                    let elapsed_time = request_times[&token].elapsed();
                    if elapsed_time >= timeout_duration {
                        sockets.remove(&token);
                        requests.remove(&token);
                        request_times.remove(&token);
                        break;
                    }
                    loop {
                        let socket;
                        match sockets.get_mut(&token) {
                            Some(s) => {
                                socket = s;
                            }

                            None => {
                                println!("sockets error");
                                sockets.remove(&token);
                                requests.remove(&token);
                                request_times.remove(&token);
                                break;
                            }
                        }
                        let read = socket.read(&mut buffer);

                        match read {
                            Ok(0) => {
                                sockets.remove(&token);
                                requests.remove(&token);
                                request_times.remove(&token);
                                break;
                            }
                            Ok(n) => {
                                let req;
                                match requests.get_mut(&token) {
                                    Some(r) => {
                                        req = r;
                                    }
                                    _ => {
                                        let r = poll.registry().deregister(socket);
                                        if r.is_err() {
                                            println!("socket deregister error {:?}", r);
                                        }
                                        sockets.remove(&token);
                                        requests.remove(&token);
                                        request_times.remove(&token);
                                        break;
                                    }
                                }

                                for b in &buffer[0..n] {
                                    req.push(*b);
                                }
                            }
                            Err(ref e) if e.kind() == io::ErrorKind::WouldBlock => {
                                break;
                            }
                            Err(_) => {
                                break;
                            }
                        }
                    }
                    let r;
                    match requests.get(&token) {
                        Some(c) => {
                            r = c;
                        }
                        _ => {
                            println!("requests error");
                            sockets.remove(&token);
                            requests.remove(&token);
                            request_times.remove(&token);
                            continue;
                        }
                    }
                    let ready = r.windows(4).any(is_double_crnl);

                    if ready {
                        let socket;
                        match sockets.get_mut(&token) {
                            Some(s) => {
                                socket = s;
                            }
                            _ => {
                                println!("sockets error");
                                sockets.remove(&token);
                                requests.remove(&token);
                                request_times.remove(&token);
                                continue;
                            }
                        }

                        let req = requests[&token].clone();
                        let parts: Vec<String> = String::from_utf8_lossy(req.as_slice())
                            .split("\r\n")
                            .map(|s| s.to_owned())
                            .collect();

                        if parts.len() > 0 {
                            let first_line = &parts[0];
                            let method_parts: Vec<&str> = first_line.split_whitespace().collect();
                            let req_str: String = String::from_utf8_lossy(&req.clone()).into();

                            // Check for cookie
                            let cookie = check_cookie(&req_str);
                            let cookie_exists = match cookie {
                                Some(cookie_value) => {
                                    println!("Cookie found: {}", cookie_value);
                                    true
                                }
                                None => {
                                    println!("No cookie found");
                                    false
                                }
                            };

                            if method_parts.len() == 3 {
                                let method = method_parts[0];
                                let path = method_parts[1];
                                let body_start = req
                                    .windows(4)
                                    .position(|window| is_double_crnl(window))
                                    .unwrap_or(req.len());
                                let body = &req[body_start..];

                                let mut port = 0;
                                match request_server.get(&event.token()) {
                                    Some(num) => {
                                        port = *num;
                                    }
                                    None => {
                                        println!("no such port exists {}", port);
                                        continue;
                                    }
                                }
                                let s_info = &server_info.ports[port];
                                match method {
                                    "GET" => handle_get(socket, path, &s_info),
                                    "POST" =>
                                        handle_post(socket, path, body, &s_info, cookie_exists),
                                    "DELETE" => handle_delete(socket, path, &s_info, cookie_exists),
                                    "HEAD" =>
                                        handle_error(socket, 405, &server_info.error_page_path),
                                    "PUT" =>
                                        handle_error(socket, 405, &server_info.error_page_path),
                                    "CONNECT" =>
                                        handle_error(socket, 405, &server_info.error_page_path),
                                    "OPTIONS" =>
                                        handle_error(socket, 405, &server_info.error_page_path),
                                    "TRACE" =>
                                        handle_error(socket, 405, &server_info.error_page_path),
                                    "PATCH" =>
                                        handle_error(socket, 405, &server_info.error_page_path),
                                    _ => handle_error(socket, 400, &server_info.error_page_path),
                                }
                            } else {
                                handle_error(socket, 400, &server_info.error_page_path);
                                println!("invalid request");
                            }
                        } else {
                            println!("lol?");
                        }

                        // Clear the request buffer and send the response
                        requests
                            .get_mut(&token)
                            .get_or_insert(&mut Vec::with_capacity(192))
                            .clear();

                        if socket.reregister(poll.registry(), token, Interest::WRITABLE).is_err() {
                            let r = poll.registry().deregister(socket);
                            if r.is_err() {
                                println!("socket deregister error {:?}", r);
                            }
                            sockets.remove(&token);
                            requests.remove(&token);
                            request_times.remove(&token);
                        }
                    }

                    // Check for timeout
                    let elapsed_time = request_times[&token].elapsed();
                    if elapsed_time >= timeout_duration {
                        sockets.remove(&token);
                        requests.remove(&token);
                        request_times.remove(&token);
                    }
                }
                token if event.is_writable() => {
                    let page = get_page("/", &server_info.error_page_path);
                    let socket: &mut TcpStream;
                    match sockets.get_mut(&token) {
                        Some(s) => {
                            socket = s;
                        }
                        _ => {
                            println!("sockets error");
                            sockets.remove(&token);
                            requests.remove(&token);
                            request_times.remove(&token);
                            continue;
                        }
                    }
                    let r = socket.write_all(page.as_bytes());
                    if r.is_err() {
                        println!("got error lol");
                        sockets.remove(&token);
                        requests.remove(&token);
                        request_times.remove(&token);
                        continue;
                    }
                    socket.reregister(poll.registry(), token, Interest::READABLE)?;

                    // Remove token/socket after handling
                    let r = poll.registry().deregister(socket);
                    if r.is_err() {
                        println!("socket deregister error {:?}", r);
                    }
                    sockets.remove(&token);
                    requests.remove(&token);
                    request_times.remove(&token);
                }
                e => println!("unreachable code lol, {:?}", e),
            };
        }
    }
}
