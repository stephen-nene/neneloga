// fn main2() {
//     println!("Hello, world!");
// }


use std::net::TcpListener;

fn main() {
    let listener = TcpListener::bind("127.0.0.1:7878").unwrap();

    for stream in listener.incoming() {
        let stream = stream.unwrap();
        print!("{:?}",stream);

        println!("Connection established!");
        return "Connection established!"
    }
}
